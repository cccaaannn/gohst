package constant

type HttpHeader string

const (
	ContentTypeHeader HttpHeader = "Content-Type"
	ContentLength     HttpHeader = "Content-Length"
	Date              HttpHeader = "Date"
	Server            HttpHeader = "Server"
	Connection        HttpHeader = "Connection"
)

func (h HttpHeader) String() string {
	return string(h)
}
