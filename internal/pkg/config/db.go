package config

import "fmt"

//DBConfig contains all configuration of database (MySQL)
type DBConfig struct {
	Host         string `default:"localhost"`
	Port         int    `default:"3306"`
	DBName       string `default:"ks_setting"`
	User         string `default:"user"`
	Password     string `default:"password"`
	ConnLifeTime int    `default:"300"`
	ConnTimeOut  int    `default:"30"`
	MaxIdleConns int    `default:"10"`
	MaxOpenConns int    `default:"80"`
	LogLevel     int    `default:"1"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.ConnTimeOut,
	)
}

func (c *DBConfig) MigrationSource() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}
