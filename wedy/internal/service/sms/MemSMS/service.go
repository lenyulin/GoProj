package MemSMS

import (
	"context"
	"fmt"
)

type MemSMS struct {
	appId    string
	signName string
}

func NewMemSMS(appId, signName string) *MemSMS {
	return &MemSMS{
		appId:    appId,
		signName: signName,
	}
}
func (m *MemSMS) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for i, number := range numbers {
		fmt.Println(fmt.Sprintf("Sending:%s, code:%s ", number, args[i]))
	}
	return nil
}
