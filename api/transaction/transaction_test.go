package transaction

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {

	t.Run("create transaction succesfully when feature toggle is enable", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"date": "2024-04-30T09:00:00.000Z",
			"amount": 1500,
			"category": "Food",
			"transaction_type": "expense",
			"note": "Lunch",
			"image_url": "https://example.com/image1.jpg"
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		date := "2024-04-30T09:00:00.000Z"
		parsedDate, _ := time.Parse(time.RFC3339, date)

		ts := Transaction{
			Date:            parsedDate,
			Amount:          1500,
			Category:        "Food",
			TransactionType: "expense",
			Note:            "Lunch",
			ImageUrl:        "https://example.com/image1.jpg",
		}

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(cStmt).WithArgs(ts.Date,ts.Amount,ts.Category,ts.TransactionType,ts.Note,ts.ImageUrl).WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateTransaction: true}

		h := New(cfg, db)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"date": "2024-04-30T09:00:00Z",
			"amount": 1500,
			"category": "Food",
			"transaction_type": "expense",
			"note": "Lunch",
			"image_url": "https://example.com/image1.jpg"
		}`, rec.Body.String())
	})
}
