package handlers

import (
	"context"
	"encoding/hex"
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
	RawTx string `json:"rawTx"`
}

type sendRawTx func(string, string) (string, error)

func Send(c *gin.Context) {

	var (
		tx       tx
		send     sendRawTx
		currency = c.Param("currency")
	)

	err := c.BindJSON(&tx)
	if err != nil {
		c.JSON(404, "bad request")
		return
	}

	switch currency {
	case "eth":
		send = sendEthBased
	case "etc":
		send = sendEthBased
	case "xlm":
		send = sendXlm
	case "btc":
		send = sendUtxoBased
	case "bch":
		send = sendUtxoBased
	case "ltc":
		send = sendUtxoBased
	}

	hash, err := send(tx.RawTx, currency)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"hash": hash})
}

func sendEthBased(rawTx, currency string) (string, error) {

	endpoint := os.Getenv(strings.ToUpper(currency))

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return "", err
	}

	rawTxBytes, err := hex.DecodeString(rawTx)

	tx := new(types.Transaction)

	rlp.DecodeBytes(rawTxBytes, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func sendXlm(rawTx string, _ string) (string, error) {
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(rawTx)
	if err != nil {
		return "", err
	}

	return resp.Hash, nil
}

func sendUtxoBased(rawTx, currency string) (string, error) {

	endpoint := os.Getenv(strings.ToUpper(currency))

	payload := strings.NewReader("data=" + rawTx)

	res, err := req.Post(endpoint, req.Header{"Content-Type": "application/x-www-form-urlencoded"}, payload)
	if err != nil {
		return "", err
	}

	result := struct {
		Data struct {
			Transaction_hash string `json:"transaction_hash"`
		} `json:"data"`
	}{}

	err = res.ToJSON(&result)
	if err != nil {
		return "", err
	}

	return result.Data.Transaction_hash, nil
}
