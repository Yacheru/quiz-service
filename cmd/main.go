package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"

	"quiz-service/init/config"
	"quiz-service/init/logger"
	"quiz-service/internal/server"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cfg := &config.ServerConfig

	if err := config.InitConfig(); err != nil {
		fmt.Println(err.Error())
		cancel()
	}

	httpLogger, err := logger.NewLogger(ctx, cfg.Debug, cfg.Log.HttpLoggerPath)
	if err != nil {
		cancel()
	}
	postgresLogger, err := logger.NewLogger(ctx, cfg.Debug, cfg.Log.PostgresLoggerPath)
	if err != nil {
		cancel()
	}
	quizLogger, err := logger.NewLogger(ctx, cfg.Debug, cfg.Log.QuizLoggerPath)
	if err != nil {
		cancel()
	}

	app, err := server.NewHTTPServer(ctx, cfg, httpLogger, postgresLogger, quizLogger)
	if err != nil {
		cancel()
	}
	quizLogger.Info("server configured")

	if app != nil {
		errs, gCtx := errgroup.WithContext(ctx)
		errs.Go(func() error {
			return app.Run()
		})

		errs.Go(func() error {
			<-gCtx.Done()
			return app.Shutdown(gCtx)
		})

		if err := errs.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel()
		}
	}

	<-ctx.Done()

	quizLogger.Info("quiz shutdown")
}
