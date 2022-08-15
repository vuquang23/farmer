CREATE TABLE IF NOT EXISTS `spot_trading_pairs` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    symbol TEXT NOT NULL,
    buy_count_allowed INT NOT NULL,
    buy_notional FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS `spot_trades` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    binance_order_id INT NOT NULL,
    spot_trading_pair_id INT NOT NULL,
    open_qty FLOAT NOT NULL,
    open_price FLOAT NOT NULL,
    close_qty FLOAT NOT NULL,
    close_price FLOAT NOT NULL,
    close_count INT NOT NULL,
    is_done BOOLEAN NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

ALTER TABLE
    `spot_trades`
ADD
    FOREIGN KEY spot_trading_pair_id REFERENCES `spot_trading_pairs` (id) ON DELETE CASCADE;