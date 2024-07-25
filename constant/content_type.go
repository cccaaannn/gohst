package constant

type ContentType string

const (
	ApplicationJson ContentType = "application/json"
	TextHtml        ContentType = "text/html"
)

func (ct ContentType) String() string {
	return string(ct)
}
