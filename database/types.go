package database

type User struct {
	ID           int64   `bson:"_id"`
	UserName     string  `bson:"user_name"`
	FirstName    string  `bson:"first_name"`
	LastName     string  `bson:"last_name"`
	WalletSet    bool    `bson:"wallet_set"`
	PublicKey    string  `bson:"public_key"`
	PrivateKey   string  `bson:"private_key"`
	PaperPnl     float64 `bson:"paper_pnl"`
	PaperBalance float64 `bson:"paper_balance"`
	PerpPnl      float64 `bson:"perp_pnl"`
	PerpBalance  float64 `bson:"perp_balance"`
}
