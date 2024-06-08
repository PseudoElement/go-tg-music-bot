package bot_manager

type BotClient struct {
	IsFirstLoad bool
	//"FIND_SIMILAR_SONGS" or "FIND_SONG_BY_KEYWORDS"
	QueryType string
	//"SEND_LIST_WITH_LINKS" or "SEND_TEXT_LIST"
	ResponseViewType string
	UserName         string
}
