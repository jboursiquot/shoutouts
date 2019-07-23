package shoutouts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// NewHandler initializes and returns a new Handler.
func NewHandler(s SQSAPI, ddb DynamoDBQuerier) *Handler {
	return &Handler{sqs: s, ddb: ddb}
}

// Handler handles incoming shoutout requests.
type Handler struct {
	sqs SQSAPI
	ddb DynamoDBQuerier
}

// Handle handles the shoutout request.
func (h *Handler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authorizedRequest(request) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	params, err := url.ParseQuery(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	cmd, err := h.parseCommand(ctx, &params)
	if err != nil {
		params := url.Values{}
		params.Set("text", "help")
		cmd, _ = h.parseCommand(ctx, &params)
	}

	rsp, err := cmd.execute(ctx)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	body, err := json.Marshal(rsp)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	headers := map[string]string{"Content-Type": "application/json"}
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

type command interface {
	execute(context.Context) (*SlackResponse, error)
}

func (h *Handler) parseCommand(ctx context.Context, params *url.Values) (command, error) {
	hc, err := parseHelpCommand(ctx, params)
	if err == nil {
		return hc, nil
	}

	sc, err := parseShoutoutCommand(ctx, params)
	if err == nil {
		sc.sqs = h.sqs
		return sc, err
	}

	lc, err := parseListCommand(ctx, params)
	if err == nil {
		lc.ddb = h.ddb
		return lc, err
	}

	return nil, fmt.Errorf("failed to parse any commands")
}

func authorizedRequest(request *events.APIGatewayProxyRequest) bool {
	var token string

	switch request.HTTPMethod {
	case http.MethodGet:
		if t, ok := request.QueryStringParameters["token"]; ok {
			token = t
		}
	case http.MethodPost:
		if params, err := url.ParseQuery(request.Body); err == nil {
			token = params.Get("token")
		}
	}

	if token == os.Getenv("SLACK_TOKEN") {
		return true
	}

	return false
}
