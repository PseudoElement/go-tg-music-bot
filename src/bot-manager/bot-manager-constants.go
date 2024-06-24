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
	FIND_SIMILAR_SONGS            = "FIND_SIMILAR_SONGS"
	FIND_SONG_BY_KEYWORDS         = "FIND_SONG_BY_KEYWORDS"
	SEND_TEXT_LIST                = "SEND_TEXT_LIST"
	SEND_LIST_WITH_LINKS          = "SEND_LIST_WITH_LINKS"
	SEND_LIST_WITH_LINKS_EXTENDED = "SEND_LIST_WITH_LINKS_EXTENDED"
	SEND_SPOTIFY_LINKS            = "SEND_SPOTIFY_LINKS"
	SEND_YOUTUBE_LINKS            = "SEND_YOUTUBE_LINKS"
	CONTACT_ADMIN                 = "CONTACT_ADMIN"
)

const (
	STAGE_MAIN_COMMAND_SELECTION         = "STAGE_MAIN_COMMAND_SELECTION"
	STAGE_LIST_TYPE_SELECTION            = "STAGE_LIST_TYPE_SELECTION"
	STAGE_SONG_NAME_INPUT                = "STAGE_SONG_NAME_INPUT"
	STAGE_LAST_COMMAND_BEFORE_SONG_INPUT = "STAGE_LAST_COMMAND_BEFORE_SONG_INPUT"
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
			tgbotapi.NewKeyboardButton("Отправить список со ссылками"),
		),
	)
	LIST_TYPE_KEYBOARD_EXTENDED = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отправить ссылки на Spotify"),
			tgbotapi.NewKeyboardButton("Отправить ссылки на YouTube"),
		),
	)
)

var COMMAND_TYPES = map[string][]string{
	FIND_SIMILAR_SONGS:            {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	FIND_SONG_BY_KEYWORDS:         {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	CONTACT_ADMIN:                 {MAIN_COMMAND, COMMAND_REQUIRED_KEYBOARD},
	SEND_LIST_WITH_LINKS:          {RESPONSE_LIST_VIEW},
	SEND_LIST_WITH_LINKS_EXTENDED: {COMMAND_REQUIRED_KEYBOARD, RESPONSE_LIST_VIEW},
	SEND_TEXT_LIST:                {RESPONSE_LIST_VIEW},
}

var COMMANDS_RUS_TO_ENG = map[string]string{
	"Найти похожие песни":            FIND_SIMILAR_SONGS,
	"Найти песню по ключевым словам": FIND_SONG_BY_KEYWORDS,
	"Отправить обычный список":       SEND_TEXT_LIST,
	"Отправить список со ссылками":   SEND_LIST_WITH_LINKS,
	"Отправить ссылки на YouTube":    SEND_YOUTUBE_LINKS,
	"Отправить ссылки на Spotify":    SEND_SPOTIFY_LINKS,
	"Связаться с администратором":    CONTACT_ADMIN,
}

var COMMANDS_LIST = []string{
	FIND_SIMILAR_SONGS,
	FIND_SONG_BY_KEYWORDS,
	SEND_LIST_WITH_LINKS,
	SEND_TEXT_LIST,
	SEND_YOUTUBE_LINKS,
	SEND_SPOTIFY_LINKS,
	SEND_LIST_WITH_LINKS_EXTENDED,
	CONTACT_ADMIN,
}

var RESPONSE_MESSAGES_FOR_COMMAND = map[string]string{
	FIND_SIMILAR_SONGS:            "Какой список тебе отправить?",
	FIND_SONG_BY_KEYWORDS:         "Какой список тебе отправить?",
	SEND_LIST_WITH_LINKS:          "Введи название песни или ключевые слова",
	SEND_TEXT_LIST:                "Введи название песни или ключевые слова",
	SEND_SPOTIFY_LINKS:            "Введи название песни или ключевые слова",
	SEND_YOUTUBE_LINKS:            "Введи название песни или ключевые слова",
	SEND_LIST_WITH_LINKS_EXTENDED: "С какими ссылками тебе отправить список?",
	CONTACT_ADMIN:                 "Telegram администратора: \n - https://t.me/BigChad",
}

var KEYBOARDS_FOR_COMMANDS = map[string]tgbotapi.ReplyKeyboardMarkup{
	FIND_SIMILAR_SONGS:            LIST_TYPE_KEYBOARD,
	FIND_SONG_BY_KEYWORDS:         LIST_TYPE_KEYBOARD,
	CONTACT_ADMIN:                 MAIN_OPTIONS_KEYBOARD,
	SEND_LIST_WITH_LINKS_EXTENDED: LIST_TYPE_KEYBOARD_EXTENDED,
}

var NEW_STAGE_AFTER_COMMAND = map[string]string{
	FIND_SIMILAR_SONGS:            STAGE_LIST_TYPE_SELECTION,
	FIND_SONG_BY_KEYWORDS:         STAGE_LIST_TYPE_SELECTION,
	SEND_LIST_WITH_LINKS_EXTENDED: STAGE_LIST_TYPE_SELECTION,
	SEND_LIST_WITH_LINKS:          STAGE_SONG_NAME_INPUT,
	SEND_TEXT_LIST:                STAGE_SONG_NAME_INPUT,
	SEND_YOUTUBE_LINKS:            STAGE_SONG_NAME_INPUT,
	SEND_SPOTIFY_LINKS:            STAGE_SONG_NAME_INPUT,
	CONTACT_ADMIN:                 STAGE_MAIN_COMMAND_SELECTION,
}
