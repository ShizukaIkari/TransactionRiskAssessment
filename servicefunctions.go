package main

import (
	"sort"

	. "transactionriskassessment/domain"

	mapset "github.com/deckarep/golang-set/v2"
)

// RelateUserToTransactions maps each unique user id to their corresponding set of transactions
func RelateUserToTransactions(transactionSlice []Transaction) TransactionsPerUserMap {
	// unique user ids in transaction list
	userIdSet := mapset.NewSet[uint]()
	// set with unique transactions, being instantiated for each new user
	var transactionSet mapset.Set[Transaction]

	// relating each user to their transactions
	userAndTransct := make(TransactionsPerUserMap)

	for transactionIndex := range transactionSlice {
		// saving position that it was in the json
		currentTransaction := &transactionSlice[transactionIndex]
		transactionSlice[transactionIndex].LineNumber = transactionIndex + 1

		// first time seeing this user id
		if !userIdSet.Contains(currentTransaction.UserId) {
			// add to the known ids set
			userIdSet.Add(currentTransaction.UserId)
			// new transaction set for new user
			transactionSet = mapset.NewSet[Transaction]()
			// the first transaction for this user in their set
			transactionSet.Add(*currentTransaction)

			// assign set for UserId key
			userAndTransct[currentTransaction.UserId] = transactionSet
		} else {
			userSet := userAndTransct[currentTransaction.UserId]
			userSet.Add(*currentTransaction)
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

	for transacIndex := range userTransactions {
		// shortening the name reference of memory space
		currentTransaction := &userTransactions[transacIndex]
		// standard risk if no match for risk rules
		risk := LOW
		if currentTransaction.DollarCentsAmount > HighRiskSingleAmount {
			risk = HIGH
		} else if currentTransaction.DollarCentsAmount > MediumRiskSingleAmount {
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

	// to update values from slice input, do not use the second return of range
	// it's a copy of the element of the slice, not a reference.
	for transacIndex := range userTransactions {
		currentTransaction := &userTransactions[transacIndex]
		totalAmount += currentTransaction.DollarCentsAmount

		risk := LOW
		if totalAmount > MediumRiskTotalAmount && totalAmount <= HighRiskTotalAmount {
			risk = MEDIUM
		} else if totalAmount > HighRiskTotalAmount {
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
func allTransactionsRisk(allTransactions []Transaction) RiskRateResults {
	var calcRiskRates RiskRateResults
	var risksSlice []string

	// sorting transactions by input position to make result position match
	sort.Sort(TransactionsByPosition(allTransactions))

	for _, transac := range allTransactions {
		risksSlice = append(risksSlice, transac.RiskRate.String())
	}
	calcRiskRates.RiskRates = risksSlice

	return calcRiskRates
}

// Receives the input transactions, returns a slice with each transaction risk level ordered by transaction line number
func CheckTransactions(userTransactions TransactionsPerUserMap) RiskRateResults {
	var transactionsSlice []Transaction

	// key, value - no use for key, so _ instead. Transactions is the copy of the set
	for _, transactions := range userTransactions {
		transactSlice := transactions.ToSlice()
		// sorting before calculating makes totalAmountRisk deliver consistent results
		sort.Sort(TransactionsByPosition(transactSlice))
		riskPerSingleAmount(transactSlice)
		riskPerTotalAmount(transactSlice)
		riskPerMultipleCards(transactSlice)

		transactionsSlice = append(transactionsSlice, transactSlice...)
	}

	return allTransactionsRisk(transactionsSlice)
}
