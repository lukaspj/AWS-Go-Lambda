package lambdahandler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type LambdaRequest struct {
	Body           []byte `json:"body"`
	PathParameters map[string]string
}

type LambdaResponse struct {
	Error             error               `json:"error"`
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              []byte              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

type Func func(context.Context, *LambdaRequest) *LambdaResponse
type Middleware func(Func) Func

type LambdaHandler struct {
	Middleware []Middleware
}

func New() *LambdaHandler {
	return &LambdaHandler{Middleware: make([]Middleware, 0)}
}

func (h *LambdaHandler) Use(middleware ...Middleware) *LambdaHandler {
	h.Middleware = append(h.Middleware, middleware...)
	return h
}

func (h *LambdaHandler) Handle(fn Func) func(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		fn := fn
		for i := range h.Middleware {
			// Reverse loop
			middleware := h.Middleware[len(h.Middleware)-1-i]
			fn = middleware(fn)
		}
		return toAwsResponse(fn(ctx, fromAwsRequest(ctx, request)))
	}
}

func ErrorResponse(statusCode int, err error) *LambdaResponse {
	return &LambdaResponse{
		Error:           err,
		StatusCode:      statusCode,
		Body:            []byte(err.Error()),
		IsBase64Encoded: false,
	}
}

func fromAwsRequest(ctx context.Context, request *events.APIGatewayProxyRequest) *LambdaRequest {
	return &LambdaRequest{
		Body:           []byte(request.Body),
		PathParameters: request.PathParameters,
	}
}

func toAwsResponse(response *LambdaResponse) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode:        response.StatusCode,
		Headers:           response.Headers,
		MultiValueHeaders: response.MultiValueHeaders,
		Body:              string(response.Body),
		IsBase64Encoded:   response.IsBase64Encoded,
	}, response.Error
}

func internalServerError(err error) *LambdaResponse {
	return ErrorResponse(http.StatusInternalServerError, err)
}
