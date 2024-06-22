package bot_manager

type BotClient struct {
	IsFirstLoad bool
	//"FIND_SIMILAR_SONGS" or "FIND_SONG_BY_KEYWORDS" or "CONTACT_ADMIN"
	MainCommandSelected string
	//"SEND_LIST_WITH_LINKS" or "SEND_TEXT_LIST"
	ResponseViewType string
	UserName         string
	//STAGE_QUERY_TYPE_SELECTION or STAGE_LIST_TYPE_SELECTION or STAGE_SONG_NAME_INPUT
	Stage string
}
