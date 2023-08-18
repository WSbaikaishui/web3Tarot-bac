package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"net/http"
	"strconv"
	apiErr "web3Tarot-backend/errors"

	"web3Tarot-backend/service/transaction"
	"web3Tarot-backend/setting"
	"web3Tarot-backend/util"
)

func Recharge(ctx *gin.Context) {
	wa := ctx.Value(util.AuthKey).(string)

	var param transaction.ChargeData
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	if param.Address != wa {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter("invalid address"))
		return
	}

	if transaction.IsTransactionExist(param.TxHash) {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter("transaction already used"))
		return
	}

	err := CheckParam(&param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	err = transaction.ReCharge(param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func CheckParam(param *transaction.ChargeData) error {
	client := rpc.NewClient(setting.AppSetting.Neo3NodeUrl)

	tx := client.GetApplicationLog(param.TxHash)
	notification := tx.Result.Executions[0].Notifications[0]
	if tx.HasError() {
		return fmt.Errorf("analyse tx err :%v", tx.Error)
	}

	if notification.EventName == "Transfer" {
		if v, ok := notification.State.Value.([]interface{}); ok {
			var result []map[string]interface{}
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					result = append(result, m)
				}
			}
			if len(result) != 3 {
				return fmt.Errorf("invalid notification")
			}
			var result2 []string
			for _, item := range result {
				if m, ok := item["value"].(string); ok {
					result2 = append(result2, m)
				}
			}

			sender, err := crypto.Base64Decode(result2[0])
			sender2, err := helper.UInt160FromString(helper.BytesToHex(helper.ReverseBytes(sender)))
			sender3 := crypto.ScriptHashToAddress(sender2, helper.DefaultAddressVersion)
			if err != nil {
				return err
			}

			recipient, err := crypto.Base64Decode(result2[1])
			recipient2, err := helper.UInt160FromString(helper.BytesToHex(helper.ReverseBytes(recipient)))
			if err != nil {
				return err
			}
			recipient3 := crypto.ScriptHashToAddress(recipient2, helper.DefaultAddressVersion)

			amount, err := strconv.Atoi(result2[2])
			if err != nil {
				return err
			}
			if sender3 == param.Address && recipient3 == setting.AppSetting.ContractAddress && amount == param.Amount {
				return nil
			} else {
				return fmt.Errorf("invalid params")
			}
		}
	} else {
		return fmt.Errorf("this is not a transfer GAS transaction")
	}
	return nil
}
