package db

import "time"

type Trade struct {
	UUID          string    `db:"uuid"`
	SerialCode    string    `db:"serial_code"`
	Status        string    `db:"status"`
	StockCode     string    `db:"stock_code"`
	Lot           int       `db:"lot"`
	AvgBuyPrice   float64   `db:"avg_buy_price"`
	AvgSellPrice  float64   `db:"avg_sell_price"`
	DividenAmount *float64  `db:"dividen_amount"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
