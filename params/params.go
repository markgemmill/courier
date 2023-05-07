package params

type CourierParams struct {
	Host     string
	Port     int
	User     string
	Password string
}

type EnvelopParams struct {
	SendFrom string
	ReplyTo  string
	SendTo   []string
	SendCc   []string
	SendBcc  []string
}

type MessageParams struct {
	HighPriority bool
	Subject      string
	TextMessage  string
	HtmlMessage  string
	TemplateType string
	TemplateData map[string]any
	Attachments  []string
}

type Parameters struct {
	CourierParams
	EnvelopParams
	MessageParams
}

func SetMessage(p *Parameters, message string, html bool) {
	if html {
		p.HtmlMessage = message
	} else {
		p.TextMessage = message
	}
}

// SetTemplateData is a helper function for translating a mapping of type T map[string]T
// to map[string]any
func SetTemplateData[T any](p *Parameters, m map[string]T) {
	p.TemplateData = make(map[string]interface{}, len(m))
	for k, v := range m {
		p.TemplateData[k] = v
	}
}
