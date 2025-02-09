package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHandles struct {
	pool         *pgxpool.Pool
	domainsTable string
	handlesTable string
}

func NewPostgresHandlesProvider(config *pgxpool.Config, handlesTable string, domainsTable string) (*PostgresHandles, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return &PostgresHandles{}, err
	}

	pg := &PostgresHandles{pool, handlesTable, domainsTable}

	healthy, status := pg.IsHealthy(context.Background())

	if !healthy {
		return &PostgresHandles{}, errors.New(status)
	}

	// TODO: Check can access tables

	return pg, nil
}

func (pg *PostgresHandles) GetDecentralizedIDForHandle(ctx context.Context, handle Handle) (DecentralizedID, error) {
	connection, err := pg.pool.Acquire(ctx)

	defer connection.Release()

	if err != nil {
		return "", err
	}

	canProvide, err := pg.CanProvideForDomain(ctx, handle.Domain)

	if err != nil {
		return "", err
	}

	if !canProvide {
		return "", &CannotGetHandelsFromDomainError{domain: handle.Domain}
	}

	var did DecentralizedID

	err = connection.QueryRow(
		ctx,
		"select did from handles_server_active_handles where handle = $1",
		handle.String(),
	).Scan(&did)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return did, nil
}

func (pg *PostgresHandles) CanProvideForDomain(ctx context.Context, domain Domain) (bool, error) {
	connection, err := pg.pool.Acquire(ctx)

	defer connection.Release()

	if err != nil {
		return false, err
	}

	exists := false

	err = connection.QueryRow(ctx, "select exists(select 1 from handles_server_active_domains where domain = $1)", domain).Scan(&exists)

	return exists, err
}

func (pg *PostgresHandles) IsHealthy(ctx context.Context) (bool, string) {
	connection, err := pg.pool.Acquire(ctx)

	if err != nil {
		return false, err.Error()
	}

	defer connection.Release()

	err = connection.Ping(ctx)

	if err != nil {
		return false, err.Error()
	}

	return true, "Connected to database"
}
