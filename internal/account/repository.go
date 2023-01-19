package account

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/models/account"
)

type (
	AccountRepository interface {
		Register(ctx context.Context, params account.Account) (int64, error)
		FindByEmail(ctx context.Context, email string) (account.Account, error)
		FindByID(ctx context.Context, id int64) (account.Account, error)
		Update(ctx context.Context, id int64, params account.Account) error
		Delete(ctx context.Context, id int64) error
	}

	accountRepositoryImpl struct {
		db        *sql.DB
		tableName string
	}
)

func NewAccountRepositoryImpl(db *sql.DB, tableName string) AccountRepository {
	return &accountRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}

func (ar *accountRepositoryImpl) Register(ctx context.Context, params account.Account) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s(name, password, email, address, created_at) VALUES (?, ?, ?, ?, ?)`, ar.tableName)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.Name,
		params.Password,
		params.Email,
		params.Address,
		params.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return 0, exception.ErrInternalServer
	}

	ID, _ := result.LastInsertId()

	return ID, nil
}

func (ar *accountRepositoryImpl) FindByEmail(ctx context.Context, email string) (account.Account, error) {
	var user account.Account
	query := fmt.Sprintf(`SELECT id, name, password, email, address, created_at, update_at FROM %s WHERE email = ?`, ar.tableName)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return user, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Address,
		&user.CreatedAt,
		&user.UpdateAt,
	)
	if err != nil {
		log.Println(err)
		return user, exception.ErrNotFound
	}

	return user, nil
}

func (ar *accountRepositoryImpl) FindByID(ctx context.Context, id int64) (account.Account, error) {
	var user account.Account
	query := fmt.Sprintf(`SELECT id, name, password, email, address, created_at, update_at FROM %s WHERE id = ?`, ar.tableName)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return user, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Address,
		&user.CreatedAt,
		&user.UpdateAt,
	)
	if err != nil {
		log.Println(err)
		return user, exception.ErrNotFound
	}

	return user, nil
}

func (ar *accountRepositoryImpl) Update(ctx context.Context, id int64, params account.Account) error {
	query := fmt.Sprintf(`UPDATE %s SET name = ?, password = ?, email = ?, address = ?, update_at = ? WHERE id = %d`, ar.tableName, id)
	stmt, err := ar.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		params.Name,
		params.Password,
		params.Email,
		params.Address,
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

func (ar *accountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = %d`, ar.tableName, id)
	stmt, err := ar.db.PrepareContext(ctx, query)
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
