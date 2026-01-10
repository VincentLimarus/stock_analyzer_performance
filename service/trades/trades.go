package trades

import (
	constant "VincentLimarus/stock-analyzer-performance/constant/tradeConst"
	db "VincentLimarus/stock-analyzer-performance/model/db/trades"
	dto "VincentLimarus/stock-analyzer-performance/model/dto/http/trades"
	"VincentLimarus/stock-analyzer-performance/repository"
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Trade struct {
	repoRegistry repository.IRegistry
}

func NewTrade(
	repoRegistry repository.IRegistry,
) *Trade {
	return &Trade{
		repoRegistry: repoRegistry,
	}
}

type ITrade interface {
	UpsertTradeData(ctx context.Context, payload []dto.TradeData) error
}

func (r *Trade) UpsertTradeData(ctx context.Context, payload []dto.TradeData) error {
	if len(payload) == 0 {
		return nil
	}

	txFunc := func(tx *sqlx.Tx) error {
		upsertPayloads := make([]db.Trade, 0, len(payload))

		for _, item := range payload {
			status := strings.ToUpper(strings.Split(item.Status, " ")[0])
			if status != constant.TradeClosedStatus {
				continue
			}

			dividendAmount := 0.0
			if item.DividenAmount != nil && *item.DividenAmount != 0.0 {
				dividendAmount = *item.DividenAmount
			}

			upsertPayloads = append(upsertPayloads, db.Trade{
				SerialCode:    item.SerialCode,
				Status:        item.Status,
				StockCode:     item.Code,
				Lot:           item.Lot,
				AvgBuyPrice:   item.AveragePrice.Buy,
				AvgSellPrice:  item.AveragePrice.Sell,
				DividenAmount: &dividendAmount,
			})
		}

		if len(upsertPayloads) == 0 {
			return nil
		}

		return r.repoRegistry.GetTrade().UpsertTradeData(ctx, tx, upsertPayloads)
	}

	return r.repoRegistry.GetTx().WithTx(ctx, txFunc, nil)
}
