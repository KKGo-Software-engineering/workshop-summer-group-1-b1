package spender

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/constanst"
	"github.com/KKGo-Software-engineering/workshop-summer/api/transaction"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Spender struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

const (
	cStmt = `INSERT INTO spender (name, email) VALUES ($1, $2) RETURNING id;`
)

func (h handler) Create(c echo.Context) error {
	if !h.flag.EnableCreateSpender {
		return c.JSON(http.StatusForbidden, "create new spender feature is disabled")
	}

	logger := mlog.L(c)
	ctx := c.Request().Context()
	var sp Spender
	err := c.Bind(&sp)
	if err != nil {
		logger.Error(constanst.NonIntError, zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, sp.Name, sp.Email).Scan(&lastInsertId)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	sp.ID = lastInsertId
	return c.JSON(http.StatusCreated, sp)
}

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT id, name, email FROM spender`)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var sps []Spender
	for rows.Next() {
		var sp Spender
		err := rows.Scan(&sp.ID, &sp.Name, &sp.Email)
		if err != nil {
			logger.Error(constanst.ScanError, zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sps = append(sps, sp)
	}

	return c.JSON(http.StatusOK, sps)
}

func (h handler) Get(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		logger.Error(constanst.NonIntError)
		return c.JSON(http.StatusBadRequest, constanst.NonIntError)
	}

	row := h.db.QueryRowContext(ctx, `SELECT id, name, email FROM spender WHERE id=$1`, id)
	if row.Err() != nil {
		logger.Error(constanst.QueryError, zap.Error(row.Err()))
		return c.JSON(http.StatusNotFound, row.Err())
	}

	var sp Spender
	err := row.Scan(&sp.ID, &sp.Name, &sp.Email)
	if err != nil {
		logger.Error(constanst.ScanError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, sp)
}

func (h handler) Update(c echo.Context) error {
	if !h.flag.EnableUpdateSpender {
		return c.JSON(http.StatusForbidden, "update spender feature is disabled")
	}

	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		logger.Error(constanst.NonIntError)
		return c.JSON(http.StatusBadRequest, constanst.NonIntError)
	}

	var sp Spender
	err := c.Bind(&sp)
	if err != nil {
		logger.Error(constanst.BadRequestBody, zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = h.db.ExecContext(ctx, `UPDATE spender SET name=$1, email=$2 WHERE id=$3`, sp.Name, sp.Email, id)
	if err != nil {
		logger.Error("update error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("update successfully", zap.String("id", id))
	return c.JSON(http.StatusOK, "update successfully")
}

func (h handler) GetTransactions(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		logger.Error(constanst.NonIntError)
		return c.JSON(http.StatusBadRequest, constanst.NonIntError)
	}

	rows, err := h.db.QueryContext(ctx, `SELECT id, spender_id, date, amount, category, transaction_type, note, image_url FROM transaction WHERE spender_id=$1`, id)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var ts []transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		err := rows.Scan(&t.ID, &t.SpenderID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl)
		if err != nil {
			logger.Error(constanst.ScanError, zap.Error(err))
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

func (h handler) GetSummary(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		logger.Error(constanst.NonIntError)
		return c.JSON(http.StatusBadRequest, constanst.NonIntError)
	}

	rows, err := h.db.QueryContext(ctx, `SELECT amount, transaction_type FROM transaction WHERE spender_id=$1`, id)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var totalIncome, totalExpenses, currentBalance float64
	for rows.Next() {
		var amount float64
		var transactionType string
		err := rows.Scan(&amount, &transactionType)
		if err != nil {
			logger.Error(constanst.ScanError, zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if transactionType == "income" {
			totalIncome += amount
		} else {
			totalExpenses += amount
		}
	}
	currentBalance = totalIncome - totalExpenses
	return c.JSON(http.StatusOK, map[string]interface{}{
		"summary": map[string]float64{
			"total_income":    totalIncome,
			"total_expenses":  totalExpenses,
			"current_balance": currentBalance,
		},
	})
}
