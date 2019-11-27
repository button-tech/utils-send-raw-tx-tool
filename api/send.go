package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	"github.com/button-tech/logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/stellar/go/clients/horizon"
	"github.com/valyala/fasthttp"
)

type tx struct {
	Data     string `json:"data"`
	Currency string `json:"currency,omitempty"`
}

type wavesDataTxToSubmit struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type utxoBasedPostDataTxToSubmit struct {
	Data struct {
		TransactionHash string `json:"transaction_hash"`
	} `json:"data"`
}

type utxoBasedGetDataTxToSubmit struct {
	Result string `json:"result"`
}

type h struct {
	Hash string `json:"hash"`
}

type broadcastResponse struct {
	Code int    `json:"code"`
	Hash string `json:"hash"`
	Ok   bool   `json:"ok"`
}

type xrpDataTxToSubmit struct {
	Method string `json:"method"`
	Params []struct {
		TxBlob string `json:"tx_blob"`
	} `json:"params"`
}

type xrpSentTxInfo struct {
	Result struct {
		EngineResult        string `json:"engine_result"`
		EngineResultCode    int    `json:"engine_result_code"`
		EngineResultMessage string `json:"engine_result_message"`
		Status              string `json:"status"`
		TxJSON              struct {
			Fee  string `json:"Fee"`
			Hash string `json:"hash"`
		} `json:"tx_json"`
	} `json:"result"`
}

type tronDataTxToSubmit struct {
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Amount       int    `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
}

type tronSentResult struct {
	Result bool `json:"result"`
}

type cosmosSentTxInfo struct {
	CheckTx struct {
		Code      int      `json:"code"`
		Data      string   `json:"data"`
		Log       string   `json:"log"`
		GasUsed   int      `json:"gas_used"`
		GasWanted int      `json:"gas_wanted"`
		Info      string   `json:"info"`
		Tags      []string `json:"tags"`
	} `json:"check_tx"`
	DeliverTx struct {
		Code      int      `json:"code"`
		Data      string   `json:"data"`
		Log       string   `json:"log"`
		GasUsed   int      `json:"gas_used"`
		GasWanted int      `json:"gas_wanted"`
		Info      string   `json:"info"`
		Tags      []string `json:"tags"`
	} `json:"deliver_tx"`
	Hash   string `json:"hash"`
	Height int    `json:"height"`
}

const submitMethod = "submit"

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
		logger.Error("send", err.Error())
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
	case "BNB":
		send = sendBnB
	case "XRP":
		send = sendXRP
	case "TRON":
		send = sendTron
	case "COSMOS":
		send = sendCosmos
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

	var tx types.Transaction
	if err = rlp.DecodeBytes(rawTxBytes, &tx); err != nil {
		return "", err
	}

	if err = c.SendTransaction(context.Background(), &tx); err != nil {
		return "", err
	}
	hexedHash := tx.Hash().Hex()

	return hexedHash, nil
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

	var r wavesDataTxToSubmit
	if err = res.ToJSON(&r); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}

	if len(r.Message) != 0 {
		return "", errors.New(r.Message)
	}
	hash := r.ID
	return hash, nil
}

// for utxo based
func sendDataPOST(data, endpoint string) (string, error) {
	rq := req.New()
	payload := strings.NewReader("data=" + data)
	res, err := rq.Post(endpoint, req.Header{"Content-Type": "application/x-www-form-urlencoded"}, payload)
	if err != nil {
		return "", err
	}

	if res.Response().StatusCode != 200 {
		return "", errors.New("Invalid transaction")
	}

	var r utxoBasedPostDataTxToSubmit
	if err = res.ToJSON(&r); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}
	hash := r.Data.TransactionHash
	return hash, nil
}

func sendDataGET(data, endpoint string) (string, error) {
	rq := req.New()
	res, err := rq.Get(endpoint + "/sendtx/" + data)
	if err != nil {
		return "", err
	}

	if res.Response().StatusCode != 200 {
		return "", errors.Wrap(errors.New("invalid transaction"), "sendDataGET")
	}

	var r utxoBasedGetDataTxToSubmit
	if err = res.ToJSON(&r); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}
	hash := r.Result
	return hash, nil
}

func sendBnB(data, currency string) (string, error) {
	e := os.Getenv(currency)

	rq := req.New()
	resp, err := rq.Post(e, req.Header{"Content-type": "text/plain"}, data)
	if err != nil {
		return "", err
	}

	broadcasted := make([]broadcastResponse, 1)
	if err := resp.ToJSON(&broadcasted); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}
	hash := broadcasted[0].Hash
	return hash, nil
}

func sendXRP(data, currency string) (string, error) {
	broadcasted, err := submitXRPTx(data, currency)
	if err != nil {
		return "", err
	}

	if err := checkSubmitXRPTxStatus(broadcasted); err != nil {
		return "", err
	}

	hash := broadcasted.Result.TxJSON.Hash
	return hash, nil
}

func submitXRPTx(data, currency string) (*xrpSentTxInfo, error) {
	e := os.Getenv(currency)

	rq := req.New()
	resp, err := rq.Post(e, req.BodyJSON(xrpTxToSubmit(data)))
	if err != nil {
		return nil, errors.Wrap(err, "submitXRPTxRequest")
	}

	if resp.Response().StatusCode != 200 {
		return nil, errors.Wrap(errors.New("StatusCodeNotOk"), "Request to Ripple")
	}

	var info xrpSentTxInfo
	if err := resp.ToJSON(&info); err != nil {
		return nil, errors.Wrap(err, "XRPtoJSON")
	}
	return &info, nil
}

func xrpTxToSubmit(txBlob string) *xrpDataTxToSubmit {
	return &xrpDataTxToSubmit{
		Method: submitMethod,
		Params: []struct {
			TxBlob string `json:"tx_blob"`
		}{
			{
				TxBlob: txBlob,
			},
		},
	}
}

func checkSubmitXRPTxStatus(info *xrpSentTxInfo) error {
	if info.Result.Status == "error" {
		return errors.Wrap(errors.New("ResponseStatusError"), "RippleAPI")
	}
	return nil
}

func sendTron(data, currency string) (string, error) {
	ok, err := submitTronTx(data, currency)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("tx reverse")
	}

	var t tronDataTxToSubmit
	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return "", errors.Wrap(err, "sendTron")
	}

	hash := t.TxID
	return hash, nil
}

func submitTronTx(data, currency string) (bool, error) {
	e := os.Getenv(currency)

	rq := req.New()
	resp, err := rq.Post(e, req.BodyJSON(&data))
	if err != nil {
		return false, errors.Wrap(err, "submitTronTx")
	}

	var r tronSentResult
	if err = resp.ToJSON(&r); err != nil {
		return false, errors.Wrap(err, "XRPtoJSON")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return false, errors.Wrap(errors.New("statusResponseNotOk"), "submitTronTx")
	}
	ok := r.Result

	return ok, nil
}

func sendCosmos(data, currency string) (string, error) {
	return submitCosmosTx(data, currency)
}

func submitCosmosTx(data, currency string) (string, error) {
	e := os.Getenv(currency)

	rq := req.New()
	resp, err := rq.Post(e, req.BodyJSON(&data))
	if err != nil {
		return "", errors.Wrap(err, "submitCosmosTx")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return "", errors.Wrap(errors.New("responseStatusNotOk"), "submitCosmosTx")
	}

	var info cosmosSentTxInfo
	if err = resp.ToJSON(&info); err != nil {
		return "", errors.Wrap(err, "COSMOStoJSON")
	}
	hash := info.Hash

	return hash, nil
}
