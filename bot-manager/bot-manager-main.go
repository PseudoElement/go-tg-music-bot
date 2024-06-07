package bot_manager

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/go-tg-music-bot/ai"
	shazam_api "github.com/pseudoelement/go-tg-music-bot/shazam-api"
	"github.com/pseudoelement/go-tg-music-bot/types"
	"github.com/pseudoelement/go-tg-music-bot/utils"
)

type BotManager struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	//SHAZAM_API_SERVICE or CHAT_GPT_SERVICE
	activeMusicService string
	musicApiServices   map[string]types.MusicApiService
	clients            map[int64]*BotClient
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
	bm.selectMusicService(useChatGPT)

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

func (bm *BotManager) selectMusicService(useChatGPT bool) {
	if useChatGPT {
		bm.activeMusicService = CHAT_GPT_SERVICE
	} else {
		bm.activeMusicService = SHAZAM_API_SERVICE
	}
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
				user := bm.clients[userId]
				msg = bm.handleQuery(update, user)
				msg.ReplyMarkup = MAIN_OPTIONS_KEYBOARD
			}
			msg.ReplyToMessageID = update.Message.MessageID
			bm.Bot().Send(msg)
		}
	}
}

func (bm *BotManager) isKeyboardCommand(text string) bool {
	return utils.Includes(COMMANDS_LIST, COMMANDS_RUS_TO_ENG[text])
}

func (bm *BotManager) sendGreetingMessage(update tgbotapi.Update, userId int64) {
	userName := bm.clients[userId].UserName
	text := fmt.Sprintf(`
        Привет, %s. Выбирай интересующую тебя опцию внизу в меню. Я попробую помочь с твоим запросом!
    `, userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bm.clients[userId].IsFirstLoad = false
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = MAIN_OPTIONS_KEYBOARD
	bm.Bot().Send(msg)
}

func (bm *BotManager) handleClientsConfig(userId int64, userName string) {
	_, ok := bm.clients[userId]
	if !ok {
		bm.clients[userId] = &BotClient{
			IsFirstLoad: true,
			NeedLinks:   false,
			QueryType:   FIND_SIMILAR_SONGS,
			UserName:    userName,
		}
	}
}

func (bm *BotManager) handleKeyboardCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch COMMANDS_RUS_TO_ENG[update.Message.Text] {
	case FIND_SIMILAR_SONGS:
		text := RESPONSE_MESSAGES_ON_COMMAND[FIND_SIMILAR_SONGS]
		keyboard := KEYBOARDS_ON_COMMAND[FIND_SIMILAR_SONGS]
		bm.clients[update.Message.From.ID].QueryType = FIND_SIMILAR_SONGS
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
	case FIND_SONG_BY_KEYWORDS:
		text := RESPONSE_MESSAGES_ON_COMMAND[FIND_SONG_BY_KEYWORDS]
		keyboard := KEYBOARDS_ON_COMMAND[FIND_SONG_BY_KEYWORDS]
		bm.clients[update.Message.From.ID].QueryType = FIND_SONG_BY_KEYWORDS
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
	case SEND_LIST_WITH_LINKS:
		mainCommandSelected := bm.clients[update.Message.From.ID].QueryType
		text := RESPONSE_MESSAGES_ON_COMMAND[mainCommandSelected]
		bm.clients[update.Message.From.ID].NeedLinks = true
		keyboard := KEYBOARDS_ON_COMMAND[SEND_LIST_WITH_LINKS]
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
	case SEND_TEXT_LIST:
		mainCommandSelected := bm.clients[update.Message.From.ID].QueryType
		text := RESPONSE_MESSAGES_ON_COMMAND[mainCommandSelected]
		bm.clients[update.Message.From.ID].NeedLinks = false
		keyboard := KEYBOARDS_ON_COMMAND[SEND_LIST_WITH_LINKS]
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
	default:
		fmt.Println("Keyboard command not found!")
	}
	return msg
}

func (bm *BotManager) handleQuery(update tgbotapi.Update, user *BotClient) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var response string
	var err error
	if user.QueryType == FIND_SIMILAR_SONGS {
		response, err = bm.musicApiServices[bm.activeMusicService].QuerySimilarSongs(update.Message.Text, false, user.NeedLinks)
		response = "Вот подборка из 20 похожих песен: \n" + response
	} else if user.QueryType == FIND_SONG_BY_KEYWORDS {
		response, err = bm.musicApiServices[bm.activeMusicService].QuerySongByKeyWords(update.Message.Text, user.NeedLinks)
		response = "Наиболее вероятные совпадения по запросу: \n" + response
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Произошла ошибка, попробуй еще раз! Текст ошибки(для разработчика): %s", err.Error())
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
	}

	return msg
}
