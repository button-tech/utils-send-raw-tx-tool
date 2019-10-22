package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
)

func SendGrams(c *gin.Context) {
	var (
		tx tx
	)

	err := c.BindJSON(&tx)
	if err != nil {
		c.JSON(404, "bad request")
		return
	}

	if len(tx.Data) == 0 {
		c.JSON(404, "bad request")
		return
	}

	err = sendTon(tx.Data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "query has been sent to the network"})
}

func sendTon(data string) error {

	workdir := os.Getenv("WORKDIR")

	stdout, err := exec.Command(workdir+"wrappers/send_grams.py", data, workdir).Output()
	if err != nil {
		return err
	}

	if string(stdout) == "error\n" {
		return errors.New("Failed")
	}

	return nil
}
