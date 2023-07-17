package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	_ "github.com/go-sql-driver/mysql"
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

func (r *Repository) Follow(ctx context.Context, userId types.UserId, followId types.UserId) error {
	// when user follows someone, insert new row
	row, err := r.db.ExecContext(
		ctx,
		"INSERT INTO Followers (user_id, following_id) VALUES (?, ?)", userId, followId,
	)
	if err != nil {
		return err
	}

	// if no rows are affected, return error
	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Unfollow(ctx context.Context, userId types.UserId, followId types.UserId) error {
	// when user unfollows someone, delete row
	row, err := r.db.ExecContext(ctx, "DELETE FROM Followers WHERE user_id = ? AND following_id = ?", userId, followId)
	if err != nil {
		return err
	}

	// if no rows are affected, return error
	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

// helper function to retrieve users
func get(ctx context.Context, r *Repository, userId types.UserId, retrieveType int) ([]types.UserId, error) {
	var query string
	// GetUserFollowers or GetFollowingUser
	switch retrieveType {
	case 0:
		query = "SELECT following_id FROM Followers WHERE user_id = ?"
	case 1:
		query = "SELECT user_id FROM Followers WHERE following_id = ?"
	}
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// retrieve all rows -> all users
	var res []types.UserId
	// iterate over result
	for rows.Next() {
		var id types.UserId
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("there is no data for your request")
	}
	return res, nil
}

func (r *Repository) GetUserFollowers(ctx context.Context, userId types.UserId) ([]types.UserId, error) {
	followers, err := get(ctx, r, userId, 0)
	return followers, err
}

func (r *Repository) GetFollowingUser(ctx context.Context, userId types.UserId) ([]types.UserId, error) {
	followers, err := get(ctx, r, userId, 1)
	return followers, err
}
