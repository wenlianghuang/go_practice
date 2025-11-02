package security

import (
	"sync"
	"time"
)

// TokenStore defines the contract for storing and querying revoked JWT IDs.
type TokenStore interface {
	Revoke(jti string, expiresAt time.Time) error
	IsRevoked(jti string) (bool, error)
}

// InMemoryTokenStore provides a lightweight in-memory TokenStore implementation.
// It is intended for development or testing scenarios where persistence is not required.
type InMemoryTokenStore struct {
	mu      sync.RWMutex
	revoked map[string]time.Time
}

// NewInMemoryTokenStore returns a fresh in-memory store with no revoked tokens.
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		revoked: make(map[string]time.Time),
	}
}

// Revoke marks the provided token ID as revoked until the supplied expiration time.
// A write lock is used because the method mutates the shared map.
func (s *InMemoryTokenStore) Revoke(jti string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revoked[jti] = expiresAt
	// Opportunistically clean up stale entries whenever we add a new one.
	s.gcLocked()
	return nil
}

// IsRevoked reports whether the token ID has been revoked and has not yet expired.
// The read lock allows concurrent lookups without blocking other readers.
func (s *InMemoryTokenStore) IsRevoked(jti string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, ok := s.revoked[jti]
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		// If the revocation has expired, it is effectively no longer revoked.
		return false, nil
	}
	return true, nil
}

// gcLocked removes expired revocation entries; callers must already hold the lock.
func (s *InMemoryTokenStore) gcLocked() {
	now := time.Now()
	for token, exp := range s.revoked {
		if now.After(exp) {
			// Remove entries that have expired to avoid unbounded growth.
			delete(s.revoked, token)
		}
	}
}
