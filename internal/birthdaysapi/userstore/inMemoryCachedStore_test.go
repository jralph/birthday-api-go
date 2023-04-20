package userstore

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

func TestInMemoryCachedStore(t *testing.T) {
	t.Run("in memory store can put items", func(t *testing.T) {
		mockUserStore := &MockUserStore{
			Users: map[string]*User{},
		}

		inMemoryCachedStore := NewInMemoryCachedStore(mockUserStore, time.Second*1)

		user := User{
			Username:    "joe",
			DateOfBirth: DateOfBirth(time.Now()),
		}

		err := inMemoryCachedStore.Put(&user)
		if err != nil {
			t.Errorf("error putting user: %s", err)
		}

		saved, _ := mockUserStore.Get("joe")

		if diff := cmp.Diff(*saved, user, cmpopts.IgnoreUnexported(DateOfBirth{})); diff != "" {
			t.Errorf("NewInMemoryCachedStore{} mismatch (-want +got):\n%s", diff)
		}

		if mockUserStore.PutCalls != 1 {
			t.Errorf("unexpected number of put calls to parent store, got %d want %d", mockUserStore.PutCalls, 1)
		}

		if mockUserStore.GetCalls != 1 {
			t.Errorf("unexpected number of get calls to parent store, got %d want %d", mockUserStore.GetCalls, 1)
		}
	})

	t.Run("in memory store can get items", func(t *testing.T) {
		mockUserStore := &MockUserStore{
			Users: map[string]*User{},
		}

		inMemoryCachedStore := NewInMemoryCachedStore(mockUserStore, time.Second*1)

		user := User{
			Username:    "joe",
			DateOfBirth: DateOfBirth(time.Now()),
		}

		err := inMemoryCachedStore.Put(&user)
		if err != nil {
			t.Errorf("error putting user: %s", err)
		}

		cached, _ := inMemoryCachedStore.Get("joe")

		if mockUserStore.PutCalls != 1 {
			t.Errorf("unexpected number of put calls to parent store, got %d want %d", mockUserStore.PutCalls, 1)
		}

		if mockUserStore.GetCalls != 0 {
			t.Errorf("unexpected number of get calls to parent store, got %d want %d", mockUserStore.GetCalls, 0)
		}

		saved, _ := mockUserStore.Get("joe")

		// Check user from both stores matches
		if diff := cmp.Diff(*saved, *cached, cmpopts.IgnoreUnexported(DateOfBirth{})); diff != "" {
			t.Errorf("NewInMemoryCachedStore{} mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("in memory store expires items", func(t *testing.T) {
		mockUserStore := &MockUserStore{
			Users: map[string]*User{},
		}

		inMemoryCachedStore := NewInMemoryCachedStore(mockUserStore, time.Second*1)

		user := User{
			Username:    "joe",
			DateOfBirth: DateOfBirth(time.Now()),
		}

		err := inMemoryCachedStore.Put(&user)
		if err != nil {
			t.Errorf("error putting user: %s", err)
		}

		time.Sleep(time.Millisecond * 1100)

		cached, _ := inMemoryCachedStore.Get("joe")

		if mockUserStore.PutCalls != 1 {
			t.Errorf("unexpected number of put calls to parent store, got %d want %d", mockUserStore.PutCalls, 1)
		}

		// At this point, we expect a call to be made to the parent store as the cache will have expired
		if mockUserStore.GetCalls != 1 {
			t.Errorf("unexpected number of get calls to parent store, got %d want %d", mockUserStore.GetCalls, 1)
		}

		saved, _ := mockUserStore.Get("joe")

		// Check user from both stores matches
		if diff := cmp.Diff(*saved, *cached, cmpopts.IgnoreUnexported(DateOfBirth{})); diff != "" {
			t.Errorf("NewInMemoryCachedStore{} mismatch (-want +got):\n%s", diff)
		}

		// Check user from memory store matches sent user
		if diff := cmp.Diff(*cached, user, cmpopts.IgnoreUnexported(DateOfBirth{})); diff != "" {
			t.Errorf("NewInMemoryCachedStore{} mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("in memory store can get items that are not cached", func(t *testing.T) {
		mockUserStore := &MockUserStore{
			Users: map[string]*User{},
		}

		inMemoryCachedStore := NewInMemoryCachedStore(mockUserStore, time.Second*1)

		user := User{
			Username:    "joe",
			DateOfBirth: DateOfBirth(time.Now()),
		}

		err := mockUserStore.Put(&user)

		if err != nil {
			t.Errorf("error putting user: %s", err)
		}

		cached, _ := inMemoryCachedStore.Get("joe")

		if mockUserStore.PutCalls != 1 {
			t.Errorf("unexpected number of put calls to parent store, got %d want %d", mockUserStore.PutCalls, 1)
		}

		if mockUserStore.GetCalls != 1 {
			t.Errorf("unexpected number of get calls to parent store, got %d want %d", mockUserStore.GetCalls, 1)
		}

		saved, _ := mockUserStore.Get("joe")

		// Check user from both stores matches
		if diff := cmp.Diff(*saved, *cached, cmpopts.IgnoreUnexported(DateOfBirth{})); diff != "" {
			t.Errorf("NewInMemoryCachedStore{} mismatch (-want +got):\n%s", diff)
		}
	})
}
