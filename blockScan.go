package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

func getShardID(shard *ton.BlockIDExt) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func getNotSeenShards(ctx context.Context, api ton.APIClientWrapped, shard *ton.BlockIDExt, shardLastSeqno map[string]uint32) (ret []*ton.BlockIDExt, err error) {
	if no, ok := shardLastSeqno[getShardID(shard)]; ok && no == shard.SeqNo {
		return nil, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, fmt.Errorf("get block data: %w", err)
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, fmt.Errorf("get parent blocks (%d:%x:%d): %w", shard.Workchain, uint64(shard.Shard), shard.Shard, err)
	}

	for _, parent := range parents {
		ext, err := getNotSeenShards(ctx, api, parent, shardLastSeqno)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ext...)
	}

	ret = append(ret, shard)
	return ret, nil
}


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	adnlHost := os.Getenv("ADNLHOST")
	adnlPort := os.Getenv("ADNLPORT")
	adnlKey := os.Getenv("ADNLKEY")
	client := liteclient.NewConnectionPool()

	// cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	// if err != nil {
	// 	log.Fatalln("get config err: ", err.Error())
	// 	return
	// }
	err = client.AddConnection(context.Background(), fmt.Sprintf("%s:%s", adnlHost, adnlPort), adnlKey)
	if err != nil {
		log.Fatalf("AddConnection: %s", err.Error())
	}
	// err = client.AddConnectionsFromConfig(context.Background(), cfg)
	// if err != nil {
	// 	log.Fatalln("connection err: ", err.Error())
	// 	return
	// }

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	// api.SetTrustedBlockFromConfig(cfg)

	log.Println("checking proofs since config init block, it may take near a minute...")

	master, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	log.Println("master proofs chain successfully verified, all data is now safe and trusted!")

	ctx := api.Client().StickyContext(context.Background())

	shardLastSeqno := map[string]uint32{}

	firstShards, err := api.GetBlockShardsInfo(ctx, master)
	if err != nil {
		log.Fatalln("get shards err:", err.Error())
		return
	}
	for _, shard := range firstShards {
		shardLastSeqno[getShardID(shard)] = shard.SeqNo
	}

	for {
		log.Printf("scanning %d master block...\n", master.SeqNo)

		currentShards, err := api.GetBlockShardsInfo(ctx, master)
		if err != nil {
			log.Fatalln("get shards err:", err.Error())
			return
		}

		var newShards []*ton.BlockIDExt
		for _, shard := range currentShards {
			notSeen, err := getNotSeenShards(ctx, api, shard, shardLastSeqno)
			if err != nil {
				log.Fatalln("get not seen shards err:", err.Error())
				return
			}
			shardLastSeqno[getShardID(shard)] = shard.SeqNo
			newShards = append(newShards, notSeen...)
		}
		newShards = append(newShards, master)

		var txList []*tlb.Transaction

		for _, shard := range newShards {
			log.Printf("scanning block %d of shard %x in workchain %d...", shard.SeqNo, uint64(shard.Shard), shard.Workchain)

			var fetchedIDs []ton.TransactionShortInfo
			var after *ton.TransactionID3
			var more = true

			for more {
				fetchedIDs, more, err = api.WaitForBlock(master.SeqNo).GetBlockTransactionsV2(ctx, shard, 100, after)
				if err != nil {
					log.Fatalln("get tx ids err:", err.Error())
					return
				}

				if more {
					after = fetchedIDs[len(fetchedIDs)-1].ID3()
				}

				for _, id := range fetchedIDs {
					tx, err := api.GetTransaction(ctx, shard, address.NewAddress(0, byte(shard.Workchain), id.Account), id.LT)
					if err != nil {
						log.Fatalln("get tx data err:", err.Error())
						return
					}
					txList = append(txList, tx)
				}
			}
		}

		for i, transaction := range txList {
			log.Println(i, transaction.String())
		}

		if len(txList) == 0 {
			log.Printf("no transactions in %d block\n", master.SeqNo)
		}

		master, err = api.WaitForBlock(master.SeqNo+1).LookupBlock(ctx, master.Workchain, master.Shard, master.SeqNo+1)
		if err != nil {
			log.Fatalln("get masterchain info err: ", err.Error())
			return
		}
	}
}
