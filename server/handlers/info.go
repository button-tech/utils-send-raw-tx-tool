package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type txDescription struct {
	Currency     string `json:"currency"`
	SignedTxType string `json:"signedTxType"`
	Example      tx     `json:"example"`
}

type descriptions struct {
	Info []txDescription `json:"info"`
}

var info = descriptions{Info: []txDescription{
	{"ETH Based(ETH, ETC)", "hex", tx{"f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772", "ETH"}},
	{"XLM", "base64", tx{"AAAAAIoNf6rpFNnSNVDkeb+HA8fYqSEbiEO9ltU1qWHtNia2AAAAZAFZFRMAAABHAAAAAQAAAAAAAAAAAAAJGKvhY4IAAAABAAAADUJVVFRPTiBXYWxsZXQAAAAAAAABAAAAAAAAAAEAAAAAig1/qukU2dI1UOR5v4cDx9ipIRuIQ72W1TWpYe02JrYAAAAAAAAAAAAAAGQAAAAAAAAAAe02JrYAAABAsuCOvG7ncDNM2J2xsxJJDZvVzT8eiFRgRCR8xAa1xvWde8kkiTq8IET7av3feEb2h3rMM+q9o+zx+2A2exkxBA==", "XLM"}},
	{"UTXO Based(BTC, BCH, LTC)", "hex", tx{"0100000001e84509c3a8fb1ee3685c96f12d5ffacf608960b8cda2cabe269a7f1ddfe153c2010000006b483045022100a70718365a3e29f05be8033f5b659d06e530fe03ad5d93f9a989ce59746493850220369b1c1997436c28775d9b8f93f46988143c9e0361644e57ea87936d9be6774c012102bcbe5228dc72dd3babff5a159bbbf49a515ccfc37952af194316abfa249db0acffffffff02e8030000000000001976a9143ac991a4209cf14c980c726339148ed470464fff88ac0a490100000000001976a9143ac991a4209cf14c980c726339148ed470464fff88ac00000000", "BTC"}},
	{"WAVES", "json", tx{"{\"senderPublicKey\":\"5sgLhwTbDZhUDuVJoM4uxnAz9AXiiX5v5zLEePqnz73F\",\"recipient\":\"address:3PDn2Sqwdz7Zbj6PJcNniRYKdLR3U3DJabR\",\"assetId\":\"\",\"amount\":1000,\"feeAssetId\":\"\",\"fee\":100000,\"attachment\":\"\",\"timestamp\":1567605589000,\"signature\":\"4zvtuqJh5AWZzzkuouh1ypXmEAKPciRZyzqoB7e86ycKi6k7R5XfSKkmiAXYrb6DWh7sNGNBAMp8pTWEEqD26xDu\",\"type\":4}", "WAVES"}},
}}

var result []byte

func init() {
	var err error
	result, err = json.MarshalIndent(&info, "", " ")
	if err != nil {
		log.Fatal(err)
	}
}

func GetInfo(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)
	_, err := c.Writer.Write(result)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
	}
}
