package main

import (
	"net/http"

	"github.com/vv-sam/otus-project/server/internal/handlers"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
)

func main() {
	ar := repository.NewJsonRepository[*agent.Info]("D:\\otus-data", "agents")
	cr := repository.NewJsonRepository[*configuration.Factorio]("D:\\otus-data", "configurations")
	tr := repository.NewJsonRepository[*task.Task]("D:\\otus-data", "tasks")

	ah := handlers.NewAgents(ar, &services.Validator{})
	ch := handlers.NewConfiguration(cr, &services.Validator{})
	th := handlers.NewTasks(tr, &services.Validator{})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/agents", ah.GetAll)
	mux.HandleFunc("GET /api/agents/{id}", ah.GetById)
	mux.HandleFunc("POST /api/agents", ah.Post)
	mux.HandleFunc("PUT /api/agents/{id}", ah.Put)
	mux.HandleFunc("DELETE /api/agents/{id}", ah.Delete)

	mux.HandleFunc("GET /api/configurations", ch.GetAll)
	mux.HandleFunc("GET /api/configurations/{id}", ch.GetById)
	mux.HandleFunc("POST /api/configurations", ch.Post)
	mux.HandleFunc("PUT /api/configurations/{id}", ch.Put)
	mux.HandleFunc("DELETE /api/configurations/{id}", ch.Delete)

	mux.HandleFunc("GET /api/tasks", th.GetAll)
	mux.HandleFunc("GET /api/tasks/{id}", th.GetById)
	mux.HandleFunc("POST /api/tasks", th.Post)
	mux.HandleFunc("PUT /api/tasks/{id}", th.Put)
	mux.HandleFunc("DELETE /api/tasks/{id}", th.Delete)

	http.ListenAndServe(":8080", mux)
}
