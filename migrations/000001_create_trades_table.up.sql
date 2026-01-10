CREATE TABLE trades (
    id SERIAL PRIMARY KEY,
    
    -- Core fields (indexed & queryable)
    stock_code VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    lot INT NOT NULL CHECK (lot > 0),
    
    -- Average Price
    avg_buy_price DECIMAL(15,2) NOT NULL,
    avg_sell_price DECIMAL(15,2),
    percentage_change DECIMAL(10,4),
    
    -- Gain/Loss
    gross_gain_loss DECIMAL(15,2),
    fee DECIMAL(15,2) NOT NULL DEFAULT 0,
    net_gain_loss DECIMAL(15,2),
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
