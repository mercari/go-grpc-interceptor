package i18n

import (
	"github.com/mercari/go-grpc-interceptor/acceptlang"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var defaultLanguage = "en"

func SetDefaultLanguage(lang string) {
	defaultLanguage = lang
}

var _ grpc.UnaryServerInterceptor = UnaryServerInterceptor

type tfuncKey struct{}

func UnaryServerInterceptor(origctx context.Context, origreq interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return acceptlang.UnaryServerInterceptor(origctx, origreq, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		acceptLangs := acceptlang.FromContext(ctx)
		tfunc := HandleI18n(acceptLangs)
		ctx = context.WithValue(ctx, tfuncKey{}, tfunc)
		return handler(ctx, req)
	})
}

func HandleI18n(acceptLangs acceptlang.AcceptLanguages) i18n.TranslateFunc {
	langs := acceptLangs.Languages()
	langs = append(langs, defaultLanguage)
	return i18n.MustTfunc(langs[0], langs[1:]...)
}

func MustTfunc(ctx context.Context) i18n.TranslateFunc {
	tfunc, ok := ctx.Value(tfuncKey{}).(i18n.TranslateFunc)
	if !ok {
		panic("could not find TranslateFunc from context")
	}
	return tfunc
}
