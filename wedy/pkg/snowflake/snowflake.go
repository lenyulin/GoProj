package snowflake

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrGenerateSnowFlakeError = errors.New("Generate SnowFlake Number Error")
)

type Generator struct {
	AccessToken string
}

func Generate() (int64, error) {
	var tk Generator
	tk.AccessToken = "shiwozaifangwenya"
	tkJson, err := json.Marshal(tk)
	url := "http://14.103.175.18:12345/tools/snackFlow/generate"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tkJson))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return -1, ErrGenerateSnowFlakeError
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) == "Internal Server Error" {
		return -1, ErrGenerateSnowFlakeError
	}
	var tmp int64
	err = binary.Read(bytes.NewBuffer(body), binary.BigEndian, &tmp)
	if err != nil {
		return -1, err
	}
	return tmp, nil
}
