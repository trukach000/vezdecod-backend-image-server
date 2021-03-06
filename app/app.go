package app

import (
	"backend-image-server/pkg/config"
	"backend-image-server/pkg/database"
	"backend-image-server/pkg/httpext"
	"backend-image-server/pkg/redisclient"
	"backend-image-server/pkg/swagger"
	"net/http"
	"os"

	chilogrus "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

func Setup() *chi.Mux {

	r := chi.NewRouter()
	log := logrus.New()

	cfg := config.Get()

	db, err := database.InitDatabase(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	if err != nil {
		logrus.Errorf("can't connect to database: %s", err)
		os.Exit(1)
	}

	redisClient, err := redisclient.InitRedisClient(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPass,
		cfg.RedisDatabaseName,
		"30s",
	)

	if err != nil {
		logrus.Errorf("can't connect to redis: %s", err)
		os.Exit(1)
	}

	r.Use(
		chilogrus.Logger("logger", log),
		chim.Recoverer,
		chim.NoCache,
		database.NewDatabaseMiddleware(db).Attach,
		redisclient.NewRedisMiddleware(redisClient).Attach,
	)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", swagger.WrapSwagger)

	// frontend files
	httpext.ServeFile(r, "/site", "./site/index.html")
	httpext.ServeDir(r, "/site/*", http.Dir("./site"))

	// Upload and get API
	r.Post("/upload", UploadImage)
	r.Get("/get/{id}", GetImage)

	//Compare API
	r.Post("/compare", CompareImage)

	return r
}
