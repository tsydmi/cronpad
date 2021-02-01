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
	"github.com/ts-dmitry/cronpad/backend/service/report"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type RestServer struct {
	authenticator *AuthService
	httpServer    *http.Server

	tagHandlers        tagHandlers
	projectTagHandlers projectTagHandlers
	baseTagHandlers    baseTagHandlers

	dayHandlers          dayHandlers
	eventHandlers        eventHandlers
	projectHandlers      projectHandlers
	adminProjectHandlers adminProjectHandlers
	userHandlers         userHandlers
	reportsHandlers      reportsHandlers
}

type uuidProvider struct {
}

func (p *uuidProvider) New() string {
	return uuid.New().String()
}

func CreateRestServer(database *mongo.Database, jwtAuth *JwtAuthService, keycloakUrl string) *RestServer {
	uuidProvider := &uuidProvider{}
	dayStore := repository.CreateDayStore(database, uuidProvider)
	tagStore := repository.CreateTagStore(database, uuidProvider)
	projectStore := repository.CreateProjectStore(database, uuidProvider)
	userService := service.CreateUserService(keycloakUrl, projectStore)
	validator := CreateValidator()

	return &RestServer{
		authenticator: &AuthService{jwtAuthService: jwtAuth},

		tagHandlers:        tagHandlers{tagStore: tagStore, projectStore: projectStore, validator: validator},
		projectTagHandlers: projectTagHandlers{tagStore: tagStore, projectStore: projectStore, validator: validator},
		baseTagHandlers:    baseTagHandlers{store: tagStore, validator: validator},

		dayHandlers:          dayHandlers{store: dayStore},
		eventHandlers:        eventHandlers{service: service.CreateEventService(dayStore, uuidProvider), validator: validator},
		projectHandlers:      projectHandlers{store: projectStore, userService: userService},
		adminProjectHandlers: adminProjectHandlers{store: projectStore, validator: validator, userService: userService},
		userHandlers:         userHandlers{service: userService},
		reportsHandlers:      reportsHandlers{service: report.CreateReportService(dayStore, tagStore, projectStore), projectStore: projectStore},
	}
}

func (s *RestServer) Run() {
	useSsl := false

	if useSsl {
		log.Fatalf("Application does not support HTTPS.")
	} else {
		s.httpServer = s.makeHTTPServer(9000, s.routes())
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf(err.Error())
		}
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

		// Project Manager
		r.Route("/manager", func(routeManager chi.Router) {
			routeManager.Use(s.authenticator.HasRole(projectManagerRole))

			routeManager.Get("/project-reports/{id}", s.reportsHandlers.projectReport) //TODO check if user assigned to the project here!
			routeManager.Get("/projects/{id}/users", s.projectHandlers.users)

			routeManager.Post("/tags", s.projectTagHandlers.create)
			routeManager.Put("/tags/{id}", s.projectTagHandlers.update)
			routeManager.Delete("/tags/{id}", s.projectTagHandlers.delete)
		})

		// Admin
		r.Route("/admin", func(routeAdmin chi.Router) {
			routeAdmin.Use(s.authenticator.HasRole(adminRole))

			routeAdmin.Route("/projects", func(routeAdminProject chi.Router) {
				routeAdminProject.Post("/", s.adminProjectHandlers.create)
				routeAdminProject.Post("/search", s.adminProjectHandlers.search)
				routeAdminProject.Put("/{id}", s.adminProjectHandlers.update)
				routeAdminProject.Delete("/{id}", s.adminProjectHandlers.delete)
			})

			routeAdmin.Route("/base-tags", func(routeUser chi.Router) {
				routeUser.Post("/", s.baseTagHandlers.create)
				routeUser.Put("/{id}", s.baseTagHandlers.update)
				routeUser.Delete("/{id}", s.baseTagHandlers.delete)
			})

			routeAdmin.Get("/users", s.userHandlers.findAll)
			routeAdmin.Post("/user-reports", s.reportsHandlers.userReport)
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
