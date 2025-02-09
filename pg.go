package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHandles struct {
	pool *pgxpool.Pool
}

func NewPostgresHandlesProvider(config *pgxpool.Config) (*PostgresHandles, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return &PostgresHandles{}, err
	}

	pg := &PostgresHandles{pool}

	healthy, status := pg.IsHealthy()

	if !healthy {
		return &PostgresHandles{}, errors.New(status)
	}

	return pg, nil
}

// TODO: Should pass in request context?
func (pg *PostgresHandles) GetDecentralizedIDForHandle(handle Handle) (DecentralizedID, error) {
	connection, err := pg.pool.Acquire(context.Background())

	defer connection.Release()

	if err != nil {
		return "", err
	}

	var did DecentralizedID

	err = connection.QueryRow(
		context.Background(),
		"select did from handles_server_active_handles where handle = $1",
		handle.String(),
	).Scan(&did)

	if err != nil {
		return "", err
	}

	return did, nil
}

func (pg *PostgresHandles) CanProvideForDomain(domain Domain) (bool, error) {
	connection, err := pg.pool.Acquire(context.Background())

	defer connection.Release()

	if err != nil {
		return false, err
	}

	exists := false

	err = connection.QueryRow(context.Background(), "select exists(select 1 from handles_server_active_domains where domain = $1)", domain).Scan(&exists)

	return exists, err
}

func (pg *PostgresHandles) IsHealthy() (bool, string) {
	connection, err := pg.pool.Acquire(context.Background())

	if err != nil {
		return false, err.Error()
	}

	defer connection.Release()

	err = connection.Ping(context.Background())

	if err != nil {
		return false, err.Error()
	}

	return true, "Connected to database"
}
