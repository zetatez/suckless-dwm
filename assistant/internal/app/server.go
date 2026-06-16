package app

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"assistant/internal/app/module"
	"assistant/internal/app/modules/health"
	"assistant/internal/app/modules/svc"
	"assistant/internal/bootstrap/psl"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context) error {
	logger := psl.GetLogger()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		logger.WithFields(map[string]interface{}{
			"error":  err,
			"stack":  string(debug.Stack()),
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Error("panic recovered")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	modules := []module.Module{
		health.NewHealthModule(),
		svc.NewModule(),
	}

	api := r.Group("/api")
	{
		for _, m := range modules {
			logger.WithFields(map[string]interface{}{"module": m.Name(), "prefix": "/api/" + m.Name()}).Info("registering module")
			group := api.Group("/" + m.Name())
			moduleMiddleware := m.Middleware()
			if len(moduleMiddleware) > 0 {
				group.Use(moduleMiddleware...)
			}
			m.Register(group)
		}
	}

	cfg := psl.GetConfig().App
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.WithFields(map[string]interface{}{"address": addr}).Info("server running")

	srv := &http.Server{Addr: addr, Handler: r}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.WithFields(map[string]interface{}{"reason": ctx.Err().Error()}).Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		logger.WithFields(map[string]interface{}{"reason": err.Error()}).Info("server error received")
		return err
	}
}
