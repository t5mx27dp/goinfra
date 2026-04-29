package graceful

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

type testApp struct {
	onSetUp    func(ctx context.Context, makeGo MakeGoroutine) error
	onTearDown func()
}

func (app *testApp) SetUp(ctx context.Context, makeGo MakeGoroutine) error {
	if app.onSetUp != nil {
		return app.onSetUp(ctx, makeGo)
	}
	return nil
}

func (app *testApp) TearDown() {
	if app.onTearDown != nil {
		app.onTearDown()
	}
}

func testSendStopSignal(controller *Controller) {
	i := rand.Intn(len(controller.signals))
	controller.signal <- controller.signals[i]
}

func TestController(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		expected := 2
		actual := 0

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			actual++
			return nil
		}

		app.onTearDown = func() {
			actual++
		}

		controller := NewController(app)

		err := controller.Run()
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	t.Run("Run_MakeGoroutine", func(t *testing.T) {
		expected := 2
		actual := 0

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			makeGo(func() {
				select {
				case <-ctx.Done():
					actual++
				}
			})

			return nil
		}

		app.onTearDown = func() {
			actual++
		}

		controller := NewController(app)

		err := controller.Run()
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	t.Run("Run_Listen", func(t *testing.T) {
		expected := 2
		actual := 0

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			actual++
			return nil
		}

		app.onTearDown = func() {
			actual++
		}

		controller := NewController(app)

		controller.Listen([]os.Signal{syscall.SIGTERM})

		go testSendStopSignal(controller)

		err := controller.Run()
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	t.Run("Run_WithContext", func(t *testing.T) {
		expected := 2
		actual := 0

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			actual++
			return nil
		}

		app.onTearDown = func() {
			actual++
		}

		controller := NewController(app)

		controller.WithContext(ctx, cancel)

		controller.Listen([]os.Signal{syscall.SIGTERM})

		go testSendStopSignal(controller)

		err := controller.Run()
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	t.Run("Run_WithContext_Cancel", func(t *testing.T) {
		expected := 2
		actual := 0

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			actual++
			return nil
		}

		app.onTearDown = func() {
			actual++
		}

		controller := NewController(app)

		controller.WithContext(ctx, cancel)

		controller.Listen([]os.Signal{syscall.SIGTERM})

		cancel()

		err := controller.Run()
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	t.Run("Run_SetUpError", func(t *testing.T) {
		expected := errors.New("err")

		app := &testApp{}

		app.onSetUp = func(ctx context.Context, makeGo MakeGoroutine) error {
			return expected
		}

		controller := NewController(app)

		err := controller.Run()
		require.NotNil(t, err)

		require.Equal(t, expected, err)
	})
}
