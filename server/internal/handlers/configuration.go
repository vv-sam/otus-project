package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/services"
)

type configurationRepository interface {
	Get(id uuid.UUID) (*configuration.Factorio, error)
	GetAll() ([]*configuration.Factorio, error)
	Add(configuration *configuration.Factorio) error
	Update(id uuid.UUID, configuration *configuration.Factorio) error
	Delete(id uuid.UUID) error
}

type Configuration struct {
	r configurationRepository
	v *services.Validator
}

func NewConfiguration(r configurationRepository, v *services.Validator) *Configuration {
	return &Configuration{r: r, v: v}
}

func (c *Configuration) GetById(w http.ResponseWriter, r *http.Request) {
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

	configuration, err := c.r.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(configuration)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (c *Configuration) GetAll(w http.ResponseWriter, r *http.Request) {
	configuration, err := c.r.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(configuration)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (c *Configuration) Post(w http.ResponseWriter, r *http.Request) {
	var configuration configuration.Factorio
	if err := json.NewDecoder(r.Body).Decode(&configuration); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal configuration: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if !c.v.IsValid(&configuration) {
		http.Error(w, "invalid configuration", http.StatusBadRequest)
		return
	}

	if err := c.r.Add(&configuration); err != nil {
		http.Error(w, fmt.Errorf("failed to add configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Configuration) Put(w http.ResponseWriter, r *http.Request) {
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

	var configuration configuration.Factorio
	if err := json.NewDecoder(r.Body).Decode(&configuration); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal configuration: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if !c.v.IsValid(&configuration) {
		http.Error(w, "invalid configuration", http.StatusBadRequest)
		return
	}

	if err := c.r.Update(uuid, &configuration); err != nil {
		http.Error(w, fmt.Errorf("failed to update configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Configuration) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.r.Delete(uuid); err != nil {
		http.Error(w, fmt.Errorf("failed to delete configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
