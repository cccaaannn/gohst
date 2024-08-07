package constant

import "strconv"

type HTTPStatusCode int

const (
	ContinueStatus                      HTTPStatusCode = 100
	SwitchingProtocolsStatus            HTTPStatusCode = 101
	ProcessingStatus                    HTTPStatusCode = 102
	EarlyHintsStatus                    HTTPStatusCode = 103
	OkStatus                            HTTPStatusCode = 200
	CreatedStatus                       HTTPStatusCode = 201
	AcceptedStatus                      HTTPStatusCode = 202
	NonAuthoritativeInformationStatus   HTTPStatusCode = 203
	NoContentStatus                     HTTPStatusCode = 204
	ResetContentStatus                  HTTPStatusCode = 205
	PartialContentStatus                HTTPStatusCode = 206
	MultiStatusStatus                   HTTPStatusCode = 207
	AlreadyReportedStatus               HTTPStatusCode = 208
	ImUsedStatus                        HTTPStatusCode = 226
	MultipleChoicesStatus               HTTPStatusCode = 300
	MovedPermanentlyStatus              HTTPStatusCode = 301
	FoundStatus                         HTTPStatusCode = 302
	SeeOtherStatus                      HTTPStatusCode = 303
	NotModifiedStatus                   HTTPStatusCode = 304
	UseProxyStatus                      HTTPStatusCode = 305
	SwitchProxyStatus                   HTTPStatusCode = 306
	TemporaryRedirectStatus             HTTPStatusCode = 307
	PermanentRedirectStatus             HTTPStatusCode = 308
	BadRequestStatus                    HTTPStatusCode = 400
	UnauthorizedStatus                  HTTPStatusCode = 401
	PaymentRequiredStatus               HTTPStatusCode = 402
	ForbiddenStatus                     HTTPStatusCode = 403
	NotFoundStatus                      HTTPStatusCode = 404
	MethodNotAllowedStatus              HTTPStatusCode = 405
	NotAcceptableStatus                 HTTPStatusCode = 406
	ProxyAuthenticationRequiredStatus   HTTPStatusCode = 407
	RequestTimeoutStatus                HTTPStatusCode = 408
	ConflictStatus                      HTTPStatusCode = 409
	GoneStatus                          HTTPStatusCode = 410
	LengthRequiredStatus                HTTPStatusCode = 411
	PreconditionFailedStatus            HTTPStatusCode = 412
	PayloadTooLargeStatus               HTTPStatusCode = 413
	UriTooLongStatus                    HTTPStatusCode = 414
	UnsupportedMediaTypeStatus          HTTPStatusCode = 415
	RangeNotSatisfiableStatus           HTTPStatusCode = 416
	ExpectationFailedStatus             HTTPStatusCode = 417
	ImATeapotStatus                     HTTPStatusCode = 418
	MisdirectedRequestStatus            HTTPStatusCode = 421
	UnprocessableEntityStatus           HTTPStatusCode = 422
	LockedStatus                        HTTPStatusCode = 423
	FailedDependencyStatus              HTTPStatusCode = 424
	TooEarlyStatus                      HTTPStatusCode = 425
	UpgradeRequiredStatus               HTTPStatusCode = 426
	PreconditionRequiredStatus          HTTPStatusCode = 428
	TooManyRequestsStatus               HTTPStatusCode = 429
	RequestHeaderFieldsTooLargeStatus   HTTPStatusCode = 431
	UnavailableForLegalReasonsStatus    HTTPStatusCode = 451
	InternalServerErrorStatus           HTTPStatusCode = 500
	NotImplementedStatus                HTTPStatusCode = 501
	BadGatewayStatus                    HTTPStatusCode = 502
	ServiceUnavailableStatus            HTTPStatusCode = 503
	GatewayTimeoutStatus                HTTPStatusCode = 504
	HttpVersionNotSupportedStatus       HTTPStatusCode = 505
	VariantAlsoNegotiatesStatus         HTTPStatusCode = 506
	InsufficientStorageStatus           HTTPStatusCode = 507
	LoopDetectedStatus                  HTTPStatusCode = 508
	NotExtendedStatus                   HTTPStatusCode = 510
	NetworkAuthenticationRequiredStatus HTTPStatusCode = 511
)

