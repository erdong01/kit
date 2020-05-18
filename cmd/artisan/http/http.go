package http

import (
	"context"
	"os"
	"os/signal"
	"github.com/erDong01/micro-kit/cmd/artisan/http/errGroup"
	"github.com/erDong01/micro-kit/cmd/artisan/http/scHttp"
	"github.com/erDong01/micro-kit/cmd/artisan/http/studentHttp"
	"github.com/erDong01/micro-kit/cmd/artisan/http/tmrHttp"
	"github.com/erDong01/micro-kit/internal/log"
	"time"
)

func New() {
	scHttp := scHttp.Start()
	tmrHttp := tmrHttp.Start()
	studentHttp := studentHttp.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tmrHttp.Shutdown(ctx); err != nil {
		log.Fatal("Server02 Shutdown:", err)
	}
	if err := scHttp.Shutdown(ctx); err != nil {
		log.Fatal("Server01 Shutdown:", err)
	}
	if err := studentHttp.Shutdown(ctx); err != nil {
		log.Fatal("Server03 Shutdown:", err)
	}

	log.Info("Server exiting")
	if err := errGroup.ErrGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}
