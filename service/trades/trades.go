package trades

import (
	constant "VincentLimarus/stock-analyzer-performance/constant/tradeConst"
	db "VincentLimarus/stock-analyzer-performance/model/db/trades"
	dto "VincentLimarus/stock-analyzer-performance/model/dto/http/trades"
	"VincentLimarus/stock-analyzer-performance/repository"
	"context"
	"math"
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
	GetListTradeByYears(ctx context.Context, years []string) ([]dto.TradeData, error)
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

func (r *Trade) GetListTradeByYears(ctx context.Context, years []string) ([]dto.TradeData, error) {
	tradeRecords, err := r.repoRegistry.GetTrade().GetListTradeByYears(ctx, years)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TradeData, 0, len(tradeRecords))
	for _, record := range tradeRecords {
		var dividenAmount *float64
		if record.DividenAmount != nil {
			dividenAmount = record.DividenAmount
		}

		percentageChange := math.Round((record.AvgSellPrice - record.AvgBuyPrice) / record.AvgBuyPrice * 10000) / 100
		gross := (float64(record.AvgBuyPrice) * float64(record.Lot) * 100) - (float64(record.AvgSellPrice) * float64(record.Lot) * 100)
		fee := -1 * (float64(record.AvgBuyPrice) * float64(record.Lot) * 100 * constant.StockbitBrokerFeePercentage) - (float64(record.AvgSellPrice) * float64(record.Lot) * 100 * constant.StockbitBrokerSellFeePercentage)
		net := gross + fee

		result = append(result, dto.TradeData{
			SerialCode: record.SerialCode,
			Status:     record.Status,
			Code:       record.StockCode,
			Lot:        record.Lot,
			AveragePrice: dto.AveragePrice{
				Buy:  record.AvgBuyPrice,
				Sell: record.AvgSellPrice,
				PercentageChange: percentageChange,
			},
			GainLoss: dto.GainLoss{
				Gross: gross,
				Fee:   fee,
				Net:   net,
			},
			DividenAmount: dividenAmount,
		})
	}

	return result, nil
}