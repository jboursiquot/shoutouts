package shoutouts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NewCallback returns a new Callback.
func NewCallback() *Callback {
	return &Callback{}
}

// Callback calls back Slack to relay a response to the user that requested the shoutout.
type Callback struct{}

// Call performs the HTTP POST operation to Slack using the callback URL attached to the Shoutout.
func (c *Callback) Call(ctx context.Context, shoutout *Shoutout) error {
	r := SlackResponse{
		ResponseType: SlackResponseTypeInChannel,
		Attachments: []*SlackAttachment{
			&SlackAttachment{
				Title: fmt.Sprintf(
					"%s shoutout to @%s from @%s for",
					shoutout.Kind.Name,
					shoutout.RecipientName,
					shoutout.SenderName,
				),
				ThumbURL: shoutout.Kind.ThumbURL,
				Text:     fmt.Sprintf("%s", shoutout.Comment),
			},
		},
	}

	client := &http.Client{Timeout: time.Second * 10}
	url := shoutout.ResponseURL

	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to marshal shoutout: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to POST to response URL: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected server response to be 200 OK, got %v", res.Status)
	}

	return nil
}
