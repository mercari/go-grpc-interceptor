package acceptlang

import (
	"sort"
	"strconv"
	"strings"

	multiint "github.com/mercari/go-grpc-interceptor/multiinterceptor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// DefaultAcceptLangKey is metadata key name for accept language
var DefaultAcceptLangKey = "accept-language"

var _ grpc.UnaryServerInterceptor = UnaryServerInterceptor

type alKey struct{}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	acceptLangs := HandleAcceptLanguage(ctx)
	ctx = context.WithValue(ctx, alKey{}, acceptLangs)
	return handler(ctx, req)
}

func StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := stream.Context()
	acceptLangs := HandleAcceptLanguage(ctx)
	ctx = context.WithValue(ctx, alKey{}, acceptLangs)
	stream = multiint.NewServerStreamWithContext(stream, ctx)
	return handler(srv, stream)
}

func FromContext(ctx context.Context) AcceptLanguages {
	al, ok := ctx.Value(alKey{}).(AcceptLanguages)
	if !ok || al == nil {
		return []AcceptLanguage{}
	}
	return al
}

func HandleAcceptLanguage(ctx context.Context) AcceptLanguages {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	header, ok := md[DefaultAcceptLangKey]
	if !ok || len(header) == 0 {
		return nil
	}

	acceptLangHeader := header[0]
	acceptLangHeaderSlice := strings.Split(acceptLangHeader, ",")

	acceptLangs := make(AcceptLanguages, len(acceptLangHeaderSlice))
	for i, lang := range acceptLangHeaderSlice {
		lang = strings.TrimSpace(lang)
		qualSlice := strings.Split(lang, ";q=")
		if len(qualSlice) == 2 {
			qual, err := strconv.ParseFloat(qualSlice[1], 32)
			if err != nil {
				acceptLangs[i] = newAcceptLanguage(qualSlice[0], 1)
			} else {
				acceptLangs[i] = newAcceptLanguage(qualSlice[0], float32(qual))
			}
		} else {
			acceptLangs[i] = newAcceptLanguage(lang, 1)
		}
	}

	sort.Sort(sort.Reverse(acceptLangs))
	return acceptLangs
}

type AcceptLanguage struct {
	Language string
	Quality  float32
}

func newAcceptLanguage(lang string, qual float32) AcceptLanguage {
	return AcceptLanguage{Language: lang, Quality: qual}
}

type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Languages() []string {
	langs := make([]string, len(al))
	for i := range al {
		langs[i] = al[i].Language
	}
	return langs
}

func (al AcceptLanguages) Len() int {
	return len(al)
}

func (al AcceptLanguages) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

func (al AcceptLanguages) Less(i, j int) bool {
	return al[i].Quality < al[j].Quality
}
