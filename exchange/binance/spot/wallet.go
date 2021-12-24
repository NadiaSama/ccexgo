package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
)

type (
	DepositHisrecReq struct {
		*binance.RestReq
	}

	DepositHisrecRecord struct {
		Amount        string `json:"amount"`
		Coin          string `json:"coin"`
		Network       string `json:"network"`
		Status        int    `json:"status"`
		Address       string `json:"address"`
		AddressTag    string `json:"addressTag"`
		TxID          string `json:"txId"`
		InsertTime    int64  `json:"insertTime"`
		TransferType  int    `json:"transferType"`
		UnlockConfirm int    `json:"unlockConfirm"`
		ConfirmTimes  string `json:"confirmTimes"`
	}
)

const (
	DepositHisrecEndPoint = "/sapi/v1/capital/deposit/hisrec"
)

func NewDepositHisrecReq() *DepositHisrecReq {
	return &DepositHisrecReq{
		binance.NewRestReq(),
	}
}

func (rc *RestClient) DepositHisrec(ctx context.Context, req *DepositHisrecReq) ([]DepositHisrecRecord, error) {
	var ret []DepositHisrecRecord
	if err := rc.GetRequest(ctx, DepositHisrecEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "reqeust deposit hisrec fail")
	}
	return ret, nil
}
