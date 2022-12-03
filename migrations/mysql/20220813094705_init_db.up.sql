CREATE TABLE IF NOT EXISTS `spot_workers` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(16) NOT NULL,
    unit_buy_allowed INT NOT NULL,
    unit_notional FLOAT NOT NULL,
    capital FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS `spot_trades` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(16) NOT NULL,
    side VARCHAR(8) NOT NULL,
    binance_order_id INT NOT NULL,
    spot_worker_id INT NOT NULL,
    qty VARCHAR(64) NOT NULL,
    cummulative_quote_qty VARCHAR(64) NOT NULL,
    price FLOAT NOT NULL,
    ref INT,
    is_done BOOLEAN NOT NULL,
    unit_bought INT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

ALTER TABLE
    `spot_trades`
ADD
    FOREIGN KEY (spot_worker_id) REFERENCES `spot_workers` (id) ON DELETE CASCADE;

CREATE INDEX Sym_Side_IsDone_Time ON `spot_trades`(spot_worker_id, side, is_done, created_at);

CREATE TABLE IF NOT EXISTS `history_spot_trades` (
    id INT PRIMARY KEY,
    symbol VARCHAR(16) NOT NULL,
    side VARCHAR(8) NOT NULL,
    binance_order_id INT NOT NULL,
    qty VARCHAR(64) NOT NULL,
    cummulative_quote_qty VARCHAR(64) NOT NULL,
    price FLOAT NOT NULL,
    ref INT,
    unit_bought INT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
