ALTER TABLE
    `spot_trades` DROP INDEX Sym_Side_IsDone_Time;

DROP TABLE IF EXISTS `spot_trades`;

DROP TABLE IF EXISTS `spot_trading_pairs`;