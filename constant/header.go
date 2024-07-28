package constant

type HttpHeader string

const (
	ContentTypeHeader   HttpHeader = "Content-Type"
	ContentLengthHeader HttpHeader = "Content-Length"
	DateHeader          HttpHeader = "Date"
	ServerHeader        HttpHeader = "Server"
	ConnectionHeader    HttpHeader = "Connection"
)

func (h HttpHeader) String() string {
	return string(h)
}
