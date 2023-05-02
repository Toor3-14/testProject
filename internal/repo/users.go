package repo
 
import (
	"github.com/jackc/pgx/v4/pgxpool"
	"errors"
	"context"
	"fmt"
)
 
type UserRepo struct {
	DB *pgxpool.Pool
}
func (r *UserRepo) InitDB() error {
	// Init users
	rows, err := r.DB.Query(context.Background(), `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE  table_schema = 'public'
		AND    table_name   = 'users'
	);`)
	if err != nil {
		return err
	}
	if rows.Next() == true {
		values, err := rows.Values()
		if err != nil {
			return err
		}

		if b, ok := values[0].(bool); ok == true && b == false {
			_, err := r.DB.Query(context.Background(), `CREATE TABLE users (
				id BIGSERIAL PRIMARY KEY,
				name CHARACTER VARYING(30) UNIQUE NOT NULL,
				age INTEGER NOT NULL
			);`)
			if err != nil {
				return err
			}
		}else if ok == false {
			return errors.New("Error while type assertion 'isExist Table' var from 'init users table'")
		}
	}
	return nil
}
func (r *UserRepo) Insert(name string, age int) (int, error) {
	query := fmt.Sprintf("INSERT INTO users (name,age) VALUES('%s',%d) RETURNING id", name, age)
	rows, err := r.DB.Query(context.Background(), query)
	if err != nil {
	  return 0, err
	}
	if rows.Next() == true {
		values, err := rows.Values()
		if err != nil {
			return 0, err
		}
		return int(values[0].(int64)), nil
	}else {
		return 0, errors.New("The username already exists")
	}
}
// GetById, insert data from db to model, decode model to json ...
// ... 