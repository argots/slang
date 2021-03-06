package eval

var _ Value = &lazy{}

// Lazy lazily evaluates a function when someone attempts to fetch the
// value.  It caches the value once it has been calculated.
func Lazy(fn func() Valuable) Valuable {
	return &lazy{fn: fn}
}

type lazy struct {
	v          Value
	fn         func() Valuable
	inProgress bool
}

func (l *lazy) Value() Value {
	if l.inProgress {
		return NewError(NewString("recursion"))
	}
	if l.v == nil {
		l.inProgress = true
		defer func() {
			l.inProgress = false
		}()
		l.v = l.fn().Value()
	}
	return l.v
}

func (l *lazy) Type() string {
	return l.Value().Type()
}

func (l *lazy) Get(v Valuable) Valuable {
	return l.Value().Get(v)
}

func (l *lazy) Code() Code {
	return l.Value().Code()
}
