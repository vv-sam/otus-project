package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("item not found")
)

type uniqueObject interface {
	GetId() uuid.UUID
}

type JsonRepository[T uniqueObject] struct {
	m     sync.RWMutex
	items []T
	path  string
}

func NewJsonRepository[T uniqueObject](dir string, entity string) *JsonRepository[T] {
	r := &JsonRepository[T]{}
	r.initialize(dir, entity)
	return r
}

func (r *JsonRepository[T]) initialize(dir string, entity string) error {
	if dir == "" {
		return fmt.Errorf("directory has to be set")
	}

	r.path = path.Join(dir, entity+".json")
	return r.loadRepository()
}

func (r *JsonRepository[T]) Get(id uuid.UUID) (T, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	for _, item := range r.items {
		if item.GetId() == id {
			return item, nil
		}
	}

	// Думаю, что отсутствие записи с нужным id не является ошибкой, но в будущем может случится переезд на sql/nosql/s3/etc базу данных,
	// тогда нужно будет возвращать ошибку
	var res T
	return res, nil
}

func (r *JsonRepository[T]) GetAll() ([]T, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	res := make([]T, len(r.items))
	copy(res, r.items)
	return res, nil
}

func (r *JsonRepository[T]) Add(item T) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.items = append(r.items, item)
	return r.saveRepository()
}

func (r *JsonRepository[T]) Update(id uuid.UUID, item T) error {
	r.m.Lock()
	defer r.m.Unlock()

	for i, item := range r.items {
		if item.GetId() == id {
			r.items[i] = item
			return r.saveRepository()
		}
	}

	// А тут, я считаю, это ошибка, так как запись должна существовать для обновления данных
	return ErrNotFound
}

func (r *JsonRepository[T]) Delete(id uuid.UUID) error {
	r.m.Lock()
	defer r.m.Unlock()

	for i, item := range r.items {
		if item.GetId() == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return r.saveRepository()
		}
	}

	return ErrNotFound
}

func (r *JsonRepository[T]) loadRepository() error {
	data, err := loadRawData(r.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &r.items)
}

func (r *JsonRepository[T]) saveRepository() error {
	data, err := json.Marshal(r.items)
	if err != nil {
		return err
	}
	return saveRawData(r.path, data)
}

func loadRawData(p string) ([]byte, error) {
	file, err := os.OpenFile(p, os.O_RDONLY, 0400)
	defer file.Close()

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return []byte{}, err
	}

	return io.ReadAll(file)
}

func saveRawData(p string, data []byte) error {
	file, err := os.OpenFile(p, os.O_WRONLY, 0200)
	defer file.Close()

	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		file, err = os.Create(p)
	}

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
