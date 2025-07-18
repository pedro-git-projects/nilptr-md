package httpext

import "net/http"

type HeaderKey string
type HeaderValue string

const (
	ContentType                   HeaderKey = "Content-Type"
	Accept                        HeaderKey = "Accept"
	Authorization                 HeaderKey = "Authorization"
	IdempotencyKey                HeaderKey = "Idempotency-Key"
	RequestID                     HeaderKey = "X-Request-ID"
	CacheControl                  HeaderKey = "Cache-Control"
	AcceptEncoding                HeaderKey = "Accept-Encoding"
	Link                          HeaderKey = "Link"
	StrictTransportSecurity       HeaderKey = "Strict-Transport-Security"
	ContentTypeOptions            HeaderKey = "X-Content-Type-Options"
	ReferrerPolicy                HeaderKey = "Referrer-Policy"
	PermissionsPolicy             HeaderKey = "Permissions-Policy"
	AccessControlAllowOrigin      HeaderKey = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials HeaderKey = "Access-Control-Allow-Credentials"
	AccessControlAllowHeaders     HeaderKey = "Access-Control-Allow-Headers"
)

const (
	ApplicationJSON          HeaderValue = "application/json"
	TextHTML                 HeaderValue = "text/html"
	TextMarkdown             HeaderValue = "text/markdown"
	BearerScheme             HeaderValue = "Bearer %s" // format with token
	NoStore                  HeaderValue = "no-store"
	MaxAge60Public           HeaderValue = "max-age=60, public"
	EncodingGzipBr           HeaderValue = "gzip, br"
	ETagWeak                 HeaderValue = `W/"%s"` // format with hash
	IfNoneMatchAny           HeaderValue = "*"
	LinkRelNext              HeaderValue = `<%s>; rel="next"` // format with URL
	LinkRelPrev              HeaderValue = `<%s>; rel="prev"`
	HSTS                     HeaderValue = "max-age=31536000; includeSubDomains"
	NoSniff                  HeaderValue = "nosniff"
	ReferrerNoReferrer       HeaderValue = "no-referrer"
	PermissionsDefault       HeaderValue = "geolocation=(), microphone=()"
	CORSAllowOriginAll       HeaderValue = "*" // add domain
	CORSAllowCredentialsTrue HeaderValue = "true"
	CORSAllowHeadersDefault  HeaderValue = "Content-Type, Authorization, X-Request-ID"
)

func (k HeaderKey) Add(h http.Header, v HeaderValue) {
	h.Add(string(k), string(v))
}
