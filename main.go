package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/walnuts1018/s3-oauth2-proxy/router/handler"
	"github.com/walnuts1018/s3-oauth2-proxy/router"
	appConfig "github.com/walnuts1018/s3-oauth2-proxy/config"
	"github.com/walnuts1018/s3-oauth2-proxy/infrastructure/auth"
	"github.com/walnuts1018/s3-oauth2-proxy/infrastructure/s3"
	"github.com/walnuts1018/s3-oauth2-proxy/logger"
	"github.com/walnuts1018/s3-oauth2-proxy/tracer"
	"github.com/walnuts1018/s3-oauth2-proxy/usecase"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	cfg, err := appConfig.Load()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := logger.CreateLogger(cfg.LogLevel, cfg.LogType)
	slog.SetDefault(logger)

	closeTracer, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create tracer provider", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer closeTracer()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load aws config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

	authRepo, err := auth.NewAuthRepository(cfg.OIDC)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create auth repository", slog.String("error", err.Error()))
		os.Exit(1)
	}

	s3Repo := s3.NewS3Repository(awsCfg, cfg.S3.Bucket)

	authUsecase := usecase.NewAuthUsecase(authRepo)
	proxyUsecase := usecase.NewProxyUsecase(s3Repo)

	authHandler := handler.NewAuthHandler(authUsecase)
	proxyHandler := handler.NewProxyHandler(proxyUsecase)

	e := router.NewRouter(&cfg.App, authHandler, proxyHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: e,
	}

	go func() {
		slog.Info("Server is running", slog.String("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	stop()
	slog.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
