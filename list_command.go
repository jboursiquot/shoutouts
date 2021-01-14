package shoutouts

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

type listCommand struct {
	userid string
	ddb    DynamoDBQuerier
}

// parseListCommand parses the list command given the key=value pairs in the params map
// 	token=PjsIpgv7fGCKHB6bK4EAx43U
// 	team_id=T0001
// 	team_domain=example
// 	channel_id=C2147483705
// 	channel_name=test
// 	user_id= <@U030XRXJ2|jboursiquot>
// 	user_name=Steve
// 	command=/shoutout help
// 	text=list @username
// 	response_url=https://hooks.slack.com/commands/1234/5678
func parseListCommand(ctx context.Context, params *url.Values) (*listCommand, error) {
	cmdRe, err := regexp.Compile("(list)\\s(<@[a-zA-Z0-9_|]*>)")
	if err != nil {
		return nil, errors.New("failed to compile list command regex")
	}

	cmdSubmatches := cmdRe.FindAllStringSubmatch(params.Get("text"), -1)
	if cmdSubmatches == nil {
		return nil, errors.New("not a list command")
	}

	rec := cmdSubmatches[0][2]
	recRe, err := regexp.Compile("@([a-zA-Z0-9]*)")
	if err != nil {
		return nil, errors.New("failed to compile userid regex")
	}

	recSubmatches := recRe.FindAllStringSubmatch(rec, -1)
	if recSubmatches == nil {
		return nil, errors.New("failed to obtain recipient")
	}

	return &listCommand{userid: recSubmatches[0][1]}, nil
}

func (c *listCommand) execute(ctx context.Context) (*SlackResponse, error) {
	l := NewLister(c.ddb)
	list, err := l.List(ctx, c.userid)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}

	var attachments []*SlackAttachment
	for _, s := range list {
		a := &SlackAttachment{}
		a.Title = fmt.Sprintf("%s shoutout from <@%s>", s.Kind.Abbrev, s.SenderID)
		a.Text = fmt.Sprintf("\"%s\"", s.Comment)
		a.ThumbURL = s.Kind.ThumbURL
		attachments = append(attachments, a)
	}

	res := SlackResponse{
		ResponseType: SlackResponseTypeInChannel,
		Attachments:  attachments,
	}

	if attachments == nil {
		res.Text = "None"
	}

	return &res, nil
}
