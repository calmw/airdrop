package blockchain

import (
	"airdrop/binding"
	"airdrop/db"
	"airdrop/model"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	TokenURI   = "https://ipfs.io/ipfs/QmR3zuDAMDYtsFvmtCyqt1LUtgi6UXVEbNpnFq4RyNSWW6?filename=tt_community_new_user.json"
	Contract   = ""
	PrivateKey = ""
	Rpc        = "https://polygon-bor.publicnode.com"
	ChainId    = 137
	//ChainId = 80001
)

type UserNft struct {
	User    string
	TokenId int
	TxHash  string
}

func Client() *ethclient.Client {
	client, err := ethclient.Dial(Rpc)
	if err != nil {
		panic("dail failed")
	}
	return client
}

func AwardItem(address string, tokenId int) (error, model.UserNft) {
	cli := Client()
	AirDrop, err := binding.NewAirDropToken(common.HexToAddress(Contract), cli)
	if err != nil {
		log.Println(err)
		return err, model.UserNft{}
	}
	privateKeyEcdsa, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Println("privateKey  ", "err", err)
		return err, model.UserNft{}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyEcdsa, big.NewInt(int64(ChainId)))
	if err != nil {
		log.Panicln(err)
		return err, model.UserNft{}
	}

	transactOpts := bind.TransactOpts{
		From:      auth.From,
		Nonce:     nil,
		Signer:    auth.Signer, // Method to use for signing the transaction (mandatory)
		Value:     big.NewInt(0),
		GasPrice:  nil,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  0,
		Context:   context.Background(),
		NoSend:    false, // Do all transact steps but do not send the transaction
	}

	item, err := AirDrop.AwardItem(&transactOpts, common.HexToAddress(address), TokenURI)
	if err != nil {
		log.Println(err, "++", tokenId)
		return err, model.UserNft{}
	}
	log.Println(address, "-----", item.Hash(), "-----", tokenId)
	SaveToDB(model.UserNft{
		User:    address,
		TokenId: tokenId,
		IsMint:  1,
		TxHash:  item.Hash().String(),
	})
	return nil, model.UserNft{
		User:    address,
		TokenId: tokenId,
		TxHash:  item.Hash().String(),
	}
}

func GetUser() (error, []model.UserNft) {
	var userNfts []model.UserNft
	err := db.Mysql.Table("user_nft").Where("is_mint=0").Find(&userNfts).Debug().Error
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userNfts
		} else {
			return err, []model.UserNft{}
		}
	}

	return nil, userNfts
}

func SaveToDB(userNft model.UserNft) {
	whereCondition := fmt.Sprintf("user='%s'", userNft.User)
	err := db.Mysql.Table("user_nft").Where(whereCondition).UpdateColumns(&userNft).Debug().Error
	if err != nil {
		log.Println(err)
	}
}

func GetTokenIdStart() (error, int64) {
	var count int64
	err := db.Mysql.Table("user_nft").Where("is_mint>0").Count(&count).Debug().Error
	if err != nil {
		return err, count
	}

	return nil, count
}

func DumpFile(filename string, userNfts UserNft) {
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	data, err := json.MarshalIndent(userNfts, "", "\t") // 带缩进的美化版
	if err != nil {
		panic(err)
	}
	fp.Write(data)
}

func ImportToDb() {
	accountFile, err := os.Open("./account/account.txt") // 注意文件最后一行需要有换行，否则只能读到倒数第二行
	if err != nil {
		log.Println("os.Open error ", err)
		return
	}
	defer accountFile.Close()
	buf := bufio.NewReader(accountFile)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("end of file")
				break
			} else {
				log.Printf("read file err:%s", err.Error())
				break
			}
		}

		address := strings.TrimRight(string(line), "\n")
		userNft := model.UserNft{
			User: address,
		}
		fmt.Println(userNft, 123)
		err = db.Mysql.Table("user_nft").Order("id desc").Create(&userNft).Debug().Error
		if err != nil {
			log.Println(err)
		}
	}
}
