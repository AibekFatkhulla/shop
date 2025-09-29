package main

import (
	"context"
	"fmt"

	"github.com/aibekfatkhulla/shop/config"
	"github.com/aibekfatkhulla/shop/internal/repository"
	"github.com/aibekfatkhulla/shop/internal/server"
	"github.com/aibekfatkhulla/shop/internal/service"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	ctx := context.Background()

	_ = godotenv.Load("local.env")
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	fmt.Printf("cfg: %+v\n", cfg)
	pg, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.Db))
	if err != nil {
		panic(err)
	}
	defer pg.Close(ctx)

	repo := repository.NewRepository(pg)
	svc := service.NewService(repo)
	srv := server.NewServer(svc)

	err = srv.Run(":8080")
	if err != nil {
		panic(err)
	}
}
