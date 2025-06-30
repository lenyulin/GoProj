package web

type VideoVO struct {
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	AuthorId   int64  `json:"author_id,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	Id         int64  `json:"id,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`
}
