package bot_manager

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/go-tg-music-bot/ai"
	shazam_api "github.com/pseudoelement/go-tg-music-bot/shazam-api"
	"github.com/pseudoelement/go-tg-music-bot/types"
	"github.com/pseudoelement/go-tg-music-bot/utils"
)

type BotManager struct {
	bot              *tgbotapi.BotAPI
	updates          tgbotapi.UpdatesChannel
	musicApiServices map[string]types.MusicApiService
	useChatGPT       bool
	clients          map[int64]*BotClient
}

func NewBotManager() *BotManager {
	botManager := &BotManager{}
	botManager.init()
	return botManager
}

func (bm *BotManager) Bot() *tgbotapi.BotAPI {
	return bm.bot
}

func (bm *BotManager) init() {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		panic("BOT_TOKEN doesn't exist!")
	}

	useChatGPT := bm.needUseChatGPT()
	bm.useChatGPT = useChatGPT

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bm.bot = bot

	bot.Debug = true
	fmt.Printf("Authorized on account %s!\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 100
	bm.updates = bot.GetUpdatesChan(u)

	var chatGPT *ai.ChatGPT
	if useChatGPT {
		chatGPT, err = ai.NewChatGPTService()
		if err != nil {
			panic(err)
		}
	}

	shazamApi, err := shazam_api.NewShazamApiService()
	if err != nil {
		panic(err)
	}

	bm.musicApiServices = map[string]types.MusicApiService{
		SHAZAM_API_SERVICE: shazamApi,
		CHAT_GPT_SERVICE:   chatGPT,
	}
	bm.clients = make(map[int64]*BotClient)
}

func (bm *BotManager) needUseChatGPT() bool {
	useAiEnv := os.Getenv("USE_AI")
	var useChatGPT bool
	if len(useAiEnv) == 0 {
		useChatGPT = false
	} else {
		boolean, err := strconv.ParseBool(useAiEnv)
		if err != nil {
			fmt.Println("Incorrect format of .env var USE_AI!")
			useChatGPT = false
		}
		useChatGPT = boolean
	}

	return useChatGPT
}

func (bm *BotManager) Broadcast() {
	for update := range bm.updates {
		// fmt.Println("Video from TG: ", update.Message.Video)
		// fmt.Println("Photo from TG: ", update.Message.Photo)
		userId := update.Message.From.ID
		userName := update.Message.From.UserName
		bm.handleClientsConfig(userId, userName)
		isNewUser := bm.clients[userId].IsFirstLoad
		if update.Message != nil {
			if isNewUser {
				bm.sendGreetingMessage(update, userId)
				continue
			}
			var msg tgbotapi.MessageConfig
			if bm.isKeyboardCommand(update.Message.Text) {
				msg = bm.handleKeyboardCommand(update)
			} else {
				actionType := bm.clients[userId].ActionType
				msg = bm.handleQuery(update, actionType)
			}
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = botKeys
			bm.Bot().Send(msg)
		}
	}
}

func (bm *BotManager) isKeyboardCommand(text string) bool {
	return keyboardCommandsRusToEng[text] == FIND_SIMILAR_SONGS || keyboardCommandsRusToEng[text] == FIND_SONG_BY_KEYWORDS
}

func (bm *BotManager) sendGreetingMessage(update tgbotapi.Update, userId int64) {
	userName := bm.clients[userId].UserName
	text := fmt.Sprintf(`
        Привет, %s. Выбирай интересующую тебя опцию внизу в меню. Я попробую помочь с твоим запросом!
    `, userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bm.clients[userId].IsFirstLoad = false
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = botKeys
	bm.Bot().Send(msg)
}

func (bm *BotManager) handleClientsConfig(userId int64, userName string) {
	_, ok := bm.clients[userId]
	if !ok {
		bm.clients[userId] = &BotClient{
			IsFirstLoad: true,
			ActionType:  FIND_SIMILAR_SONGS,
			UserName:    userName,
		}
	}
}

func (bm *BotManager) handleKeyboardCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch keyboardCommandsRusToEng[update.Message.Text] {
	case FIND_SIMILAR_SONGS:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название песни")
		bm.clients[update.Message.From.ID].ActionType = FIND_SIMILAR_SONGS
	case FIND_SONG_BY_KEYWORDS:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ключевые слова из песни(отрывок из текста, часть названия и т.д.)")
		bm.clients[update.Message.From.ID].ActionType = FIND_SONG_BY_KEYWORDS
	default:
		fmt.Println("Keyboard command not found!")
	}
	return msg
}

func (bm *BotManager) handleQuery(update tgbotapi.Update, actionType string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var response string
	var err error
	if actionType == FIND_SIMILAR_SONGS {
		if bm.useChatGPT {
			response, err = bm.musicApiServices[CHAT_GPT_SERVICE].QuerySimilarSongs(update.Message.Text, false)
		} else {
			response, err = bm.musicApiServices[SHAZAM_API_SERVICE].QuerySimilarSongs(update.Message.Text, false)
		}
	} else if actionType == FIND_SONG_BY_KEYWORDS {
		err = utils.MethodNotImplemented()
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Некорректный запрос, попробуй еще раз! Текст ошибки(для разработчика): %s", err.Error())
		if strings.HasPrefix(err.Error(), utils.MethodNotImplemented().Error()) {
			errorMsg = "На данный момент опция `Найти песню по ключевым слова` не добавлена в приложение Musician-bot! Но скоро все будет :)"
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
	}

	return msg
}
