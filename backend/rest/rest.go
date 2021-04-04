package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
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

	tagHandlers     tagHandlers
	dayHandlers     dayHandlers
	eventHandlers   eventHandlers
	projectHandlers projectHandlers
	managerHandlers managerHandlers
	adminHandlers   adminHandlers
}

type uuidProvider struct {
}

func (p *uuidProvider) New() string {
	return uuid.New().String()
}

func CreateRestServer(database *mongo.Database, jwtAuth *JwtAuthService, keycloakUrl string) *RestServer {
	uuidProvider := &uuidProvider{}
	validator := CreateValidator()

	dayStore := repository.CreateDayStore(database, uuidProvider)
	tagStore := repository.CreateTagStore(database, uuidProvider)
	projectStore := repository.CreateProjectStore(database, uuidProvider)

	userService := service.CreateUserService(keycloakUrl, projectStore)
	reportService := service.CreateReportService(dayStore, tagStore, projectStore)
	eventService := service.CreateEventService(dayStore, uuidProvider)

	return &RestServer{
		authenticator: &AuthService{authenticator: jwtAuth},

		tagHandlers:     tagHandlers{tagStore: tagStore, projectStore: projectStore, validator: validator},
		dayHandlers:     dayHandlers{store: dayStore},
		eventHandlers:   eventHandlers{service: eventService, validator: validator},
		projectHandlers: projectHandlers{store: projectStore, userService: userService},

		adminHandlers: adminHandlers{
			validator:     validator,
			projectStore:  projectStore,
			baseTagStore:  tagStore,
			userService:   userService,
			reportService: reportService,
		},

		managerHandlers: managerHandlers{
			validator:     validator,
			projectStore:  projectStore,
			tagStore:      tagStore,
			userService:   userService,
			reportService: reportService,
		},
	}
}

func (s *RestServer) Run() error {
	useSsl := false

	if useSsl {
		log.Fatalf("Application does not support HTTPS.")
	} else {
		s.httpServer = s.makeHTTPServer(9000, s.routes())
		err := s.httpServer.ListenAndServe()
		if err != nil {
			return err
		}
	}

	defer s.Shutdown()

	return nil
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

	router.Get("/api/health", func(writer http.ResponseWriter, request *http.Request) {
		render.Status(request, http.StatusOK)
	})

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

			routeManager.Get("/project-reports/{id}", s.managerHandlers.getProjectReport)
			routeManager.Get("/projects/{id}/users", s.managerHandlers.getProjectUsers)

			routeManager.Post("/tags", s.managerHandlers.createTag)
			routeManager.Put("/tags/{id}", s.managerHandlers.updateTag)
			routeManager.Delete("/tags/{id}", s.managerHandlers.deleteTag)
		})

		// Admin
		r.Route("/admin", func(routeAdmin chi.Router) {
			routeAdmin.Use(s.authenticator.HasRole(adminRole))

			routeAdmin.Route("/projects", func(routeAdminProject chi.Router) {
				routeAdminProject.Post("/", s.adminHandlers.createProject)
				routeAdminProject.Post("/search", s.adminHandlers.findProject)
				routeAdminProject.Put("/{id}", s.adminHandlers.updateProject)
				routeAdminProject.Delete("/{id}", s.adminHandlers.deleteProject)
			})

			routeAdmin.Route("/base-tags", func(routeUser chi.Router) {
				routeUser.Post("/", s.adminHandlers.createBaseTag)
				routeUser.Put("/{id}", s.adminHandlers.updateBaseTag)
				routeUser.Delete("/{id}", s.adminHandlers.deleteBaseTag)
			})

			routeAdmin.Get("/users", s.adminHandlers.findAllUser)
			routeAdmin.Post("/user-reports", s.adminHandlers.userReport)
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
