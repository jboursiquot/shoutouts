package shoutouts

import (
	"context"
	"errors"
	"net/url"
)

const (
	helpKindUnspecified = "unspecified"
	helpKindUsage       = "usage"
	helpKindValues      = "values"
)

type helpCommand struct {
	kind string
}

// parseHelp attempts to extract a "help" command from the incoming parameters.
func parseHelpCommand(ctx context.Context, params *url.Values) (*helpCommand, error) {
	text := params.Get("text")
	if text == "help" || text == "" {
		return &helpCommand{kind: helpKindUnspecified}, nil
	}

	if text == "help usage" {
		return &helpCommand{kind: helpKindUsage}, nil
	}

	if text == "help values" {
		return &helpCommand{kind: helpKindValues}, nil
	}

	return nil, errors.New("not a help command")
}

// Execute executes the command and returns the appropriate Slack response to relay.
func (c *helpCommand) execute(ctx context.Context) (*SlackResponse, error) {
	var r SlackResponse
	switch c.kind {
	case helpKindUnspecified:
		fallthrough
	case helpKindUsage:
		r = SlackResponse{
			ResponseType: SlackResponseTypeEphemeral,
			Text:         "This command allows you to give shoutouts to your teammates for exemplifying the core values of the organization.",
			Attachments: []*SlackAttachment{
				{
					Title: "Get help",
					Text:  "`/shoutout help`",
				},
				{
					Title: "Show Core Values",
					Text:  "`/shoutout help values`",
				},
				{
					Title: "Give a shoutout",
					Text:  "`/shoutout <@username> <core-value-abbrev> <comment>`",
				},
				{
					Title: "List shoutouts for a user",
					Text:  "`/shoutout list <@username>`",
				},
			},
		}
	case helpKindValues:
		r = SlackResponse{
			ResponseType: SlackResponseTypeEphemeral,
			Attachments: []*SlackAttachment{
				{Title: "Innovative Thinking | IT", Text: descIT, ThumbURL: thumbIT},
				{Title: "Results Focus | RF", Text: descRF, ThumbURL: thumbRF},
				{Title: "Team First | TF", Text: descTF, ThumbURL: thumbTF},
			},
		}
	}
	return &r, nil
}
