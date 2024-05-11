package transactions

import (
	"database/sql"
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}
type Transaction struct {
	ID              int       `json:"id"`
	Date            string `json:"date"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	TransactionType string    `json:"trancsationType"`
	Note            string    `json:"note"`
	ImageURL        string    `json:"imageURL"`
	UserID          int       `json:"userID"`
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}
func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT id, date, amount,category,trancsation_type,note,image_url,user_id FROM transaction`)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var ts []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageURL, &t.UserID)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ts = append(ts, t)
	}

	return c.JSON(http.StatusOK, ts)
}
