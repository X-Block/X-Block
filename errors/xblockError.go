package errors

type xblockError struct {
	errmsg string
	callstack *CallStack
	root error
	code ErrCode
}

