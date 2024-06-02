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
)

var botKeys = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Найти похожие песни"),
		tgbotapi.NewKeyboardButton("Найти песню по ключевым словам"),
	),
)

var keyboardCommandsRusToEng = map[string]string{
	"Найти похожие песни":            FIND_SIMILAR_SONGS,
	"Найти песню по ключевым словам": FIND_SONG_BY_KEYWORDS,
}
