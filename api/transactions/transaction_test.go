package transactions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTransation(t *testing.T) {
	t.Run("get all transaction succesfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"ID", "Date", "Amount", "Category", "TransactionType", "Note", "ImageURL", "UserID"}).
			AddRow(1, "2024-05-11T16:47:48.6967838+07:00", 123.45, "Food", "Expense", "Dinner with friends", "http://example.com/receipt.jpg", 100).
			AddRow(2, "2024-05-11T16:47:48.6967838+07:00", 150.00, "Transport", "Expense", "Monthly train ticket", "http://example.com/ticket.jpg", 100)
		mock.ExpectQuery(`SELECT id, date, amount,category,trancsation_type,note,image_url,user_id FROM transaction`).WillReturnRows(rows)
		cfg := config.FeatureFlag{EnableCreateSpender: false}

		h := New(cfg, db)
		err := h.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)



		assert.JSONEq(t, `[{"id":1,"date":"2024-05-11T16:47:48.6967838+07:00","amount":123.45,"category":"Food","trancsationType":"Expense","note":"Dinner with friends","imageURL":"http://example.com/receipt.jpg","userID":100},
		{"id":2,"date":"2024-05-11T16:47:48.6967838+07:00","amount":150,"category":"Transport","trancsationType":"Expense","note":"Monthly train ticket","imageURL":"http://example.com/ticket.jpg","userID":100}]`, rec.Body.String())
	})

	t.Run("get all transaction failed on database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(`SELECT id, date, amount,category,trancsation_type,note,image_url,user_id FROM transaction`).WillReturnError(assert.AnError)

		h := New(config.FeatureFlag{}, db)
		err := h.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
