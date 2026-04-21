// Package temporal contains the temporal client
package temporal

import (
	"context"
	"log"

	"GrandExchange/internal/config"
	"GrandExchange/internal/models"
	"GrandExchange/internal/service"
	"GrandExchange/temporal/activities"
	"GrandExchange/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type TemporalClient interface {
	StartSearchWorkflow(context.Context, string) ([]models.SearchItem, error)
}

type ItemTemporalClient struct {
	Client      client.Client
	itemService service.ItemService
}

func NewTemporalClient(service service.ItemService) (TemporalClient, worker.Worker, error) {
	t, err := client.Dial(client.Options{
		HostPort: config.TemporalHostPort,
	})
	if err != nil {
		return nil, nil, err
	}

	w := worker.New(t, OsrsItemsQueue, worker.Options{})

	activities := activities.NewOsrsActivities(service)

	w.RegisterWorkflow(workflows.SearchWorkflow)
	w.RegisterActivity(activities)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatal(err)
		}
	}()

	return &ItemTemporalClient{Client: t, itemService: service}, w, nil
}

func (t *ItemTemporalClient) StartSearchWorkflow(ctx context.Context, itemName string) ([]models.SearchItem, error) {
	we, err := t.Client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        "search-" + itemName,
			TaskQueue: "ItemQueue",
		},
		workflows.SearchWorkflow,
		itemName,
	)
	if err != nil {
		return nil, err
	}

	var response models.SearchActivityResponse
	err = we.Get(ctx, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}
