package exn

// Wrap wraps err, adding the string context to its message. If
// err is nil, returns nil instead.
func Wrap(context string, err error) error {
	if err != nil {
		err = wrappedErr{
			context: context,
			err:     err,
		}
	}
	return err
}

// WrapThrow(th, context, err) is equivalent to th(Wrap(context, err)),
// i.e. if err != nil, it throws with an error wrapping err and supplying
// the additional context.
func WrapThrow(th Thrower, context string, err error) {
	th(Wrap(context, err))
}

// wrapper error used by Wrap.
type wrappedErr struct {
	context string
	err     error
}

func (e wrappedErr) Error() string {
	return e.context + ": " + e.err.Error()
}

func (e wrappedErr) Unwrap() error {
	return e.err
}
