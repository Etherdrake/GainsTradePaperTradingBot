package stringbuilders

import (
	"HootTelegram/database"
	"fmt"
)

func BuildWalletMainString(guid int64) (string, error) {
	keys, err := database.GetKeys(guid)
	if err != nil {
		return "", err
	}

	publicKey := "`" + keys.PublicKey + "`"
	privateKey := "`" + keys.PrivateKey + "`"

	// You would need to replace the 0.00 balances with actual values.
	// These could be fetched similarly to how the keys are fetched.
	balanceMatic := "0.00 MATIC"
	balanceArb := "0.00 ARB"
	balanceEth := "0.00 ETH"

	return fmt.Sprintf("Publickey: %s\n\nPrivatekey: %s\n\nBalance MATIC: %s\nBalance ARB: %s\nBalance ETH: %s",
		publicKey, privateKey, balanceMatic, balanceArb, balanceEth), nil
}
