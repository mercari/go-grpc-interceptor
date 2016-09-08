package requestdump

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/transport"
)

func init() {
	marshaler = jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		OrigName:     true,
	}
}

var marshaler jsonpb.Marshaler

type protoMessage struct {
	msg proto.Message
}

func (m protoMessage) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, m.msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type metadataMarshaller metadata.MD

func (mm metadataMarshaller) MarshalLog(kv zap.KeyValue) error {
	for k, v := range mm {
		// TODO: array handling
		if len(v) == 0 {
			continue
		}

		value := v[0]
		if strings.HasSuffix(k, "-bin") {
			zap.Base64(k, []byte(value)).AddTo(kv)
		} else {
			kv.AddString(k, value)
		}
	}
	return nil
}

func dump(ctx context.Context, opts options, logger zap.Logger, info *grpc.UnaryServerInfo, request bool, msg interface{}, err error) {
	direction := "request"
	code := zap.Skip()
	header := zap.Skip()
	trailer := zap.Skip()
	if request {
		if md, ok := metadata.FromContext(ctx); ok {
			header = zap.Marshaler("header", metadataMarshaller(md))
		}
	} else {
		direction = "response"
		code = zap.String("code", grpc.Code(err).String())

		// NOTE: not sure this is public API, will be broken in futre
		if stream, ok := transport.StreamFromContext(ctx); ok {
			md := stream.Trailer()
			trailer = zap.Marshaler("trailer", metadataMarshaller(md))
		}

		// NOTE: there is no API to get response header
	}

	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}

	fields := []zap.Field{
		zap.String("method", info.FullMethod),
		zap.String("direction", direction),
		zap.String("addr", addr),
		code,
		header,
		trailer,
	}

	if err != nil {
		logger.Info("request dump",
			zap.Nest(opts.rootKey,
				append(fields,
					zap.String("error", err.Error()),
				)...,
			),
		)
		return
	}

	protoMsg, ok := msg.(proto.Message)
	if !ok {
		logger.Info("request dump",
			zap.Nest(opts.rootKey,
				append(fields,
					zap.String("error", fmt.Sprintf("not proto.Message: %v", msg)),
				)...,
			),
		)
		return
	}

	logger.Info("request dump",
		zap.Nest(opts.rootKey,
			append(fields,
				zap.Object("message", protoMessage{msg: protoMsg}),
			)...,
		),
	)
	return
}
