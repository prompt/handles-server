package main

import (
	"context"
	"fmt"
	"strings"
)

type DecentralizedID string

type Hostname string

type Domain string

type Username string

type Handle struct {
	Domain   Domain
	Username Username
}

func (handle Handle) String() string {
	return strings.ToLower(fmt.Sprintf("%s.%s", handle.Username, handle.Domain))
}

type ProvidesDecentralizedIDs interface {
	GetDecentralizedIDForHandle(ctx context.Context, handle Handle) (DecentralizedID, error)
	CanProvideForDomain(ctx context.Context, domain Domain) (bool, error)
	IsHealthy(ctx context.Context) (bool, string)
}

type DecentralizedIDNotFoundError struct {
	handle Handle
}

func (e DecentralizedIDNotFoundError) Error() string {
	return fmt.Sprintf("No DID found for %s", e.handle.String())
}

type CannotGetHandelsFromDomainError struct {
	domain Domain
}

func (e CannotGetHandelsFromDomainError) Error() string {
	return fmt.Sprintf("Domain %s is not supported by this server", e.domain)
}

func (e *CannotGetHandelsFromDomainError) Is(target error) bool {
	_, ok := target.(*CannotGetHandelsFromDomainError)
	return ok
}
