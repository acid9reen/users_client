package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestCache_New(t *testing.T) {
	New()
}

func TestCache_StoreValidUser(t *testing.T) {
	cache := New()

	user := &UserReadModel{1, "TestName", "TestLastName"}
	cache.Set(user)

	if retrieved := cache.Get(1); user != retrieved {
		t.Fatalf("User in cache differs from given one `%v` (given) != `%v` (retrieved)", user, retrieved)
	}
}

func TestCache_UserCleanedOnExpire(t *testing.T) {
	lifetime := time.Millisecond * 10
	cache := New(SetLifetime(lifetime))

	user := &UserReadModel{1, "TestName", "TestLastName"}
	cache.Set(user)
	time.Sleep(lifetime * 2)

	if cache.Get(1) != nil {
		t.Fatalf("User was not cleaned after expire")
	}
}

func BenchmarkCache_GetSetOperations(b *testing.B) {
	numUsers := 100_000
	cache := New(SetLifetime(time.Millisecond * 10))
	users := make([]*UserReadModel, numUsers)
	for i, v := range users {
		users[i] = &UserReadModel{uint64(i), "TestName" + fmt.Sprint(v), "TestLastName" + fmt.Sprint(v)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go cache.Set(users[rand.Intn(numUsers)])
		go cache.Get(uint64(rand.Intn(numUsers)))
	}
}
