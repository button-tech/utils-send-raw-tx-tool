package handlers

import (
	"github.com/gin-gonic/gin"
)

type txDescription struct {
	Currency string `json:"currency"`
	SignedTxType string `json:"signedTxType"`
	Example tx `json:"example"`
}

type descriptions struct {
	Data []txDescription `json:"data"`
}

var info = descriptions{Data:[]txDescription{
	{"ETH","hex", tx{"f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"}},
	{"ETC","hex", tx{"f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"}},
	{"XLM","base64", tx{"AAAAAIoNf6rpFNnSNVDkeb+HA8fYqSEbiEO9ltU1qWHtNia2AAAAZAFZFRMAAABHAAAAAQAAAAAAAAAAAAAJGKvhY4IAAAABAAAADUJVVFRPTiBXYWxsZXQAAAAAAAABAAAAAAAAAAEAAAAAig1/qukU2dI1UOR5v4cDx9ipIRuIQ72W1TWpYe02JrYAAAAAAAAAAAAAAGQAAAAAAAAAAe02JrYAAABAsuCOvG7ncDNM2J2xsxJJDZvVzT8eiFRgRCR8xAa1xvWde8kkiTq8IET7av3feEb2h3rMM+q9o+zx+2A2exkxBA=="}},
}}

func GetInfo(c *gin.Context){
	c.JSON(200, &info)
}