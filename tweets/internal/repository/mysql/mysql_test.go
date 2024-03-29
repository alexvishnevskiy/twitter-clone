package mysql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	_, _, err = repo.Put(ctx, types.UserId(1), "some content", nil, nil)
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
	if err = repo.DeletePost(ctx, types.TweetId(1)); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repository{db}
	ctx := context.Background()
	// what we want
	curTime := time.Now()
	mediaUrl := "url"
	retweetId := types.TweetId(2)
	want := model.Tweet{
		TweetId:   types.TweetId(1),
		UserId:    types.UserId(1),
		RetweetId: &retweetId,
		MediaUrl:  &mediaUrl,
		Content:   "content",
		CreatedAt: curTime,
	}

	testCases := []struct {
		name  string
		query string
	}{
		{
			name:  "GetByTweet",
			query: "^SELECT \\* FROM Tweets WHERE tweet_id IN \\(\\?\\)$",
		},
		{
			name:  "GetByUser",
			query: "^SELECT \\* FROM Tweets WHERE user_id IN \\(\\?\\)$",
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				// Create rows to return
				rows := sqlmock.NewRows([]string{"user_id", "tweet_id", "retweet_id", "content", "media_url", "created_at"}).
					AddRow(1, 1, 2, "content", "url", curTime.Format(layout))
				// Set expectation
				mock.ExpectQuery(tc.query).
					WithArgs(1).
					WillReturnRows(rows)

				var (
					res []model.Tweet
					err error
				)
				// Call the get method
				switch tc.name {
				case "GetByTweet":
					res, err = repo.GetByTweet(ctx, types.TweetId(1))
				case "GetByUser":
					res, err = repo.GetByUser(ctx, types.UserId(1))
				}

				if err != nil {
					t.Errorf("error was not expected while getting tweet: %s", err)
				}
				// We make sure that all expectations were met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
				// compare structures
				if diff := cmp.Diff(want, res[0], cmpopts.IgnoreFields(model.Tweet{}, "CreatedAt")); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}
