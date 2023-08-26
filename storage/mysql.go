package storage

import (
	"fmt"
	"time"

	"tempo/config"
	"tempo/helper"
	"tempo/repository/mysqlrepo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	mysql_migrate "github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetMySqlDb() *gorm.DB {
	dbName := config.Instance().DB.Database
	return MySqlDbConn(&dbName)
}

func MySqlDbConn(dbName *string) *gorm.DB {
	dbURL := getMysqlUrl(dbName)

	cfg := config.Instance()

	logrusLogger := gorm_logrus.New()
	logrusLogger.LogMode(logger.Silent)
	logrusLogger.Debug = false
	if cfg.DB.Debug {
		logrusLogger.Debug = true
		logrusLogger.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{
		CreateBatchSize: 500,
		Logger:          logrusLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dbURL))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dbURL))
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Instance().DB.MaxConnLifeTime))
	sqlDB.SetMaxOpenConns(config.Instance().DB.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.Instance().DB.MaxIdleConnections)

	return db
}

func CreateMySqlDb(dbName string) error {
	dbConn := MySqlDbConn(nil)
	defer func(dbConn *gorm.DB) {
		if sqlDB, err := dbConn.DB(); err != nil {
			logrus.Warnf("Error when get db connection: %s", err)
		} else {
			err = sqlDB.Close()
			if err != nil {
				logrus.WithError(err).Warning("Error when closing mysql db")
			}
		}
	}(dbConn)
	return dbConn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
}

func MigrateMysqlDb(dbConn *gorm.DB, migrationFolder *string, rollback bool, versionToForce int) error {
	dbConfig := config.Instance().DB

	var validMigrationFolder = dbConfig.Migrations.Path
	if helper.Val(migrationFolder) != "" {
		validMigrationFolder = *migrationFolder
	}

	if validMigrationFolder == "" {
		return fmt.Errorf("empty migration folder")
	}
	logrus.Infof("Migration folder: %s", validMigrationFolder)

	db, err := dbConn.DB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	driver, err := mysql_migrate.WithInstance(db, &mysql_migrate.Config{})
	if err != nil {
		logrus.WithError(err).Warning("Error when instantiating driver")
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+validMigrationFolder,
		dbConfig.Client,
		driver)
	if err != nil {
		logrus.WithError(err).Warning("Error when instantiating migrate")
		return err
	}
	if rollback {
		logrus.Info("About to Rolling back 1 step")
		err = m.Steps(-1)
	} else if versionToForce != -1 {
		logrus.Info(fmt.Sprintf("About to force version %d", versionToForce))
		err = m.Force(versionToForce)
	} else {
		logrus.Info("About to run migration")
		err = m.Up()
	}
	if err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}

	if sqlDB, err := db.DB(); err != nil {
		logrus.Warnf("Error when get db connection: %s", err)
	} else {
		err = sqlDB.Close()
		if err != nil {
			logrus.Warnf("Error when closing db: %s", err)
		}
	}
}

func TruncateNonRefTables(db *gorm.DB) error {
	models := []interface{}{
		mysqlrepo.User{},
	}
	for _, v := range models {
		err := db.Statement.Parse(v)
		if err != nil {
			return err
		}

		tableName := db.Statement.Schema.Table
		err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func getMysqlUrl(dbName *string) string {
	dbConfig := config.Instance().DB

	dbNameTmp := ""
	if dbName != nil {
		dbNameTmp = *dbName
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true&multiStatements=true", dbConfig.Username,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbNameTmp)
}
