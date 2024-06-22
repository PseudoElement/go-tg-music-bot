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
		isCommandRequired := bm.clients[userId].Stage != STAGE_SONG_NAME_INPUT
		if update.Message != nil {
			if isNewUser {
				bm.sendGreetingMessage(update, userId)
				continue
			}
			var msg tgbotapi.MessageConfig
			if bm.isKeyboardCommand(update.Message.Text) {
				msg = bm.handleKeyboardCommand(update)
			} else if isCommandRequired {
				msg = bm.handleCommandRequiredWarning(update)
			} else {
				user := bm.clients[userId]
				msg = bm.handleQuery(update, user)
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
			IsFirstLoad:         true,
			MainCommandSelected: FIND_SIMILAR_SONGS,
			ResponseViewType:    SEND_TEXT_LIST,
			UserName:            userName,
			Stage:               STAGE_QUERY_TYPE_SELECTION,
		}
	}
}

func (bm *BotManager) handleCommandRequiredWarning(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо выбрать одну из доступных команд!")
	return msg
}

func (bm *BotManager) handleKeyboardCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch COMMANDS_RUS_TO_ENG[update.Message.Text] {
	case FIND_SIMILAR_SONGS:
		msg = bm.getResponseMessage(update, FIND_SIMILAR_SONGS)
	case FIND_SONG_BY_KEYWORDS:
		msg = bm.getResponseMessage(update, FIND_SONG_BY_KEYWORDS)
	case SEND_LIST_WITH_LINKS:
		msg = bm.getResponseMessage(update, SEND_LIST_WITH_LINKS)
	case SEND_TEXT_LIST:
		msg = bm.getResponseMessage(update, SEND_TEXT_LIST)
	case CONTACT_ADMIN:
		msg = bm.getResponseMessage(update, CONTACT_ADMIN)
	default:
		fmt.Println("Keyboard command not found!")
	}
	return msg
}

func (bm *BotManager) getResponseMessage(update tgbotapi.Update, command string) tgbotapi.MessageConfig {
	needKeyboard := utils.Includes(COMMAND_TYPES[command], COMMAND_REQUIRED_KEYBOARD)
	needSetResponseListView := utils.Includes(COMMAND_TYPES[command], RESPONSE_LIST_VIEW)
	isMainCommand := utils.Includes(COMMAND_TYPES[command], MAIN_COMMAND)

	text := RESPONSE_MESSAGES_FOR_COMMAND[command]
	bm.clients[update.Message.From.ID].Stage = NEW_STAGE_ON_COMMAND[command]
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	if isMainCommand {
		bm.clients[update.Message.From.ID].MainCommandSelected = command
	}
	if needKeyboard {
		mainCommandSelected := bm.clients[update.Message.From.ID].MainCommandSelected
		keyboard := KEYBOARDS_FOR_MAIN_COMMAND[mainCommandSelected]
		msg.ReplyMarkup = keyboard
	}
	if needSetResponseListView {
		bm.clients[update.Message.From.ID].ResponseViewType = command
	}

	return msg
}

func (bm *BotManager) handleQuery(update tgbotapi.Update, user *BotClient) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var response string
	var err error
	if user.MainCommandSelected == FIND_SIMILAR_SONGS {
		switch user.ResponseViewType {
		case SEND_TEXT_LIST:
			response, err = bm.musicApiServices[bm.activeMusicService].QuerySimilarSongs(update.Message.Text, false)
		case SEND_LIST_WITH_LINKS:
			response, err = bm.musicApiServices[bm.activeMusicService].QuerySimilarSongsLinks(update.Message.Text)
		}
		response = "Вот подборка похожих песен: \n" + response
	} else if user.MainCommandSelected == FIND_SONG_BY_KEYWORDS {
		switch user.ResponseViewType {
		case SEND_TEXT_LIST:
			response, err = bm.musicApiServices[bm.activeMusicService].QuerySongByKeyWords(update.Message.Text)
		case SEND_LIST_WITH_LINKS:
			response, err = bm.musicApiServices[bm.activeMusicService].QuerySongByKeyWordsLinks(update.Message.Text)
		}
		response = "Наиболее вероятные совпадения по запросу: \n" + response
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Произошла ошибка, попробуй еще раз! Текст ошибки(для разработчика): %s", err.Error())
		if utils.IndexOfSubstring(err.Error(), utils.SimilarSongsNotFound().Error()) != -1 {
			errorMsg = err.Error()
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
	}

	msg.ReplyMarkup = MAIN_OPTIONS_KEYBOARD
	return msg
}
