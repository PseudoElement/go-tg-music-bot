package bot_manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	SHAZAM_API_SERVICE = "SHAZAM_API_SERVICE"
	CHAT_GPT_SERVICE   = "CHAT_GPT_SERVICE"
)

const (
	YOUTUBE_LINK_SEARCHER = "YOUTUBE_LINK_SEARCHER"
	SPOTIFY_LINK_SEARCHER = "SPOTIFY_LINK_SEARCHER"
)

const (
	FIND_SIMILAR_SONGS    = "FIND_SIMILAR_SONGS"
	FIND_SONG_BY_KEYWORDS = "FIND_SONG_BY_KEYWORDS"
	SEND_TEXT_LIST        = "SEND_TEXT_LIST"
	SEND_LIST_WITH_LINKS  = "SEND_LIST_WITH_LINKS"
	CONTACT_ADMIN         = "CONTACT_ADMIN"
)

const (
	STAGE_MAIN_COMMAND_SELECTION = "STAGE_MAIN_COMMAND_SELECTION"
	STAGE_LIST_TYPE_SELECTION    = "STAGE_LIST_TYPE_SELECTION"
	STAGE_SONG_NAME_INPUT        = "STAGE_SONG_NAME_INPUT"
)

const (
	MAIN_COMMAND              = "MAIN_COMMAND"
	RESPONSE_LIST_VIEW        = "RESPONSE_LIST_VIEW"
	COMMAND_REQUIRED_KEYBOARD = "COMMAND_REQUIRED_KEYBOARD"
)

var (
	MAIN_OPTIONS_KEYBOARD = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Найти похожие песни"),
			tgbotapi.NewKeyboardButton("Найти песню по ключевым словам"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Связаться с администратором"),
		),
	)
	LIST_TYPE_KEYBOARD = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отправить обычный список"),
			tgbotapi.NewKeyboardButton("Отправить список с ссылками на YouTube"),
		),
	)
)

var COMMAND_TYPES = map[string][]string{
	FIND_SIMILAR_SONGS:    {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	FIND_SONG_BY_KEYWORDS: {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	CONTACT_ADMIN:         {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	SEND_LIST_WITH_LINKS:  {RESPONSE_LIST_VIEW},
	SEND_TEXT_LIST:        {RESPONSE_LIST_VIEW},
}

var COMMANDS_RUS_TO_ENG = map[string]string{
	"Найти похожие песни":                    FIND_SIMILAR_SONGS,
	"Найти песню по ключевым словам":         FIND_SONG_BY_KEYWORDS,
	"Отправить обычный список":               SEND_TEXT_LIST,
	"Отправить список с ссылками на YouTube": SEND_LIST_WITH_LINKS,
	"Связаться с администратором":            CONTACT_ADMIN,
}

var COMMANDS_LIST = []string{
	FIND_SIMILAR_SONGS,
	FIND_SONG_BY_KEYWORDS,
	SEND_LIST_WITH_LINKS,
	SEND_TEXT_LIST,
	CONTACT_ADMIN,
}

var RESPONSE_MESSAGES_FOR_COMMAND = map[string]string{
	FIND_SIMILAR_SONGS:    "Введи название песни",
	FIND_SONG_BY_KEYWORDS: "Введи ключевые слова из песни(отрывок из текста, часть названия и т.д.)",
	SEND_LIST_WITH_LINKS:  "Какой список тебе отправить?",
	SEND_TEXT_LIST:        "Какой список тебе отправить?",
	CONTACT_ADMIN:         "Telegram администратора: \n - https://t.me/BigChad",
}

var KEYBOARDS_FOR_MAIN_COMMAND = map[string]tgbotapi.ReplyKeyboardMarkup{
	FIND_SIMILAR_SONGS:    LIST_TYPE_KEYBOARD,
	FIND_SONG_BY_KEYWORDS: LIST_TYPE_KEYBOARD,
	CONTACT_ADMIN:         MAIN_OPTIONS_KEYBOARD,
}

var NEW_STAGE_ON_COMMAND = map[string]string{
	FIND_SIMILAR_SONGS:    STAGE_LIST_TYPE_SELECTION,
	FIND_SONG_BY_KEYWORDS: STAGE_LIST_TYPE_SELECTION,
	SEND_LIST_WITH_LINKS:  STAGE_SONG_NAME_INPUT,
	SEND_TEXT_LIST:        STAGE_SONG_NAME_INPUT,
	CONTACT_ADMIN:         STAGE_MAIN_COMMAND_SELECTION,
}
