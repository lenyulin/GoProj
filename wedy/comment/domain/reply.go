package domain

type Reply struct {
	Id             int64           `json:"comment_id"`
	Content        string          `json:"content"`
	ReplyTo        int64           `json:"parent_comment_id"`
	Ctime          int64           `json:"ctime"`
	Like           int64           `json:"like"`
	Dislike        int64           `json:"dislike"`
	Picture        string          `json:"picture"`
	ReplyToUser    ReplyToUser     `json:"reply_user"`
	MentionedUsers []MentionedUser `json:"mentioned_users"`
}
