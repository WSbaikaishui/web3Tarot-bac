package util

import (
	"encoding/json"
	"errors"
	"strconv"
)

func GetTxAmount(address string, txHash string) (int, error) {
	url := "https://testnet.toncenter.com/api/v2/getTransactions?address=EQDG5ALYstmhiMhpSN9j6tEHmLq2Owx7yyfTH1W15wS0NqZH&limit=100&to_lt=0&archival=true&api_key=0e682c2f19e7e6970846fb50d975e8095e292db2b55afb64750a5470a1582431"
	res, err := Get(url)
	if err != nil {
		return 0, err
	}
	var re Transaction
	err = json.Unmarshal(res, &re)
	value := 0
	for _, v := range re.Result {
		if len(v.OutMsgs) == 0 {
			if v.TransactionID.Hash == txHash && v.InMsg.Source == address {
				val, _ := strconv.Atoi(v.InMsg.Value)
				value = val
				break
			}
		}
	}
	if value == 0 {
		return 0, errors.New("无效的tx")
	}

	return value, nil
}

// 创建与 JSON 结构相匹配的 Go 结构体
type Transaction struct {
	OK     bool `json:"ok"`
	Result []struct {
		Type    string `json:"@type"`
		Address struct {
			Type           string `json:"@type"`
			AccountAddress string `json:"account_address"`
		} `json:"address"`
		Utime         int64  `json:"utime"`
		Data          string `json:"data"`
		TransactionID struct {
			Type string `json:"@type"`
			Lt   string `json:"lt"`
			Hash string `json:"hash"`
		} `json:"transaction_id"`
		Fee        string `json:"fee"`
		StorageFee string `json:"storage_fee"`
		OtherFee   string `json:"other_fee"`
		InMsg      struct {
			Type        string `json:"@type"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
			Value       string `json:"value"`
			FwdFee      string `json:"fwd_fee"`
			IhrFee      string `json:"ihr_fee"`
			CreatedLt   string `json:"created_lt"`
			BodyHash    string `json:"body_hash"`
			MsgData     struct {
				Type string `json:"@type"`
				Text string `json:"text"`
			} `json:"msg_data"`
			Message string `json:"message"`
		} `json:"in_msg"`
		OutMsgs []struct {
			Type        string `json:"@type"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
			Value       string `json:"value"`
			FwdFee      string `json:"fwd_fee"`
			IhrFee      string `json:"ihr_fee"`
			CreatedLt   string `json:"created_lt"`
			BodyHash    string `json:"body_hash"`
			MsgData     struct {
				Type      string `json:"@type"`
				Body      string `json:"body"`
				InitState string `json:"init_state"`
			} `json:"msg_data"`
			Message string `json:"message"`
		} `json:"out_msgs"`
	} `json:"result"`
}
