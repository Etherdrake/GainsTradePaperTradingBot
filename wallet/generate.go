package wallet

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"
)

func GenerateWallet() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	privKeyString := hexutil.Encode(privateKeyBytes)[2:]

	// PUBKEY

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	pubKeyString := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return strings.ToLower(pubKeyString), privKeyString
}
