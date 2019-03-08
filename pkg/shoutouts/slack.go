package shoutouts

const (
	// SlackResponseTypeInChannel is the in_channel response.
	SlackResponseTypeInChannel = "in_channel"

	// SlackResponseTypeEphemeral is the ephemeral response visible only to the interacting user.
	SlackResponseTypeEphemeral = "ephemeral"
)

// SlackField is part of the SlackAttachment.
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SlackAction is part of the SlackAttachment.
type SlackAction struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SlackAttachment is a part of a rich response to a Slack message.
type SlackAttachment struct {
	Title          string         `json:"title"`
	Fields         []*SlackField  `json:"fields,omitempty"`
	AuthorName     string         `json:"author_name,omitempty"`
	AuthorIcon     string         `json:"author_icon,omitempty"`
	ImageURL       string         `json:"image_url,omitempty"`
	ThumbURL       string         `json:"thumb_url,omitempty"`
	Text           string         `json:"text,omitempty"`
	Fallback       string         `json:"fallback,omitempty"`
	CallbackID     string         `json:"callback_id,omitempty"`
	Color          string         `json:"color,omitempty"`
	AttachmentType string         `json:"attachment_type,omitempty"`
	Actions        []*SlackAction `json:"actions,omitempty"`
}

// SlackResponse is the response sent to Slack after capturing the request.
type SlackResponse struct {
	ResponseType string             `json:"response_type"`
	Text         string             `json:"text"`
	Attachments  []*SlackAttachment `json:"attachments"`
}
