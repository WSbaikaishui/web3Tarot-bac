package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	apiErr "web3Tarot-backend/errors"
)

type BaseResp struct {
	apiErr.ErrorInfo
	Data interface{} `json:"data"`
}

func EncodeError(err error) string {
	var (
		code      int
		errDetail apiErr.ErrorInfo
	)

	if errInfo, ok := err.(apiErr.ErrorInfo); ok {
		code = errInfo.StatusCode()
		errDetail = errInfo
	} else {
		errInfo := apiErr.ErrInternal(err.Error())
		code = errInfo.StatusCode()
		errDetail = errInfo
	}

	return fmt.Sprintf("%d: %d", errDetail.Code, code) + errDetail.Message
}

func EncodeResp(gCtx *gin.Context, data interface{}) {
	resp := BaseResp{}
	resp.Message = "ok"
	if data != nil {
		resp.Data = data
	} else {
		resp.Data = struct{}{}
	}
	gCtx.JSON(http.StatusOK, resp)
}

func EncodeBytes(gCtx *gin.Context, bs []byte) {
	if len(bs) == 0 {
		gCtx.Data(http.StatusOK, "application/json", []byte(`{"code":0, "message":"ok", "data":{}}`))
		return
	}
	data := []byte(`{"code":0, "message":"ok", "data":`)
	data = append(data, bs...)
	data = append(data, '}')
	gCtx.Data(http.StatusOK, "application/json", data)
}

func IsToday(t time.Time) bool {
	return time.Now().Format("2006-01-02") == t.Local().Format("2006-01-02")
}

func CalSeconds(str string) (int64, error) {
	t, err := time.ParseInLocation("15:04:05", str, time.Local)
	if err != nil {
		return 0, err
	}
	return int64(t.Hour()*3600 + t.Minute()*60 + t.Second()), nil
}

var (
	serverErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "errors_total",
			Help:      "Count number of errors that happen during RESTful processing on server side.",
		},
		[]string{"method", "uri", "err_code", "err_message"},
	)
)

func init() {
	prometheus.MustRegister(serverErrors)
}
