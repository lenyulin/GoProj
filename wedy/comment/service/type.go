package service

import "errors"

type Comment struct {
	Id      int64  `json:"id"`
	Vid     int64  `json:"vid"`
	Content string `json:"content"`
	User    int64  `json:"user"`
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

var (
	ErrGenerateSnowFlakeError = errors.New("Generate SnowFlake Number Error")
)
