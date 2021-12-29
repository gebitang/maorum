package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"time"
)

func (t *TrainClient) handleClaim(ctx context.Context, userId string) {
	now := time.Now().Format("2006-01-02")
	traceId := bot.UniqueConversationId(userId, now)
	trace, err := bot.ReadTransferByTrace(ctx, traceId, t.Config.ClientId, t.Config.SessionId, t.Config.PrivateKey)
	if err != nil {
		in := &bot.TransferInput{
			AssetId:     cnbAssetId,
			RecipientId: userId,
			Amount:      number.FromString("1"),
			TraceId:     traceId,
			Memo:        "test from bot",
		}

		transfer, e := bot.CreateTransfer(ctx, in, t.Config.ClientId, t.Config.SessionId, t.Config.PrivateKey, t.Config.Pin, t.Config.PinToken)
		if e != nil {
			mErr := &bot.Error{}
			eb, _ := json.Marshal(e)
			json.Unmarshal(eb, mErr)
			// {"status":202,"code":20125,"description":"Transfer has been paid by someone else."}
			if mErr.Code == 20125 {
				t.sendTextMsg(ctx, userId, "keystore已经被其他应用使用")
			}
			// {"status":202,"code":20117,"description":"Insufficient balance."}
			if mErr.Code == 20117 {
				t.sendTextMsg(ctx, userId, "余额不足，请先转账或打赏CNB")
				transferAction := fmt.Sprintf("mixin://transfer/%s", t.Config.ClientId)
				t.Client.SendAppButton(ctx, bot.UniqueConversationId(userId, t.Config.ClientId), userId, "打赏", transferAction, "#1DDA99")
			}
			return
		}
		tt, _ := json.MarshalIndent(transfer, "", "  ")
		fmt.Println("transfer result: ", string(tt))
		return
	}

	if len(trace.SnapshotId) > 0 {
		t.sendTextMsg(ctx, userId, "今日已领取，请明天再来。")
		return
	}
}

func (t *TrainClient) handleDonate(ctx context.Context, userId string) {
	transferAction := fmt.Sprintf("mixin://transfer/%s", t.Config.ClientId)
	t.Client.SendAppButton(ctx, bot.UniqueConversationId(userId, t.Config.ClientId), userId, "点我打赏", transferAction, "#000000")
}

func (t *TrainClient) helpMsgWithInfo(ctx context.Context, userId, info string) {
	t.sendTextMsg(ctx, userId, info+helpMsg)
	t.Client.SendAppButton(ctx, bot.UniqueConversationId(userId, t.Config.ClientId), userId, "签到", "input:/claim", "#1DDA99")
	t.Client.SendAppButton(ctx, bot.UniqueConversationId(userId, t.Config.ClientId), userId, "打赏", "input:/donate", "#f05d5d")
}

func (t *TrainClient) handleAssets(ctx context.Context, userId, data string) bool {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	botMsg := bot.MessageView{
		ConversationId: uniqueCid,
		UserId:         userId,
	}
	if isValidUUID(data) {
		asset, err := readNetworkAsset(ctx, data)
		if err != nil {
			return false
		}
		b, _ := json.MarshalIndent(asset, "", "  ")
		content := fmt.Sprintf("```json\n%s\n```", string(b))
		t.Client.SendPost(ctx, botMsg, content)
	} else {
		assets, err := bot.AssetSearch(ctx, data)
		if err != nil {
			return false
		}
		if len(assets) > 0 {
			t.sendTextMsg(ctx, userId, assets[0].AssetId)
			b, _ := json.MarshalIndent(assets, "", "  ")
			content := fmt.Sprintf("```json\n%s\n```", string(b))
			t.Client.SendPost(ctx, botMsg, content)
		}
	}

	return true
}

func (t *TrainClient) sendTextMsg(ctx context.Context, userId, content string) {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	t.Client.SendMessage(ctx, uniqueCid, userId, uuid.New().String(), bot.MessageCategoryPlainText, content, "")
}

func (t *TrainClient) handleUser(ctx context.Context, userId, data string) bool {
	user, err := bot.GetUser(ctx, data, t.Config.ClientId, t.Config.SessionId, t.Config.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)

	err = t.Client.SendContact(ctx, uniqueCid, userId, user.UserId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	transferAction := fmt.Sprintf("mixin://transfer/%s", user.UserId)
	label := fmt.Sprintf("\ntransfer to %s\n", user.FullName)
	if data != user.UserId {
		t.sendTextMsg(ctx, userId, user.UserId)
	}

	err = t.Client.SendAppButton(ctx, uniqueCid, userId, label, transferAction, "#1DDA99")
	if err != nil {
		fmt.Println(err)
		return false
	}
	encode, err := qrcode.Encode(transferAction, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err)
		return false
	}

	b, err := t.sendPlainImage(ctx, userId, "image/jpeg", encode, 300, 300)
	return b
}

func (t *TrainClient) sendPlainImage(ctx context.Context, userId, imgType string, encode []byte, w, h int) (bool, error) {
	uniqueCid := bot.UniqueConversationId(userId, t.Config.ClientId)
	attachment, err := bot.CreateAttachment(ctx, t.Config.ClientId, t.Config.SessionId, t.Config.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	err = uploadAttachmentTo(attachment.UploadUrl, encode)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	img := &ImageMessage{
		AttachmentID: attachment.AttachmentId,
		MimeType:     "image/jpeg",
		Width:        w,
		Height:       h,
		Size:         len(encode),
		Thumbnail:    base64.StdEncoding.EncodeToString(encode),
	}
	byteImg, err := json.Marshal(img)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	err = t.Client.SendMessage(ctx, uniqueCid, userId, uuid.New().String(), bot.MessageCategoryPlainImage, string(byteImg), "")
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, err
}
