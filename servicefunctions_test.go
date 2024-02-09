package main

import (
	"testing"
	. "transactionriskassessment/domain"

	"github.com/stretchr/testify/assert"

	mapset "github.com/deckarep/golang-set/v2"
)

/** Mocks for test cases*/
var transactionsListMock = []Transaction{
	{TransactionId: 1, UserId: 1, DollarCentsAmount: 200000, IdCardUsed: 1, LineNumber: 1},
	{TransactionId: 2, UserId: 1, DollarCentsAmount: 600000, IdCardUsed: 1, LineNumber: 2},
	{TransactionId: 3, UserId: 3, DollarCentsAmount: 1100000, IdCardUsed: 1, LineNumber: 3},
	{TransactionId: 4, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 2, LineNumber: 4},
	{TransactionId: 5, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 3, LineNumber: 5},
	{TransactionId: 6, UserId: 2, DollarCentsAmount: 100000, IdCardUsed: 4, LineNumber: 6},
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
		transactions []Transaction
	}
	// test case definition
	tests := []struct {
		name string
		args args
		want TransactionsPerUserMap
	}{
		// test cases
		{
			name: "should return a map with three keys for each user and transactions' set for each key",
			args: args{transactionsListMock},
			want: expectedMap,
		},
	}
	// test cases execution
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelateUserToTransactions(tt.args.transactions)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGreaterRisk(t *testing.T) {
	type args struct {
		riskLevel1 RiskLevel
		riskLevel2 RiskLevel
	}
	tests := []struct {
		name string
		args args
		want RiskLevel
	}{
		{
			name: "should return HIGH",
			args: args{LOW, HIGH},
			want: HIGH,
		},
		{
			name: "should return MEDIUM",
			args: args{MEDIUM, LOW},
			want: MEDIUM,
		},
		{
			name: "should return LOW",
			args: args{LOW, LOW},
			want: LOW,
		},
		{
			name: "should return HIGH",
			args: args{HIGH, MEDIUM},
			want: HIGH,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := greaterRisk(test.args.riskLevel1, test.args.riskLevel2)
			assert.Equal(t, test.want, got)
		})
	}

}

func TestRiskPerSingleAmount(t *testing.T) {
	tests := []struct {
		name string
		args []Transaction
		// this function updates the risk field in the transactions passed
		want []Transaction
	}{
		{
			name: "should return low for values <= $5000 transactions",
			args: []Transaction{
				{DollarCentsAmount: 312000},
				{DollarCentsAmount: 499999},
				{DollarCentsAmount: 0},
			},
			want: []Transaction{
				{RiskRate: LOW, DollarCentsAmount: 312000},
				{RiskRate: LOW, DollarCentsAmount: 499999},
				{RiskRate: LOW, DollarCentsAmount: 0},
			},
		},
		{
			name: "should return medium for values between $5000.01 and $10000 transactions",
			args: []Transaction{
				{DollarCentsAmount: 999999},
				{DollarCentsAmount: 500001},
				{DollarCentsAmount: 700000},
			},
			want: []Transaction{
				{RiskRate: MEDIUM, DollarCentsAmount: 999999},
				{RiskRate: MEDIUM, DollarCentsAmount: 500001},
				{RiskRate: MEDIUM, DollarCentsAmount: 700000},
			},
		},
		{
			name: "should return high for values greater than $10000 transactions",
			args: []Transaction{
				{DollarCentsAmount: 5000000},
				{DollarCentsAmount: 1500000},
				{DollarCentsAmount: 1100000},
			},
			want: []Transaction{
				{RiskRate: HIGH, DollarCentsAmount: 5000000},
				{RiskRate: HIGH, DollarCentsAmount: 1500000},
				{RiskRate: HIGH, DollarCentsAmount: 1100000},
			},
		},
		{
			name: "should return low, medium and high",
			args: []Transaction{
				{DollarCentsAmount: 500000},
				{DollarCentsAmount: 1000000},
				{DollarCentsAmount: 1000001},
			},
			want: []Transaction{
				{RiskRate: LOW, DollarCentsAmount: 500000},
				{RiskRate: MEDIUM, DollarCentsAmount: 1000000},
				{RiskRate: HIGH, DollarCentsAmount: 1000001},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			riskPerSingleAmount(test.args)
			// should've been updated
			got := test.args
			assert.Equal(t, test.want, got)

		})
	}

}

