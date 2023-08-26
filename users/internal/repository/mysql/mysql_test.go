package mysql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/users/pkg/model"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestRepository_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	user := model.User{
		1,
		"user",
		"alex",
		"ivanov",
		"email",
		"password",
	}

	mock.ExpectExec("INSERT INTO User").
		WithArgs(user.UserId, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Register(ctx, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		t.Errorf("Error was not expecting while register: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	mock.ExpectExec("DELETE FROM User WHERE user_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err = repo.Delete(ctx, types.UserId(1)); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_RetrievePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	mail := "somemail@mail.com"
	password := "password"
	row := sqlmock.NewRows([]string{"password"}).
		AddRow(password)

	mock.ExpectQuery("SELECT password FROM User WHERE email = ?").
		WithArgs(mail).
		WillReturnRows(row)

	_, retrievedPassword, err := repo.RetrievePassword(ctx, mail)
	if err != nil {
		t.Errorf("Failed to retrieve password: %s", err)
	}
	if diff := cmp.Diff(password, retrievedPassword); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	userData := model.User{UserId: 1, Email: "somemail@gmail.com"}

	mock.ExpectExec("UPDATE User SET email = 'somemail@gmail.com' WHERE user_id = ?").
		WithArgs(userData.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(ctx, userData)
	if err != nil {
		t.Errorf("Error updating data table: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
