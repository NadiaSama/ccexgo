package deribit

import "github.com/shopspring/decimal"

type (
	AccountSummaryResp struct {
		Balance                    decimal.Decimal `json:"balance"`
		OptionsSessionUPL          decimal.Decimal `json:"options_session_upl"`
		DepositAddress             string          `json:"deposit_address"`
		OptionGamma                decimal.Decimal `json:"option_gamma"`
		OptionTheta                decimal.Decimal `json:"option_theta"`
		UserName                   string          `json:"user"`
		Equity                     decimal.Decimal `json:"equity"`
		Type                       string          `json:"type"`
		Currency                   string          `json:"currency"`
		DeltaTotal                 decimal.Decimal `json:"delta_total"`
		FuturesSessionRPL          decimal.Decimal `json:"futures_session_rpl"`
		PortfolioMarginingEnabled  bool            `json:"portfolio_margining_enabled"`
		TotalPL                    decimal.Decimal `json:"total_pl"`
		MarginBalance              decimal.Decimal `json:"margin_balance"`
		TfaEnabled                 bool            `json:"tfa_enabled"`
		OptionsSessionRPL          decimal.Decimal `json:"options_session_rpl"`
		OptionsDelta               decimal.Decimal `json:"options_delta"`
		FuturesPL                  decimal.Decimal `json:"futures_pl"`
		ID                         int             `json:"id"`
		SessionUPL                 decimal.Decimal `json:"session_upl"`
		AvailableWithdrawalFunds   decimal.Decimal `json:"available_withdrawal_funds"`
		CreationTimestmap          int64           `json:"creation_timestamp"`
		OptionsPL                  decimal.Decimal `json:"options_pl"`
		SystemName                 string          `json:"system_name"`
		InitialMargin              decimal.Decimal `json:"initial_margin"`
		ProjectedInitialMargin     decimal.Decimal `json:"projected_initial_margin"`
		MaintenanceMargin          decimal.Decimal `json:"maintenance_margin"`
		ProjectedMaintenanceMargin decimal.Decimal `json:"projected_maintenance_margin"`
		SessinRPL                  decimal.Decimal `json:"session_rpl"`
		InteruserTransfersEnabled  bool            `json:"interuser_transfers_enabled"`
		OptionsVega                decimal.Decimal `json:"options_vega"`
		ProjectedDeltaTotal        decimal.Decimal `json:"projectd_delta_total"`
		Email                      string          `json:"email"`
		FuturesSessionUPL          decimal.Decimal `json:"futures_session_upl"`
		AvailableFunds             decimal.Decimal `json:"available_funds"`
		OptionsValue               decimal.Decimal `json:"options_value"`
	}

	AccountSummaryReq struct {
		AuthToken
		Currency string `json:"currency"`
		Extended bool   `json:"extended"`
	}
)

const (
	AccountSummaryMethod = "private/get_account_summary"
)
