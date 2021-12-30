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
	"github.com/MixinNetwork/bot-api-go-client"
	"log"
	"net/http"
)

const (
	Debug            = false
	TimeFormat       = "2006-01-02 15:04:05"
	DefineTimeFormat = "2006-01-02 15:04"
	cnbAssetId       = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	helpMsg          = "\n1. 支持用户查询，请发送 user_id | identity_number\n  2. 支持资产查询，请发送 asset_id | symbol\n  3. 支持每日领取 1cnb，请发送 /claim 或点击签到\n  4. 支持打赏，请发送 /donate 或点击打赏"
)

var (
	mars       *TrainClient
	httpClient = &http.Client{}
	helpMap    []string
)

type (
	ImageMessage struct {
		AttachmentID string `json:"attachment_id,omitempty"`
		MimeType     string `json:"mime_type,omitempty"`
		Width        int    `json:"width,omitempty"`
		Height       int    `json:"height,omitempty"`
		Size         int    `json:"size,omitempty"`
		Thumbnail    string `json:"thumbnail,omitempty"`
	}

	MixinBlazeHandler func(ctx context.Context, msg bot.MessageView, clientID string) error

	TrainClient struct {
		Ctx    context.Context
		Config *config.Config
		Client *bot.BlazeClient
	}
)

func (f MixinBlazeHandler) OnMessage(ctx context.Context, msg bot.MessageView, clientID string) error {
	return f(ctx, msg, clientID)
}

func (f MixinBlazeHandler) OnAckReceipt(ctx context.Context, msg bot.MessageView, clientID string) error {
	indent, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	if Debug {
		log.Println("ack Message...", string(indent))
	}
	return nil
}

func (f MixinBlazeHandler) SyncAck() bool {
	return true
}

func NewClient(ctx context.Context, config *config.Config, client *bot.BlazeClient) *TrainClient {
	mars = &TrainClient{
		Ctx:    ctx,
		Config: config,
		Client: client,
	}
	return mars
}

func Handler(ctx context.Context, botMsg bot.MessageView, clientID string) error {
	if Debug {
		marshal, _ := json.MarshalIndent(botMsg, "", "  ")
		fmt.Println("msg data: ", string(marshal))
	}
	decodeBytes, _ := base64.StdEncoding.DecodeString(botMsg.Data)

	if botMsg.Category == bot.MessageCategorySystemAccountSnapshot {
		ss := &bot.Snapshot{}
		json.Unmarshal(decodeBytes, ss)
		asset, _ := readNetworkAsset(ctx, ss.AssetId)
		con := fmt.Sprintf("打赏的%s %s 已收到，感谢支持。", ss.Amount, asset.Symbol)
		mars.sendTextMsg(ctx, botMsg.UserId, con)
		return nil
	}

	if botMsg.Category != bot.MessageCategoryPlainText {
		mars.sendTextMsg(ctx, botMsg.UserId, "仅支持文本信息")
		return nil
	}

	//----------------mao rum----------------
	if botMsg.Category == bot.MessageCategoryPlainImage {
		img := &ImageMessage{}
		json.Unmarshal(decodeBytes, img)
		mars.downloadAttachment(ctx, botMsg.UserId, img)
		return nil
	}

	data := string(decodeBytes)
	if data == "rum" {
		mars.drinkRum(ctx, botMsg.UserId, data)
		return nil
	}

	if data == "mur" {
		mars.readLatestRum(ctx, botMsg.UserId, data)
		return nil
	}

	if data == "mao" {
		mars.drawMaoPic(ctx, botMsg.UserId, data)
		return nil
	}

	if data == "gao" {
		mars.sendFiveItems(ctx, botMsg.UserId, data)
		return nil
	}

	if data == "GAO" {
		mars.sendGaoItems(ctx, botMsg.UserId, data)
		return nil
	}

	if matchMaoRum(data) {
		mars.sendOneItems(ctx, botMsg.UserId, data)
		return nil
	}

	b, d := isDailyItem(data)
	if b {
		if d.Min > 120 {
			mars.sendTextMsg(ctx, botMsg.UserId, "时长最长120分钟，注意休息。")
			return nil
		}
		addOneItem(ctx, botMsg.UserId, d, mars)
		return nil
	}

	if b, d = isDefinedFormat(data); b {
		if d.Min > 120 {
			mars.sendTextMsg(ctx, botMsg.UserId, "时长最长120分钟，注意休息。")
			return nil
		}
		addOneItem(ctx, botMsg.UserId, d, mars)
		return nil
	}

	//------------------bot demo-----------------
	mars.sendTextMsg(ctx, botMsg.UserId, "To Be Continued...")

	return nil
}

func addOneItem(ctx context.Context, userId string, d *dbm.DailyItem, c *TrainClient) {
	err := models.DailyItemStore.Create(ctx, d)
	if err != nil {
		fmt.Println("create item err ", err)
		return
	}
	c.sendTextMsg(ctx, userId, fmt.Sprintf("恭喜完成%s%d分钟", bubble.ItemMap[d.ItemComm].Name, d.Min))
}

func botDemo(ctx context.Context, botMsg bot.MessageView, data string) error {

	if isHelpInfo(data) {
		mars.helpMsgWithInfo(ctx, botMsg.UserId, "")
		return nil
	}

	if data == "/claim" {
		mars.handleClaim(ctx, botMsg.UserId)
		return nil
	}

	if data == "/donate" {
		mars.handleDonate(ctx, botMsg.UserId)
		return nil
	}

	if isValidUUID(data) {
		a := mars.handleAssets(ctx, botMsg.UserId, data)
		b := mars.handleUser(ctx, botMsg.UserId, data)
		if !a && !b {
			mars.helpMsgWithInfo(ctx, botMsg.UserId, "指令输入不正确")
		}
		return nil
	}

	if isNumber(data) {
		a := mars.handleUser(ctx, botMsg.UserId, data)
		if !a {
			mars.helpMsgWithInfo(ctx, botMsg.UserId, "指令输入不正确")
		}
		return nil
	} else {
		a := mars.handleAssets(ctx, botMsg.UserId, data)
		if !a {
			mars.helpMsgWithInfo(ctx, botMsg.UserId, "指令输入不正确")
		}
	}
	return nil
}
