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

///

// package transaction

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreateTransaction(t *testing.T) {

// 	t.Run("create transaction succesfully when feature toggle is enable", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		date := "2024-04-30T09:00:00.000Z"
// 		parsedDate, _ := time.Parse(time.RFC3339, date)

// 		ts := Transaction{
// 			Date:            parsedDate,
// 			SenderID:        1,
// 			Amount:          1500,
// 			Category:        "Food",
// 			TransactionType: "expense",
// 			Note:            "Lunch",
// 			ImageUrl:        "https://example.com/image1.jpg",
// 		}

// 		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
// 		mock.ExpectQuery(cStmt).WithArgs(ts.SenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl).WillReturnRows(row)
// 		cfg := config.FeatureFlag{EnableCreateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Create(c)

// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 		assert.JSONEq(t, `{
// 			"id": 1,
// 			"sender_id": 1,
// 			"date": "2024-04-30T09:00:00Z",
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`, rec.Body.String())
// 	})

// 	t.Run("create transaction failed when feature toggle is disable", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		cfg := config.FeatureFlag{EnableCreateTransaction: false}

// 		h := New(cfg, nil)
// 		err := h.Create(c)

// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusForbidden, rec.Code)
// 	})

// 	t.Run("create transaction failed when bad request body", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ bad request body }`))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		cfg := config.FeatureFlag{EnableCreateTransaction: true}

// 		h := New(cfg, nil)
// 		err := h.Create(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("get transaction successfully", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		date := "2024-04-30T09:00:00.000Z"
// 		parsedDate, _ := time.Parse(time.RFC3339, date)

// 		ts := Transaction{
// 			ID:              1,
// 			SenderID:        1,
// 			Date:            parsedDate,
// 			Amount:          1500,
// 			Category:        "Food",
// 			TransactionType: "expense",
// 			Note:            "Lunch",
// 			ImageUrl:        "https://example.com/image1.jpg",
// 		}

// 		row := sqlmock.NewRows([]string{"id", "sender_id", "date", "amount", "category", "transaction_type", "note", "image_url"}).AddRow(ts.ID, ts.SenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl)
// 		mock.ExpectQuery("SELECT * FROM transaction WHERE id = $1").WithArgs(ts.ID).WillReturnRows(row)

// 		cfg := config.FeatureFlag{EnableCreateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Get(c)

// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		assert.JSONEq(t, `{
// 			"id": 1,
// 			"sender_id": 1,
// 			"date": "2024-04-30T09:00:00Z",
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`, rec.Body.String())
// 	})

// 	t.Run("get transaction failed when bad request id", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("bad id")

// 		cfg := config.FeatureFlag{EnableCreateTransaction: true}

// 		h := New(cfg, nil)
// 		err := h.Get(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("get transaction failed when query error", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		ts := Transaction{
// 			ID: 1,
// 		}

// 		mock.ExpectQuery("SELECT * FROM transaction WHERE id = $1").WithArgs(ts.ID).WillReturnError(assert.AnError)

// 		cfg := config.FeatureFlag{EnableCreateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Get(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction failed when feature toggle is disable", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: false}

// 		h := New(cfg, nil)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusForbidden, rec.Code)
// 	})

// 	t.Run("update transaction failed when bad request body", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{ bad request body }`))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, nil)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction failed when query error", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		ts := Transaction{
// 			ID: 1,
// 		}

// 		mock.ExpectQuery("SELECT * FROM transaction WHERE id = $1").WithArgs(ts.ID).WillReturnError(assert.AnError)

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction failed when wrong id", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction successfully", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		date := "2024-04-30T09:00:00.000Z"
// 		parsedDate, _ := time.Parse(time.RFC3339, date)

// 		ts := Transaction{
// 			ID:              1,
// 			Date:            parsedDate,
// 			SenderID:        1,
// 			Amount:          1500,
// 			Category:        "Food",
// 			TransactionType: "expense",
// 			Note:            "Lunch",
// 			ImageUrl:        "https://example.com/image1.jpg",
// 		}

// 		mock.ExpectExec("UPDATE transaction SET sender_id = $1, date = $2, amount = $3, category = $4, transaction_type = $5, note = $6, image_url = $7 WHERE id = $8").WithArgs(ts.SenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl, ts.ID).WillReturnResult(sqlmock.NewResult(1, 1))

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	})

// 	t.Run("update transaction failed when query error", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")

// 		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		date := "2024-04-30T09:00:00.000Z"
// 		parsedDate, _ := time.Parse(time.RFC3339, date)

// 		ts := Transaction{
// 			ID:              1,
// 			Date:            parsedDate,
// 			SenderID:        1,
// 			Amount:          1500,
// 			Category:        "Food",
// 			TransactionType: "expense",
// 			Note:            "Lunch",
// 			ImageUrl:        "https://example.com/image1.jpg",
// 		}

// 		mock.ExpectExec("UPDATE transaction SET sender_id = $1, date = $2, amount = $3, category = $4, transaction_type = $5, note = $6, image_url = $7 WHERE id = $8").WithArgs(ts.SenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl, ts.ID).WillReturnError(assert.AnError)

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})


// 	t.Run("update transaction failed when wrong id", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("bad id")

// 		db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 		defer db.Close()

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, db)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction failed when bad request body", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{ bad request body }`))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")
// 		cfg := config.FeatureFlag{EnableUpdateTransaction: true}

// 		h := New(cfg, nil)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("update transaction failed when feature toggle is disable", func(t *testing.T) {
// 		e := echo.New()
// 		defer e.Close()

// 		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{
// 			"date": "2024-04-30T09:00:00.000Z",
// 			"sender_id":1,
// 			"amount": 1500,
// 			"category": "Food",
// 			"transaction_type": "expense",
// 			"note": "Lunch",
// 			"image_url": "https://example.com/image1.jpg"
// 		}`))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/transactions/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues("1")

// 		cfg := config.FeatureFlag{EnableUpdateTransaction: false}

// 		h := New(cfg, nil)
// 		err := h.Update(c)

// 		assert.NoError(t, err)
// 	})


// }
