package main

import (
	// "fmt"
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
	Transaction_id      uint      `json:"id"`
	User_id             uint      `json:"user_id"`
	Dollar_cents_amount int       `json:"amount_us_cents"`
	Id_card_used        uint      `json:"card_id"`
	Risk_rate           RiskLevel `json:"transaction_risk"`
}

type RiskRate struct {
	Risk_rates []string `json:"risk_ratings"`
}

// implementing custom sorting function for transaction

func (transactions TransactionsByID) Len() int {
	return len(transactions)
}
func (transactions TransactionsByID) Swap(i, j int) {
	transactions[i], transactions[j] = transactions[j], transactions[i]
}
func (transactions TransactionsByID) Less(i, j int) bool {
	return transactions[i].Transaction_id < transactions[j].Transaction_id
}

type TransactionsInput struct {
	InputTransactions []Transaction `json:"transactions"`
}

/*
*  This function maps each unique user id to their corresponding set of transactions
 */
func relateUserToTransactions(transacts []Transaction) TransactionsPerUserMap {
	// unique user ids in transaction list
	userIdSet := mapset.NewSet[uint]()
	// set with unique transactions, being instantiated for each new user
	var transactionSet mapset.Set[Transaction]

	// relating each user to their transactions
	userAndTransct := make(TransactionsPerUserMap)

	for _, transact := range transacts {
		// first time seeing this user id
		if !userIdSet.Contains(transact.User_id) {
			// add to the known ids set
			userIdSet.Add(transact.User_id)
			// new transaction set for new user
			transactionSet = mapset.NewSet[Transaction]()
			// the first transaction for this user in their set
			transactionSet.Add(transact)

			// assign set for UserId key
			userAndTransct[transact.User_id] = transactionSet
		} else {
			userSet := userAndTransct[transact.User_id]
			userSet.Add(transact)
		}
	}
	return userAndTransct
}

// Checks which risk is greater, returning the ... greater ... one.
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

	for _, transac := range userTransactions {
		if transac.Dollar_cents_amount > highRiskAmount {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, HIGH)
		} else if transac.Dollar_cents_amount > mediumRiskAmount {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, MEDIUM)
		} else {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, LOW)
		}
	}
}

// should I give this function the responsibility to filter the map/transactions for a given user id?
// for now I'm determining that the transaction list has only  one disctinct user_id
func riskPerTotalAmount(userTransactions []Transaction) {
	// in US cents
	var totalAmount int
	mediumRiskAmount := 1000000
	highRiskAmount := 2000000

	for _, transac := range userTransactions {
		totalAmount += transac.Dollar_cents_amount

		if totalAmount > mediumRiskAmount && totalAmount <= highRiskAmount {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, MEDIUM)
		} else if totalAmount > highRiskAmount {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, HIGH)
		} else {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, LOW)
		}
	}

}

// So ... I think I'll merge this function with the per total amout to avoid using another for loop
// BUT if I do this, I think I'll be breaking the single responsibility principle
func riskPerMultipleCards(userTransactions []Transaction) {
	cardIdSet := mapset.NewSet[uint]()

	for _, transac := range userTransactions {
		if !cardIdSet.Contains(transac.Id_card_used) && cardIdSet.Cardinality() < 1 {
			// add to the known ids set
			cardIdSet.Add(transac.Id_card_used)
			// first transaction to be analyzed has only one card used. so risk level is low
			transac.Risk_rate = greaterRisk(transac.Risk_rate, LOW)
		} else if !cardIdSet.Contains(transac.Id_card_used) && cardIdSet.Cardinality() == 1 {
			cardIdSet.Add(transac.Id_card_used)
			// now user has used more than 1 card
			transac.Risk_rate = greaterRisk(transac.Risk_rate, MEDIUM)
		} else if !cardIdSet.Contains(transac.Id_card_used) && cardIdSet.Cardinality() >= 2 {
			transac.Risk_rate = greaterRisk(transac.Risk_rate, HIGH)
		}
	}

}

func allTransactionsRisk(allTransactions []Transaction) RiskRate {
	var calcRiskRates RiskRate
	// sorting transactions by ID to maintain the order of input *
	// * this fails if the input isn't in transaction Id order. Future feature to use the json parsing to keep the input order in a field
	sort.Sort(TransactionsByID(allTransactions))
	var risksSlice []string

	for _, transac := range allTransactions {
		risksSlice = append(risksSlice, transac.Risk_rate.String())
	}
	calcRiskRates.Risk_rates = risksSlice

	return calcRiskRates
}

// Receives the input transactions, returns a slice with each transaction risk level ordered by transactionId
func checkTransactions(userTransactions TransactionsPerUserMap) RiskRate {
	var transactionsSlice []Transaction

	// key, value - no use for key, so _ instead. Transactions is a set
	for _, transactions := range userTransactions {
		transactSlice := transactions.ToSlice()
		riskPerSingleAmount(transactSlice)
		riskPerTotalAmount(transactSlice)
		riskPerMultipleCards(transactSlice)

		transactionsSlice = append(transactionsSlice, transactSlice...)
	}

	return allTransactionsRisk(transactionsSlice)
}
