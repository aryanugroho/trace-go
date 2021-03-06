package trace

import "net/http"

const (
	// HeaderSpanID is the header name in requests for the SpanID.
	HeaderSpanID = "Trace-Span"
	// HeaderTraceID is the header name in requests for the TraceID.
	HeaderTraceID = "Trace-Trace"
)

// NewSpanIDFromRequest generates a new SpanID from a request.
// If there is no span information in the request, a root SpanID will be generated.
func NewSpanIDFromRequest(req *http.Request) SpanID {
	parentSpanID := ParseSpanID(req)
	if parentSpanID.IsEmpty() {
		return NewRootSpanID()
	}
	return NewSpanID(parentSpanID)
}

// ParseSpanID returns the SpanID embedded in the request.
// If the ID parsing fails, the ID is set to 0.
func ParseSpanID(req *http.Request) SpanID {
	spanIDStr := req.Header.Get(HeaderSpanID)
	traceIDStr := req.Header.Get(HeaderTraceID)
	spanID, _ := ParseID(spanIDStr)
	traceID, _ := ParseID(traceIDStr)
	return SpanID{
		Trace: traceID,
		Span:  spanID,
	}
}

// SetHeaders sets the Trace HTTP headers from a SpanID
// Returns a new request as convention:
// https://groups.google.com/forum/#!topic/golang-nuts/-j6p12SSpXI
func SetHeaders(req *http.Request, s SpanID) *http.Request {
	// https://groups.google.com/forum/#!topic/golang-nuts/-j6p12SSpXI
	req2 := copyRequest(req)
	req2.Header.Set(HeaderSpanID, s.Span.String())
	req2.Header.Set(HeaderTraceID, s.Trace.String())
	return req2
}

func copyRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
