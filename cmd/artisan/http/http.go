package http

import (
	"context"
	"os"
	"os/signal"
	"github.com/erDong01/micro-kit/cmd/artisan/http/errGroup"
	"github.com/erDong01/micro-kit/cmd/artisan/http/studentHttp"
	"time"
)

func New() {
	studentHttp := studentHttp.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := studentHttp.Shutdown(ctx); err != nil {
	}

	if err := errGroup.ErrGroup.Wait(); err != nil {
	}
}
