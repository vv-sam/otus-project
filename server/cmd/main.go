package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/vv-sam/otus-project/server/cmd/docs"
	"github.com/vv-sam/otus-project/server/internal/handlers"
	"github.com/vv-sam/otus-project/server/internal/middleware"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
)

// @title			Otus-Project
// @version		1.0
// @description	This is a simple API for the Otus-Project
// @host			localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	user := flag.String("user", "admin", "admin username")
	pass := flag.String("password", "", "admin password")
	flag.Parse()

	if *pass == "" {
		fmt.Println("password is required")
		os.Exit(1)
	}

	as := services.NewUsers(map[string]string{
		*user: *pass,
	})

	am := middleware.NewAuthMiddleware(as)

	ar := repository.NewJsonRepository[*agent.Info]("D:\\otus-data", "agents")
	cr := repository.NewJsonRepository[*configuration.Factorio]("D:\\otus-data", "configurations")
	tr := repository.NewJsonRepository[*task.Task]("D:\\otus-data", "tasks")

	ah := handlers.NewAgents(ar, &services.Validator{})
	ch := handlers.NewConfiguration(cr, &services.Validator{})
	th := handlers.NewTasks(tr, &services.Validator{})
	au := handlers.NewAuth(as)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	mux.HandleFunc("POST /api/auth/login", au.Login)

	mux.HandleFunc("GET /api/agents", ah.GetAll)
	mux.HandleFunc("GET /api/agents/{id}", ah.GetById)
	mux.Handle("POST /api/agents", am.Authenticate(ah.Post))
	mux.Handle("PUT /api/agents/{id}", am.Authenticate(ah.Put))
	mux.Handle("DELETE /api/agents/{id}", am.Authenticate(ah.Delete))

	mux.HandleFunc("GET /api/configurations", ch.GetAll)
	mux.HandleFunc("GET /api/configurations/{id}", ch.GetById)
	mux.Handle("POST /api/configurations", am.Authenticate(ch.Post))
	mux.Handle("PUT /api/configurations/{id}", am.Authenticate(ch.Put))
	mux.Handle("DELETE /api/configurations/{id}", am.Authenticate(ch.Delete))

	mux.HandleFunc("GET /api/tasks", th.GetAll)
	mux.HandleFunc("GET /api/tasks/{id}", th.GetById)
	mux.Handle("POST /api/tasks", am.Authenticate(th.Post))
	mux.Handle("PUT /api/tasks/{id}", am.Authenticate(th.Put))
	mux.Handle("DELETE /api/tasks/{id}", am.Authenticate(th.Delete))

	http.ListenAndServe(":8080", mux)
}
