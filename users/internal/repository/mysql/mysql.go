package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/users/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func New(driverName string, dataSourceName string) (*Repository, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

func (r *Repository) Register(
	ctx context.Context,
	nickname string,
	firstname string,
	lastname string,
	email string,
	password string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO User (nickname, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)",
		nickname, firstname, lastname, email, password,
	)
	return err
}

// outputs password for email address
func (r *Repository) RetrievePassword(
	ctx context.Context,
	email string,
) (string, error) {
	var password string
	rows, err := r.db.QueryContext(ctx, "SELECT password FROM User WHERE email = ?", email)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	userid types.UserId,
) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM User WHERE user_id = ?", userid)
	return err
}

func (r *Repository) Update(
	ctx context.Context,
	userData model.User,
) error {
	var conditions []string

	// iterate over fields
	v := reflect.ValueOf(userData)
	for i := 0; i < v.NumField(); i++ {
		length := len(v.Type().Field(i).Tag)
		fieldName := string(v.Type().Field(i).Tag[6 : length-1])
		fieldValue := v.Field(i).Interface()
		if fieldValue != "" && fieldName != "user_id" {
			conditions = append(conditions, fmt.Sprintf("%s = '%s'", fieldName, fieldValue))
		}
	}

	setStatement := strings.Join(conditions, ", ")
	execStatement := fmt.Sprintf("UPDATE User SET %s WHERE user_id = ?", setStatement)
	_, err := r.db.ExecContext(ctx, execStatement, userData.UserId)
	return err
}
