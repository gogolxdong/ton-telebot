package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	tele "gopkg.in/telebot.v3"

	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/jetton"

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
			Data: "buy_token",
		},
		{
			Text: "卖出Jetton",
			Data: "sell_token",
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

func initUser(c tele.Context, db *gorm.DB) {
	adnlHost := os.Getenv("ADNLHOST")
	adnlPort := os.Getenv("ADNLPORT")
	adnlKey := os.Getenv("ADNLKEY")
	u := c.Sender()
	tx := db.Find(&user, "uid=?", u.ID)
	if tx.Error != nil {
		fmt.Println("Find", tx.Error)
		c.Send("欢迎使用TON钱包机器人！", walletMenu(nil))
	}
	fmt.Println(adnlHost, adnlPort, adnlKey)
	err := client.AddConnection(context.Background(), fmt.Sprintf("%s:%s", adnlHost, adnlPort), adnlKey)
	if err != nil {
		log.Fatalf("AddConnection: %s", err.Error())
	}

	seed := strings.Split(user[0].Mnemonic, " ")
	w, err := wallet.FromSeed(api, seed, wallet.V4R2)
	if err != nil {
		log.Fatalf("FromSeed: %s", err.Error())
	}
	walletClient = w
	fmt.Println("walletClient:", walletClient)
}

var client *liteclient.ConnectionPool = liteclient.NewConnectionPool()
var api ton.APIClientWrapped = ton.NewAPIClient(client).WithRetry()
var stonRouter *address.Address = address.MustParseAddr("EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC")

var walletClient *wallet.Wallet
var user []model.User

var snipeJetton string
var botToken string
var dbUser string
var dbPassword string
var dbname string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	botToken = os.Getenv("TOKEN")
	dbUser = os.Getenv("DBUSER")
	dbPassword = os.Getenv("DBPASSWORD")
	dbname = os.Getenv("DBNAME")
	pref := tele.Settings{
		Token:  botToken,
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

	// sqlDB, err := db.DB()
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Hour)

	b.Handle("/start", func(c tele.Context) error {
		initUser(c, db)
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
				PrivateKey: privateKey, Mnemonic: strings.Join(seed, " "), CreateTime: time.Now(), UpdateTime: time.Now()})
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

		} else if callBack == "buy_token" {
			c.Send("请输入购买数量：")
			var amount tlb.Coins
			b.Handle(tele.OnText, func(c tele.Context) error {
				var err error
				amount = tlb.MustFromDecimal(c.Text(), 9)
				log.Printf("walletClient: %v", walletClient)
				walletAddr := walletClient.Address()
				snipeJettonAddress, err := address.ParseAddr(snipeJetton)
				if err != nil {
					log.Fatalf("ParseAddr: %s", err.Error())
				}
				token := jetton.NewJettonMasterClient(api, snipeJettonAddress)
				fmt.Println(client)
				ctx := client.StickyContext(context.Background())

				tokenWallet, err := token.GetJettonWallet(ctx, walletAddr)
				if err != nil {
					log.Fatalf("GetJettonWallet: %s", err.Error())
				}

				tokenBalance, err := tokenWallet.GetBalance(ctx)
				if err != nil {
					log.Fatalf("GetBalance: %s", err.Error())
				}
				log.Println("our jetton balance:", tokenBalance.String())

				to := address.MustParseAddr(user[0].JettonWalletAddress)
				transferPayload, err := tokenWallet.BuildTransferPayloadV2(to, to, amount, tlb.ZeroCoins, nil, nil)
				if err != nil {
					log.Fatalf("BuildTransferPayloadV2: %s", err.Error())
				}

				msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)

				tx, _, err := walletClient.SendWaitTransaction(ctx, msg)
				if err != nil {
					log.Fatalf("SendWaitTransaction: %s", err.Error())
				}
				c.Send(fmt.Sprintf("交易发送：%s", base64.StdEncoding.EncodeToString(tx.Hash)))
				return nil
			})

		} else if callBack == "sell_token" {
			c.Send("请输入售出数量：")
			// var amount tlb.Coins
			b.Handle(tele.OnText, func(c tele.Context) error {
				var err error
				// amount = tlb.MustFromDecimal(c.Text(), 9)

				walletAddr := walletClient.Address()
				snipeJettonAddress, err := address.ParseAddr(snipeJetton)
				if err != nil {
					log.Fatalf("ParseAddr: %s", err.Error())
				}
				token := jetton.NewJettonMasterClient(api, snipeJettonAddress)

				ctx := client.StickyContext(context.Background())

				tokenWallet, err := token.GetJettonWallet(ctx, walletAddr)
				if err != nil {
					log.Fatalf("GetJettonWallet: %s", err.Error())
				}

				tokenBalance, err := tokenWallet.GetBalance(ctx)
				if err != nil {
					log.Fatalf("GetBalance: %s", err.Error())
				}
				log.Println("our jetton balance:", tokenBalance.String())
				body := cell.BeginCell().
					MustStoreUInt(0x25938561, 32).    // swap op code
					MustStoreUInt(rand.Uint64(), 64). // query id
					MustStoreAddr(stonRouter).
					MustStoreRef(
						cell.BeginCell().
							MustStoreBigCoins(tlb.MustFromTON("1.521").Nano()).
							EndCell(),
					).EndCell()

				message := &wallet.Message{
					Mode: 1,
					InternalMessage: &tlb.InternalMessage{
						Bounce:  true,
						DstAddr: stonRouter,
						Amount:  tlb.MustFromTON("0.03"),
						Body:    body,
					},
				}
				block, err := api.CurrentMasterchainInfo(context.Background())
				if err != nil {
					panic(err)
				}

				response, err := api.RunGetMethod(
					context.Background(),
					block,
					stonRouter,
					"get_pool_address",
					map[string]*address.Address{
						"token0": address.MustParseAddr("EQA2kCVNwVsil2EM2mB0SkXytxCqQjS4mttjDpnXmwG9T6bO"),
						"token1": address.MustParseAddr("EQA2kCVNwVsil2EM2mB0SkXytxCqQjS4mttjDpnXmwG9T6bO"),
					},
				)
				poolAddress := response.MustCell(0).BeginParse().MustLoadAddr()

				if err != nil {
					panic(fmt.Errorf("failed to get pool address: %w", err))
				}

				fmt.Printf("Pool Address: %s\n", poolAddress)

				tx, _, err := walletClient.SendWaitTransaction(ctx, message)
				if err != nil {
					log.Fatalf("SendWaitTransaction:%s", err.Error())
				}
				c.Send(fmt.Sprintf("交易发送：%s", base64.StdEncoding.EncodeToString(tx.Hash)))
				return nil
			})
		}
		return nil
	})

	b.Start()
}
