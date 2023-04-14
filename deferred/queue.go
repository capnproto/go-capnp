package deferred

type Queue []func()

func (q *Queue) Run() {
	funcs := *q
	for i, f := range funcs {
		if f != nil {
			f()
			funcs[i] = nil
		}
	}
}

func (q *Queue) Defer(f func()) {
	*q = append(*q, f)
}
