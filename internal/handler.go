package internal

import (
	"fmt"
	"net/http"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/mercari/go-circuitbreaker"
	"context"
	"time"
	"Transaction-manager/models"
	"log/slog"
)

func ConfigRouter(producer *kafka.Producer, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "MAIN PAGE")
	})

	accountsGroup := r.Group("/accounts")
	{
		accountsGroup.GET("/:userID", func(c *gin.Context) {
			GetAccount(c, logger)
		})
		accountsGroup.POST("/transfer/:userID", func(c *gin.Context) {
			Transfer(c, producer, logger)
		})
	}

	return r
}

func GetAccount(c *gin.Context, logger *slog.Logger) {
	userID := c.Param("userID")
	c.String(http.StatusOK, "User ID: %s", userID)
}

func Transfer(c *gin.Context, producer *kafka.Producer, logger *slog.Logger) {
	userID := c.Param("userID")

	var request models.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipientId := request.RecipientID
	amount := request.Amount
	valute := request.Valute

	message := strings.Join([]string{userID, recipientId, amount, valute}, ",")
	
	const topic = "TransactionManagerTopic"

	cb := circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(5*time.Second),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncFailureRate(3, 0.5)),
	)
	
	ctx := context.Background()
	result, err := cb.Do(ctx, func() (interface{}, error) {
		return sendMessage(producer, topic, message)
	})

	if err != nil {
		fmt.Printf("Failed to send message: %s\n", err)
		c.String(http.StatusInternalServerError, "Failed to send message")
		return
	}

	fmt.Println("message was successfully sent:", result)
	c.String(http.StatusOK, "message was successfully sent")

}