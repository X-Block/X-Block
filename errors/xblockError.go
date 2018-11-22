package errors

type xblockError struct {
	errmsg string
	callstack *CallStack
	root error
	code ErrCode
}

func (e xblockError) Error() string {
	return e.errmsg
}

