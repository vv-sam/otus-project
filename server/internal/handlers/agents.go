package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/services"
)

type agentRepository interface {
	Get(id uuid.UUID) (*agent.Info, error)
	GetAll() ([]*agent.Info, error)
	Add(agent *agent.Info) error
	Update(id uuid.UUID, agent *agent.Info) error
	Delete(id uuid.UUID) error
}

type Agents struct {
	r agentRepository
	v *services.Validator
}

func NewAgents(r agentRepository, v *services.Validator) *Agents {
	return &Agents{r: r, v: v}
}

func (a *Agents) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to parse id: %w", err).Error(), http.StatusBadRequest)
		return
	}

	agent, err := a.r.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(agent)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal agent: %w", err).Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a *Agents) GetAll(w http.ResponseWriter, r *http.Request) {
	agents, err := a.r.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(agents)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal agents: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a *Agents) Post(w http.ResponseWriter, r *http.Request) {
	var agent agent.Info
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal agent: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if !a.v.IsValid(agent) {
		http.Error(w, "invalid agent", http.StatusBadRequest)
		return
	}

	if err := a.r.Add(&agent); err != nil {
		http.Error(w, fmt.Errorf("failed to add agent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Agents) Put(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to parse id: %w", err).Error(), http.StatusBadRequest)
		return
	}

	var agent agent.Info
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal agent: %w", err).Error(), http.StatusBadRequest)
		return
	}
	agent.AgentId = uuid

	if !a.v.IsValid(agent) {
		http.Error(w, "invalid agent", http.StatusBadRequest)
	}

	if err := a.r.Update(uuid, &agent); err != nil {
		http.Error(w, fmt.Errorf("failed to update agent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Agents) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to parse id: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := a.r.Delete(uuid); err != nil {
		http.Error(w, fmt.Errorf("failed to delete agent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
