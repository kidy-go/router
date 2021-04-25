// status_code.go kee > 2021/02/17

package router

import (
	"net/http"
)

const (
	// 2xx
	StatusOK                          = 200
	StatusCreated                     = 201
	StatusAccepted                    = 202
	StatusNonAuthoritativeInformation = 203
	StatusNoContent                   = 204
	StatusResetContent                = 205
	StatusPartialContent              = 206
	StatusMultiStatus                 = 207
	StatusAlreadyReported             = 208

	// 3xx
	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307
	StatusPermanentRedirect = 308

	// 4xx
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthenticationRequired  = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusAsTeapot                     = 418 // 反爬虫

	// 4xx-other
	StatusPageExpired                      = 419
	StatusBlockedByWindowsParentalControls = 450
	StatusInvalidToken                     = 498
	StatusTokenRequired                    = 499

	// 5xx
	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusBandwidthLimitExceeded        = 509
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
	StatusInvalidSSLCertificate         = 526
	StatusSiteOverloaded                = 529
	StatusSiteFrozen                    = 530
	StatusNetworkReadTimeout            = 598
)

var unofficialStatusText = map[int]string{
	StatusPageExpired:                      "Page Expired",
	StatusBlockedByWindowsParentalControls: "Blocked by Windows Parental Controls",
	StatusInvalidToken:                     "Invalid Token",
	StatusTokenRequired:                    "Token Required",
	StatusBandwidthLimitExceeded:           "Bandwidth Limit Exceeded",
	StatusInvalidSSLCertificate:            "Invalid SSL Certificate",
	StatusSiteOverloaded:                   "Site is overloaded",
	StatusSiteFrozen:                       "Site is frozen",
	StatusNetworkReadTimeout:               "Network read timeout error",
}

func StatusText(code int) (text string) {
	if text = http.StatusText(code); text != "" {
		return
	}

	text = unofficialStatusText[code]
	return
}
