package shoutouts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSAPI is the minimal interface needed to enqueue a shoutout.
type SQSAPI interface {
	SendMessageWithContext(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error)
}

type shoutoutCommand struct {
	shoutout *Shoutout
	sqs      SQSAPI
}

// parseShoutoutCommand parses the shoutout command given the key=value pairs in the params map
// 	token=PjsIpgv7fGCKHB6bK4EAx43U
// 	team_id=T0001
// 	team_domain=example
// 	channel_id=C2147483705
// 	channel_name=test
// 	user_id= <@U030XRXJ2|jboursiquot>
// 	user_name=Steve
// 	command=/shoutout help
// 	text=@username core-value-abbrev comment
// 	response_url=https://hooks.slack.com/commands/1234/5678
func parseShoutoutCommand(ctx context.Context, params *url.Values) (*shoutoutCommand, error) {
	log.Printf("%#v", params)
	cmdRe := regexp.MustCompile("(<@[a-zA-Z0-9_|]*>)(\\s{1})(it|IT|rf|RF|tf|TF){1}(\\s{1})(.*)")
	cmdSubmatches := cmdRe.FindAllStringSubmatch(params.Get("text"), -1)
	if cmdSubmatches == nil {
		return nil, errors.New("not a shoutout command")
	}

	k, err := lookupShoutoutKind(strings.ToUpper(cmdSubmatches[0][3]))
	if err != nil {
		return nil, err
	}

	rec := cmdSubmatches[0][1]
	recRe := regexp.MustCompile("[\\p{L}\\d-_.]+")
	recSubmatches := recRe.FindAllStringSubmatch(rec, -1)
	if recSubmatches == nil {
		return nil, errors.New("failed to obtain recipient")
	}

	s := New()
	s.Kind = k
	s.SenderID = params.Get("user_id")
	s.SenderName = params.Get("user_name")
	s.RecipientID = recSubmatches[0][0]
	s.RecipientName = recSubmatches[1][0]
	s.TeamID = params.Get("team_id")
	s.Comment = cmdSubmatches[0][5]
	s.ResponseURL = params.Get("response_url")
	c := shoutoutCommand{shoutout: s}

	return &c, nil
}

func (c *shoutoutCommand) execute(ctx context.Context) (*SlackResponse, error) {
	bs, err := json.Marshal(c.shoutout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shoutout: %s", err)
	}

	in := sqs.SendMessageInput{
		MessageBody: aws.String(string(bs)),
		QueueUrl:    aws.String(os.Getenv("QUEUE_URL")),
	}

	var res *SlackResponse

	if _, err := c.sqs.SendMessageWithContext(ctx, &in); err != nil {
		res = &SlackResponse{
			ResponseType: SlackResponseTypeEphemeral,
			Text:         "Uh oh. There was a problem sending your shoutout. Try again later.",
		}
		return res, fmt.Errorf("failed to publish shoutout: %s", err)
	}

	res = &SlackResponse{
		ResponseType: SlackResponseTypeEphemeral,
		Text:         "You sent a shoutout!",
	}

	return res, nil
}

func lookupShoutoutKind(ks string) (*ShoutoutKind, error) {
	for k, v := range shoutoutKinds {
		if strings.ToLower(ks) == k {
			return v, nil
		}
	}
	return nil, fmt.Errorf("shoutout kind '%s' is not valid", ks)
}
