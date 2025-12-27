package pipe

type Pipe[T any] struct {
	id uint32
}

func New[T any](opts ...Option) *Pipe[T] {
	p := &Pipe[T]{}
	for _, fn := range opts {
		fn(p)
	}
	return p
}

func (p *Pipe[T]) Send(in T) (err error) {
	id = p.id
	if err = setData(in); err != nil {
		return
	}
	send()
	return
}

func (p *Pipe[T]) Recv() (out T, err error) {
	id = p.id
	recv()
	err = getData(&out)
	return
}

func (p *Pipe[T]) setID(id uint32) {
	p.id = id
}

type Option func(t target)
type target interface {
	setID(id uint32)
}

func WithID(id uint32) Option {
	return func(t target) {
		t.setID(id)
	}
}
