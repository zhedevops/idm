package employee

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) FindById(id int64) (employee Entity, err error) {
	err = r.db.Get(&employee, "SELECT * FROM employee WHERE id = $1", id)
	return
}

func (r *Repository) Create(e *Entity) error {
	query := "INSERT INTO employee (name) VALUES ($1) RETURNING id, name"

	// В PostgreSQL Get выполнит запрос и сразу вернёт вставленную запись
	return r.db.Get(e, query, e.Name)
}

func (r *Repository) CreateNamed(e *Entity) error {
	query := `
		INSERT INTO employee (name)
		VALUES (:name)
		RETURNING id
	`

	// Используем sqlx.NamedQuery, чтобы подставить значения по тегам struct
	rows, err := r.db.NamedQuery(query, e)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&e.Id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) FindAll() (employees []Entity, err error) {
	query := "SELECT id, name, created_at, updated_at FROM employee ORDER BY id"
	err = r.db.Select(&employees, query)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *Repository) FilterByIDs(ids []int64) (employees []Entity, err error) {
	query, args, err := sqlx.In("SELECT * FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	err = r.db.Select(&employees, query, args...)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *Repository) DeleteById(id int64) (int64, error) {
	res, err := r.db.Exec("DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r *Repository) DeleteByIds(ids []int64) (int64, error) {
	query, args, err := sqlx.In("DELETE FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return 0, err
	}
	query = r.db.Rebind(query)
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r *Repository) BeginTransaction() (tx *sqlx.Tx, err error) {
	return r.db.Beginx()
}

func (r *Repository) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	err = tx.Get(
		&isExists,
		"SELECT EXISTS (SELECT 1 FROM employee WHERE name = $1)",
		name,
	)
	return isExists, err
}

func (r *Repository) CreateTx(tx *sqlx.Tx, request CreateRequest) (employeeId int64, err error) {
	var e = request.ToEntity()
	query := `INSERT INTO employee (name) VALUES (:name) RETURNING id`
	res, err := tx.NamedQuery(query, e)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	if res.Next() {
		if err := res.Scan(&employeeId); err != nil {
			return 0, err
		}
	}
	return employeeId, err
}
