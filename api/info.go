package api

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

type txDescription struct {
	Currency     string `json:"currency"`
	SignedTxType string `json:"signedTxType"`
	Example      tx     `json:"example"`
}

type handlersDescriptions struct {
	Send      []txDescription `json:"/send"`
	SendGrams txDescription   `json:"/sendGrams"`
}

var info = handlersDescriptions{Send: []txDescription{
	{"ETH Based(ETH, ETC)", "hex", tx{"f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772", "ETH"}},
	{"XLM", "base64", tx{"AAAAAIoNf6rpFNnSNVDkeb+HA8fYqSEbiEO9ltU1qWHtNia2AAAAZAFZFRMAAABHAAAAAQAAAAAAAAAAAAAJGKvhY4IAAAABAAAADUJVVFRPTiBXYWxsZXQAAAAAAAABAAAAAAAAAAEAAAAAig1/qukU2dI1UOR5v4cDx9ipIRuIQ72W1TWpYe02JrYAAAAAAAAAAAAAAGQAAAAAAAAAAe02JrYAAABAsuCOvG7ncDNM2J2xsxJJDZvVzT8eiFRgRCR8xAa1xvWde8kkiTq8IET7av3feEb2h3rMM+q9o+zx+2A2exkxBA==", "XLM"}},
	{"UTXO Based(BTC, BCH, LTC)", "hex", tx{"0100000001e84509c3a8fb1ee3685c96f12d5ffacf608960b8cda2cabe269a7f1ddfe153c2010000006b483045022100a70718365a3e29f05be8033f5b659d06e530fe03ad5d93f9a989ce59746493850220369b1c1997436c28775d9b8f93f46988143c9e0361644e57ea87936d9be6774c012102bcbe5228dc72dd3babff5a159bbbf49a515ccfc37952af194316abfa249db0acffffffff02e8030000000000001976a9143ac991a4209cf14c980c726339148ed470464fff88ac0a490100000000001976a9143ac991a4209cf14c980c726339148ed470464fff88ac00000000", "BTC"}},
	{"WAVES", "json", tx{`{"senderPublicKey":"5sgLhwTbDZhUDuVJoM4uxnAz9AXiiX5v5zLEePqnz73F","recipient":"address:3PDn2Sqwdz7Zbj6PJcNniRYKdLR3U3DJabR","assetId":"","amount":1000,"feeAssetId":"","fee":100000,"attachment":"","timestamp":1567605589000,"signature":"4zvtuqJh5AWZzzkuouh1ypXmEAKPciRZyzqoB7e86ycKi6k7R5XfSKkmiAXYrb6DWh7sNGNBAMp8pTWEEqD26xDu","type":4}`, "WAVES"}},
	{"BNB", "json", tx{"0x31033757117cf38040ab70485f7e247c75eb8b2074305f01a4cf8b41a3e940fe", "BNB"}},
	{"XRP", "json", tx{"1200002280000000240000000361D4838D7EA4C6800000000000000000000000000055534400000000004B4E9C06F24296074F7BC48F92A97916C6DC5EA968400000000000000A732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB74473045022100D184EB4AE5956FF600E7536EE459345C7BBCF097A84CC61A93B9AF7197EDB98702201CEA8009B7BEEBAA2AACC0359B41C427C1C5B550A4CA4B80CF2174AF2D6D5DCE81144B4E9C06F24296074F7BC48F92A97916C6DC5EA983143E9D4A2B8AA0780F682D136F7A56D6724EF53754", "XRP"}},
	{"TRON", "json", tx{`{"signature:["97c825b41c77de2a8bd65b3df55cd4c0df59c307c0187e42321dcc1cc455ddba583dd9502e17cfec5945b34cad0511985a6165999092a6dec84c2bdd97e649fc01"],"txID":"454f156bf1256587ff6ccdbc56e64ad0c51e4f8efea5490dcbc720ee606bc7b8","raw_data":{"contract":[{"parameter":{"value":{"amount":1000,"owner_address":"41e552f6487585c2b58bc2c9bb4492bc1f17132cd0","to_address":"41d1e7a6bc354106cb410e65ff8b181c600ff14292"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"267e","ref_block_hash":"9a447d222e8de9f2","expiration":1530893064000,"timestamp":1530893006233}}`, "TRON"}},
	{"COSMOS", "json", tx{`{"tx": {"msg": ["string"],"fee": {"gas": "string","amount": [{"denom": "stake","amount": "50"}]},"memo": "string","signature": {"signature": "MEUCIQD02fsDPra8MtbRsyB1w7bqTM55Wu138zQbFcWx4+CFyAIge5WNPfKIuvzBZ69MyqHsqD8S1IwiEp+iUb6VSdtlpgY=","pub_key": {"type": "tendermint/PubKeySecp256k1","value": "Avz04VhtKJh8ACCVzlI8aTosGy0ikFXKIVHQ3jKMrosH"},"account_number": "0","sequence": "0"}},"mode": "block"}`, "COSMOS"}},

	// todo: need to approve
	{"ALGORAND", "binary", tx{"", "ALGORAND"}},
},
	SendGrams: txDescription{"TON", "hex", tx{Data: "B5EE9C724101020100A50001CF88015F22047CC96B03CA2F33C363340E4550317A32D4C6BF58773E799BD6ECB67F2607DFEAA63832274E15FCB4768E6E5C5A31D7D0F9380B3CFEDC31E7DECC1FA065D13A79909867EA366E2B7DCBFF0B05CF91D7B986383B9A46106270DF6F079C6058000000180C010070620072379E8B9816F48D2DC83D6F9995D602C590F305D64893A9BFB581B0260A59EB21DCD65000000000000000000000000000005445535437D8FAE2"}},
}

var infoResponse []byte

func infoHandler(ctx *routing.Context) error {
	ctx.SetContentType("application/json")
	if _, err := ctx.Write(infoResponse); err != nil {
		return err
	}
	return nil
}
