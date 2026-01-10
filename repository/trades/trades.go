package trades

import (
	db "VincentLimarus/stock-analyzer-performance/model/db/trades"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Trade struct {
	db *sqlx.DB
}

func NewTrade(db *sqlx.DB) *Trade {
	return &Trade{
		db: db,
	}
}

type ITrade interface {
	UpsertTradeData(ctx context.Context, tx *sqlx.Tx, payload []db.Trade) error
}

func (r *Trade) UpsertTradeData(ctx context.Context, tx *sqlx.Tx, payload []db.Trade) error {
	if len(payload) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(payload)*7)
	valueClauses := make([]string, 0, len(payload))
	paramIndex := 1
	for _, trade := range payload {
		valueClauses = append(valueClauses,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				paramIndex, paramIndex+1, paramIndex+2,
				paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6))

		args = append(args,
			trade.SerialCode,
			trade.StockCode,
			trade.Status,
			trade.Lot,
			trade.AvgBuyPrice,
			trade.AvgSellPrice,
			trade.DividenAmount,
		)
		paramIndex += 7
	}

	query := fmt.Sprintf(UpsertTradeQuery, strings.Join(valueClauses, ", "))
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Inserted %d rows", rowsAffected)

	return nil
}
