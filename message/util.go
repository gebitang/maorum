package message

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func isNumber(u string) bool {
	_, err := strconv.Atoi(u)
	return err == nil
}

func readNetworkAsset(ctx context.Context, name string) (*bot.Asset, error) {
	body, err := bot.Request(ctx, "GET", "/network/assets/"+name, nil, "")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data  *bot.Asset `json:"data"`
		Error bot.Error  `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return resp.Data, nil
}

func isHelpInfo(info string) bool {
	for _, v := range helpMap {
		if v == strings.TrimSpace(info) {
			return true
		}
	}
	return false
}

func uploadAttachmentTo(uploadURL string, file []byte) error {
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(file))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("x-amz-acl", "public-read")
	req.Header.Add("Content-Length", strconv.Itoa(len(file)))

	resp, err := httpClient.Do(req)
	if resp != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	return nil
}

func toMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func millisecondToTime(millis int64) time.Time {
	millisecond := int64(time.Millisecond)
	return time.Unix(0, millis*millisecond)
}
