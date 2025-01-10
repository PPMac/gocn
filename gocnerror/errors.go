package gocnerror

type GocnError struct {
	err    error
	ErrFuc ErrorFuc
}

func Default() *GocnError {
	return &GocnError{}
}

func (e *GocnError) Error() string {
	return e.err.Error()
}

func (e *GocnError) Put(err error) {
	e.check(err)
}

func (e *GocnError) check(err error) {
	if err != nil {
		e.err = err
		panic(e)
	}
}

type ErrorFuc func(msError *GocnError)

//暴露一个方法 让用户自定义

func (e *GocnError) Result(errFuc ErrorFuc) {
	e.ErrFuc = errFuc
}
func (e *GocnError) ExecResult() {
	e.ErrFuc(e)
}
