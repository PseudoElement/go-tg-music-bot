package bot_manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	SHAZAM_API_SERVICE = "SHAZAM_API_SERVICE"
	CHAT_GPT_SERVICE   = "CHAT_GPT_SERVICE"
)

const (
	FIND_SIMILAR_SONGS    = "FIND_SIMILAR_SONGS"
	FIND_SONG_BY_KEYWORDS = "FIND_SONG_BY_KEYWORDS"
	SEND_TEXT_LIST        = "SEND_TEXT_LIST"
	SEND_LIST_WITH_LINKS  = "SEND_LIST_WITH_LINKS"
)

var (
	MAIN_OPTIONS_KEYBOARD = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Найти похожие песни"),
			tgbotapi.NewKeyboardButton("Найти песню по ключевым словам"),
		),
	)
	LIST_TYPE_KEYBOARD = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отправить обычный список"),
			tgbotapi.NewKeyboardButton("Отправить список с ссылками на YouTube"),
		),
	)
	EMPTY_KEYBOARD = tgbotapi.NewReplyKeyboard()
)

var COMMANDS_RUS_TO_ENG = map[string]string{
	"Найти похожие песни":                    FIND_SIMILAR_SONGS,
	"Найти песню по ключевым словам":         FIND_SONG_BY_KEYWORDS,
	"Отправить обычный список":               SEND_TEXT_LIST,
	"Отправить список с ссылками на YouTube": SEND_LIST_WITH_LINKS,
}

var COMMANDS_LIST = []string{
	FIND_SIMILAR_SONGS,
	FIND_SONG_BY_KEYWORDS,
	SEND_LIST_WITH_LINKS,
	SEND_TEXT_LIST,
}

var RESPONSE_MESSAGES_ON_COMMAND = map[string]string{
	FIND_SIMILAR_SONGS:    "Введи название песни",
	FIND_SONG_BY_KEYWORDS: "Введи ключевые слова из песни(отрывок из текста, часть названия и т.д.)",
	SEND_LIST_WITH_LINKS:  "Какой список тебе отправить?",
	SEND_TEXT_LIST:        "Какой список тебе отправить?",
}

var KEYBOARDS_ON_COMMAND = map[string]tgbotapi.ReplyKeyboardMarkup{
	FIND_SIMILAR_SONGS:    LIST_TYPE_KEYBOARD,
	FIND_SONG_BY_KEYWORDS: LIST_TYPE_KEYBOARD,
	SEND_LIST_WITH_LINKS:  EMPTY_KEYBOARD,
	SEND_TEXT_LIST:        EMPTY_KEYBOARD,
}
