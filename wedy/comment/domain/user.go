package domain

type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type MentionedUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ReplyToUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
