package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/repository"
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

// @Summary Get agent by id
// @Description Get agent by id
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 200 {object} agent.Info
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/agents/{id} [get]
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

	if agent == nil {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(agent)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal agent: %w", err).Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// @Summary Get all agents
// @Description Get all agents
// @Tags agents
// @Accept json
// @Produce json
// @Success 200 {array} agent.Info
// @Failure 500 {object} error
// @Router /api/agents [get]
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

// @Summary Create a new agent
// @Description Create a new agent
// @Tags agents
// @Accept json
// @Produce json
// @Param agent body agent.Info true "Agent"
// @Success 201 {object} agent.Info
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/agents [post]
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

// @Summary Update an agent
// @Description Update an agent
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Param agent body agent.Info true "Agent"
// @Success 200 {object} agent.Info
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/agents/{id} [put]
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
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "agent not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Errorf("failed to update agent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete an agent
// @Description Delete an agent
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 200 {object} agent.Info
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/agents/{id} [delete]
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
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "agent not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Errorf("failed to delete agent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
