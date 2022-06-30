package requestdump

import (
	"go.uber.org/zap"
)

type Option interface {
	apply(*options)
}

type optionApplyer func(*options)

func (a optionApplyer) apply(opt *options) {
	a(opt)
}

type options struct {
	enabled bool
	logger  zap.Logger
	rootKey string
}

func newOptions() options {
	return options{
		enabled: true,
		rootKey: "request_dump",
	}
}

func Disable() Option {
	return optionApplyer(func(opt *options) {
		opt.enabled = false
	})
}

func Zap(logger zap.Logger) Option {
	return optionApplyer(func(opt *options) {
		opt.logger = logger
	})
}

func RootKey(key string) Option {
	return optionApplyer(func(opt *options) {
		opt.rootKey = key
	})
}
