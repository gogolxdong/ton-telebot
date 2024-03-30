package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	tele "gopkg.in/telebot.v3"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"

	"gorm.io/gen"
	"gorm.io/gorm"

	"tontelebot/model"
)

func walletMenu(user []model.User) *tele.ReplyMarkup {
	var button [][]tele.InlineButton
	var jettonButton = []tele.InlineButton{
		{
			Text: "狙击合约",
			Data: "snipe_token",
		},
		{
			Text: "购买Jetton",
			Data: "buy_ton",
		},
		{
			Text: "卖出Jetton",
			Data: "sell_ton",
		},
	}
	if len(user) == 0 {
		button = append(button, []tele.InlineButton{
			{
				Text: "创建钱包",
				Data: "create_wallet",
			},
		},
			jettonButton,
		)
	}
	for _, u := range user {
		button = append(button, []tele.InlineButton{
			{
				Text: fmt.Sprintf("导出私钥%s", u.JettonWalletAddress),
				Data: "export_private_key",
			},
			{
				Text: fmt.Sprintf("导出助记词%s", u.JettonWalletAddress),
				Data: "export_mnemonic",
			},
		}, jettonButton)
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: button,
	}
}
func genModel(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "models",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()

}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	token := os.Getenv("TOKEN")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbname := os.Getenv("DBNAME")
	adnlHost := os.Getenv("ADNLHOST")
	adnlPort := os.Getenv("ADNLPORT")
	adnlKey := os.Getenv("ADNLKEY")
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// genModel(db)

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	var walletClient *wallet.Wallet
	var user []model.User
	var client *liteclient.ConnectionPool
	var api ton.APIClientWrapped
	var snipeJetton string
	b.Handle("/start", func(c tele.Context) error {
		u := c.Sender()
		tx := db.Find(&user, "uid=?", u.ID)
		if tx.Error != nil {
			fmt.Println(tx.Error)
			return c.Send("欢迎使用TON钱包机器人！", walletMenu(nil))
		}
		client = liteclient.NewConnectionPool()

		err := client.AddConnection(context.Background(), fmt.Sprintf("%s:%s", adnlHost, adnlPort), adnlKey)
		if err != nil {
			log.Fatal(err)
		}

		api = ton.NewAPIClient(client).WithRetry()
		seed := strings.Split(user[0].Mnemonic, " ")
		w, err := wallet.FromSeed(api, seed, wallet.V4R2)
		if err != nil {
			log.Fatal(err)
		}
		walletClient = w
		return c.Send("欢迎使用TON钱包机器人！", walletMenu(user))

	})

	b.Handle(tele.OnCallback, func(c tele.Context) error {
		u := c.Sender()

		var callBack = c.Callback().Data
		if callBack == "create_wallet" {
			seed := wallet.NewSeed()
			w, err := wallet.FromSeed(api, seed, wallet.V4R2)
			if err != nil {
				log.Fatal(err)
			}
			walletClient = w
			privateKey := base64.StdEncoding.EncodeToString(walletClient.PrivateKey())
			tx := db.Save(&model.User{UID: u.ID, UserName: u.Username, JettonWalletAddress: w.Address().String(),
				PrivateKey: privateKey, Mnemonic: strings.Join(seed, " "), CreateTime: time.Now()})
			if tx.Error != nil {
				log.Fatal(tx.Error)
			}
			c.Send(fmt.Sprintf("钱包已创建,地址为: %s\n助记词: %s\n私钥：%s\n", w.Address(), seed, privateKey))
		} else if callBack == "export_private_key" {
			privateKey := base64.StdEncoding.EncodeToString(walletClient.PrivateKey())
			c.Send(fmt.Sprintf("私钥己导出: %s\n", privateKey))

		} else if callBack == "export_mnemonic" {
			c.Send(user[0].Mnemonic)

		} else if callBack == "snipe_token" {
			c.Send("请输入合约地址：")

			b.Handle(tele.OnText, func(c tele.Context) error {
				snipeJetton = c.Text()
				c.Send(fmt.Sprintf("你输入的合约地址是：%s", snipeJetton))
				return nil
			})

		} else if callBack == "buy_ton" {
			c.Send(fmt.Sprintf("购买%s", snipeJetton))

		} else if callBack == "sell_ton" {
			c.Send(fmt.Sprintf("卖出%s", snipeJetton))

		}
		return nil
	})

	b.Start()
}
