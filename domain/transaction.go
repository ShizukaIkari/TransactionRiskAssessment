package domain

import (
	mapset "github.com/deckarep/golang-set/v2"
)

// used for enum
type RiskLevel uint

// used for sort interface
type TransactionsByPosition []Transaction

// alias for type to better legibility
type TransactionsPerUserMap map[uint]mapset.Set[Transaction]

// risk enum
const (
	LOW RiskLevel = iota
	MEDIUM
	HIGH
)

// constant thresholds to decide transaction risk
const (
	MediumRiskSingleAmount = 500000
	HighRiskSingleAmount   = 1000000
	MediumRiskTotalAmount  = 1000000
	HighRiskTotalAmount    = 2000000
)

// magic string function to convert the statuses to string
func (risk RiskLevel) String() string {
	switch risk {
	case LOW:
		return "low"
	case MEDIUM:
		return "medium"
	case HIGH:
		return "high"
	}

	return "unknown"
}

type Transaction struct {
	TransactionId     uint      `json:"id"`
	UserId            uint      `json:"user_id"`
	DollarCentsAmount int       `json:"amount_us_cents"`
	IdCardUsed        uint      `json:"card_id"`
	RiskRate          RiskLevel `json:"transaction_risk"`
	LineNumber        int       `json:"-"`
}

type RiskRateResults struct {
	RiskRates []string `json:"risk_ratings"`
}

// implementing custom sorting function for transaction

func (transactions TransactionsByPosition) Len() int {
	return len(transactions)
}
func (transactions TransactionsByPosition) Swap(i, j int) {
	transactions[i], transactions[j] = transactions[j], transactions[i]
}
func (transactions TransactionsByPosition) Less(i, j int) bool {
	return transactions[i].LineNumber < transactions[j].LineNumber
}

type TransactionsInput struct {
	InputTransactions []Transaction `json:"transactions"`
}
