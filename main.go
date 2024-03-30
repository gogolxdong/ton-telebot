package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	tele "gopkg.in/telebot.v3"

	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	token := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/dbname")
	if err != nil {
		log.Fatal(err)
		return
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	var walletClient *wallet.Wallet
	var walletMenu = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				tele.InlineButton{
					Text: "创建钱包",
					Data: "create_wallet",
				},
				tele.InlineButton{
					Text: "导出私钥",
					Data: "export_private_key",
				},
			},
			{
				tele.InlineButton{
					Text: "狙击合约",
					Data: "snipe_token",
				},
				tele.InlineButton{
					Text: "购买TON",
					Data: "buy_ton",
				},
				tele.InlineButton{
					Text: "卖出TON",
					Data: "sell_ton",
				},
			},
		},
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("欢迎使用TON钱包机器人!点击下方按钮开始创建你的钱包:", walletMenu)

	})
	b.Handle(tele.OnCallback, func(c tele.Context) error {
		var callBack = c.Callback().Data
		if callBack == "create_wallet" {
			client := liteclient.NewConnectionPool()

			err := client.AddConnection(context.Background(), "", "")
			if err != nil {
				log.Fatal(err)
			}

			api := ton.NewAPIClient(client).WithRetry()
			seed := wallet.NewSeed()
			w, err := wallet.FromSeed(api, seed, wallet.V4R2)
			if err != nil {
				log.Fatal(err)
			}
			walletClient = w
			privateKey := w.PrivateKey()
			c.Send(fmt.Sprintf("钱包已创建,地址为: %s, %s", w.Address(), privateKey))
		} else if callBack == "export_private_key" {
			var privateKey = walletClient.PrivateKey()
			c.Send(fmt.Sprintf("钱包已创建,地址为: %s", privateKey))

		} else if callBack == "snipe_token" {

		} else if callBack == "buy_ton" {

		} else if callBack == "sell_ton" {

		}
		return nil
	})

	b.Start()
}

