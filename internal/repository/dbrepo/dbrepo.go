package dbrepo

import (
	"database/sql"
	"github.com/sokolovss/BNBsite/internal/config"
	"github.com/sokolovss/BNBsite/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgreRepo(conn *sql.DB, config *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		config,
		conn,
	}

}
