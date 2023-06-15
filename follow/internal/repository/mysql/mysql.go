package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/follow/pkg/model"
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

func (r *Repository) Follow(ctx context.Context, userId model.UserId, followId model.UserId) error {
	row, err := r.db.ExecContext(
		ctx,
		"INSERT INTO Followers (user_id, following_id) VALUES (?, ?)", userId, followId,
	)
	if err != nil {
		return err
	}

	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Unfollow(ctx context.Context, userId model.UserId, followId model.UserId) error {
	row, err := r.db.ExecContext(ctx, "DELETE FROM Followers WHERE user_id = ? AND following_id = ?", userId, followId)
	if err != nil {
		return err
	}

	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func get(ctx context.Context, r *Repository, userId model.UserId, retrieveType int) ([]model.UserId, error) {
	var query string
	switch retrieveType {
	case 0:
		query = "SELECT user_id FROM Followers WHERE user_id = ?"
	case 1:
		query = "SELECT following_id FROM Followers WHERE following_id = ?"
	}
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.UserId
	// iterate over result
	for rows.Next() {
		var id model.UserId
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

func (r *Repository) GetUserFollowers(ctx context.Context, userId model.UserId) ([]model.UserId, error) {
	followers, err := get(ctx, r, userId, 0)
	return followers, err
}

func (r *Repository) GetFollowingUser(ctx context.Context, userId model.UserId) ([]model.UserId, error) {
	followers, err := get(ctx, r, userId, 1)
	return followers, err
}
