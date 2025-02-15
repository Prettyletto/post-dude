package repository

import (
	"database/sql"
	"fmt"

	"github.com/Prettyletto/post-dude/cmd/internal/db"
	"github.com/Prettyletto/post-dude/cmd/internal/model"
)

type CollectionRepository interface {
	SaveCollection(colllection *model.Collection) error
	FindAllCollections() ([]model.Collection, error)
	FindCollectionById(id int) (*model.Collection, error)
	UpdateCollection(id int, collection *model.Collection) error
	DeleteCollection(id int) error
}

type collectionRepository struct {
	database *db.DataBase
}

func NewCollectionRepository(database *db.DataBase) CollectionRepository {
	return &collectionRepository{database: database}
}

func (r *collectionRepository) SaveCollection(collection *model.Collection) error {
	query := `INSERT INTO collections (name) VALUES(?)`

	_, err := r.database.DB.Exec(query, collection.Name)
	if err != nil {
		return fmt.Errorf("Failed to insert collection: %w", err)
	}

	return nil
}

func (r *collectionRepository) FindAllCollections() ([]model.Collection, error) {
	rows, err := r.database.DB.Query(`SELECT id,name FROM collections`)
	if err != nil {
		return nil, fmt.Errorf("error querying collections: %w", err)
	}
	defer rows.Close()

	var collections []model.Collection

	for rows.Next() {
		var coll model.Collection
		if err := rows.Scan(&coll.ID, &coll.Name); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		collections = append(collections, coll)
	}
	return collections, nil
}

func (r *collectionRepository) FindCollectionById(id int) (*model.Collection, error) {
	query := `SELECT * FROM collections WHERE id = ?`
	row := r.database.DB.QueryRow(query, id)

	var collection model.Collection
	err := row.Scan(&collection.ID, &collection.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no collection found")
		}
		return nil, fmt.Errorf("Error on query")
	}

	return &collection, nil
}

func (r *collectionRepository) UpdateCollection(id int, collection *model.Collection) error {
	query := `UPDATE collections SET name = ? where id = ?`
	result, err := r.database.DB.Exec(query, collection.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update collection")
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows: %w", err)
	}
	if updatedRows == 0 {
		return fmt.Errorf("no collection found with id %d", collection.ID)
	}

	return nil
}

func (r *collectionRepository) DeleteCollection(id int) error {
	query := `DELETE FROM collections WHERE id = ?`
	result, err := r.database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete collection")
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to deleted rows: %w", err)
	}
	if deletedRows == 0 {
		return fmt.Errorf("no collection found with id %d", id)
	}

	return nil
}
