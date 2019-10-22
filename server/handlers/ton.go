package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"strings"
)

var workdir = os.Getenv("WORKDIR")

func SigningMessageHash(c *gin.Context) {

	params := struct {
		DestinationAddress string `json:"destinationAddress"`
		Seqno              string `json:"seqno"`
		Amount             string `json:"amount"`
	}{}

	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(401, "bad request")
		return
	}

	stdout, err := exec.Command(workdir+
		"liteclient-build/crypto/fift", "-I",
		workdir+"lite-client/crypto/fift/lib/",
		"-s",
		workdir+"wrappers/signing_message_hash.fif",
		params.DestinationAddress,
		params.Seqno,
		params.Amount).Output()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash := strings.TrimSuffix(string(stdout), " ")
	hash = strings.Replace(hash, "'", "\"", -1)

	if len(hash) != 77 {
		c.JSON(400, "bad request")
		return
	}

	c.JSON(200, gin.H{"hash": hash})
}

func SendGrams(c *gin.Context) {
	var (
		tx tx
	)

	err := c.BindJSON(&tx)
	if err != nil {
		c.JSON(400, "bad request")
		return
	}

	if len(tx.Data) == 0 {
		c.JSON(400, "bad request")
		return
	}

	err = sendGr(tx.Data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "query has been sent to the network"})
}

func sendGr(data string) error {

	stdout, err := exec.Command(workdir+"wrappers/send_grams.py", data, workdir).Output()
	if err != nil {
		return err
	}

	if string(stdout) == "error\n" {
		return errors.New("Failed")
	}

	return nil
}
