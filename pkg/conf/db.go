// Package conf contains functions to configure logging, the database and router
package conf

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	// database drivers
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	log "github.com/sirupsen/logrus"

	"github.com/epip-io/go-birthday-calendar/pkg/models"
)

// ConfigureDatabase performs DB migrations, including table and field creation
func ConfigureDatabase(cfg *models.DBConfig, logger *log.Entry) (*gorm.DB, error) {
	if cfg.Conn == "" {
		switch strings.ToLower(cfg.Engine) {
		case "mysql":
			if cfg.Port == 0 {
				cfg.Port = 3306
			}
			cfg.Conn = cfg.User + ":" + cfg.Pass + "@(" + cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port) + ")/" + cfg.Name + "?charset=utf8&parseTime=True&loc=Local"

		case "mssql":
			if cfg.Port == 0 {
				cfg.Port = 1433
			}
			cfg.Conn = "host=" + cfg.Host + " port=" + fmt.Sprintf("%d", cfg.Port) + " user=" + cfg.User + " password=" + cfg.Pass + " dbname=" + cfg.Name

		case "postgres":
			if cfg.Port == 0 {
				cfg.Port = 1433
			}
			cfg.Conn = "sqlserver://" + cfg.User + ":" + cfg.Pass + "@" + cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port) + "?database=" + cfg.Name

		case "sqlite3":

			cfg.Conn = cfg.Name
		}
	}

	logger.WithFields(log.Fields{
		"engine":      cfg.Engine,
		"conn_string": cfg.Conn,
	}).Info("connecting to database")
	db, err := gorm.Open(strings.ToLower(cfg.Engine), cfg.Conn)
	db.SetLogger(logger)
	if err != nil {
		return nil, err
	}

	logger.WithFields(log.Fields{
		"table": "persons",
	}).Debug("auto migrating table")
	db.AutoMigrate(&models.Person{})

	return db, nil
}
