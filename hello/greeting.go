package hello

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func GreetSomeone(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var notifResult string
	err := workflow.ExecuteActivity(ctx, SendNotification, name, "welcome!").Get(ctx, &notifResult)
	if err != nil {
		return "", err
	}

	var paymentResult string
	err = workflow.ExecuteActivity(ctx, ProcessPayment, 150.0, "USD").Get(ctx, &paymentResult)
	if err != nil {
		return "", err
	}

	return "Hello " + name + "! " + notifResult + " | " + paymentResult, nil
}
