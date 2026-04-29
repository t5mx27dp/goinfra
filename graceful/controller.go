package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
)

var ErrControllerIsRunning = errors.New("controller is running")
var ErrAlreadySetSignals = errors.New("already set signals")

type MakeGoroutine func(fn func())

type Controller struct {
	app App

	ctx    context.Context
	cancel context.CancelFunc

	wg *sync.WaitGroup

	signal  chan os.Signal
	signals []os.Signal

	mutex   sync.Mutex
	running bool
}

func NewController(app App) *Controller {
	return &Controller{
		app: app,
	}
}

// WithContext 用于设置 Context。
func (c *Controller) WithContext(ctx context.Context, cancel context.CancelFunc) {
	c.ctx = ctx
	c.cancel = cancel
}

// Listen 用于添加要监听的信号，例如 []os.Signal{syscall.SIGINT, syscall.SIGTERM}。
func (c *Controller) Listen(signals []os.Signal) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.signal != nil {
		return ErrAlreadySetSignals
	}

	c.signals = signals
	c.signal = make(chan os.Signal, 1)

	return nil
}

func (c *Controller) Run() error {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return ErrControllerIsRunning
	}
	c.running = true
	c.mutex.Unlock()

	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.cancel == nil {
		c.ctx, c.cancel = context.WithCancel(c.ctx)
	}

	c.wg = &sync.WaitGroup{}

	defer func() {
		c.cancel()
		c.wg.Wait()
		c.app.TearDown()
	}()

	err := c.app.SetUp(c.ctx, c.MakeGoroutine)
	if err != nil {
		return err
	}

	if c.signal != nil {
		signal.Notify(c.signal, c.signals...)
		defer signal.Stop(c.signal)

		select {
		case <-c.ctx.Done():
		case <-c.signal:
		}
	}

	return nil
}

func (c *Controller) MakeGoroutine(fn func()) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		fn()
	}()
}
