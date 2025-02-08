package main

import (
	"fmt"
)

type DecentralizedID string

type Hostname string

type Domain string

type Username string

type Handle struct {
	Hostname Hostname
	Domain   Domain
	Username Username
}

func (handle Handle) String() string {
	return fmt.Sprintf("%s.%s", handle.Username, handle.Domain)
}

type ResolvesHandlesToDecentralizedIDs interface {
	ResolveHandleToDID(handle Handle) (DecentralizedID, error)
	CanResolveHandlesAtDomain(domain Domain) bool
	IsResolverHealthy() (bool, string)
}

type DecentralizedIDNotFoundError struct {
	handle Handle
}

func (e DecentralizedIDNotFoundError) Error() string {
	return fmt.Sprintf("No DID found for %s", e.handle.Hostname)
}

type CannotResolveDomainError struct {
	handle Handle
}

func (e CannotResolveDomainError) Error() string {
	return fmt.Sprintf("Domain %s is not supported by this server", e.handle.Domain)
}

type MapOfDids = map[Hostname]DecentralizedID

type MapOfDomains = map[Domain]bool

type InMemoryResolver struct {
	dids      MapOfDids
	domains   MapOfDomains
	isHealthy bool
}

func NewInMemoryResolver(dids MapOfDids, domains MapOfDomains) *InMemoryResolver {
	return &InMemoryResolver{dids, domains, true}
}

func (handles *InMemoryResolver) ResolveHandleToDID(handle Handle) (DecentralizedID, error) {
	if !handles.CanResolveHandlesAtDomain(handle.Domain) {
		return "", &CannotResolveDomainError{handle: handle}
	}

	if did, found := handles.dids[handle.Hostname]; found {
		return did, nil
	}

	return "", DecentralizedIDNotFoundError{handle}
}

func (handles *InMemoryResolver) CanResolveHandlesAtDomain(domain Domain) bool {
	return handles.domains[domain]
}

func (handles *InMemoryResolver) IsResolverHealthy() (bool, string) {
	if handles.isHealthy {
		return true, fmt.Sprintf("Available with %d resolvable handles for %d domains", len(handles.dids), len(handles.domains))
	}

	return false, "Not healthy"
}

func (handles *InMemoryResolver) SetHealthy(isHealthy bool) {
	handles.isHealthy = isHealthy
}
