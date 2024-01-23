package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Transaction struct {
// 	Transaction_id      uint `json:"id"`
// 	User_id             uint `json:"user_id"`
// 	Dollar_cents_amount int  `json:"amount_us_cents"`
// 	Id_card_used        uint `json:"card_id"`
// 	Risk_rate           int  `json:"transaction_risk"`
// }

// type TransactionsInput struct {
// 	InputTransactions []Transaction `json:"transactions"`
// }

func assessTransactions(context *gin.Context) {
	var newTransactionsList TransactionsInput

	if err := context.BindJSON(&newTransactionsList); err != nil {
		return
	}

	// fmt.Println(newTransactionsList)
	// fmt.Println("Input Transactions array:", newTransactionsList.InputTransactions)
	var mappedUserTransac = relateUserToTransactions(newTransactionsList.InputTransactions)
	fmt.Println(mappedUserTransac)
	resultantRatings := checkTransactions(mappedUserTransac)
	// call function that process and returns transaction status
	context.IndentedJSON(http.StatusOK, resultantRatings)
}

func main() {
	// server to run the API
	router := gin.Default()
	router.POST("/check_transactions", assessTransactions)
	router.Run("localhost:9090")
}
