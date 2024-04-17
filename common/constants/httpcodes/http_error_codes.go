package httpcodes

const (
	ErrBadRequest                  = 400
	ErrUnauthorized                = 401
	ErrPaymentRequired             = 402
	ErrForbidden                   = 403
	ErrNotFound                    = 404
	ErrMethodNotAllowed            = 405
	ErrNotAcceptable               = 406
	ErrProxyAuthRequired           = 407
	ErrRequestTimeout              = 408
	ErrConflict                    = 409
	ErrGone                        = 410
	ErrLengthRequired              = 411
	ErrPreconditionFailed          = 412
	ErrPayloadTooLarge             = 413
	ErrURITooLong                  = 414
	ErrUnsupportedMediaType        = 415
	ErrRangeNotSatisfiable         = 416
	ErrExpectationFailed           = 417
	ErrTeapot                      = 418
	ErrMisdirectedRequest          = 421
	ErrUnprocessableEntity         = 422
	ErrLocked                      = 423
	ErrFailedDependency            = 424
	ErrTooEarly                    = 425
	ErrUpgradeRequired             = 426
	ErrPreconditionRequired        = 428
	ErrTooManyRequests             = 429
	ErrRequestHeaderFieldsTooLarge = 431
	ErrUnavailableForLegalReasons  = 451

	ErrInternalServerError           = 500
	ErrNotImplemented                = 501
	ErrBadGateway                    = 502
	ErrServiceUnavailable            = 503
	ErrGatewayTimeout                = 504
	ErrHTTPVersionNotSupported       = 505
	ErrVariantAlsoNegotiates         = 506
	ErrInsufficientStorage           = 507
	ErrLoopDetected                  = 508
	ErrNotExtended                   = 510
	ErrNetworkAuthenticationRequired = 511
)
