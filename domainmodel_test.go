package main

import (
	// "github.com/stretchr/testify/assert"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

var transactionsListMock = []Transaction{
	{Transaction_id: 1, User_id: 1, Dollar_cents_amount: 200000, Id_card_used: 1},
	{Transaction_id: 2, User_id: 1, Dollar_cents_amount: 600000, Id_card_used: 1},
	{Transaction_id: 3, User_id: 3, Dollar_cents_amount: 1100000, Id_card_used: 1},
	{Transaction_id: 4, User_id: 2, Dollar_cents_amount: 100000, Id_card_used: 2},
	{Transaction_id: 5, User_id: 2, Dollar_cents_amount: 100000, Id_card_used: 3},
	{Transaction_id: 6, User_id: 2, Dollar_cents_amount: 100000, Id_card_used: 4},
}

var expectedMap = map[uint]mapset.Set[Transaction]{}

func TestRelateUserToTransactions(t *testing.T) {
	type args struct {
		transacts []Transaction
	}
	tests := []struct {
		name string
		args args
		want TransactionsPerUserMap
	}{
		{
			name: "should return a map with three keys for each user and transactions' set for each key",
			args: transactionsMock,
			want: expectedMap,
		},
	}
}

// func relateUserToTransactions(transacts []Transaction) TransactionsPerUserMap {
// func greaterRisk(currentRisk, newRisk RiskLevel) RiskLevel {
// func riskPerSingleAmount(userTransactions []Transaction) {

// func riskPerTotalAmount(userTransactions []Transaction) {

// func riskPerMultipleCards(userTransactions []Transaction) {
// func allTransactionsRisk(allTransactions []Transaction) RiskRate {
// func checkTransactions(userTransactions TransactionsPerUserMap) RiskRate {
