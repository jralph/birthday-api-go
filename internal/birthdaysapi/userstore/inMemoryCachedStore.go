package userstore

import (
	"sync"
	"time"
)

// CachedUser is a simple implementation of a cached user object with an expiry
type CachedUser struct {
	User
	Expiry time.Time
}

// InMemoryCachedStore is a rudimentary in-memory caching wrapper for a store (or chain of stores)
// It provides the ability to cache objects for a given duration
type InMemoryCachedStore struct {
	store    UserStore
	cache    map[string]CachedUser
	duration time.Duration
	lock     sync.Mutex
}

// NewInMemoryCachedStore handles creation of a new InMemoryCachedStore object and handles its configuration
// It returns a configured and ready to use InMemoryCachedStore
func NewInMemoryCachedStore(store UserStore, duration time.Duration) *InMemoryCachedStore {
	return &InMemoryCachedStore{
		store:    store,
		cache:    map[string]CachedUser{},
		duration: duration,
	}
}

// Put is responsible for saving a given user
// It will simply do nothing if the user in the cache is the same as the user requested to save
// This avoids extra calls to the upstream store when not necessary and also refreshes the cache expiry
func (s *InMemoryCachedStore) Put(user *User) error {
	s.lock.Lock()
	cached, exists := s.cache[user.Username]
	s.lock.Unlock()
	if exists {
		cachedDob := time.Time(cached.DateOfBirth).Format("2006-01-02")
		newDob := time.Time(user.DateOfBirth).Format("2006-01-02")

		if cachedDob == newDob {
			// Refresh the cache
			s.lock.Lock()
			cached.Expiry = time.Now().Add(s.duration)
			s.lock.Unlock()
			return nil
		}
	}

	err := s.store.Put(user)
	if err != nil {
		// Cached data is different to the requested new data and there was an error
		// Remove the user from the cache
		if exists {
			s.lock.Lock()
			delete(s.cache, user.Username)
			s.lock.Unlock()
		}
		return err
	}

	// Save the newly saved user to the cache
	s.lock.Lock()
	s.cache[user.Username] = CachedUser{
		User:   *user,
		Expiry: time.Now().Add(s.duration),
	}
	s.lock.Unlock()

	return nil
}

// Get handles fetching a user from the store or cache
// If a cached user is returned, it is refreshed and returned
// If the user is not found in the cache, the store is queried and the user is cached if it exists
func (s *InMemoryCachedStore) Get(username string) (*User, error) {
	s.lock.Lock()
	cached, exists := s.cache[username]
	s.lock.Unlock()

	// If the cache is still alive, return it, if not, delete it and continue
	if exists && cached.Expiry.After(time.Now()) {
		// Refresh the cache
		s.lock.Lock()
		cached.Expiry = time.Now().Add(s.duration)
		s.lock.Unlock()
		return &cached.User, nil
	} else if exists {
		s.lock.Lock()
		delete(s.cache, username)
		s.lock.Unlock()
	}

	// Fetch, then cache, the user
	user, err := s.store.Get(username)
	if err != nil {
		return nil, err
	} else if user == nil {
		// User does not exist
		// If implementing a delete function at any point, a cache clear will need to be done here
		return nil, nil
	}

	s.lock.Lock()
	s.cache[user.Username] = CachedUser{
		User:   *user,
		Expiry: time.Now().Add(s.duration),
	}
	s.lock.Unlock()
	return user, nil
}
