package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	cfg "farmer/internal/pkg/config"
)

type Migration interface {
	MigrateUp(up int) error
	MigrateDown(down int) error
	Close() (error, error)
}

type MySQLMigration struct {
	m *migrate.Migrate
}

func NewMigration(dir string) (Migration, error) {
	db, _ := sql.Open("mysql", cfg.Instance().DB.MigrationSource())
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		dir,
		"mysql",
		driver,
	)

	if err != nil {
		return nil, err
	}

	return &MySQLMigration{
		m,
	}, nil

}

func (t *MySQLMigration) MigrateUp(up int) error {
	if up == 0 {
		fmt.Println("haa ")
		return t.m.Up()
	} else {
		return t.m.Steps(up)
	}
}

func (t *MySQLMigration) MigrateDown(down int) error {
	if down == 0 {
		return t.m.Down()
	} else {
		return t.m.Steps(-down)
	}
}

func (t *MySQLMigration) Close() (error, error) {
	return t.m.Close()
}
