package shoutouts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

// NewSanitizingHandler initializes and returns a new SanitizingHandler.
func NewSanitizingHandler(s SQSAPI, ddb DynamoDBQuerier, logger *logrus.Logger, sanitizerServiceEndpoint string) *SanitizingHandler {
	h := SanitizingHandler{
		endpoint:   sanitizerServiceEndpoint,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	h.sqs = s
	h.ddb = ddb
	h.log = logger
	return &h
}

// SanitizingHandler handles incoming shoutout requests.
type SanitizingHandler struct {
	Handler
	endpoint   string
	httpClient *http.Client
}

// Handle handles the shoutout request after sanitizinng it.
func (h *SanitizingHandler) Handle(ctx context.Context, request *events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	bs, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		h.log.WithError(err).Error("failed to decode request body")
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	params, err := url.ParseQuery(string(bs))
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	if !h.authorizedRequest(params.Get("token")) {
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	cmd, err := h.parseCommand(ctx, &params)
	if err != nil {
		params := url.Values{}
		params.Set("text", "help")
		cmd, _ = h.parseCommand(ctx, &params)
	}

	if sc, ok := cmd.(*shoutoutCommand); ok {
		h.log.WithField("shoutout", sc).Infoln("shoutout command received")
		if err := sc.sanitize(ctx, h.endpoint, h.httpClient); err != nil {
			h.log.WithError(err).Errorln("failed to sanitize shoutout")
			return events.LambdaFunctionURLResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}
	}

	rsp, err := cmd.execute(ctx)
	if err != nil {
		h.log.Println(err)
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	body, err := json.Marshal(rsp)
	if err != nil {
		h.log.Println(err)
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	headers := map[string]string{"Content-Type": "application/json"}
	return events.LambdaFunctionURLResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func (h SanitizingHandler) authorizedRequest(token string) bool {
	return token == os.Getenv("SLACK_TOKEN")
}
