package main

import (
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
)

// section emulating enum in go
type RiskLevel uint
type TransactionsByID []Transaction
type TransactionsPerUserMap map[uint]mapset.Set[Transaction]

const (
	LOW RiskLevel = iota
	MEDIUM
	HIGH
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
}

type RiskRate struct {
	RiskRates []string `json:"risk_ratings"`
}

// implementing custom sorting function for transaction

func (transactions TransactionsByID) Len() int {
	return len(transactions)
}
func (transactions TransactionsByID) Swap(i, j int) {
	transactions[i], transactions[j] = transactions[j], transactions[i]
}
func (transactions TransactionsByID) Less(i, j int) bool {
	return transactions[i].TransactionId < transactions[j].TransactionId
}

type TransactionsInput struct {
	InputTransactions []Transaction `json:"transactions"`
}

// RelateUserToTransactions maps each unique user id to their corresponding set of transactions
func RelateUserToTransactions(transacts []Transaction) TransactionsPerUserMap {
	// unique user ids in transaction list
	userIdSet := mapset.NewSet[uint]()
	// set with unique transactions, being instantiated for each new user
	var transactionSet mapset.Set[Transaction]

	// relating each user to their transactions
	userAndTransct := make(TransactionsPerUserMap)

	for _, transact := range transacts {
		// first time seeing this user id
		if !userIdSet.Contains(transact.UserId) {
			// add to the known ids set
			userIdSet.Add(transact.UserId)
			// new transaction set for new user
			transactionSet = mapset.NewSet[Transaction]()
			// the first transaction for this user in their set
			transactionSet.Add(transact)

			// assign set for UserId key
			userAndTransct[transact.UserId] = transactionSet
		} else {
			userSet := userAndTransct[transact.UserId]
			userSet.Add(transact)
		}
	}
	return userAndTransct
}

// Checks which risk is greater, returning the greater one.
func greaterRisk(currentRisk, newRisk RiskLevel) RiskLevel {
	if currentRisk < newRisk {
		return newRisk
	}
	return currentRisk
}

// Analyzes the amount of each transaction, updating their Risk Level according to the risk rules
func riskPerSingleAmount(userTransactions []Transaction) {
	// in US cents
	mediumRiskAmount := 500000
	highRiskAmount := 1000000

	for transacIndex := range userTransactions {
		// shortening the name reference of memory space
		currentTransaction := &userTransactions[transacIndex]
		// standard risk if no match for risk rules
		risk := LOW
		if currentTransaction.DollarCentsAmount > highRiskAmount {
			risk = HIGH
		} else if currentTransaction.DollarCentsAmount > mediumRiskAmount {
			risk = MEDIUM
		}
		currentTransaction.RiskRate = greaterRisk(userTransactions[transacIndex].RiskRate, risk)
	}
}

// Sums up the dollar amount while iterating through user transactions
// updating their Risk Level according to the risk rules
func riskPerTotalAmount(userTransactions []Transaction) {
	// in US cents
	var totalAmount int
	mediumRiskAmount := 1000000
	highRiskAmount := 2000000

	// to update values from slice input, do not use the second return of range
	// it's a copy of the element of the slice, not a reference.
	for transacIndex := range userTransactions {
		currentTransaction := &userTransactions[transacIndex]
		totalAmount += currentTransaction.DollarCentsAmount

		risk := LOW
		if totalAmount > mediumRiskAmount && totalAmount <= highRiskAmount {
			risk = MEDIUM
		} else if totalAmount > highRiskAmount {
			risk = HIGH
		}

		currentTransaction.RiskRate = greaterRisk(currentTransaction.RiskRate, risk)
	}
}

// Checks how many different cards are in the user transaction set
// updating each transaction Risk Level according to the risk rules
func riskPerMultipleCards(userTransactions []Transaction) {
	cardIdSet := mapset.NewSet[uint]()

	for transacIndex := range userTransactions {
		currentTransaction := &userTransactions[transacIndex]
		isCardInSet := cardIdSet.Contains(currentTransaction.IdCardUsed)

		if !isCardInSet {
			cardIdSet.Add(currentTransaction.IdCardUsed)
		}

		risk := LOW
		if cardIdSet.Cardinality() <= 1 {
			risk = LOW
		} else if cardIdSet.Cardinality() == 2 {
			risk = MEDIUM
		} else {
			risk = HIGH
		}

		currentTransaction.RiskRate = greaterRisk(currentTransaction.RiskRate, risk)
	}
}

// Receives a slice with all transactions that were sent to the API
// with risks calculated and returns a slice with resultant risk in strings
// sorted by transaction id * might change
func allTransactionsRisk(allTransactions []Transaction) RiskRate {
	var calcRiskRates RiskRate
	// sorting transactions by ID to maintain the order of input *
	// * this fails if the input isn't in transaction Id order. Future feature to use the json parsing to keep the input order in a field
	sort.Sort(TransactionsByID(allTransactions))
	var risksSlice []string

	for _, transac := range allTransactions {
		risksSlice = append(risksSlice, transac.RiskRate.String())
	}
	calcRiskRates.RiskRates = risksSlice

	return calcRiskRates
}

// Receives the input transactions, returns a slice with each transaction risk level ordered by transactionId
func CheckTransactions(userTransactions TransactionsPerUserMap) RiskRate {
	var transactionsSlice []Transaction

	// key, value - no use for key, so _ instead. Transactions is the copy of the set
	for _, transactions := range userTransactions {
		transactSlice := transactions.ToSlice()
		riskPerSingleAmount(transactSlice)
		riskPerTotalAmount(transactSlice)
		riskPerMultipleCards(transactSlice)

		transactionsSlice = append(transactionsSlice, transactSlice...)
	}

	return allTransactionsRisk(transactionsSlice)
}
