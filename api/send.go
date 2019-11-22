package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/imroc/req"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/stellar/go/clients/horizon"
	"github.com/valyala/fasthttp"
)

type tx struct {
	Data     string `json:"data"`
	Currency string `json:"currency,omitempty"`
}

type h struct {
	Hash string `json:"hash"`
}

type sendRawTx func(string, string) (string, error)

func sendHandler(ctx *routing.Context) error {
	var tx tx
	if err := json.Unmarshal(ctx.PostBody(), &tx); err != nil {
		return err
	}

	currency := strings.ToUpper(tx.Currency)
	send := sendBased(currency)
	if send == nil {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, "incorrect currency")
		return nil
	}

	hash, err := send(tx.Data, currency)
	if err == nil {
		respondWithJSON(ctx, fasthttp.StatusOK, hash)
		return nil
	}
	hash, err = send(tx.Data, "RESERVE_"+currency)
	if err != nil {
		respondWithJSON(ctx, fasthttp.StatusInternalServerError, err.Error())
		return nil
	}

	respondWithJSON(ctx, fasthttp.StatusOK, h{Hash: hash})
	return nil
}

func sendBased(currency string) (send sendRawTx) {
	switch currency {
	case "ETH", "ETC":
		send = sendEthBased
	case "XLM":
		send = sendXlm
	case "BTC", "BCH", "LTC":
		send = sendUtxoBased
	case "WAVES":
		send = sendWaves
	default:
		send = nil
	}
	return
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
		return "", errors.New("invalid transaction")
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
