package asset

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	. "XBlock/cli/common"
	"XBlock/account"
	. "XBlock/common"
	. "XBlock/core/asset"
	"XBlock/core/contract"
	"XBlock/core/signature"
	"XBlock/core/transaction"
	"XBlock/net/httpjsonrpc"

	"github.com/urfave/cli"
)

const (
	RANDBYTELEN    = 4
	REFERTXHASHLEN = 64
)

func newContractContextWithoutProgramHashes(data signature.SignableData) *contract.ContractContext {
	return &contract.ContractContext{
		Data:       data,
		Codes:      make([][]byte, 1),
		Parameters: make([][][]byte, 1),
	}
}

func openWallet(name string, passwd []byte) account.Client {
	if name == account.WalletFileName {
		fmt.Println("Using default wallet: ", account.WalletFileName)
	}
	wallet := account.Open(name, passwd)
	if wallet == nil {
		fmt.Println("Failed to open wallet: ", name)
		os.Exit(1)
	}
	return wallet
}

func getUintHash(programHashStr, assetHashStr string) (Uint160, Uint256, error) {
	programHashHex, err := hex.DecodeString(programHashStr)
	if err != nil {
		fmt.Println("Decoding program hash string failed")
		return Uint160{}, Uint256{}, err
	}
	var programHash Uint160
	if err := programHash.Deserialize(bytes.NewReader(programHashHex)); err != nil {
		fmt.Println("Deserialization program hash failed")
		return Uint160{}, Uint256{}, err
	}
	assetHashHex, err := hex.DecodeString(assetHashStr)
	if err != nil {
		fmt.Println("Decoding asset hash string failed")
		return Uint160{}, Uint256{}, err
	}
	var assetHash Uint256
	if err := assetHash.Deserialize(bytes.NewReader(assetHashHex)); err != nil {
		fmt.Println("Deserialization asset hash failed")
		return Uint160{}, Uint256{}, err
	}
	return programHash, assetHash, nil
}

func signTransaction(signer *account.Account, tx *transaction.Transaction) error {
	signature, err := signature.SignBySigner(tx, signer)
	if err != nil {
		fmt.Println("SignBySigner failed")
		return err
	}
	transactionContract, err := contract.CreateSignatureContract(signer.PubKey())
	if err != nil {
		fmt.Println("CreateSignatureContract failed")
		return err
	}
	transactionContractContext := newContractContextWithoutProgramHashes(tx)
	if err := transactionContractContext.AddContract(transactionContract, signer.PubKey(), signature); err != nil {
		fmt.Println("AddContract failed")
		return err
	}
	tx.SetPrograms(transactionContractContext.GetPrograms())
	return nil
}

func makeRegTransaction(admin, issuer *account.Account, name string, value Fixed64) (string, error) {
	asset := &Asset{name, byte(0x00), AssetType(Share), UTXO}
	transactionContract, err := contract.CreateSignatureContract(admin.PubKey())
	if err != nil {
		fmt.Println("CreateSignatureContract failed")
		return "", err
	}
	tx, _ := transaction.NewRegisterAssetTransaction(asset, value, issuer.PubKey(), transactionContract.ProgramHash)
	tx.Nonce = uint64(rand.Int63())
	if err := signTransaction(issuer, tx); err != nil {
		fmt.Println("sign regist transaction failed")
		return "", err
	}
	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialize registtransaction failed")
		return "", err
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}

func makeIssueTransaction(issuer *account.Account, programHashStr, assetHashStr string, value Fixed64) (string, error) {
	programHash, assetHash, err := getUintHash(programHashStr, assetHashStr)
	if err != nil {
		return "", err
	}
	issueTxOutput := &transaction.TxOutput{
		AssetID:     assetHash,
		Value:       value,
		ProgramHash: programHash,
	}
	outputs := []*transaction.TxOutput{issueTxOutput}
	tx, _ := transaction.NewIssueAssetTransaction(outputs)
	tx.Nonce = uint64(rand.Int63())
	if err := signTransaction(issuer, tx); err != nil {
		fmt.Println("sign issue transaction failed")
		return "", err
	}
	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialization of issue transaction failed")
		return "", err
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}

