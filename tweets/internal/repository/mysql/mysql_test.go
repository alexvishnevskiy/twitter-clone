package mysql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"testing"
	"time"
)

func TestRepository_Put(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO Tweets").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "some content", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Put(ctx, model.UserId(1), "some content", nil, nil)
	if err != nil {
		t.Errorf("error was not expected while inserting tweet: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_DeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	mock.ExpectExec("DELETE FROM Tweets WHERE tweet_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the Delete method
	if err = repo.DeletePost(ctx, model.TweetId(1)); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetByTweet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	// Create rows to return
	rows := sqlmock.NewRows([]string{"user_id", "tweet_id", "retweet_id", "content", "media_url", "created_at"}).
		AddRow(1, 1, 2, "content", "url", time.Now().Format(layout))
	// Set expectation
	mock.ExpectQuery("^SELECT \\* FROM Tweets WHERE tweet_id IN \\(\\?\\)$").
		WithArgs(1).
		WillReturnRows(rows)

	// Call the get method
	res, err := repo.GetByTweet(ctx, model.TweetId(1))
	if err != nil {
		t.Errorf("error was not expected while getting tweet: %s", err)
	}
	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	testId := model.TweetId(1)
	userId := model.UserId(1)
	testUrl := "url"
	if res[0].TweetId == testId && res[0].UserId == userId && res[0].MediaUrl == &testUrl && res[0].Content == "content" {
		t.Fatal("unexpected results")
	}
}

func TestRepository_GetByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()

	// Create rows to return
	rows := sqlmock.NewRows([]string{"user_id", "tweet_id", "retweet_id", "content", "media_url", "created_at"}).
		AddRow(1, 1, 2, "content", "url", time.Now().Format(layout))
	// Set expectation
	mock.ExpectQuery("^SELECT \\* FROM Tweets WHERE user_id IN \\(\\?\\)$").
		WithArgs(1).
		WillReturnRows(rows)

	// Call the get method
	res, err := repo.GetByUser(ctx, model.UserId(1))
	if err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}
	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	testId := model.TweetId(1)
	userId := model.UserId(1)
	testUrl := "url"
	if res[0].TweetId == testId && res[0].UserId == userId &&
		res[0].MediaUrl == &testUrl && res[0].Content == "content" {
		t.Fatal("unexpected results")
	}
}
