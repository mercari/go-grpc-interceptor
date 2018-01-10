package acceptlang

import (
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func newMetadataContext(ctx context.Context, val string) context.Context {
	md := metadata.Pairs(DefaultAcceptLangKey, val)
	return metadata.NewIncomingContext(ctx, md)
}

func TestHandleAcceptLanguage(t *testing.T) {
	table := []string{
		"da, en-gb;q=0.8, en;q=0.7",
		"da ,en-gb;q=0.8 ,en;q=0.7",
		"en;q=0.7, en-gb;q=0.8, da",
	}

	for i := range table {
		ctx := context.Background()
		acceptLangs := HandleAcceptLanguage(newMetadataContext(ctx, table[i]))

		t.Logf("header: %s", table[i])
		if got, want := len(acceptLangs), 3; got != want {
			t.Fatalf("expect len() = %d, but got %d", want, got)
		}

		al := acceptLangs[0]
		if got, want := al.Language, "da"; got != want {
			t.Fatalf("expect language = %q, but got %q", want, got)
		}
		if got, want := al.Quality, float32(1); got != want {
			t.Fatalf("expect quality = %f, but got %f", want, got)
		}

		al = acceptLangs[1]
		if got, want := al.Language, "en-gb"; got != want {
			t.Fatalf("expect language = %q, but got %q", want, got)
		}
		if got, want := al.Quality, float32(0.8); got != want {
			t.Fatalf("expect quality = %f, but got %f", want, got)
		}

		al = acceptLangs[2]
		if got, want := al.Language, "en"; got != want {
			t.Fatalf("expect language = %q, but got %q", want, got)
		}
		if got, want := al.Quality, float32(0.7); got != want {
			t.Fatalf("expect quality = %f, but got %f", want, got)
		}
	}
}

func TestHandleAcceptLanguageOrder(t *testing.T) {
	header := "en-gb, da, en"
	ctx := context.Background()
	acceptLangs := HandleAcceptLanguage(newMetadataContext(ctx, header))

	if got, want := len(acceptLangs), 3; got != want {
		t.Fatalf("expect len() = %d, but got %d", want, got)
	}

	al := acceptLangs[0]
	if got, want := al.Language, "en-gb"; got != want {
		t.Fatalf("expect language = %q, but got %q", want, got)
	}
	al = acceptLangs[1]
	if got, want := al.Language, "da"; got != want {
		t.Fatalf("expect language = %q, but got %q", want, got)
	}
	al = acceptLangs[2]
	if got, want := al.Language, "en"; got != want {
		t.Fatalf("expect language = %q, but got %q", want, got)
	}
}

type testServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *testServerStream) Context() context.Context {
	return ss.ctx
}

func (ss *testServerStream) SendMsg(m interface{}) error {
	return nil
}

func (f *testServerStream) RecvMsg(m interface{}) error {
	return nil
}

func TestUnaryServer(t *testing.T) {
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		acceptLangs := FromContext(ctx)
		if got, want := len(acceptLangs), 1; got != want {
			t.Fatalf("expect len() = %d, but got %d", want, got)
		}
		al := acceptLangs[0]
		if got, want := al.Language, "en"; got != want {
			t.Fatalf("expect language = %q, but got %q", want, got)
		}

		return "output", nil
	}

	ctx := context.Background()
	ctx = newMetadataContext(ctx, "en")
	_, err := UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamServerWithoutRequestID(t *testing.T) {
	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		acceptLangs := FromContext(stream.Context())
		if got, want := len(acceptLangs), 1; got != want {
			t.Fatalf("expect len() = %d, but got %d", want, got)
		}
		al := acceptLangs[0]
		if got, want := al.Language, "en"; got != want {
			t.Fatalf("expect language = %q, but got %q", want, got)
		}

		return nil
	}
	testService := struct{}{}
	ctx := context.Background()
	ctx = newMetadataContext(ctx, "en")
	testStream := &testServerStream{ctx: ctx}

	err := StreamServerInterceptor(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
