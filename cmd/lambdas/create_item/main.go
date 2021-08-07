package main

import (
	"BryrupTeater.Backend/pkg/lambdahandler"
	"BryrupTeater.Backend/pkg/testtype"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type TestEvent struct {
	Name string `json:"name"`
}

type TestResponse struct {
	Id string `json:"id"`
}

func CreateItem(ctx context.Context, request *lambdahandler.LambdaRequest) *lambdahandler.LambdaResponse {
	var event TestEvent
	var err error
	err = json.Unmarshal(request.Body, &event)

	if err != nil {
		err := fmt.Errorf("failed to unmarshal event due to %w", err)
		log.Error(err.Error())
		return lambdahandler.ErrorResponse(http.StatusBadRequest, err)
	}

	cfg := lambdahandler.GetAwsConfig(ctx)
	repo := testtype.NewDynamoRepo(cfg, "test_table")
	id, err := repo.Create(ctx, &testtype.TestItem{
		Name: event.Name,
	})
	if err != nil {
		return lambdahandler.ErrorResponse(http.StatusBadRequest, err)
	}

	body, err := json.Marshal(TestResponse{Id: id})
	if err != nil {
		return lambdahandler.ErrorResponse(http.StatusBadRequest, err)
	}

	return &lambdahandler.LambdaResponse{
		StatusCode: http.StatusOK,
		Body: body,
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	lambda.Start(
		lambdahandler.New().
			Use(lambdahandler.LogEntryExitContext).
			Use(lambdahandler.AwsConfig).
			Handle(CreateItem),
	)
}
