package deribit

import "github.com/shopspring/decimal"

type (
	TradeResult struct {
		Amount          decimal.Decimal `json:"amount"`
		BlockTradeID    string          `json:"block_trade_id"`
		Direction       string          `json:"direction"`
		Fee             decimal.Decimal `json:"fee"`
		FeeCurrency     string          `json:"fee_currency"`
		IndexPrice      decimal.Decimal `json:"index_price"`
		InstrumentName  string          `json:"instrument_name"`
		IV              decimal.Decimal `json:"iv"`
		Label           string          `json:"label"`
		Liquidation     string          `json:"liquidation"`
		Liquidity       string          `json:"liquidity"`
		MarkPrice       decimal.Decimal `json:"mark_price"`
		MatchingID      string          `json:"matching_id"`
		OrderID         string          `json:"order_id"`
		OrderType       string          `json:"order_type"`
		PostOnly        bool            `json:"post_only"`
		Price           decimal.Decimal `json:"price"`
		ProfitLoss      decimal.Decimal `json:"profit_loss"`
		ReduceOnly      bool            `json:"reduce_only"`
		SelfTrade       bool            `json:"self_trade"`
		State           string          `json:"state"`
		TickDirection   int             `json:"tick_direction"`
		Timestamp       int64           `json:"timestamp"`
		TradeID         string          `json:"trade_id"`
		TradeSeq        int64           `json:"trade_seq"`
		UnderlyingPrice decimal.Decimal `json:"underlying_price"`
	}

	TradeResp struct {
		HasMore bool          `json:"has_more"`
		Trades  []TradeResult `json:"trades"`
	}

	GetUserTradesByCurrencyReq struct {
		AuthToken
		*clientReq
	}
)

const (
	GetUserTradesByCurrency = "private/get_user_trades_by_currency"
)

func NewGetUserTradesByCurrencyReq(currency string) *GetUserTradesByCurrencyReq {
	cr := newClientReq()
	cr.addField("currency", currency)

	return &GetUserTradesByCurrencyReq{
		clientReq: cr,
	}
}

func (gr *GetUserTradesByCurrencyReq) MarshalJSON() ([]byte, error) {
	gr.addField("access_token", gr.AuthToken.AccessToken)
	return gr.clientReq.MarshalJSON()
}

func (gr *GetUserTradesByCurrencyReq) Kind(k string) *GetUserTradesByCurrencyReq {
	gr.addField("kind", k)
	return gr
}

func (gr *GetUserTradesByCurrencyReq) StartID(id string) *GetUserTradesByCurrencyReq {
	gr.addField("start_id", id)
	return gr
}

func (gr *GetUserTradesByCurrencyReq) EndID(id string) *GetUserTradesByCurrencyReq {
	gr.addField("end_id", id)
	return gr
}

func (gr *GetUserTradesByCurrencyReq) Count(ct int) *GetUserTradesByCurrencyReq {
	gr.addField("count", ct)
	return gr
}

func (gr *GetUserTradesByCurrencyReq) IncludeOld(io bool) *GetUserTradesByCurrencyReq {
	gr.addField("include_old", io)
	return gr
}

func (gr *GetUserTradesByCurrencyReq) Sorting(sr string) *GetUserTradesByCurrencyReq {
	gr.addField("sorting", sr)
	return gr
}
