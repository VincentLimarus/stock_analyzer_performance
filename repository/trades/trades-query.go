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
			id, 
			stock_code, 
			status, 
			lot, 
			avg_buy_price, 
			avg_sell_price, 
			dividen_amount, 
			created_at, 
			updated_at
		FROM trades
		WHERE EXTRACT(YEAR FROM created_at) = $1
		ORDER BY created_at ASC
	`

	GetAllListTradeQuery = `
		SELECT 
			id, 
			stock_code, 
			status, 
			lot, 
			avg_buy_price, 
			avg_sell_price, 
			dividen_amount, 
			created_at, 
			updated_at
		FROM trades
		ORDER BY created_at DESC
	`
)
