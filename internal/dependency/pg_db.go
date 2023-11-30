package dependency

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lil-oren/rest/internal/constant"

	"github.com/jmoiron/sqlx"
)

func NewPGDB(config Config, logger Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(constant.ConnectionStringTemplate,
		config.PostgreDB.DBHost,
		config.PostgreDB.DBUser,
		config.PostgreDB.DBPass,
		config.PostgreDB.DBName,
		config.PostgreDB.DBPort,
	))

	if err != nil {
		logger.Fatalf("Failed connect to DB %v", err)
		return nil, err
	}

	logger.Infof("Successfully connect to database", nil)

	return db, nil
}
