package lambdahandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)

func LogEntryExitContext(fn Func) Func {
	return func(ctx context.Context, request *LambdaRequest) *LambdaResponse {
		lc, ok := lambdacontext.FromContext(ctx)
		if !ok {
			err := errors.New("failed to get lambdacontext")
			return internalServerError(err)
		}

		log.WithFields(
			log.Fields{
				"function-name":    lambdacontext.FunctionName,
				"function-version": lambdacontext.FunctionVersion,
				"memory-limit":     fmt.Sprintf("%dMB", lambdacontext.MemoryLimitInMB),
				"log-group":        lambdacontext.LogGroupName,
				"log-stream":       lambdacontext.LogStreamName,
				"request-id":       lc.AwsRequestID,
				"function-arn":     lc.InvokedFunctionArn,
			},
		).Info("Entered Lambda")

		resp := fn(ctx, request)

		log.Info("Exited Lambda")

		return resp
	}
}

var AwsConfigCtxKey = "AWS_CONFIG_CTX_KEY"

func GetAwsConfig(ctx context.Context) aws.Config {
	// Assume it's set, and if it's not it should be handled in middleware
	return ctx.Value(AwsConfigCtxKey).(aws.Config)
}

func AwsConfig(fn Func) Func {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "eu-north-1"
		return nil
	})

	return func(ctx context.Context, request *LambdaRequest) *LambdaResponse {
		if err != nil {
			return internalServerError(err)
		}

		return fn(context.WithValue(ctx, AwsConfigCtxKey, cfg), request)
	}
}