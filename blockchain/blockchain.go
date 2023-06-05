package blockchain

import (
	"airdrop/binding"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

const (
	TokenURI   = "https://ipfs.io/ipfs/QmP2soBpShrzLf5K51vtgFfLF3egaBGnyEibZq6jskGU5A?filename=tt_community.json"
	Contract   = "0x30674B79c036C8222e85d0AB72F2CB47f446d1cB"
	PrivateKey = "ce79cca4e66c331b5ca2945a5e267c2b74daf7132a7d5031889cf153aa11b853"
)

type UserNft struct {
	User    string
	TokenId int
	TxHash  string
}

func Client() *ethclient.Client {
	client, err := ethclient.Dial("https://rpc-mumbai.maticvigil.com")
	//client, err := ethclient.Dial("http://polygon.drpc.org")
	if err != nil {
		panic("dail failed")
	}
	return client
}

func AwardItem(address string, tokenId int) (error, UserNft) {
	cli := Client()
	AirDrop, err := binding.NewAirDropToken(common.HexToAddress(Contract), cli)
	if err != nil {
		log.Println(err)
		return err, UserNft{}
	}
	privateKeyEcdsa, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Println("privateKey  ", "err", err)
		return err, UserNft{}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyEcdsa, big.NewInt(int64(80001)))
	if err != nil {
		log.Panicln(err)
		return err, UserNft{}
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
		log.Println(err)
		return err, UserNft{}
	}
	log.Println(address, "-----", item.Hash(), "-----", tokenId)
	return nil, UserNft{
		address,
		tokenId,
		item.Hash().String(),
	}
}
