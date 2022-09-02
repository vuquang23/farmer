CREATE TABLE IF NOT EXISTS `spot_workers` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    symbol TEXT NOT NULL,
    unit_buy_allowed INT NOT NULL,
    unit_notional FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS `spot_trades` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    side VARCHAR(5) NOT NULL,
    binance_order_id INT NOT NULL,
    spot_worker_id INT NOT NULL,
    qty FLOAT NOT NULL,
    cummulative_quote_quantity FLOAT NOT NULL,
    ref INT,
    is_done BOOLEAN NOT NULL,
    unit_bought INT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

ALTER TABLE
    `spot_trades`
ADD
    FOREIGN KEY spot_worker_id REFERENCES `spot_workers` (id) ON DELETE CASCADE;