func TestRiskPerTotalAmount(t *testing.T) {
	tests := []struct {
		name string
		args []Transaction
		// this function updates the risk field in the transactions passed
		want []Transaction
	}{
		{
			name: "The accumulating values should return LOW, MEDIUM and HIGH",
			args: []Transaction{
				{DollarCentsAmount: 1000000},
				{DollarCentsAmount: 500000},
				{DollarCentsAmount: 500010},
			},
			want: []Transaction{
				{RiskRate: LOW, DollarCentsAmount: 1000000},
				{RiskRate: MEDIUM, DollarCentsAmount: 500000},
				{RiskRate: HIGH, DollarCentsAmount: 500010},
			},
		},
		{
			name: "Accumulating values varying, should return HIGH, MEDIUM and LOW",
			args: []Transaction{
				// should not update a higher risk to a medium
				{DollarCentsAmount: 1100000, RiskRate: HIGH},
				// medium since 16k is between 10k and 20k
				{DollarCentsAmount: 500000},
				// this shouldn't be possible, but currently is. so risklevel should be low (for the total amount rules)
				{DollarCentsAmount: -1000000},
			},
			want: []Transaction{
				{RiskRate: HIGH, DollarCentsAmount: 1100000},
				{RiskRate: MEDIUM, DollarCentsAmount: 500000},
				{RiskRate: LOW, DollarCentsAmount: -1000000},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			riskPerTotalAmount(test.args)
			// should've been updated
			got := test.args
			assert.Equal(t, test.want, got)
		})
	}
}

func TestRiskPerMultipleCards(t *testing.T) {
	tests := []struct {
		name string
		args []Transaction
		// this function updates the risk field in the transactions passed
		want []Transaction
	}{
		{
			name: "Same cardId, should return low 3x",
			args: []Transaction{
				{IdCardUsed: 1},
				{IdCardUsed: 1},
				{IdCardUsed: 1},
			},
			want: []Transaction{
				{RiskRate: LOW, IdCardUsed: 1},
				{RiskRate: LOW, IdCardUsed: 1},
				{RiskRate: LOW, IdCardUsed: 1},
			},
		},
		{
			name: "Two distinct cardIds, should return low, medium and medium",
			args: []Transaction{
				{IdCardUsed: 1},
				{IdCardUsed: 2},
				{IdCardUsed: 1},
			},
			want: []Transaction{
				{RiskRate: LOW, IdCardUsed: 1},
				{RiskRate: MEDIUM, IdCardUsed: 2},
				{RiskRate: MEDIUM, IdCardUsed: 1},
			},
		},
		{
			name: "Three distinct cardIds, should return low, medium, high",
			args: []Transaction{
				{IdCardUsed: 1},
				{IdCardUsed: 2},
				{IdCardUsed: 3},
			},
			want: []Transaction{
				{RiskRate: LOW, IdCardUsed: 1},
				{RiskRate: MEDIUM, IdCardUsed: 2},
				{RiskRate: HIGH, IdCardUsed: 3},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			riskPerMultipleCards(test.args)
			// should've been updated
			got := test.args
			assert.Equal(t, test.want, got)
		})
	}
}

func TestAllTransactionsRisk(t *testing.T) {
	tests := []struct {
		name string
		args []Transaction
		want RiskRateResults
	}{
		{
			name: `Should return after sorting LineNumber "medium", "high", "low", "low" `,
			args: []Transaction{
				{LineNumber: 3, RiskRate: LOW},
				{LineNumber: 1, RiskRate: MEDIUM},
				{LineNumber: 5, RiskRate: LOW},
				{LineNumber: 2, RiskRate: HIGH},
			},
			want: RiskRateResults{RiskRates: []string{"medium", "high", "low", "low"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := allTransactionsRisk(test.args)
			// should've been updated
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCheckTransactions(t *testing.T) {
	tests := []struct {
		name string
		args TransactionsPerUserMap
		want RiskRateResults
	}{
		{
			name: "evaluating 8 transactions from 3 users, should return string slice with risk ratings ordered by line number ",
			args: map[uint]mapset.Set[Transaction]{
				10: mapset.NewSet([]Transaction{
					{TransactionId: 12, UserId: 10, DollarCentsAmount: 2200000, IdCardUsed: 1, LineNumber: 5}, //high per rule 4
					{TransactionId: 9, UserId: 10, DollarCentsAmount: 10000, IdCardUsed: 1, LineNumber: 2},    // evaluated first, so low (no rule match)
				}...),

				5: mapset.NewSet([]Transaction{
					{TransactionId: 13, UserId: 5, DollarCentsAmount: 500600, IdCardUsed: 321, LineNumber: 1}, // medium per rule 1
					{TransactionId: 2, UserId: 5, DollarCentsAmount: 10000, IdCardUsed: 121, LineNumber: 3},   // medium per rule 5
					{TransactionId: 5, UserId: 5, DollarCentsAmount: 10000, IdCardUsed: 132, LineNumber: 7},   // high per rule 6
				}...),

				31: mapset.NewSet([]Transaction{
					{TransactionId: 4, UserId: 31, DollarCentsAmount: 5000000, IdCardUsed: 22, LineNumber: 8}, // high per rules 2, 4 (1 and 5 match, but are lower priority).
					{TransactionId: 3, UserId: 31, DollarCentsAmount: 1200000, IdCardUsed: 21, LineNumber: 4}, // high per rule 2
					{TransactionId: 11, UserId: 31, DollarCentsAmount: 1000, IdCardUsed: 20, LineNumber: 6},   // medium per rules 3, 5
				}...),
			},
			want: RiskRateResults{RiskRates: []string{"medium", "low", "medium", "high", "high", "medium", "high", "high"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CheckTransactions(test.args)
			// should've been updated
			assert.Equal(t, test.want, got)
		})
	}
}
