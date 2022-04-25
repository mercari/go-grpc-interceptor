package xrequestid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

// Logs the incoming request with the new request id. full.const
// FullMethod is the full RPC method string, i.e., /package.service/method.
func logRequestWithID(requestData interface{}, requestID, fullMethod string) {
	methodPath := strings.Split(fullMethod, "/")
	if len(methodPath) > 1 {
		logrus.WithFields(logrus.Fields{
			"Request Data":    fmt.Sprintf("%+v", requestData),
			"Request ID":      requestID,
			"Package.Service": methodPath[0],
			"Method Name":     methodPath[1],
		}).Infof("Request ID appended to request")
	} else {
		logrus.WithFields(logrus.Fields{
			"Request Data":     fmt.Sprintf("%+v", requestData),
			"Request ID":       requestID,
			"Full Method Path": fullMethod,
		}).Infof("Request ID appended to request")
	}
}
