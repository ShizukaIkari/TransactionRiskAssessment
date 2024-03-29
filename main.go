package main

import (
	"net/http"
	"transactionriskassessment/domain"

	"github.com/gin-gonic/gin"
)

/**
* assessTransactions receives a JSON from request body, calls RelateUserToTransactions to map
* user ids to its respective transactions and finally calls CheckTransactions to assess the risk
* for each transaction according to requirement's rules, returning the JSON with their risks
 */
func AssessTransactions(context *gin.Context) {
	var newTransactionsList domain.TransactionsInput

	// if API input is not what's expected, end the function
	if err := context.BindJSON(&newTransactionsList); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call API logic
	mappedUserTransac := RelateUserToTransactions(newTransactionsList.InputTransactions)
	resultantRatings := CheckTransactions(mappedUserTransac)

	// call function that process and returns transaction risks
	context.IndentedJSON(http.StatusOK, resultantRatings)
}

func main() {
	// server to run the API
	router := gin.Default()
	router.POST("/check_transactions", AssessTransactions)
	router.Run("localhost:9090")
}
