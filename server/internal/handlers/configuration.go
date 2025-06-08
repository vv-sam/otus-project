package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/history"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
)

type configurationRepository interface {
	Get(id uuid.UUID) (*configuration.Factorio, error)
	GetAll() ([]*configuration.Factorio, error)
	Add(configuration *configuration.Factorio) error
	Update(id uuid.UUID, configuration *configuration.Factorio) error
	Delete(id uuid.UUID) error
	GetHistory() ([]history.Log[*configuration.Factorio], error)
}

type Configuration struct {
	r configurationRepository
	v *services.Validator
}

func NewConfiguration(r configurationRepository, v *services.Validator) *Configuration {
	return &Configuration{r: r, v: v}
}

// @Summary Get configuration by id
// @Description Get configuration by id
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {object} configuration.Factorio
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/configurations/{id} [get]
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

	if configuration == nil {
		http.Error(w, "configuration not found", http.StatusNotFound)
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

// @Summary Get all configurations
// @Description Get all configurations
// @Tags configurations
// @Accept json
// @Produce json
// @Success 200 {array} configuration.Factorio
// @Failure 500 {object} error
// @Router /api/configurations [get]
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

// @Summary Create a new configuration
// @Description Create a new configuration
// @Tags configurations
// @Accept json
// @Produce json
// @Param configuration body configuration.Factorio true "Configuration"
// @Security BearerAuth
// @Success 201
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/configurations [post]
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

// @Summary Update a configuration
// @Description Update a configuration
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param configuration body configuration.Factorio true "Configuration"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/configurations/{id} [put]
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
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "configuration not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Errorf("failed to update configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a configuration
// @Description Delete a configuration
// @Tags configurations
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/configurations/{id} [delete]
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
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "configuration not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Errorf("failed to delete configuration: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get history
// @Description Get history
// @Tags configurations
// @Accept json
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} error
// @Router /api/configurations/history [get]
func (c *Configuration) GetHistory(w http.ResponseWriter, r *http.Request) {
	history, err := c.r.GetHistory()
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get history: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(history)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal history: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
