// status_code.go kee > 2021/02/17

package router

import (
	"net/http"
)

const (
	// 2xx
	StatusCodeOK                          = 200
	StatusCodeCreated                     = 201
	StatusCodeAccepted                    = 202
	StatusCodeNonAuthoritativeInformation = 203
	StatusCodeNoContent                   = 204
	StatusCodeResetContent                = 205
	StatusCodePartialContent              = 206
	StatusCodeMultiStatus                 = 207
	StatusCodeAlreadyReported             = 208

	// 3xx
	StatusCodeMultipleChoices   = 300
	StatusCodeMovedPermanently  = 301
	StatusCodeFound             = 302
	StatusCodeSeeOther          = 303
	StatusCodeNotModified       = 304
	StatusCodeUseProxy          = 305
	StatusCodeTemporaryRedirect = 307
	StatusCodePermanentRedirect = 308

	// 4xx
	StatusCodeBadRequest                   = 400
	StatusCodeUnauthorized                 = 401
	StatusCodePaymentRequired              = 402
	StatusCodeForbidden                    = 403
	StatusCodeNotFound                     = 404
	StatusCodeMethodNotAllowed             = 405
	StatusCodeNotAcceptable                = 406
	StatusCodeProxyAuthenticationRequired  = 407
	StatusCodeRequestTimeout               = 408
	StatusCodeConflict                     = 409
	StatusCodeGone                         = 410
	StatusCodeLengthRequired               = 411
	StatusCodePreconditionFailed           = 412
	StatusCodeRequestEntityTooLarge        = 413
	StatusCodeRequestURITooLong            = 414
	StatusCodeUnsupportedMediaType         = 415
	StatusCodeRequestedRangeNotSatisfiable = 416
	StatusCodeExpectationFailed            = 417
	StatusCodeAsTeapot                     = 418 // 反爬虫

	// 4xx-other
	StatusCodePageExpired                      = 419
	StatusCodeBlockedByWindowsParentalControls = 450
	StatusCodeInvalidToken                     = 498
	StatusCodeTokenRequired                    = 499

	// 5xx
	StatusCodeInternalServerError           = 500
	StatusCodeNotImplemented                = 501
	StatusCodeBadGateway                    = 502
	StatusCodeServiceUnavailable            = 503
	StatusCodeGatewayTimeout                = 504
	StatusCodeHTTPVersionNotSupported       = 505
	StatusCodeVariantAlsoNegotiates         = 506
	StatusCodeInsufficientStorage           = 507
	StatusCodeLoopDetected                  = 508
	StatusCodeBandwidthLimitExceeded        = 509
	StatusCodeNotExtended                   = 510
	StatusCodeNetworkAuthenticationRequired = 511
	StatusCodeInvalidSSLCertificate         = 526
	StatusCodeSiteOverloaded                = 529
	StatusCodeSiteFrozen                    = 530
	StatusCodeNetworkReadTimeout            = 598
)

var unofficialStatusText = map[int]string{
	StatusCodePageExpired:                      "Page Expired",
	StatusCodeBlockedByWindowsParentalControls: "Blocked by Windows Parental Controls",
	StatusCodeInvalidToken:                     "Invalid Token",
	StatusCodeTokenRequired:                    "Token Required",
	StatusCodeBandwidthLimitExceeded:           "Bandwidth Limit Exceeded",
	StatusCodeInvalidSSLCertificate:            "Invalid SSL Certificate",
	StatusCodeSiteOverloaded:                   "Site is overloaded",
	StatusCodeSiteFrozen:                       "Site is frozen",
	StatusCodeNetworkReadTimeout:               "Network read timeout error",
}

func StatusText(code int) (text string) {
	if text = http.StatusText(code); text != "" {
		return
	}

	text = unofficialStatusText[code]
	return
}
