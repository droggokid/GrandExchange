// Package workflows contains all temporal item related workflows
package workflows

import (
	"time"

	"GrandExchange/internal/models"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func SearchWorkflow(ctx workflow.Context, itemName string) (models.SearchActivityResponse, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			MaximumInterval:    time.Minute,
			BackoffCoefficient: 2,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var response models.SearchActivityResponse
	err := workflow.ExecuteActivity(ctx, "SearchItemActivity", itemName).Get(ctx, &response)
	if err != nil {
		return models.SearchActivityResponse{}, err
	}

	return response, err
}