func (c HTTPStatusCode) String() string {
	return strconv.Itoa(int(c))
}

func (c HTTPStatusCode) Verb() string {
	switch c {
	case ContinueStatus:
		return "Continue"
	case SwitchingProtocolsStatus:
		return "Switching protocols"
	case ProcessingStatus:
		return "Processing"
	case EarlyHintsStatus:
		return "Early Hints"
	case OkStatus:
		return "OK"
	case CreatedStatus:
		return "Created"
	case AcceptedStatus:
		return "Accepted"
	case NonAuthoritativeInformationStatus:
		return "Non-Authoritative Information"
	case NoContentStatus:
		return "No Content"
	case ResetContentStatus:
		return "Reset Content"
	case PartialContentStatus:
		return "Partial Content"
	case MultiStatusStatus:
		return "Multi-Status"
	case AlreadyReportedStatus:
		return "Already Reported"
	case ImUsedStatus:
		return "IM Used"
	case MultipleChoicesStatus:
		return "Multiple Choices"
	case MovedPermanentlyStatus:
		return "Moved Permanently"
	case FoundStatus:
		return "Found (Previously “Moved Temporarily”)"
	case SeeOtherStatus:
		return "See Other"
	case NotModifiedStatus:
		return "Not Modified"
	case UseProxyStatus:
		return "Use Proxy"
	case SwitchProxyStatus:
		return "Switch Proxy"
	case TemporaryRedirectStatus:
		return "Temporary Redirect"
	case PermanentRedirectStatus:
		return "Permanent Redirect"
	case BadRequestStatus:
		return "Bad Request"
	case UnauthorizedStatus:
		return "Unauthorized"
	case PaymentRequiredStatus:
		return "Payment Required"
	case ForbiddenStatus:
		return "Forbidden"
	case NotFoundStatus:
		return "Not Found"
	case MethodNotAllowedStatus:
		return "Method Not Allowed"
	case NotAcceptableStatus:
		return "Not Acceptable"
	case ProxyAuthenticationRequiredStatus:
		return "Proxy Authentication Required"
	case RequestTimeoutStatus:
		return "Request Timeout"
	case ConflictStatus:
		return "Conflict"
	case GoneStatus:
		return "Gone"
	case LengthRequiredStatus:
		return "Length Required"
	case PreconditionFailedStatus:
		return "Precondition Failed"
	case PayloadTooLargeStatus:
		return "Payload Too Large"
	case UriTooLongStatus:
		return "URI Too Long"
	case UnsupportedMediaTypeStatus:
		return "Unsupported Media Type"
	case RangeNotSatisfiableStatus:
		return "Range Not Satisfiable"
	case ExpectationFailedStatus:
		return "Expectation Failed"
	case ImATeapotStatus:
		return "I’m a Teapot"
	case MisdirectedRequestStatus:
		return "Misdirected Request"
	case UnprocessableEntityStatus:
		return "Unprocessable Entity"
	case LockedStatus:
		return "Locked"
	case FailedDependencyStatus:
		return "Failed Dependency"
	case TooEarlyStatus:
		return "Too Early"
	case UpgradeRequiredStatus:
		return "Upgrade Required"
	case PreconditionRequiredStatus:
		return "Precondition Required"
	case TooManyRequestsStatus:
		return "Too Many Requests"
	case RequestHeaderFieldsTooLargeStatus:
		return "Request Header Fields Too Large"
	case UnavailableForLegalReasonsStatus:
		return "Unavailable For Legal Reasons"
	case InternalServerErrorStatus:
		return "Internal Server Error"
	case NotImplementedStatus:
		return "Not Implemented"
	case BadGatewayStatus:
		return "Bad Gateway"
	case ServiceUnavailableStatus:
		return "Service Unavailable"
	case GatewayTimeoutStatus:
		return "Gateway Timeout"
	case HttpVersionNotSupportedStatus:
		return "HTTP Version Not Supported"
	case VariantAlsoNegotiatesStatus:
		return "Variant Also Negotiates"
	case InsufficientStorageStatus:
		return "Insufficient Storage"
	case LoopDetectedStatus:
		return "Loop Detected"
	case NotExtendedStatus:
		return "Not Extended"
	case NetworkAuthenticationRequiredStatus:
		return "Network Authentication Required"
	default:
		return "Unknown"
	}
}
