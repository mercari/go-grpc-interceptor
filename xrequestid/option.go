package xrequestid

type Option interface {
	apply(*options)
}

type optionApplyer func(*options)

func (a optionApplyer) apply(opt *options) {
	a(opt)
}

type options struct {
	chainRequestID   bool
	persistRequestID bool // if true, attach request id to outgoing context
	validator        requestIDValidator
}

func ChainRequestID() Option {
	return optionApplyer(func(opt *options) {
		opt.chainRequestID = true
	})
}

func PersistRequestID() Option {
	return optionApplyer(func(opt *options) {
		opt.persistRequestID = true
	})
}

type requestIDValidator func(string) bool

// RequestIDValidator is validator function that returns true if
// request id is valid, or false if invalid.
func RequestIDValidator(validator requestIDValidator) Option {
	return optionApplyer(func(opt *options) {
		opt.validator = validator
	})
}

func defaultReqeustIDValidator(requestID string) bool {
	return true
}
