package transaction

import (
	"database/sql"
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

type Transaction struct {
	ID              int64     `json:"id"`
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
	cStmt = `INSERT INTO transaction (date , amount , category, transaction_type, note, image_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
)

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
	err = h.db.QueryRowContext(ctx, cStmt, ts.Date ,ts.Amount, ts.Category , ts.TransactionType , ts.Note, ts.ImageUrl).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	ts.ID = lastInsertId
	return c.JSON(http.StatusCreated, ts)
}
