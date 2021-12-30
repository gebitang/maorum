package main

import (
	"context"
	"fmt"
	"gebitang.com/maorum/config"
	"gebitang.com/maorum/message"
	"github.com/MixinNetwork/bot-api-go-client"
	"log"
	"time"
)

var (
	ReleaseVersion string
)

func mixinBot() {
	ctx := context.Background()
	c := config.MixinBotConfig
	for {
		client := bot.NewBlazeClient(c.ClientId, c.SessionId, c.PrivateKey)
		message.NewClient(ctx, c, client)
		if err := client.Loop(ctx, message.MixinBlazeHandler(message.Handler)); err != nil {
			log.Println("test...", err)
		}
	}
}
func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // East 8 District
	time.Local = cstZone
}

func main() {
	fmt.Println("Welcome to MaoRum", ReleaseVersion)
	mixinBot()
}
