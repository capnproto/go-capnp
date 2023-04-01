package rc

type Releaser interface {
	Release()
}

type Refcounted[T any] interface {
	Releaser
	AddRef() T
}
