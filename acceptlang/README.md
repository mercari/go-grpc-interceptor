# acceptlang

acceptlang is an grpc interceptor which parses Accept-Language from metadata like HTTP header and set the AcceptLanguage to context.

## Usage

```golang
import (
 "github.com/eltorocorp/go-grpc-request-id-interceptor/acceptlang"
 "golang.org/x/net/context"
)

func main() {
 uIntOpt := grpc.UnaryInterceptor(acceptlang.UnaryServerInterceptor)
 sIntOpt := grpc.StreamInterceptor(acceptlang.StreamServerInterceptor)
 grpc.NewServer(uIntOpt, sIntOpt)
}

func foo(ctx context.Context) {
 acceptLangs := acceptlang.FromContext(ctx)
 fmt.printf("language :%s", acceptLangs[0].Language)
}
```

## i18n with Accept-Language

Also support i18n integration with Accept-Language. github.com/nicksnyder/go-i18n is supported for now.

When you send accept language via metadata, i18n interceptor parses it and set `i18n.TranslateFunc` to context. Then use `i18n.MustTFunc(ctx)` for translactions.

```golang
import (
 "github.com/nicksnyder/go-i18n/i18n"
 grpci18n "github.com/eltorocorp/go-grpc-request-id-interceptor/acceptlang/i18n"
 "golang.org/x/net/context"
)

func main() {
 // load translation files
 i18n.LoadTranslationFile("en-us.all.json")
 i18n.LoadTranslationFile("ja-jp.all.json")

 // set default language in case of no accept language specified
 // or no valid language found
 grpci18n.SetDefaultLanguage("en")

 // use i18 interceptor. Not explicitly required acceptlang interceptor
 uIntOpt := grpc.UnaryInterceptor(grpci18n.UnaryServerInterceptor)
 grpc.NewServer(uIntOpt)
}

func foo(ctx context.Context) {
 // get TranslateFunc from context
 T := grpci18n.MustTFunc(ctx)
 fmt.printf("%s", T("hello"))
}
```
