package utils

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func CreateAddress(mnemonic string, index int64) (*ecdsa.PrivateKey, string, int64, error) {

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatalf("Failed to create wallet from mnemonic: %v", err)
		return nil, "", index, err
	}

	path := fmt.Sprintf("m/44'/60'/0'/0/%v", index)

	hdPath := hdwallet.MustParseDerivationPath(path)

	account, err := wallet.Derive(hdPath, false)
	if err != nil {
		log.Fatalf("Failed to derive account: %v", err)
		return nil, "", index, err
	}

	privKey, err := wallet.PrivateKey(account)
	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
		return nil, "", index, err
	}

	// fmt.Printf("Mnemonic: %s\n", mnemonic)
	// fmt.Printf("Private Key: 0x%s\n", privKeyHex)
	address := account.Address.Hex()

	return privKey, address, int64(index), nil
}

// SignTransaction 用助记词和索引推导账户，对交易哈希进行签名并返回签名结果
func SignTransaction(mnemonic string, index int64, digestHash []byte) ([]byte, error) {

	privKey, _, _, err := CreateAddress(mnemonic, index)
	if err != nil {
		return nil, err
	}

	// 用 go-ethereum/crypto.Sign 进行以太坊签名
	signature, err := crypto.Sign(digestHash, privKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
