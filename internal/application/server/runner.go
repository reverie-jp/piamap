package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/reverie-jp/piamap/internal/config"
	"github.com/reverie-jp/piamap/internal/platform/jwt"
	"github.com/reverie-jp/piamap/internal/platform/logger"
)

func Run() error {
	cfg := config.New()
	if err := cfg.LoadFromEnv(); err != nil {
		return fmt.Errorf("failed to load config from env: %w", err)
	}
	logger.Init(cfg)

	ctx := context.Background()
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to parse database DSN: %w", err)
	}
	poolCfg.MaxConns = cfg.Database.MaxConns
	poolCfg.MinConns = cfg.Database.MinConns
	poolCfg.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.Database.MaxConnIdleTime

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	defer db.Close()

	jwtManager := jwt.NewManager(cfg.Auth.JWTSecretKey, cfg.Auth.AccessExpiration, cfg.Auth.RefreshExpiration)
	services := initServices(cfg, db, jwtManager)

	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	for _, svc := range services {
		svc.RegisterConnectHandler(mux)
	}

	if cfg.Env == config.EnvDevelopment {
		serviceNames := make([]string, 0, len(services))
		for _, svc := range services {
			serviceNames = append(serviceNames, svc.Name)
		}
		mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(serviceNames...)))
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders: []string{
			"Connect-Protocol-Version",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Grpc-Status",
			"Grpc-Message",
		},
	})

	handler := c.Handler(h2c.NewHandler(mux, &http2.Server{
		ReadIdleTimeout: 30 * time.Second,
		PingTimeout:     15 * time.Second,
	}))

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	slog.Info("api server is running",
		slog.String("addr", srv.Addr),
		slog.String("env", string(cfg.Env)),
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	slog.Info("shutting down server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	slog.Info("server exited properly")
	return nil
}
