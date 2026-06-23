package hello

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

func SendNotification(ctx context.Context, userID string, message string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SendNotification started", "userID", userID, "time", time.Now())

	time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond)

	if rand.Float64() < 0.4 {
		logger.Error("notification service unavailable (simulated transient failure)")
		return "", errors.New("notification service unavailable (simulated transient failure)")
	}

	txnID := fmt.Sprintf("ntf-%d", rand.Int63n(100000))
	logger.Info("notification sent", "txnID", txnID)
	return txnID, nil
}

func ProcessPayment(ctx context.Context, amount float64, currency string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProcessPayment started", "amount", amount, "currency", currency)

	time.Sleep(2 * time.Second)

	if amount <= 0 {
		logger.Error("invalid amount: must be positive")
		return "", errors.New("invalid amount: must be positive")
	}

	if amount > 5000 {
		logger.Error("amount %.2f %s exceeds daily limit of 5000", amount, currency)
		return "", fmt.Errorf("amount %.2f %s exceeds daily limit of 5000", amount, currency)
	}

	if rand.Float64() < 0.25 {
		logger.Error("payment gateway timeout (simulated transient failure)")
		return "", errors.New("payment gateway timeout (simulated transient failure)")
	}

	receiptID := fmt.Sprintf("rcpt-%s-%d", currency, rand.Int63n(1000000))
	logger.Info("payment processed", "receiptID", receiptID)
	return receiptID, nil
}
