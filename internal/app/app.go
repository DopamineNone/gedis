package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/wire"
)

var (
	ProvideSet = wire.NewSet(New)
	Timeout    = 5 * time.Second
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

type App struct {
	Handler
	*sync.WaitGroup
	listener      net.Listener
	rootCtx       context.Context
	ctxCancelFunc context.CancelFunc
	shutdownCh    chan struct{}
	cleanup       sync.Once
}

func (a *App) Run() error {
	defer a.Stop()
	for {
		// get connection
		conn, err := a.listener.Accept()
		if err != nil {
			return err
		}

		// service handle
		a.Add(1)
		go func() {
			defer a.Done()
			a.Handle(a.rootCtx, conn)
		}()
	}
}

func (a *App) Stop() {
	a.cleanup.Do(func() {
		// stop incoming task
		a.listener.Close()
		a.Close()
		a.ctxCancelFunc()

		// wait working goroutines with timeout
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		go func() {
			a.Wait()
			cancel()
		}()
		<-ctx.Done()

		close(a.shutdownCh)
	})
}

func (a *App) ListenAndQuit() {
	// listen system quit signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		select {
		case <-signalChan:
			a.Stop()
		case <-a.shutdownCh: // program has been stoped by another goroutine
		}
	}()

	<-a.shutdownCh
}

func New(handler Handler, listener net.Listener) *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		Handler:       handler,
		WaitGroup:     new(sync.WaitGroup),
		listener:      listener,
		rootCtx:       ctx,
		ctxCancelFunc: cancel,
		shutdownCh:    make(chan struct{}),
	}
}
