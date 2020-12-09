package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/go-pkgz/lgr"
	R "github.com/go-pkgz/rest"
	validation "github.com/go-playground/validator/v10"
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

	tagHandlers          tagHandlers
	dayHandlers          dayHandlers
	eventHandlers        eventHandlers
	projectHandlers      projectHandlers
	adminProjectHandlers adminProjectHandlers
	userHandlers         userHandlers
}

type uuidProvider struct {
}

func (p *uuidProvider) New() string {
	return uuid.New().String()
}

func CreateRestServer(database *mongo.Database, keycloakUrl string) *RestServer {
	uuidProvider := &uuidProvider{}
	dayStore := repository.CreateDayStore(database, uuidProvider)
	projectStore := repository.CreateProjectStore(database, uuidProvider)
	validator := validation.New()

	return &RestServer{
		authenticator:        CreateAuthService(keycloakUrl),
		tagHandlers:          tagHandlers{store: repository.CreateTagStore(database, uuidProvider), validator: validator},
		dayHandlers:          dayHandlers{store: dayStore},
		eventHandlers:        eventHandlers{service: service.CreateEventService(dayStore, uuidProvider), validator: validator},
		projectHandlers:      projectHandlers{store: projectStore},
		adminProjectHandlers: adminProjectHandlers{store: projectStore, validator: validator},
		userHandlers:         userHandlers{service: service.CreateUserService(keycloakUrl)},
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

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(30 * time.Second))
		r.Use(s.authenticator.HttpMiddleware)

		r.Route("/tags", func(routeUser chi.Router) {
			routeUser.Get("/", s.tagHandlers.findAll)
			routeUser.Post("/", s.tagHandlers.create)
			routeUser.Delete("/{id}", s.tagHandlers.delete)
		})

		r.Route("/days", func(routeUser chi.Router) {
			routeUser.Get("/", s.dayHandlers.findByDateRange)
			routeUser.Get("/{date}", s.dayHandlers.findByDate)
		})

		r.Route("/events", func(routeUser chi.Router) {
			routeUser.Post("/", s.eventHandlers.create)
			routeUser.Put("/{id}", s.eventHandlers.update)
			routeUser.Delete("/{id}", s.eventHandlers.delete)
		})

		r.Get("/projects", s.projectHandlers.findAllByUser)

		// Admin
		r.Route("/admin", func(routeAdmin chi.Router) {
			// TODO add authentication middleware to allow admin roles only

			routeAdmin.Route("/projects", func(routeAdminProject chi.Router) {
				routeAdminProject.Post("/", s.adminProjectHandlers.create)
				routeAdminProject.Get("/", s.adminProjectHandlers.findAll)
				routeAdminProject.Put("/{id}", s.adminProjectHandlers.update)
				routeAdminProject.Delete("/{id}", s.adminProjectHandlers.delete)
			})

			routeAdmin.Get("/users", s.userHandlers.findAll)
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
