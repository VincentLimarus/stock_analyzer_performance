CREATE TABLE trades (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    serial_code VARCHAR(100) UNIQUE NOT NULL,
    -- Core fields (indexed & queryable)
    stock_code VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    lot INT NOT NULL CHECK (lot > 0),
    
    -- Average Price
    avg_buy_price DECIMAL(15,2) NOT NULL,
    avg_sell_price DECIMAL(15,2) NOT NULL,
    dividen_amount DECIMAL(15,2) DEFAULT 0,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
