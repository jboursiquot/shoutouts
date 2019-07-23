package shoutouts

const (
	descIT = "We enthusiastically seize possibilities and act with an entrepreneurial spirit."
	descRF = "We maintain a relentless focus on achieving excellence."
	descTF = "We achieve results as a supportive interdependent team."

	thumbIT = "https://emoji.slack-edge.com/T030XRXJ0/idea/5a620997a0eb6fe0.png"
	thumbRF = "https://emoji.slack-edge.com/T030XRXJ0/focus/fbfa7f9d78f679af.png"
	thumbTF = "https://emoji.slack-edge.com/T030XRXJ0/team/120d0991426017b8.png"
)

var shoutoutKinds map[string]*ShoutoutKind

func init() {
	shoutoutKinds = make(map[string]*ShoutoutKind, 6)
	shoutoutKinds["it"] = &ShoutoutKind{Name: "Innovative Thinking", Abbrev: "IT", Desc: descIT, ThumbURL: thumbIT}
	shoutoutKinds["rf"] = &ShoutoutKind{Name: "Results Focused", Abbrev: "RF", Desc: descRF, ThumbURL: thumbRF}
	shoutoutKinds["tf"] = &ShoutoutKind{Name: "Team First", Abbrev: "TF", Desc: descTF, ThumbURL: thumbTF}
}
