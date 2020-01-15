package api

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type status struct {
	Status string `json:"status"`
}

func sendGramsHandler(ctx *routing.Context) error {
	var tx tx
	if err := json.Unmarshal(ctx.PostBody(), &tx); err != nil {
		return err
	}

	if len(tx.Data) == 0 {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, "data is empty")
		return nil
	}

	if err := sendGr(tx.Data); err != nil {
		respondWithJSON(ctx, fasthttp.StatusInternalServerError, err.Error())
		return nil
	}

	respondWithJSON(ctx, fasthttp.StatusOK, status{Status: "query has been sent to the network"})
	return nil
}

func sendGr(data string) error {
	stdout, err := exec.Command("/app/wrappers/send_grams.py", data).Output()
	if err != nil {
		return err
	}

	if string(stdout) == "error\n" {
		return errors.Wrap(errors.New("failed"), "sendGrams")
	}
	return nil
}

func signingMsgHashHandler(ctx *routing.Context) error {
	params := struct {
		DestinationAddress string `json:"destinationAddress"`
		Seqno              string `json:"seqno"`
		Amount             string `json:"amount"`
	}{}

	if err := json.Unmarshal(ctx.PostBody(), &params); err != nil {
		return err
	}

	stdout, err := exec.Command(
		"/app/liteclient-build/crypto/fift", "-I",
		"/app/lite-client/crypto/fift/lib/",
		"-s",
		"/app/wrappers/signing_message_hash.fif",
		params.DestinationAddress,
		params.Seqno,
		params.Amount).Output()
	if err != nil {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, err.Error())
		return nil
	}

	hash := strings.TrimSuffix(string(stdout), " ")
	hash = strings.Replace(hash, "'", "\"", -1)

	if len(hash) != 77 {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, "badRequest")
		return nil
	}

	respondWithJSON(ctx, fasthttp.StatusOK, h{Hash: hash})
	return nil
}
