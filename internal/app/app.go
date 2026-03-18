package app

import (
	"context"
	"time"

	"github.com/VAGRAMCHIC/vless_reality_agent/internal/config"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/handler"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/repository"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/service"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() {
	log, err := utils.New("/var/log/ggvpn/app.log")
	if err != nil {
		log.Fatal(context.TODO(), "failed_to_init_logger", map[string]any{
			"error": err.Error(),
		})
	}
	log.Info(context.TODO(), "start_application", map[string]any{
		"curr_date": time.Now(),
	})

	cfg := config.Load()
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)

	userRepo := repository.NewUserRepository(db, log)

	userService := service.NewUserService(userRepo, log, cfg.Server, cfg.TAG)

	userH := handler.NewUserHandler(userService)

	handler := handler.NewHandler(userH, log)
	r := gin.Default()
	handler.Register(r, cfg.APIKey)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(context.TODO(), "failed_to_run_app", map[string]any{
			"error": err.Error(),
		})
	}

}
