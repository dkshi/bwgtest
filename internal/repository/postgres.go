package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	DBName   string
	Port     string
	Username string
	Password string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=%s", cfg.Host, cfg.DBName, cfg.Port, cfg.Username, cfg.Password, cfg.SSLMode)
	for {
		db, err := sqlx.Connect("postgres", connStr)
		if err == nil {
			return db, nil
		}
		logrus.Printf("Failed to connect to the database: %v. Retrying...", err)
		time.Sleep(5 * time.Second)
	}
}

func NewListener(cfg Config, name string) (*pq.Listener, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.DBName, cfg.Port, cfg.Username, cfg.Password, cfg.SSLMode)

	lis := pq.NewListener(connStr, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logrus.Printf("error listening: %s", err.Error())
		}
	})

	err := lis.Listen(name)
	if err != nil {
		return &pq.Listener{}, err
	}

	return lis, nil
}
