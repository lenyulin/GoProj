package cache

type Video struct {
	Title    string `gorm:"varchar(255)" bson:"title"`
	Content  string `gorm:"varchar(4096)" bson:"content"`
	Id       int64  `gorm:"primaryKey" bson:"id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	AuthorId int64  `gorm:"index" bson:"authorId,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
