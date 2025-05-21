package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/services"
)

type tasksRepository interface {
	Get(id uuid.UUID) (*task.Task, error)
	GetAll() ([]*task.Task, error)
	Add(task *task.Task) error
	Update(id uuid.UUID, task *task.Task) error
	Delete(id uuid.UUID) error
}

type Tasks struct {
	r tasksRepository
	v *services.Validator
}

func NewTasks(r tasksRepository, v *services.Validator) *Tasks {
	return &Tasks{r: r, v: v}
}

func (t *Tasks) GetById(w http.ResponseWriter, r *http.Request) {
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

	task, err := t.r.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (t *Tasks) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.r.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to marshal tasks: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (t *Tasks) Post(w http.ResponseWriter, r *http.Request) {
	var task task.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal task: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if !t.v.IsValid(&task) {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	if err := t.r.Add(&task); err != nil {
		http.Error(w, fmt.Errorf("failed to add task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *Tasks) Put(w http.ResponseWriter, r *http.Request) {
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

	var task task.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal task: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if !t.v.IsValid(&task) {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	if err := t.r.Update(uuid, &task); err != nil {
		http.Error(w, fmt.Errorf("failed to update task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Tasks) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := t.r.Delete(uuid); err != nil {
		http.Error(w, fmt.Errorf("failed to delete task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
