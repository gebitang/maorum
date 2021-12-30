package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gebitang.com/maorum/bubble"
	"gebitang.com/maorum/config"
	"gebitang.com/maorum/models"
	"gebitang.com/maorum/models/dbm"
	"gebitang.com/maorum/rum"
	"gebitang.com/maorum/tomato"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (t *TrainClient) drawMaoPic(ctx context.Context, userId, content string) {
	itemMap := queryItems(ctx)
	if len(itemMap) == 0 {
		t.sendTextMsg(ctx, userId, "抱歉，今天还没有定投")
		return
	}
	b, str, bb := bubble.BuildBubbleByte(itemMap)
	if b {
		t.sendTextMsg(ctx, userId, str)
		t.sendPlainImage(ctx, userId, "image/png", bb, 750, 450)
	}

}

func (t *TrainClient) readLatestRum(ctx context.Context, userId, content string) {
	if len(config.RumReadGroup) == 0 {
		t.sendTextMsg(ctx, userId, "抱歉，没有配置rum.read.group.id的信息")
		return
	}
	b, s := rum.ReadFromGroup(1)
	if !b {
		t.sendTextMsg(ctx, userId, fmt.Sprintf("获取失败 %s", s))
		return
	}
	g := make([]rum.ContentItem, 0)
	json.Unmarshal([]byte(s), &g)
	for _, c := range g {
		fmt.Println("------------------", c.Publisher, time.Unix(0, c.TimeStamp).Format(TimeFormat), c.TrxId, c.TypeUrl)
		t.sendTextMsg(ctx, userId, c.Content.Content)
		for i, img := range c.Content.Image {
			decodeString, _ := base64.StdEncoding.DecodeString(img.Content)
			if Debug {
				fmt.Println(i, " imgInfo ", img.MediaType, img.Name, "imgStr=", len(img.Content), " imgBye=", len(decodeString))
			}
			if len(decodeString) > 1024*100 {
				t.sendTextMsg(ctx, userId, "图片太大，请到Rum群组中查看")
			} else {
				t.sendPlainImage(ctx, userId, img.MediaType, decodeString, 500, 300)
			}

		}
	}
}

func (t *TrainClient) drinkRum(ctx context.Context, userId, content string) {
	itemMap := queryItems(ctx)
	if len(itemMap) == 0 {
		t.sendTextMsg(ctx, userId, "抱歉，今天还没有定投")
		return
	}
	if len(config.RumPostGroup) == 0 {
		t.sendTextMsg(ctx, userId, "抱歉，没有配置rum.post.group.id的信息")
		return
	}

	b, str, bb := bubble.BuildBubbleByte(itemMap)
	if b {
		r, s := rum.PostToGroup(bb, str)
		if r {
			t.sendTextMsg(ctx, userId, "发送成功，请稍后查看Rum上对应的群组信息")
		} else {
			t.sendTextMsg(ctx, userId, "发送失败，"+s)
		}
	}

}

func (t *TrainClient) sendOneItems(ctx context.Context, userId, content string) {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	for _, i := range bubble.ItemMap {
		if content == i.CommandName {
			buttons := generateButtons(i)
			t.Client.SendGroupAppButton(ctx, uniqueCid, userId, buttons)
			break
		}
	}
}

func (t *TrainClient) sendFiveItems(ctx context.Context, userId, content string) {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	t.Client.SendMessage(ctx, uniqueCid, userId, uuid.New().String(), bot.MessageCategoryPlainText, "恭喜，请选一味", "")

	buttons := make([]*bot.AppButtonView, 0)
	for _, i := range bubble.ItemMap {
		b := &bot.AppButtonView{
			Label:  fmt.Sprintf("%s", i.Name),
			Color:  i.HexColor,
			Action: fmt.Sprintf("input:%s", i.CommandName),
		}
		buttons = append(buttons, b)

	}
	t.Client.SendGroupAppButton(ctx, uniqueCid, userId, buttons)

}

func (t *TrainClient) sendGaoItems(ctx context.Context, userId, content string) {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	t.Client.SendMessage(ctx, uniqueCid, userId, uuid.New().String(), bot.MessageCategoryPlainText, "恭喜，请选一味", "")

	for _, i := range bubble.ItemMap {
		buttons := generateButtons(i)
		t.Client.SendGroupAppButton(ctx, uniqueCid, userId, buttons)
	}

}

