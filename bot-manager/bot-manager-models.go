package bot_manager

type BotClient struct {
	IsFirstLoad bool
	//"FIND_SIMILAR_SONGS" or "FIND_SONG_BY_KEYWORDS"
	QueryType string
	NeedLinks bool
	UserName  string
}
