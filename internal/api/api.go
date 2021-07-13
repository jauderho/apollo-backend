package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/christianselig/apollo-backend/internal/data"
	"github.com/christianselig/apollo-backend/internal/reddit"
)

type api struct {
	logger *logrus.Logger
	statsd *statsd.Client
	db     *pgxpool.Pool
	reddit *reddit.Client
	models *data.Models
}

func NewAPI(ctx context.Context, logger *logrus.Logger, statsd *statsd.Client, db *pgxpool.Pool) *api {
	reddit := reddit.NewClient(
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_CLIENT_SECRET"),
		statsd,
	)

	models := data.NewModels(ctx, db)

	return &api{logger, statsd, db, reddit, models}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/v1/health", a.healthCheckHandler)

	router.POST("/v1/device", a.upsertDeviceHandler)
	router.POST("/v1/device/:apns/account", a.upsertAccountHandler)

	return router
}