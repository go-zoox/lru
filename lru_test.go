package lru

import (
	"strings"
	"testing"
)

func TestLRU(t *testing.T) {
	cache := New(3)
	cache.Set("a", 1)
	if v, ok := cache.Get("a"); !ok || v != 1 {
		t.Errorf("cache.Get(a) = %v, %v, want 1, true", v, ok)
	}
	if cache.Len() != 1 {
		t.Errorf("cache.Len() = %v, want 1", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "a" {
		t.Errorf("cache.Keys() = %v, want a", strings.Join(cache.Keys(), ","))
	}

	cache.Set("b", 2)
	if v, ok := cache.Get("b"); !ok || v != 2 {
		t.Errorf("cache.Get(b) = %v, %v, want 2, true", v, ok)
	}
	if cache.Len() != 2 {
		t.Errorf("cache.Len() = %v, want 2", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "b,a" {
		t.Errorf("cache.Keys() = %v, want b,a", strings.Join(cache.Keys(), ","))
	}

	cache.Set("c", 3)
	if v, ok := cache.Get("c"); !ok || v != 3 {
		t.Errorf("cache.Get(c) = %v, %v, want 3, true", v, ok)
	}
	if cache.Len() != 3 {
		t.Errorf("cache.Len() = %v, want 3", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "c,b,a" {
		t.Errorf("cache.Keys() = %v, want c,b,a", strings.Join(cache.Keys(), ","))
	}

	cache.Set("d", 4)
	if v, ok := cache.Get("d"); !ok || v != 4 {
		t.Errorf("cache.Get(d) = %v, %v, want 4, true", v, ok)
	}
	if cache.Len() != 3 {
		t.Errorf("cache.Len() = %v, want 3", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "d,c,b" {
		t.Errorf("cache.Keys() = %v, want d,c,b", strings.Join(cache.Keys(), ","))
	}
	if v, ok := cache.Get("a"); ok {
		t.Errorf("cache.Get(a) = %v, %v, want nil, false", v, ok)
	}

	cache.Set("e", 5)
	if v, ok := cache.Get("e"); !ok || v != 5 {
		t.Errorf("cache.Get(e) = %v, %v, want 5, true", v, ok)
	}
	if cache.Len() != 3 {
		t.Errorf("cache.Len() = %v, want 3", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "e,d,c" {
		t.Errorf("cache.Keys() = %v, want e,d,c", strings.Join(cache.Keys(), ","))
	}
	if v, ok := cache.Get("b"); ok {
		t.Errorf("cache.Get(b) = %v, %v, want nil, false", v, ok)
	}

	cache.Delete("d")
	if v, ok := cache.Get("d"); ok {
		t.Errorf("cache.Get(d) = %v, %v, want nil, false", v, ok)
	}
	if cache.Len() != 2 {
		t.Errorf("cache.Len() = %v, want 2", cache.Len())
	}
	if strings.Join(cache.Keys(), ",") != "e,c" {
		t.Errorf("cache.Keys() = %v, want e,c", strings.Join(cache.Keys(), ","))
	}

	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("cache.Len() = %v, want 0", cache.Len())
	}
}
