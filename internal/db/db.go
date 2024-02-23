package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type PhotoEntity struct {
	Id   string
	Name string
	Data []byte
}

type Repository interface {
	Get(id string) (*PhotoEntity, error)
	GetAll() ([]*PhotoEntity, error)
	Save(entity *PhotoEntity) error
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

	var name string
	//check if table is created
	row := repo.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='photos'")
	err := row.Scan(&name)
	if err != nil {
		// if table does not exist, create it
		if errors.Is(err, sql.ErrNoRows) {
			repo.db.Exec("CREATE TABLE photos(id INTEGER PRIMARY KEY, name TEXT, data BLOB)")
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
func (repo *SQLLiteRepository) Save(e PhotoEntity) error {
	_, err := repo.db.Exec("INSERT INTO photos (name,data) VALUES (?, ?)", e.Name, e.Data)
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
func (repo *SQLLiteRepository) Get(id string) (*PhotoEntity, error) {
	rows, err := repo.db.Query("SELECT id, name, data FROM PHOTOS WHERE id =?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entity PhotoEntity
	for rows.Next() {
		err := rows.Scan(&entity.Id, &entity.Name, &entity.Data)
		if err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

// TODO: implement get all
func (repo *SQLLiteRepository) GetAll() ([]*PhotoEntity, error) {
	rows, err := repo.db.Query("SELECT id, name, data FROM photos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*PhotoEntity
	for rows.Next() {
		var entity PhotoEntity
		err := rows.Scan(&entity.Id, &entity.Name, &entity.Data)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}

func (repo *SQLLiteRepository) Delete(id string) error {
	_, err := repo.db.Exec("DELETE FROM photos WHERE id =?", id)
	if err != nil {
		fmt.Printf("Error deleting photo with id %s: %s", id, err)
		return err
	}
	return nil
}
