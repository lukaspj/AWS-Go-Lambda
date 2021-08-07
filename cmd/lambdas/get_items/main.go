package main

import (
	"BryrupTeater.Backend/pkg/lambdahandler"
	"BryrupTeater.Backend/pkg/testtype"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type ResponseItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Items   []ResponseItem `json:"items"`
}

func GetItems(ctx context.Context, request *lambdahandler.LambdaRequest) *lambdahandler.LambdaResponse {
	cfg := lambdahandler.GetAwsConfig(ctx)
	repo := testtype.NewDynamoRepo(cfg, "test_table")
	items, err := repo.List(ctx)
	if err != nil {
		return lambdahandler.ErrorResponse(http.StatusInternalServerError, err)
	}

	resp := Response{
		Items:   make([]ResponseItem, len(items)),
	}

	for i, item := range items {
		resp.Items[i] = ResponseItem{
			Id:   item.UserId,
			Name: item.Name,
		}
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return lambdahandler.ErrorResponse(http.StatusInternalServerError, err)
	}

	return &lambdahandler.LambdaResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	lambda.Start(
		lambdahandler.New().
			Use(lambdahandler.LogEntryExitContext).
			Use(lambdahandler.AwsConfig).
			Handle(GetItems),
	)
}
