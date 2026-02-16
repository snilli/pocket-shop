package main

import (
	"context"
	"fmt"
	"orm"
	"time"

	swagger "github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"pocket-shop/config"
	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/repository/db/orderrepo"
	"pocket-shop/internal/core/order/repository/db/orderreservationrepo"
	"pocket-shop/internal/core/order/service/ordersvc"
	"pocket-shop/internal/delivery/http/handler/orderhdl"
	"pocket-shop/internal/infrastructure/database"
	"pocket-shop/internal/infrastructure/ezclient"
	"pocket-shop/internal/infrastructure/server"
)

func main() {
	fx.New(
		fx.Provide(provideConfig),
		fx.Provide(server.ProvideFiberApp),
		fx.Provide(database.ProvideDatabase),
		fx.Provide(
			fx.Annotate(
				func(client *orm.Client, cfg *config.Config) order.OrderRepository {
					return orderrepo.New(client)
				},
				fx.As(new(order.OrderRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(client *orm.Client, cfg *config.Config) order.OrderReservationRepository {
					return orderreservationrepo.New(client)
				},
				fx.As(new(order.OrderReservationRepository)),
			),
		),
		fx.Provide(ezclient.NewEZClient),
		fx.Provide(
			fx.Annotate(
				ordersvc.New,
				fx.As(new(order.OrderService)),
			),
		),
		fx.Provide(orderhdl.NewHandler),
		fx.Invoke(ensureSchema),
		fx.Invoke(registerDatabaseHooks),
		fx.Invoke(runDiscoverOrderJob),
		fx.Invoke(runServer),
	).Run()
}

func provideConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cfg.SetupLogger()
	log.Info().
		Str("mode", cfg.ServerMode).
		Str("address", cfg.ServerAddr()).
		Msg("Starting server")
	return cfg, nil
}

func ensureSchema(lc fx.Lifecycle, client *orm.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return database.EnsureSchema(ctx, client)
		},
	})
}

func runDiscoverOrderJob(lc fx.Lifecycle, svc order.OrderService, cfg *config.Config) {
	var cancel context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx, cancel = context.WithCancel(context.Background())
			interval := time.Duration(cfg.DiscoverIntervalSec) * time.Second
			go func() {
				ticker := time.NewTicker(interval)
				defer ticker.Stop()
				for {
					if err := svc.DiscoverAvailable(ctx, cfg.RefSource); err != nil {
						log.Warn().Err(err).Msg("discover available EZ orders")
					}
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})
}

func registerDatabaseHooks(lc fx.Lifecycle, client *orm.Client) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Closing database connection")
			return client.Close()
		},
	})
}

func runServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	app *fiber.App,
	orderHandler *orderhdl.Handler,
) {
	if cfg.EnableSwagger {
		app.Use(swagger.New(swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
			Title:    "Order API",
			CacheAge: 3600,
		}))
		log.Info().Str("url", fmt.Sprintf("http://%s/swagger", cfg.ServerAddr())).Msg("Swagger enabled")
	}

	apiV1 := app.Group("/api/v1")
	orderHandler.RegisterRoutes(apiV1)

	log.Info().
		Str("url", fmt.Sprintf("http://%s", cfg.ServerAddr())).
		Msg("Server is ready")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
					log.Error().Err(err).Msg("Server error")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down server...")
			return app.Shutdown()
		},
	})
}
