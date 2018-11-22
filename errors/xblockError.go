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

func (e xblockError) GetErrCode()  ErrCode {
	return e.code
}

func (e xblockError) GetRoot()  error {
	return e.root
}

func (e xblockError) GetCallStack()  *CallStack {
	return e.callstack
}
