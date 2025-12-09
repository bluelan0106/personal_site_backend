package database

import (
	"fmt"

	"personal_site/config"
	"personal_site/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initMySQLDB(dsn string) (*gorm.DB, error) {
	if len(dsn) >= 8 && dsn[:8] == "mysql://" {
		dsn = dsn[8:] // Remove "mysql://" prefix if present
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.YTDataAPITokenHistory{},
		&models.BattleCatLevel{},
		&models.Reurl{},
		&models.Post{},
		&models.Tag{},
		&models.Comment{},
		&models.Reaction{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %v", err)
	}

	return db, nil
}

func initSQLiteDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.YTDataAPITokenHistory{},
		&models.BattleCatLevel{},
		&models.Reurl{},
		&models.Post{},
		&models.Tag{},
		&models.Comment{},
		&models.Reaction{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %v", err)
	}

	return db, nil
}

func InitDB() (*gorm.DB, error) {

	dsn, err := config.GetVariableAsString("DATABASE_DSN")
	if err != nil {
		return nil, fmt.Errorf("error getting DATABASE_DSN: %v", err)
	}

	switch dbType(dsn) {
	case "mysql":
		return initMySQLDB(dsn)
	case "sqlite":
		return initSQLiteDB(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type")
	}
}

func dbType(dsn string) string {
	if len(dsn) >= len("sqlite") && dsn[:len("sqlite")] == "sqlite" {
		return "sqlite"
	}
	if len(dsn) > len("mysql") && dsn[:len("mysql")] == "mysql" {
		return "mysql"
	}

	if len(dsn) > 0 {
		if dsn[0] == '/' || dsn[0] == '.' || dsn[0] == ':' {
			return "sqlite"
		}
	}
	return "not supported"
}
