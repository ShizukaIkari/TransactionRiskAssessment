package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	mapset "github.com/deckarep/golang-set/v2"
)

/** Mocks for test cases*/
var transactionsListMock = []Transaction{
	{TransactionId: 1, UserId: 1, DollarCentsAmount: 200000, IdCardUsed: 1},
	{TransactionId: 2, UserId: 1, DollarCentsAmount: 600000, IdCardUsed: 1},
	{TransactionId: 3, UserId: 3, DollarCentsAmount: 1100000, IdCardUsed: 1},
	{TransactionId: 4, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 2},
	{TransactionId: 5, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 3},
	{TransactionId: 6, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 4},
}
var expectedSetsUser1 = mapset.NewSet(transactionsListMock[:2]...)
var expectedSetsUser2 = mapset.NewSet(transactionsListMock[3:]...)
var expectedSetsUser3 = mapset.NewSet(transactionsListMock[2])

var expectedMap = map[uint]mapset.Set[Transaction]{
	1: expectedSetsUser1,
	2: expectedSetsUser2,
	3: expectedSetsUser3,
}

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
			args: args{transactionsListMock},
			want: expectedMap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelateUserToTransactions(tt.args.transacts)
			fmt.Println("got", got)
			fmt.Println("expectedMap", expectedMap)
			assert.Equal(t, tt.want, got)
		})
	}
}

// func relateUserToTransactions(transacts []Transaction) TransactionsPerUserMap {
// func greaterRisk(currentRisk, newRisk RiskLevel) RiskLevel {
// func riskPerSingleAmount(userTransactions []Transaction) {

// func riskPerTotalAmount(userTransactions []Transaction) {

// func riskPerMultipleCards(userTransactions []Transaction) {
// func allTransactionsRisk(allTransactions []Transaction) RiskRate {
// func checkTransactions(userTransactions TransactionsPerUserMap) RiskRate {
