package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRepository struct {
	data map[string]*PhotoEntity
}

func NewTestRepository() *TestRepository {
	return &TestRepository{
		data: map[string]*PhotoEntity{},
	}
}

func (repo *TestRepository) Get(id string) (*PhotoEntity, error) {
	if e, ok := repo.data[id]; ok {
		return e, nil
	}
	return nil, errors.New("not found")
}

func (repo *TestRepository) GetAll() ([]*PhotoEntity, error) {
	return nil, nil
}

func (repo *TestRepository) Save(e PhotoEntity) error {
	repo.data[e.Id] = &e
	return nil
}

func (repo *TestRepository) Delete(id string) error {
	delete(repo.data, id)
	return nil
}

func TestSave(t *testing.T) {
	repo := NewTestRepository()

	e := PhotoEntity{
		Id:   "1",
		Name: "Test Photo",
		Data: []byte("Hello World"),
	}

	err := repo.Save(e)
	assert.NoError(t, err)

	saved, err := repo.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, e.Id, saved.Id)
	assert.Equal(t, e.Name, saved.Name)
	assert.Equal(t, e.Data, saved.Data)
}

func TestGet(t *testing.T) {
	repo := NewTestRepository()

	expected := PhotoEntity{
		Id:   "1",
		Name: "Test Photo",
		Data: []byte("Hello World"),
	}

	err := repo.Save(expected)
	assert.NoError(t, err)

	saved, err := repo.Get("1")

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, saved.Id)
	assert.Equal(t, expected.Name, saved.Name)
	assert.Equal(t, expected.Data, saved.Data)
}

func TestDelete(t *testing.T) {
	repo := NewTestRepository()

	saved := PhotoEntity{
		Id:   "1",
		Name: "Test Photo",
		Data: []byte("Hello World"),
	}
	err := repo.Save(saved)
	assert.NoError(t, err)

	err = repo.Delete("1")

	assert.NoError(t, err)
	_, err = repo.Get("1")
	assert.Error(t, err)
}
