package domain

type Comment struct {
	Id             int64           `json:"comment_id"`
	VId            int64           `json:"video_id"`
	User           User            `json:"user"`
	Content        string          `json:"content"`
	Ctime          int64           `json:"ctime"`
	Like           int64           `json:"like"`
	Dislike        int64           `json:"dislike"`
	ReplyCount     int64           `json:"reply_count"`
	Picture        string          `json:"picture"`
	MentionedUsers []MentionedUser `json:"mentioned_users"`
}
