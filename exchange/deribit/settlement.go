package deribit

import "github.com/shopspring/decimal"

type (
	SettlementType string
	SettlementResp struct {
		Settlements  []Settlement `json:"settlements"`
		Continuation string       `json:"continuation"`
	}

	Settlement struct {
		Type              string          `json:"type"`
		Timestamp         int64           `json:"timestamp"`
		SessionProfitLoss decimal.Decimal `json:"session_profit_loss"`
		ProfitLoss        decimal.Decimal `json:"profit_loss"`
		Position          decimal.Decimal `json:"position"`
		MarkPrice         decimal.Decimal `json:"mark_price"`
		InstrumentName    string          `json:"instrument_name"`
		IndexPrice        decimal.Decimal `json:"index_price"`
	}

	SettlementReq struct {
		AuthToken
		InstrumentName string         `json:"instrument_name"`
		Type           SettlementType `json:"type"`
	}

	PublicSettlementByInstrumentReq struct {
		InstrumentName string         `json:"instrument_name"`
		Type           SettlementType `json:"type"`
	}
)

const (
	SettlementMethodByInstrument       = "private/get_settlement_history_by_instrument"
	PublicSettlementMethodByInstrument = "public/get_last_settlements_by_instrument"

	SettlementTypeSettlement SettlementType = "settlement"
	SettlementTypeDelivery   SettlementType = "delivery"
	SettlementTypeBankrupcty SettlementType = "bankruptcy"
)
