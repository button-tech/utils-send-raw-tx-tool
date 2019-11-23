package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/button-tech/logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/stellar/go/clients/horizon"
	"os"
	"strings"
	"time"
)

type tx struct {
	Data     string `json:"data"`
	Currency string `json:"currency,omitempty"`
}

type sendRawTx func(string, string) (string, error)

func Send(c *gin.Context) {

	start := time.Now()

	var (
		tx   tx
		send sendRawTx
	)

	err := c.BindJSON(&tx)
	if err != nil {
		logger.Error("bad request", err.Error())
		c.JSON(404, "bad request")
		return
	}

	currency := strings.ToUpper(tx.Currency)

	switch currency {
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
	case "WAVES":
		send = sendWaves
	default:
		c.JSON(400, "bad request")
		return
	}

	var hash string

	hash, err = send(tx.Data, currency)
	if err != nil {
		hash, err = send(tx.Data, "RESERVE_"+currency)
		if err != nil {
			logger.Error("send raw tx", err.Error(), logger.Params{
				"currency": currency,
			})
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(200, gin.H{"hash": hash})

	logger.LogRequest(time.Since(start), currency, "SendRawTx")
}

func sendEthBased(data, currency string) (string, error) {

	e := os.Getenv(currency)

	c, err := ethclient.Dial(e)
	if err != nil {
		return "", err
	}

	rawTxBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", nil
	}

	tx := new(types.Transaction)

	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return "", nil
	}

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

	e := os.Getenv(currency)

	var request sendRawTx

	if strings.Contains(currency, "RESERVE") {
		request = sendDataGET
	} else {
		request = sendDataPOST
	}

	result, err := request(data, e)
	if err != nil {
		return "", err
	}

	return result, nil
}

func sendWaves(data, _ string) (string, error) {

	url := os.Getenv("WAVES") + "/transactions/broadcast"

	payload := strings.NewReader(data)

	res, err := req.Post(url, req.Header{"Content-Type": "application/json"}, payload)
	if err != nil {
		return "", err
	}

	result := struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}{}

	err = res.ToJSON(&result)
	if err != nil {
		return "", err
	}

	if len(result.Message) != 0 {
		return "", errors.New(result.Message)
	}

	return result.ID, nil

}

// for utxo based
func sendDataPOST(data, endpoint string) (string, error) {

	payload := strings.NewReader("data=" + data)

	res, err := req.Post(endpoint, req.Header{"Content-Type": "application/x-www-form-urlencoded"}, payload)
	if err != nil {
		return "", err
	}

	r := struct {
		Data struct {
			Transaction_hash string `json:"transaction_hash"`
		} `json:"data"`
	}{}

	if res.Response().StatusCode != 200 {
		return "", errors.New("Invalid transaction")
	}

	err = res.ToJSON(&r)
	if err != nil {
		return "", err
	}

	return r.Data.Transaction_hash, nil
}

func sendDataGET(data, endpoint string) (string, error) {

	res, err := req.Get(endpoint + "/sendtx/" + data)
	if err != nil {
		return "", err
	}

	if res.Response().StatusCode != 200 {
		return "", errors.New("Invalid transaction")
	}

	r := struct {
		Result string `json:"result"`
	}{}

	err = res.ToJSON(&r)
	if err != nil {
		return "", err
	}

	return r.Result, nil
}
