package response

import "net/http"

// Response codes - custom codes yang map ke HTTP status
const (
	Success           = 1000
	Created           = 1001
	Accepted          = 1002
	NoContent         = 1004
	MultipleChoices   = 1300
	MovedPermanently  = 1301
	Found             = 1302
	SeeOther          = 1303
	NotModified       = 1304
	TemporaryRedirect = 1307
	PermanentRedirect = 1308

	BadRequest                   = 1400
	Unauthorized                 = 1401
	PaymentRequired              = 1402
	Forbidden                    = 1403
	NotFound                     = 1404
	MethodNotAllowed             = 1405
	NotAcceptable                = 1406
	ProxyAuthRequired            = 1407
	RequestTimeout               = 1408
	Conflict                     = 1409
	Gone                         = 1410
	LengthRequired               = 1411
	PreconditionFailed           = 1412
	RequestEntityTooLarge        = 1413
	RequestURITooLong            = 1414
	UnsupportedMediaType         = 1415
	RequestedRangeNotSatisfiable = 1416
	ExpectationFailed            = 1417
	UnprocessableEntity          = 1422
	TooManyRequests              = 1429

	ServerError        = 1500
	NotImplemented     = 1501
	BadGateway         = 1502
	ServiceUnavailable = 1503
	GatewayTimeout     = 1504
)

// httpToCustomCode maps custom response code to HTTP status code
var httpToCustomCode = map[int]int{}

func init() {
	statusMap := map[int]int{
		Success:                      http.StatusOK,
		Created:                      http.StatusCreated,
		Accepted:                     http.StatusAccepted,
		NoContent:                    http.StatusNoContent,
		MultipleChoices:              http.StatusMultipleChoices,
		MovedPermanently:             http.StatusMovedPermanently,
		Found:                        http.StatusFound,
		SeeOther:                     http.StatusSeeOther,
		NotModified:                  http.StatusNotModified,
		TemporaryRedirect:            http.StatusTemporaryRedirect,
		PermanentRedirect:            http.StatusPermanentRedirect,
		BadRequest:                   http.StatusBadRequest,
		Unauthorized:                 http.StatusUnauthorized,
		PaymentRequired:              http.StatusPaymentRequired,
		Forbidden:                    http.StatusForbidden,
		NotFound:                     http.StatusNotFound,
		MethodNotAllowed:             http.StatusMethodNotAllowed,
		NotAcceptable:                http.StatusNotAcceptable,
		ProxyAuthRequired:            http.StatusProxyAuthRequired,
		RequestTimeout:               http.StatusRequestTimeout,
		Conflict:                     http.StatusConflict,
		Gone:                         http.StatusGone,
		LengthRequired:               http.StatusLengthRequired,
		PreconditionFailed:           http.StatusPreconditionFailed,
		RequestEntityTooLarge:        http.StatusRequestEntityTooLarge,
		RequestURITooLong:            http.StatusRequestURITooLong,
		UnsupportedMediaType:         http.StatusUnsupportedMediaType,
		RequestedRangeNotSatisfiable: http.StatusRequestedRangeNotSatisfiable,
		ExpectationFailed:            http.StatusExpectationFailed,
		UnprocessableEntity:          http.StatusUnprocessableEntity,
		TooManyRequests:              http.StatusTooManyRequests,
		ServerError:                  http.StatusInternalServerError,
		NotImplemented:               http.StatusNotImplemented,
		BadGateway:                   http.StatusBadGateway,
		ServiceUnavailable:           http.StatusServiceUnavailable,
		GatewayTimeout:               http.StatusGatewayTimeout,
	}

	for customCode, httpStatus := range statusMap {
		httpToCustomCode[customCode] = httpStatus
	}
}

// GetHTTPStatus maps a custom response code to HTTP status code
func GetHTTPStatus(customCode int) int {
	if code, exists := httpToCustomCode[customCode]; exists {
		return code
	}
	return http.StatusInternalServerError
}
