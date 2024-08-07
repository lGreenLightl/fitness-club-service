package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lGreenLightl/fitness-club-service/internal/app/auth"
	logger "github.com/lGreenLightl/fitness-club-service/internal/app/logs/logrus"

	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func RunHTTPServer(createHandler func(chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	masterRouter := chi.NewRouter()
	masterRouter.Mount("/api", createHandler(apiRouter))

	logrus.Info("Starting HTTP server")
	http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), masterRouter)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		logger.NewStructuredLogger(logrus.StandardLogger()),
		middleware.Recoverer,
	)

	addAuthMiddleware(router)
	addCorsMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "no-sniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
		middleware.NoCache,
	)
}

func addAuthMiddleware(router *chi.Mux) {
	if ok, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); ok {
		router.Use(auth.HTTPMockMiddleware)
		return
	}

	var options []option.ClientOption

	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
		options = append(options, option.WithCredentialsFile("/service-account-file.json"))
	}

	firebaseApp, err := firebase.NewApp(
		context.Background(),
		&firebase.Config{ProjectID: os.Getenv("GCP")},
		options...,
	)
	if err != nil {
		logrus.Fatalf("Initializing app error: %v\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("Firebase auth client didn't create")
	}

	router.Use(auth.FirebaseHttpMiddleware{AuthClient: authClient}.HTTPMiddleware)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedOrigins:   allowedOrigins,
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	})

	router.Use(corsMiddleware.Handler)
}
