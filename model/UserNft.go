package model

type UserNft struct {
	Id      int64  `json:"-" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	User    string `json:"user" gorm:"column:user"`
	TokenId int    `json:"token_id" gorm:"column:token_id"`
	TxHash  string `json:"tx_hash" gorm:"column:tx_hash"`
	IsMint  int    `json:"is_mint" gorm:"column:is_mint"`
}
