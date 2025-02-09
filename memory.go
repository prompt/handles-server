package main

import "fmt"

type MapOfDids = map[Hostname]DecentralizedID

type MapOfDomains = map[Domain]bool

type InMemoryProvider struct {
	dids      MapOfDids
	domains   MapOfDomains
	isHealthy bool
}

func NewInMemoryProvider(dids MapOfDids, domains MapOfDomains) *InMemoryProvider {
	return &InMemoryProvider{dids, domains, true}
}

func (memory *InMemoryProvider) GetDecentralizedIDForHandle(handle Handle) (DecentralizedID, error) {
	canProvide, err := memory.CanProvideForDomain(handle.Domain)

	if err != nil {
		return "", err
	}

	if !canProvide {
		return "", &CannotGetHandelsFromDomainError{handle: handle}
	}

	if did, found := memory.dids[Hostname(handle.String())]; found {
		return did, nil
	}

	return "", DecentralizedIDNotFoundError{handle}
}

func (memory *InMemoryProvider) CanProvideForDomain(domain Domain) (bool, error) {
	return memory.domains[domain], nil
}

func (memory *InMemoryProvider) IsHealthy() (bool, string) {
	if memory.isHealthy {
		return true, fmt.Sprintf("Available with %d handles for %d domains", len(memory.dids), len(memory.domains))
	}

	return false, "Not healthy"
}

func (memory *InMemoryProvider) SetHealthy(isHealthy bool) {
	memory.isHealthy = isHealthy
}
