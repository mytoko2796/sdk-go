package httpheader

import "go.opencensus.io/plugin/ochttp/propagation/b3"

const (
	// B3 Propagation Header
	B3TraceID string = b3.TraceIDHeader
	B3SpanID  string = b3.SpanIDHeader
	B3Sampled string = b3.SampledHeader

	// HTTP Header Standard
	RequestID      string = `x-request-id`
	RequestMethod  string = `x-request-method`
	RequestScheme  string = `x-request-scheme`
	KeyServerRoute string = `x-key-server-route`
	ForwardedFor   string = `x-forwarded-for`

	// Custom HTTP Header
	SessionID     string = `x-session-id`
	UserID        string = `x-user-id`
	MerchantID    string = `x-merchant-id`
	AgentID       string = `x-agent-id`
	UserBasicInfo string = `x-user-basic-info`
	ClientID      string = `x-client-id`
	AppLang       string = `x-app-lang`
	AppDebug      string = `x-app-debug`

	// Lang Header
	LangEN  string = `EN`
	LangIDN string = `IDN`

	// UserAgent Header
	UserAgent                  string = `User-Agent`
	UserAgentHTTPClientDefault string = `SDKdefault/1.0`
	ContentAccept              string = `Accept`
	ContentType                string = `Content-Type`
	ContentJSON                string = `application/json`
	ContentXML                 string = `application/xml`
	ContentFormURLEncoded      string = `application/x-www-form-urlencoded`

	// Cache Control Header
	CacheControl        string = `Cache-Control`
	CacheNoCache        string = `no-cache`
	CacheNoStore        string = `no-store`
	CacheMustRevalidate string = `must-revalidate`

	// bpm
	BPMProcessID  string = `x-bpm-process-id`
	BPMWorkflowID string = `x-bpm-workflow-id`
	BPMInstanceID string = `x-bpm-instance-id`
	BPMJobID      string = `x-bpm-job-id`
	BPMJobType    string = `x-bpm-job-type`
)

