package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/init/config"
	"quiz-service/init/logger"
	"quiz-service/internal/repository/postgres"
	"quiz-service/internal/server/http/router"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(ctx context.Context, cfg *config.Config, httpLogger, dbLogger, quizLogger *logger.Logger) (*HTTPServer, error) {
	db, err := postgres.InitPostgresConnection(ctx, cfg, quizLogger)
	if err != nil {
		quizLogger.Error("postgres:" + err.Error())
		return nil, err
	}

	engine := setupGin(cfg.Debug)
	entry := engine.Group(cfg.Entry)
	router.InitRouterAndComponents(entry, db, cfg, httpLogger, dbLogger).Routes()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        engine,
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{server: server}, nil
}

func (s *HTTPServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func setupGin(debug bool) *gin.Engine {
	var mode = gin.ReleaseMode
	if debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logger.HTTPLogger))

	engine.LoadHTMLFiles("./web/register.html", "./web/variants.html", "./web/test.html", "./web/results.html")

	return engine
}
