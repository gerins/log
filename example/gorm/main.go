package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gerins/log"
	gormLogger "github.com/gerins/log/extension/gorm"
	middlewareLog "github.com/gerins/log/middleware/echo"
)

func main() {
	log.Init() // Using default configuration
	e := echo.New()
	db := initGormDatabase()

	// Init logging middleware
	e.Use(middlewareLog.SetLogRequest())                       // Mandatory
	e.Use(middleware.BodyDump(middlewareLog.SaveLogRequest())) // Mandatory

	// Init handler
	e.GET("/", func(c echo.Context) error {
		// Get context from echo locals.
		ctx := c.Get("ctx").(context.Context)

		var email string
		db.WithContext(ctx).Table("users").Select("email").Where("id = ?", 176).Scan(&email)

		return c.String(http.StatusOK, "Hello, visit log folder for more detail!")
	})

	e.Start("localhost:8080")
}

// initGormDatabase to PostgreSQL Database
func initGormDatabase() *gorm.DB {
	cfg := struct {
		Host         string
		Port         int
		User         string
		Pass         string
		DatabaseName string
		Pool         struct {
			MaxIdleConn     int
			MaxOpenConn     int
			MaxConnLifetime int // Minutes
		}
	}{
		Host:         "localhost",
		Port:         5434,
		User:         "root",
		Pass:         "root",
		DatabaseName: "default",
		Pool: struct {
			MaxIdleConn     int
			MaxOpenConn     int
			MaxConnLifetime int
		}{
			MaxIdleConn:     50,
			MaxOpenConn:     100,
			MaxConnLifetime: 30,
		},
	}

	address := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Pass, cfg.DatabaseName, cfg.Port)

	gormCfg := &gorm.Config{
		Logger: gormLogger.Default, // Assign new logger
	}

	db, err := gorm.Open(postgres.Open(address), gormCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	sqlDB.SetMaxIdleConns(cfg.Pool.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.Pool.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Pool.MaxConnLifetime) * time.Minute)

	log.Info("GormDB : Successfully Connected to Database")
	return db
}
