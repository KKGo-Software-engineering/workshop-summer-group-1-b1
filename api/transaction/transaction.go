package transaction

import (
	"database/sql"
	"net/http"
	"strconv"

	"time"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	ID              int64     `json:"id"`
	SpenderID       int       `json:"spender_id"`
	Date            time.Time `json:"date"`
	Amount          float32   `json:"amount"`
	Category        string    `json:"category"`
	TransactionType string    `json:"transaction_type"`
	Note            string    `json:"note"`
	ImageUrl        string    `json:"image_url"`
}

type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

const (
	cStmt = `INSERT INTO transaction ( spender_id , date , amount , category, transaction_type, note, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
)

func (h handler) Get(c echo.Context) error {
	logger := mlog.L(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error("bad request id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	var ts Transaction
	err = h.db.QueryRowContext(ctx, "SELECT * FROM transaction WHERE id = $1", id).Scan(&ts.ID, &ts.SpenderID, &ts.Date, &ts.Amount, &ts.Category, &ts.TransactionType, &ts.Note, &ts.ImageUrl)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("get successfully", zap.Int64("id", id))
	return c.JSON(http.StatusOK, ts)
}

func (h handler) GetTransactions(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		logger.Error("id is non-int")
		return c.JSON(http.StatusBadRequest, "id is non-int")
	}

	rows, err := h.db.QueryContext(ctx, `SELECT id, sender_id, date, amount, category, transaction_type, note, image_url FROM transaction WHERE sender_id=$1`, id)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var ts []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.SpenderID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ts = append(ts, t)
	}

	var totalIncome, totalExpenses, currentBalance float32
	for _, t := range ts {
		if t.TransactionType == "income" {
			totalIncome += t.Amount
		} else {
			totalExpenses += t.Amount
		}
	}
	currentBalance = totalIncome - totalExpenses

	currentPage := 1
	totalPages := len(ts)/10 + 1
	perPage := 10

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transections": ts,
		"summary": map[string]float32{
			"total_income":    totalIncome,
			"total_expenses":  totalExpenses,
			"current_balance": currentBalance,
		},

		"pagination": map[string]int{
			"current_page": currentPage,
			"total_pages":  totalPages,
			"per_page":     perPage,
		},
	})
}

func (h handler) Create(c echo.Context) error {
	if !h.flag.EnableCreateTransaction {
		return c.JSON(http.StatusForbidden, "create new transaction feature is disabled")
	}

	logger := mlog.L(c)

	var ts Transaction
	err := c.Bind(&ts)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, ts.SpenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	ts.ID = lastInsertId
	return c.JSON(http.StatusCreated, ts)
}

func (h handler) Update(c echo.Context) error {
	if !h.flag.EnableUpdateTransaction {
		return c.JSON(http.StatusForbidden, "update transaction feature is disabled")
	}

	logger := mlog.L(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	updateID := id
	if err != nil {
		logger.Error("bad request id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var ts Transaction
	err = c.Bind(&ts)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	_, err = h.db.ExecContext(ctx, "UPDATE transaction SET spender_id = $1, date = $2, amount = $3, category = $4, transaction_type = $5, note = $6, image_url = $7 WHERE id = $8", ts.SpenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl, id)
	if err != nil {
		logger.Error("exec error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("update successfully", zap.Int64("id", updateID))
	return c.JSON(http.StatusOK, ts)
}

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT * FROM transaction`)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var ts []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.SpenderID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ts = append(ts, t)
	}

	return c.JSON(http.StatusOK, ts)
}
