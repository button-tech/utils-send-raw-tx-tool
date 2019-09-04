package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/stellar/go/clients/horizon"
	"os"
	"strings"
)

type tx struct {
	Data     string `json:"data"`
	Currency string `json:"currency"`
}

type sendRawTx func(string, string) (string, error)

func Send(c *gin.Context) {

	var (
		tx   tx
		send sendRawTx
	)

	err := c.BindJSON(&tx)
	if err != nil {
		c.JSON(404, "bad request")
		return
	}

	switch tx.Currency {
	case "ETH":
		send = sendEthBased
	case "ETC":
		send = sendEthBased
	case "XLM":
		send = sendXlm
	case "BTC":
		send = sendUtxoBased
	case "BCH":
		send = sendUtxoBased
	case "LTC":
		send = sendUtxoBased
	}

	hash, err := send(tx.Data, tx.Currency)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"hash": hash})
}

func sendEthBased(data, currency string) (string, error) {

	endpoint := os.Getenv(currency)

	c, err := ethclient.Dial(endpoint)
	if err != nil {
		return "", err
	}

	rawTxBytes, err := hex.DecodeString(data)

	tx := new(types.Transaction)

	rlp.DecodeBytes(rawTxBytes, &tx)

	err = c.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func sendXlm(data, _ string) (string, error) {

	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(data)
	if err != nil {
		return "", err
	}

	return resp.Hash, nil
}

func sendUtxoBased(data, currency string) (string, error) {

	endpoint := os.Getenv(currency)

	payload := strings.NewReader("data=" + data)

	res, err := req.Post(endpoint, req.Header{"Content-Type": "application/x-www-form-urlencoded"}, payload)
	if err != nil {
		return "", err
	}

	result := struct {
		Data struct {
			Transaction_hash string `json:"transaction_hash"`
		} `json:"data"`
	}{}

	if res.Response().StatusCode != 200 {
		return "", errors.New("Invalid transaction")
	}

	err = res.ToJSON(&result)
	if err != nil {
		return "", err
	}

	return result.Data.Transaction_hash, nil
}
//
//func sendWaves(data, _ string) (string, error) {
//	jsonValue, _ := json.Marshal(data)
//
//	res, err := req.Post(os.Getenv("WAVES")+"/transactions/broadcast", req.Param{"Content-type": "application/json"}, jsonValue)
//	if err != nil {
//		return "", nil
//	}
//
//	return res.String(), nil
//}
