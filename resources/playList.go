package resources

type PlayList struct {
	ID      string   `json:"id"`
	UserID  string   `json:"user_id"`
	SongIDs []string `json:"song_ids"`
}
