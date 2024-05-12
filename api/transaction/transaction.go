package transaction

import (
	"database/sql"
	"net/http"
	"strconv"

	"time"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/constanst"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	ID              int64     `json:"id"`
	SpenderID       int       `json:"spender_id,omitempty" sql:"default:0"` //if default = 0 return nothing
	Date            time.Time `json:"date"`
	Amount          float32   `json:"amount"`
	Category        string    `json:"category"`
	TransactionType string    `json:"transaction_type,omitempty"`
	Note            string    `json:"note,omitempty"`
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
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("get successfully", zap.Int64("id", id))
	return c.JSON(http.StatusOK, ts)
}



func (h handler) Create(c echo.Context) error {
	if !h.flag.EnableCreateTransaction {
		return c.JSON(http.StatusForbidden, "create new transaction feature is disabled")
	}

	logger := mlog.L(c)

	var ts Transaction
	err := c.Bind(&ts)
	if err != nil {
		logger.Error(constanst.BadRequestBody, zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, ts.SpenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl).Scan(&lastInsertId)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	ts.ID = lastInsertId
	return c.JSON(http.StatusCreated, ts)
}

func (h handler) Update(c echo.Context) error {
	// if !h.flag.EnableUpdateTransaction {
	// 	return c.JSON(http.StatusForbidden, "update transaction feature is disabled")
	// }

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
		logger.Error(constanst.BadRequestBody, zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	_, err = h.db.ExecContext(ctx, "UPDATE transaction SET spender_id = $1, date = $2, amount = $3, category = $4, transaction_type = $5, note = $6, image_url = $7 WHERE id = $8", ts.SpenderID, ts.Date, ts.Amount, ts.Category, ts.TransactionType, ts.Note, ts.ImageUrl, id)
	if err != nil {
		logger.Error("exec error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("update successfully", zap.Int64("id", updateID))
	ts.ID = updateID
	return c.JSON(http.StatusOK, ts)
}

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT * FROM transaction`)
	if err != nil {
		logger.Error(constanst.QueryError, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var ts []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.SpenderID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl)
		if err != nil {
			logger.Error(constanst.ScanError, zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ts = append(ts, t)
	}

	return c.JSON(http.StatusOK, ts)
}
