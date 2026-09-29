package trades

var (
	UpsertTradeQuery = `
		INSERT INTO trades (
			serial_code, stock_code, status, lot, avg_buy_price, 
			avg_sell_price, dividen_amount
		) VALUES %s
		ON CONFLICT (serial_code) 
		DO UPDATE SET
			stock_code = EXCLUDED.stock_code,
			status = EXCLUDED.status,
			lot = EXCLUDED.lot,
			avg_buy_price = EXCLUDED.avg_buy_price,
			avg_sell_price = EXCLUDED.avg_sell_price,
			dividen_amount = EXCLUDED.dividen_amount,
			updated_at = NOW()
	`

	GetListTradeByYearQuery = `
		SELECT 
			uuid, 
			serial_code,
			stock_code, 
			status, 
			lot, 
			avg_buy_price, 
			avg_sell_price, 
			dividen_amount, 
			created_at, 
			updated_at
		FROM trades
		WHERE serial_code ~ '^[0-9]{4}-'
		AND split_part(serial_code, '-', 1) = ANY($1)
		ORDER BY serial_code DESC;
	`
)
