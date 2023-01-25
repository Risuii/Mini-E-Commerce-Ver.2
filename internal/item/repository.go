package item

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/models/item"
)

type (
	ItemRepository interface {
		AddItem(ctx context.Context, params item.Item) (int64, error)
		GetAllItem() ([]item.Item, error)
		FindByID(ctx context.Context, id int64) (item.Item, error)
		FindByName(ctx context.Context, name string) (item.Item, error)
		UpdateItem(ctx context.Context, id int64, params item.Item) error
		DeleteItem(ctx context.Context, id int64) error
	}

	itemRepositoryImpl struct {
		DB        *sql.DB
		tableName string
	}
)

func NewItemRepositoryImpl(db *sql.DB, tableName string) ItemRepository {
	return &itemRepositoryImpl{
		DB:        db,
		tableName: tableName,
	}
}

func (repo *itemRepositoryImpl) AddItem(ctx context.Context, params item.Item) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (storeID, name, description, quantity, created_at) VALUES (?,?,?,?,?)`, repo.tableName)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.StoreID,
		params.Name,
		params.Description,
		params.Quantity,
		params.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	ID, _ := result.LastInsertId()

	return ID, nil
}

func (repo *itemRepositoryImpl) GetAllItem() ([]item.Item, error) {
	var items []item.Item

	rows, err := repo.DB.Query(fmt.Sprintf(`SELECT id, storeID, name, description, quantity, created_at, update_at FROM %s`, repo.tableName))
	if err != nil {
		log.Println(err)
		return items, exception.ErrInternalServer
	}

	defer rows.Close()

	for rows.Next() {
		var c item.Item
		if err := rows.Scan(
			&c.ID,
			&c.StoreID,
			&c.Name,
			&c.Description,
			&c.Quantity,
			&c.CreatedAt,
			&c.UpdateAt,
		); err != nil {
			log.Println(err)
			return items, exception.ErrInternalServer
		}
		items = append(items, c)
	}

	return items, nil
}

func (repo *itemRepositoryImpl) FindByID(ctx context.Context, id int64) (item.Item, error) {
	var items item.Item

	query := fmt.Sprintf(`SELECT id, storeID, name, description, quantity, created_at, update_at FROM %s WHERE id = ?`, repo.tableName)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return items, exception.ErrInternalServer
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(
		&items.ID,
		&items.StoreID,
		&items.Name,
		&items.Description,
		&items.Quantity,
		&items.CreatedAt,
		&items.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return items, exception.ErrNotFound
	}

	return items, nil
}

func (repo *itemRepositoryImpl) FindByName(ctx context.Context, name string) (item.Item, error) {
	var items item.Item

	query := fmt.Sprintf(`SELECT id, storeID, name, description, quantity, created_at, update_at FROM %s WHERE name = ?`, repo.tableName)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return items, exception.ErrInternalServer
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, name)

	err = row.Scan(
		&items.ID,
		&items.StoreID,
		&items.Name,
		&items.Description,
		&items.Quantity,
		&items.CreatedAt,
		&items.UpdateAt,
	)

	if err != nil {
		log.Println(err)
		return items, exception.ErrNotFound
	}

	return items, nil
}

func (repo *itemRepositoryImpl) UpdateItem(ctx context.Context, id int64, params item.Item) error {
	query := fmt.Sprintf(`UPDATE %s SET name = ?, description = ? WHERE id = %d`, repo.tableName, id)
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.Name,
		params.Description,
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

func (repo *itemRepositoryImpl) DeleteItem(ctx context.Context, id int64) error {
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
