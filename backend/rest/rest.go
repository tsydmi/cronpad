package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/go-pkgz/lgr"
	R "github.com/go-pkgz/rest"
	"github.com/google/uuid"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/service"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type RestServer struct {
	authenticator *AuthService
	httpServer    *http.Server

	tagHandlers   tagHandlers
	dayHandlers   dayHandlers
	eventHandlers eventHandlers
}

type uuidProvider struct {
}

func (p *uuidProvider) New() string {
	return uuid.New().String()
}

func CreateRestServer(database *mongo.Database) *RestServer {
	uuidProvider := &uuidProvider{}
	dayStore := repository.CreateDayStore(database, uuidProvider)

	return &RestServer{
		authenticator: CreateAuthService(),
		tagHandlers:   tagHandlers{store: repository.CreateTagStore(database, uuidProvider)},
		dayHandlers:   dayHandlers{store: dayStore},
		eventHandlers: eventHandlers{service: service.CreateEventService(dayStore, uuidProvider)},
	}
}

func (s *RestServer) Run() {
	useSsl := false

	if useSsl {
		log.Fatalf("Application does not support HTTPS.")
	} else {
		s.httpServer = s.makeHTTPServer(9000, s.routes())
		s.httpServer.ListenAndServe()
	}
}

func (s *RestServer) makeHTTPServer(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
}

func (s *RestServer) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000), middleware.RealIP, R.Recoverer(log.Default()))

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)

	router.Route("/api/v1", func(rapi chi.Router) {
		rapi.Route("/tags", func(radmin chi.Router) {
			radmin.Use(middleware.Timeout(30 * time.Second))
			radmin.Use(s.authenticator.HttpMiddleware)

			radmin.Get("/", s.tagHandlers.findAll)
			radmin.Post("/", s.tagHandlers.create)
			radmin.Delete("/{id}", s.tagHandlers.delete)
		})

		rapi.Route("/days", func(radmin chi.Router) {
			radmin.Use(middleware.Timeout(30 * time.Second))
			radmin.Use(s.authenticator.HttpMiddleware)

			radmin.Get("/", s.dayHandlers.findByDateRange)
			radmin.Get("/{date}", s.dayHandlers.findByDate)
		})

		rapi.Route("/events", func(radmin chi.Router) {
			radmin.Use(middleware.Timeout(30 * time.Second))
			radmin.Use(s.authenticator.HttpMiddleware)

			radmin.Post("/", s.eventHandlers.create)
			radmin.Put("/{id}", s.eventHandlers.update)
			radmin.Delete("/{id}", s.eventHandlers.delete)
		})
	})

	return router
}

func (s *RestServer) Shutdown() {
	log.Print("[WARN] shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[DEBUG] http shutdown error, %s", err)
		}
		log.Print("[DEBUG] shutdown http server completed")
	}
	//s.lock.Unlock()
}
