package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/models/store"
)

type (
	StoreRepository interface {
		Create(ctx context.Context, params store.Store) (int64, error)
		FindByID(ctx context.Context, id int64) (store.Store, error)
		Update(ctx context.Context, id int64, params store.Store) error
		Delete(ctx context.Context, id int64) error
	}

	storeRepositoryImpl struct {
		DB        *sql.DB
		tableName string
	}
)

func NewStoreRepository(db *sql.DB, tableName string) StoreRepository {
	return &storeRepositoryImpl{
		DB:        db,
		tableName: tableName,
	}
}

func (repo *storeRepositoryImpl) Create(ctx context.Context, params store.Store) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (userID, nameStore, description, created_at) VALUES (?,?,?,?)`, repo.tableName)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		ctx,
		params.UserID,
		params.NameStore,
		params.Description,
		params.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	ID, _ := result.LastInsertId()
	return ID, nil
}

func (repo *storeRepositoryImpl) FindByID(ctx context.Context, id int64) (store.Store, error) {
	var store store.Store
	query := fmt.Sprintf(`SELECT id, userID, nameStore, description, created_at, update_at FROM %s WHERE id = ?`, repo.tableName)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return store, exception.ErrInternalServer
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(
		&store.ID,
		&store.UserID,
		&store.NameStore,
		&store.Description,
		&store.CreatedAt,
		&store.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return store, exception.ErrInternalServer
	}

	return store, nil
}

func (repo *storeRepositoryImpl) Update(ctx context.Context, id int64, params store.Store) error {
	query := fmt.Sprintf(`UPDATE %s SET nameStore = ?, description = ?, update_at = ? WHERE id = %d`, repo.tableName, id)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.NameStore,
		params.Description,
		params.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected < 1 {
		return exception.ErrNotFound
	}

	return nil
}

func (repo *storeRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = %d`, repo.tableName, id)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
	)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return exception.ErrNotFound
	}

	return nil
}
