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

type TestResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetItem(ctx context.Context, request *lambdahandler.LambdaRequest) *lambdahandler.LambdaResponse {
	id, ok := request.PathParameters["id"]
	if !ok {
		return lambdahandler.ErrorResponse(http.StatusBadRequest, fmt.Errorf("missing path parameter: id"))
	}

	cfg := lambdahandler.GetAwsConfig(ctx)
	repo := testtype.NewDynamoRepo(cfg, "test_table")
	item, err := repo.Get(ctx, id)
	if err != nil {
		return lambdahandler.ErrorResponse(http.StatusInternalServerError, err)
	}

	body, err := json.Marshal(TestResponse{Id: item.UserId, Name: item.Name})
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
			Handle(GetItem),
	)
}
