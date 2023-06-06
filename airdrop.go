package main

import (
	"airdrop/blockchain"
	"airdrop/db"
	"fmt"
	"log"
)

func main() {
	db.InitMysql()
	//blockchain.ImportToDb() // 将用户从文件导入数据库
	//os.Exit(2)
	err, users := blockchain.GetUser()
	if err != nil {
		panic("can not find user")
	}

	for _, u := range users {
		err, tokenIdStart := blockchain.GetTokenIdStart()
		if err != nil {
			panic("can not find token id")
		}
		fmt.Println(tokenIdStart, "~~~")
		err, _ = blockchain.AwardItem(u.User, int(tokenIdStart))
		if err != nil {
			log.Println("awardItem failed", u.User, err)
			continue
		}
		tokenIdStart++
	}

}