func (t *TrainClient) downloadAttachment(ctx context.Context, userId string, img *ImageMessage) error {
	att, err := bot.AttachmentShow(ctx, t.Config.ClientId, t.Config.SessionId, t.Config.PrivateKey, img.AttachmentID)
	if err != nil {
		return err
	}
	t.sendTextMsg(ctx, userId, att.ViewURL)
	imgFormat := "png"
	if strings.Contains(img.MimeType, "/") {
		imgFormat = strings.Split(img.MimeType, "/")[1]
	}
	path := img.AttachmentID + "." + imgFormat
	err = downloadToFile(att.ViewURL, path)
	if err == nil {
		cur, _ := os.Getwd()
		dir := filepath.Join(cur, "data")
		t.sendTextMsg(ctx, userId, fmt.Sprintf("%s下载到%s目录，请查看", path, dir))
	}
	return nil
}

func downloadToFile(url, name string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	path, _ := os.Getwd()
	dir := filepath.Join(path, "data")
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	return os.WriteFile(filepath.Join(dir, name), all, 0644)
}

func queryItems(ctx context.Context) map[string][]*tomato.DailyItem {
	itemMap := make(map[string][]*tomato.DailyItem, 0)
	items, err := models.DailyItemStore.FindTodayItem(ctx, todayStart(time.Now()))
	if err != nil {
		return itemMap
	}
	for _, i := range items {
		if v, found := itemMap[i.ItemComm]; found {
			v = append(v, convertItem(i))
			itemMap[i.ItemComm] = v
		} else {
			di := []*tomato.DailyItem{convertItem(i)}
			itemMap[i.ItemComm] = di
		}
	}
	return itemMap
}

func convertItem(d dbm.DailyItem) *tomato.DailyItem {
	return &tomato.DailyItem{
		T:       d.ItemType,
		Name:    bubble.ItemMap[d.ItemComm].Name,
		Stamp:   millisecondToTime(int64(d.CreatedAt)),
		Min:     d.Min,
		Comment: d.Comment,
	}
}

func generateButtons(item *bubble.TypeItem) []*bot.AppButtonView {
	buttons := make([]*bot.AppButtonView, 0)

	for i := 0; i < 5; i++ {
		min := 20 + i*5
		label := fmt.Sprintf("%s%d", item.Name, min)
		b := &bot.AppButtonView{
			Label:  label,
			Color:  item.HexColor,
			Action: fmt.Sprintf("input:%s%d,完成%s分钟", item.CommandName, min, label),
		}
		buttons = append(buttons, b)
	}

	return buttons
}

func matchMaoRum(data string) bool {
	for _, i := range bubble.ItemMap {
		if data == i.CommandName {
			return true
		}
	}
	return false
}

func isDailyItem(c string) (bool, *dbm.DailyItem) {
	d := &dbm.DailyItem{
		CreatedAt: int(toMillisecond(time.Now())),
	}
	// -D25,abcd
	if len(c) > 2 {
		it := c[:2]
		if matchMaoRum(it) {
			d.ItemType = bubble.ItemMap[it].Type
			d.ItemComm = it
			m := ""
			if strings.Contains(c, ",") {
				h := strings.Split(c, ",")[0]
				t := strings.TrimSpace(h)
				m = t[2:]
				d.Comment = strings.TrimSpace(c[len(h)+1:])
			} else {
				m = c[2:]
			}
			min, err := strconv.Atoi(m)
			if err != nil {
				return false, nil
			}
			d.Min = min
			return true, d
		}
	}
	return false, nil
}

// -D,2021-12-28 18:58,30,comment
//-D,2021-12-28 18:58,30
func isDefinedFormat(data string) (bool, *dbm.DailyItem) {
	d := &dbm.DailyItem{}
	if len(data) > 20 && strings.Contains(data, ",") && len(strings.Split(data, ",")) > 2 {
		group := strings.Split(data, ",")
		// type
		it := strings.TrimSpace(group[0])
		if !matchMaoRum(it) {
			return false, nil
		}
		// time
		t, err := time.Parse(DefineTimeFormat, strings.TrimSpace(group[1]))
		if err != nil {
			return false, nil
		}
		// min
		min, err := strconv.Atoi(strings.TrimSpace(group[2]))
		if err != nil {
			return false, nil
		}
		d.ItemType = bubble.ItemMap[it].Type
		d.ItemComm = it
		d.CreatedAt = int(toMillisecond(t))
		d.Min = min

		// comment
		if len(group) > 3 {
			d.Comment = strings.TrimSpace(group[3])
		}
		return true, d

	}
	return false, nil
}

func todayStart(now time.Time) int {
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return int(toMillisecond(t))
}
