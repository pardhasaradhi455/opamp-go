package login

import (
	"log"
	"net/http"

	"github.com/open-telemetry/opamp-go/opamp-server/login/controller"
	"github.com/open-telemetry/opamp-go/opamp-server/login/routing"

	"github.com/gorilla/mux"
)

func Start(port string) {
	loginController := &controller.LoginController{}
	StartRouter(port, loginController)
}

func StartRouter(port string, loginController *controller.LoginController) {
	router := mux.NewRouter()
	routing.AddRoutes(router, loginController)

	corsRouter := routing.EnableCors(router)
	server := http.Server{
		Addr:    ":" + port,
		Handler: corsRouter,
	}

	log.Println("stating router")
	go server.ListenAndServe()
}
