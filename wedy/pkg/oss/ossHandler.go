package oss

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
)

const Dst = "C:\\Users\\lyl69\\GolandProjects\\GoProj\\wedy\\tmp\\"

var (
	ErrUploadToOSSFailed = errors.New("upload to oss failed")
)

type OssHandler struct {
	oss *cos.Client
}

func NewOSSHandler(oss *cos.Client) *OssHandler {
	return &OssHandler{oss: oss}
}

func (hdl *OssHandler) Upload(ctx context.Context, uid int64) error {
	filepath := fmt.Sprintf("%s%d.mp4", Dst, uid)
	_, err := hdl.oss.Object.PutFromFile(context.Background(), fmt.Sprintf("%d.mp4", uid), filepath, nil)
	if err != nil {
		return ErrUploadToOSSFailed
	}
	return nil
}

func (hdl *OssHandler) Find(ctx context.Context, uid int64) error {
	return nil
}
func (hdl *OssHandler) Delete(ctx context.Context, uid int64) error {
	return nil
}
