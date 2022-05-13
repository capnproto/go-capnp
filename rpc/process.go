package rpc

import (
	"context"
	"sync"
	"time"

	"github.com/jbenet/goprocess"
	syncutil "github.com/lthibault/util/sync"
)

type process struct {
	root  goprocess.Process
	refs  refctr
	abort chan error
	err   error
}

func newProcess(c *Conn) *process {
	p := new(process)
	p.abort = make(chan error, 1)
	p.root = goprocess.Go(func(proc goprocess.Process) {
		proc.SetTeardown(p.teardown(c))
		proc.
			Go(p.sender(c)).
			Go(p.recver(c)).
			SetTeardown(p.signal(c))

		select {
		case p.err = <-p.abort:
		case <-proc.Closing():
		}
	})

	return p
}

func (p *process) sender(c *Conn) goprocess.ProcessFunc {
	return c.sender
}

func (p *process) recver(c *Conn) goprocess.ProcessFunc {
	return func(proc goprocess.Process) {
		err := c.recver(proc)
		if err == context.Canceled {
			err = nil
		}

		select {
		case p.abort <- err:
		default:
		}
	}
}

// Signal to answers that we are shutting down. They are expected to
// promptly release any references to the process that they may hold.
//
// Signal is called after the send and receive processes have returned.
// Signal returns before teardown() is called.
func (p *process) signal(c *Conn) goprocess.TeardownFunc {
	return func() error {
		syncutil.With(&c.mu, func() {
			for _, a := range c.answers {
				if a != nil && a.cancel != nil {
					a.cancel()
				}
			}
		})

		p.refs.Wait()
		return nil
	}
}

// Teardown is called after all child processes have terminated.
func (p *process) teardown(c *Conn) goprocess.TeardownFunc {
	return func() (err error) {
		c.drainQueue()

		c.mu.Lock()
		defer c.mu.Unlock()

		// Clear all tables, releasing exported clients and unfinished answers.
		exports := c.exports
		embargoes := c.embargoes
		answers := c.answers
		c.imports = nil
		c.exports = nil
		c.embargoes = nil
		c.questions = nil
		c.answers = nil

		syncutil.Without(&c.mu, func() {
			c.releaseBootstrap()
			c.releaseExports(exports)
			c.liftEmbargoes(embargoes)
			c.releaseAnswers(answers)
			c.abort(p.err)
		})

		if err = c.transport.Close(); err != nil {
			err = rpcerr.Failedf("close transport: %w", err)
		}

		return
	}
}

func (p *process) Shutdown(err error) error {
	select {
	case p.abort <- err:
	default:
	}

	return p.root.Close()
}

func (p *process) Closing() <-chan struct{} {
	return p.root.Closing()
}

func (p *process) Closed() <-chan struct{} {
	return p.root.Closed()
}

func (p *process) HandleCancel(ctx context.Context, q *question) {
	p.root.Go(func(proc goprocess.Process) {
		q.handleCancel(ctx)
	})
}

func (p *process) Go(f goprocess.ProcessFunc) { p.root.Go(f) }

// AddRef increases the refcount on the process, preventing it from
// shutting down.  Returns false if the process is already shutting
// down, indicating that the reference could not be acquired.
//
// Callers MUST hold c.mu.
func (p *process) AddRef() bool {
	select {
	case <-p.root.Closing():
		return false
	default:
		p.refs.Acquire()
		return true
	}
}

func (p *process) Release() { p.refs.Release() }

type refctr sync.WaitGroup

func (r *refctr) Acquire() { (*sync.WaitGroup)(r).Add(1) }
func (r *refctr) Release() { (*sync.WaitGroup)(r).Done() }
func (r *refctr) Wait()    { (*sync.WaitGroup)(r).Wait() }

type procCtx <-chan struct{}

func (p procCtx) Done() <-chan struct{}                 { return p }
func (procCtx) Value(key any) any                       { return nil }
func (procCtx) Deadline() (deadline time.Time, ok bool) { return }

func (p procCtx) Err() error {
	select {
	case <-p:
		return context.Canceled
	default:
		return nil
	}
}
