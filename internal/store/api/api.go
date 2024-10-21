package api

import (
	categoryController "back/internal/store/api/controllers/category"
	locationController "back/internal/store/api/controllers/location"
	locationMaterialController "back/internal/store/api/controllers/location_material"
	materialController "back/internal/store/api/controllers/material"
	transactionController "back/internal/store/api/controllers/transaction"
	unitController "back/internal/store/api/controllers/unit"
	userController "back/internal/store/api/controllers/user"
	"back/internal/store/api/helper"
	pgstore "back/internal/store/pgstore/sqlc"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type apiHandler struct {
	q *pgstore.Queries
	r *chi.Mux
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				helper.HandleError(w, "", err.Error(), http.StatusUnauthorized)
				return
			}

			options := ja.ValidateOptions()

			if token == nil || jwt.Validate(token, options...) != nil {
				helper.HandleError(w, "", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		q: q,
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/client", func(r chi.Router) {
			r.Post("/auth", userController.New(q).Auth)
		})
		r.Route("/admin", func(r chi.Router) {
			jwtSecret := os.Getenv("JWT_SECRET")
			tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
			r.Use(jwtauth.Verifier(tokenAuth), Authenticator(tokenAuth))

			r.Get("/refresh", userController.New(q).Refresh)

			r.Route("/user", func(r chi.Router) {
				r.Post("/", userController.New(q).Create)
				r.Delete("/{id}", userController.New(q).Delete)
				r.Put("/{id}", userController.New(q).Update)
				r.Get("/list", userController.New(q).List)
				r.Get("/{id}", userController.New(q).Find)
			})

			r.Route("/category", func(r chi.Router) {
				r.Post("/", categoryController.New(q).Create)
				r.Delete("/{id}", categoryController.New(q).Delete)
				r.Put("/{id}", categoryController.New(q).Update)
				r.Get("/list", categoryController.New(q).List)
				r.Get("/{id}", categoryController.New(q).Find)
				r.Get("/autocomplete", categoryController.New(q).Autocomplete)
			})

			r.Route("/location", func(r chi.Router) {
				r.Post("/", locationController.New(q).Create)
				r.Delete("/{id}", locationController.New(q).Delete)
				r.Put("/{id}", locationController.New(q).Update)
				r.Get("/list", locationController.New(q).List)
				r.Get("/{id}", locationController.New(q).Find)
				r.Get("/autocomplete", locationController.New(q).Autocomplete)
			})

			r.Route("/unit", func(r chi.Router) {
				r.Post("/", unitController.New(q).Create)
				r.Delete("/{id}", unitController.New(q).Delete)
				r.Put("/{id}", unitController.New(q).Update)
				r.Get("/list", unitController.New(q).List)
				r.Get("/{id}", unitController.New(q).Find)
				r.Get("/autocomplete", unitController.New(q).Autocomplete)
			})

			r.Route("/material", func(r chi.Router) {
				r.Post("/", materialController.New(q).Create)
				r.Delete("/{id}", materialController.New(q).Delete)
				r.Put("/{id}", materialController.New(q).Update)
				r.Get("/list", materialController.New(q).List)
				r.Get("/{id}", materialController.New(q).Find)
				r.Get("/autocomplete", materialController.New(q).Autocomplete)
			})

			r.Route("/location-material", func(r chi.Router) {
				r.Get("/list", locationMaterialController.New(q).List)
				r.Get("/{id}", locationMaterialController.New(q).Find)
				r.Get("/relation", locationMaterialController.New(q).FindRelation)
			})

			r.Route("/transaction", func(r chi.Router) {
				r.Post("/", transactionController.New(q).Create)
				r.Get("/list", transactionController.New(q).List)
				r.Get("/{id}", transactionController.New(q).Find)
			})
		})
	})

	a.r = r
	return a
}
