package server

import (
	"runtime/debug"

	"github.com/bytedance/sonic"
	fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	recovermw "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"pocket-shop/config"
)

// @title Order API
// @version 1.0
// @description API documentation for Order (instant gift card) service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// ProvideFiberApp creates and configures the Fiber application.
func ProvideFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Order API",
		ServerHeader: "Fiber",
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if code >= 500 {
				ev := log.Error().Err(err).Str("method", c.Method()).Str("path", c.Path())
				if rid := c.Locals("requestid"); rid != nil {
					if s, ok := rid.(string); ok {
						ev = ev.Str("request_id", s)
					}
				}
				ev.Msg("request failed")
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "internal_error",
				"message": err.Error(),
			})
		},
	})

	logger := log.Logger
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
		Fields: []string{
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldURL,
			fiberzerolog.FieldError,
			fiberzerolog.FieldRequestID,
			fiberzerolog.FieldIP,
		},
		Messages: []string{
			"Server error",
			"Client error",
			"Success",
		},
		Levels: []zerolog.Level{
			zerolog.ErrorLevel,
			zerolog.WarnLevel,
			zerolog.InfoLevel,
		},
		Next: func(c fiber.Ctx) bool {
			path := c.Path()
			return path == "/livez"
		},
	}))

	app.Use(recovermw.New(recovermw.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			stack := debug.Stack()
			log.Error().
				Interface("panic", e).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("stack", string(stack)).
				Msg("panic recovered")
		},
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	app.Get("/livez", healthcheck.New())

	return app
}
