package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grpc-service-ref/internal/config"
	"grpc-service-ref/internal/domain/models"
	"grpc-service-ref/internal/storage"
	e "grpc-service-ref/pkg"
)

type Storage struct {
	db *gorm.DB
}

func NewDsn(dbCfg config.PostgresConfig) (string, error) {
	const op = "storage.postgresql.NewDsn"

	host := dbCfg.Host
	user := dbCfg.User
	password := dbCfg.Password
	name := dbCfg.DBName
	port := dbCfg.Port

	if host == "" || user == "" || password == "" || name == "" {
		return "", e.Err(op, storage.ErrNotSetDBParameter)
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	return dsn, nil
}

func New(dbCfg config.PostgresConfig) (*Storage, error) {
	const op = "storage.postgresql.New"

	dsn, err := NewDsn(dbCfg)

	if err != nil {
		return nil, e.Err(op, err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, e.Err(op, err)
	}

	err = db.AutoMigrate(&models.URL{})
	if err != nil {
		return nil, e.Err(op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	const op = "storage.postgresql.Close"

	db, err := s.db.DB()
	if err != nil {
		return e.Err(op, err)
	}

	err = db.Close()
	if err != nil {
		return e.Err(op, err)
	}

	return nil
}
