package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Entity struct {
	Id   string
	Name string
}

type Repository interface {
	Get(id string) (*Entity, error)
	GetAll() ([]*Entity, error)
	Save(entity *Entity) error
	Delete(id string) error
}

type SQLLiteRepository struct {
	db sql.DB
}

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}
	return db
}

func NewSQLLiteRepository() *SQLLiteRepository {
	repo := &SQLLiteRepository{
		db: *GetDb(),
	}

	//check if table is created
	row := repo.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='entities'")
	err := row.Scan()
	if err != nil {
		// if table does not exist, create it
		if errors.Is(err, sql.ErrNoRows) {
			repo.db.Exec("CREATE TABLE entities (id TEXT PRIMARY KEY, name TEXT)")
		} else {
			panic(err)
		}
	}

	return repo
}

/*
Saves an entity to the SQLLite repository.
@param {Entity} e - The entity to be saved.
@returns {error} - An error if the save operation fails, otherwise null.
*/
func (repo *SQLLiteRepository) Save(e Entity) error {
	_, err := repo.db.Exec("INSERT INTO entities (id, name) VALUES (?,?)", e.Id, e.Name)
	if err != nil {
		return err
	}
	return nil
}

/*
*
Retrieves an entity from the SQLLiteRepository based on the given id.
@param {string} id - The id of the entity to retrieve.
@returns {Entity} - The retrieved entity, along with an error if the operation fails.
*/
func (repo *SQLLiteRepository) Get(id string) (*Entity, error) {
	rows, err := repo.db.Query("SELECT id, name FROM entities WHERE id =?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entity Entity
	for rows.Next() {
		err := rows.Scan(&entity.Id, &entity.Name)
		if err != nil {
			return nil, err
		}
	}
	return &entity, nil
}

// TODO: implement get all
func (repo *SQLLiteRepository) GetAll() ([]*Entity, error) {
	return nil, nil
}

// TODO: implement delete
func (repo *SQLLiteRepository) Delete(id string) error {
	return nil
}
