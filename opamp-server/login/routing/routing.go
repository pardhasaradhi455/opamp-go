package routing

import (
	"context"
	"net/http"

	"github.com/open-telemetry/opamp-go/opamp-server/login/service"
	"github.com/open-telemetry/opamp-go/opamp-server/login/controller"
	"github.com/open-telemetry/opamp-go/opamp-server/login/models"

	"github.com/gorilla/mux"
)

var basepath = "/api"

func GetRoutes(loginController *controller.LoginController) []models.RouteDef {
	return []models.RouteDef{
		{
			Path:            basepath + "/roles",
			Method:          "GET",
			HandlerFunction: loginController.GetRoles,
			Protected:       true,
		},
		{
			Path:            basepath + "/agents",
			Method:          "GET",
			HandlerFunction: loginController.GetAgents,
			Protected:       true,
		},
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := service.ValidateToken(authHeader)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), service.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AddRoutes(router *mux.Router, loginController *controller.LoginController) {

	for _, route := range GetRoutes(loginController) {

		handler := route.HandlerFunction

		if route.Protected {
			handler = AuthMiddleware(handler)
		}

		router.HandleFunc(route.Path, handler).Methods(route.Method)
	}
}

func EnableCors(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
