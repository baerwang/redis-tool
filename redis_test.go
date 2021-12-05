package main

import (
	"github.com/stretchr/testify/assert"
)

import (
	"testing"
)

func TestRedisPrefix(t *testing.T) {
	rc := Config{
		Host:     "127.0.0.1",
		Password: "xxxx",
		Prefix:   "test",
		Port:     6379,
	}
	err := LoadRedisSession(rc)
	assert.NoError(t, err)

	conn := GetSession()
	set, err := conn.Set("test", "test111")
	assert.NoError(t, err)
	t.Log(set)
	get, err := conn.Get("test")
	assert.NoError(t, err)
	t.Log(get)
	del, err := conn.Del("test")
	assert.NoError(t, err)
	t.Log(del)
}

func TestRedisNotPrefix(t *testing.T) {
	rc := Config{
		Host:     "127.0.0.1",
		Password: "xxxx",
		Prefix:   "test",
		Port:     6379,
	}

	err := LoadRedisSession(rc)
	assert.NoError(t, err)

	conn := GetSession().Whether(true)
	set, err := conn.Set("test", "test111")
	assert.NoError(t, err)
	t.Log(set)
	get, err := conn.Get("test")
	assert.NoError(t, err)
	t.Log(get)
	del, err := conn.Del("test")
	assert.NoError(t, err)
	t.Log(del)
}
