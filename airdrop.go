package main

import (
	"airdrop/blockchain"
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var userNfts []blockchain.UserNft
	accountFile, err := os.Open("./account/account.txt") // 注意文件最后一行需要有换行，否则只能读到倒数第二行
	if err != nil {
		log.Println("os.Open error ", err)
		return
	}
	defer accountFile.Close()
	buf := bufio.NewReader(accountFile)
	tokenIdStart := 0
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
		err, userNft := blockchain.AwardItem(address, tokenIdStart)
		if err != nil {
			log.Println("awardItem failed", address, err)
		}
		userNfts = append(userNfts, userNft)
		tokenIdStart++
	}
	// Write to file
	DumpFile("./data/nft.json", userNfts)

}

func DumpFile(filename string, userNfts []blockchain.UserNft) {
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
