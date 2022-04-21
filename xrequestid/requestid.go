package xrequestid

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "x-request-id"

func HandleRequestID(ctx context.Context, validator requestIDValidator) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[DefaultXRequestIDKey]
	if !ok || len(header) == 0 {
		return newRequestID()
	}

	requestID := header[0]
	if requestID == "" {
		return newRequestID()
	}

	if !validator(requestID) {
		return newRequestID()
	}

	return requestID
}

func HandleRequestIDChain(ctx context.Context, validator requestIDValidator) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[DefaultXRequestIDKey]
	if !ok || len(header) == 0 {
		return newRequestID()
	}

	requestID := header[0]
	if requestID == "" {
		return newRequestID()
	}

	if !validator(requestID) {
		return newRequestID()
	}

	return fmt.Sprintf("%s,%s", requestID, newRequestID())
}

func newRequestID() string {
	return uuid.NewString()
}
