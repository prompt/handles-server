package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHandles struct {
	pool         *pgxpool.Pool
	didsTable    string
	domainsTable string
}

func NewPostgresHandlesProvider(config *pgxpool.Config, didsTable string, domainsTable string) (*PostgresHandles, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return &PostgresHandles{}, err
	}

	pg := &PostgresHandles{pool, didsTable, domainsTable}

	healthy, status := pg.IsHealthy(context.Background())

	if !healthy {
		return &PostgresHandles{}, errors.New(status)
	}

	canAccessTables, err := pg.canAccessTables(context.Background())

	if err != nil {
		return &PostgresHandles{}, err
	}

	if !canAccessTables {
		return &PostgresHandles{}, errors.New("cannot access tables")
	}

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

	query := fmt.Sprintf(
		"select did from %s where LOWER(handle) = LOWER($1)",
		pgx.Identifier{pg.didsTable}.Sanitize(),
	)

	err = connection.QueryRow(ctx, query, handle.String()).Scan(&did)

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

	query := fmt.Sprintf(
		"select exists(select 1 from %s where domain = $1)",
		pgx.Identifier{pg.domainsTable}.Sanitize(),
	)

	err = connection.QueryRow(ctx, query, domain).Scan(&exists)

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

func (pg *PostgresHandles) canAccessTables(ctx context.Context) (bool, error) {
	connection, err := pg.pool.Acquire(ctx)

	if err != nil {
		return false, err
	}

	defer connection.Release()

	canAccessTables := false

	query := fmt.Sprintf(
		"select true from %s union select true from %s",
		pgx.Identifier{pg.didsTable}.Sanitize(),
		pgx.Identifier{pg.domainsTable}.Sanitize(),
	)

	err = connection.QueryRow(ctx, query).Scan(&canAccessTables)

	return canAccessTables, err
}
