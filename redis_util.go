package tool

import "github.com/gomodule/redigo/redis"

func (r *Session) Set(key, value string) (bool, error) {
	s, err := redis.String(r.Do(SET, key, value))
	if err != nil {
		return false, err
	}
	return s == "OK", nil
}

func (r *Session) Get(key string) (string, error) {
	return redis.String(r.Do(GET, key))
}

func (r *Session) Del(key string) (bool, error) {
	return redis.Bool(r.Do(DEL, key))
}

func (r *Session) StringExists(key string) (bool, error) {
	return redis.Bool(r.Do(GET, key))
}

func (r *Session) SetExpire(key, value string, expire int) (bool, error) {
	exists, err := r.Set(key, value)
	if err != nil {
		return false, err
	}
	if exists {
		return redis.Bool(r.Do(EXPIRE, key, expire))
	}
	return false, nil
}

func (r *Session) HSetSend(hashKey, key string, value interface{}) (bool, error) {
	return redis.Bool(r.Do(HSET, hashKey, key, value))
}

func (r *Session) HDel(hashKey string, key ...string) (bool, error) {
	return redis.Bool(r.Do(HDEL, redis.Args{}.Add(hashKey).AddFlat(key)...))
}

func (r *Session) HSetExpire(hashKey, key string, value interface{}, expire int) (bool, error) {
	exists, err := r.HSetSend(hashKey, key, value)
	if err != nil {
		return false, err
	}
	if exists {
		return redis.Bool(r.Do(EXPIRE, hashKey, expire))
	}
	return false, nil
}

func (r *Session) HaShGet(hashKey, key string) ([]byte, error) {
	return redis.Bytes(r.Do(HGET, hashKey, key))
}

func (r *Session) HashAll(hashKey string) ([]interface{}, error) {
	return redis.Values(r.Do(HGETALL, hashKey))
}

func (r *Session) HashAllMap(hashKey string) (map[string]string, error) {
	return redis.StringMap(r.Do(HGETALL, hashKey))
}

func (r *Session) Exists(key string) (bool, error) {
	return redis.Bool(r.Do(EXISTS, key))
}

func (r *Session) HExists(hashKey, key string) (bool, error) {
	return redis.Bool(r.Do(HEXISTS, hashKey, key))
}

func (r *Session) HExistsAndSet(hashKey, key string, value interface{}) (bool, error) {
	exists, err := r.HExists(hashKey, key)
	if err != nil {
		return false, err
	}
	if !exists {
		return r.HSetSend(hashKey, key, value)
	}
	return false, nil
}

func (r *Session) LpushSend(key string, value interface{}) (bool, error) {
	return redis.Bool(r.Do(LPUSH, key, value))
}

func (r *Session) RPUSHSend(key string, value interface{}) (bool, error) {
	return redis.Bool(r.Do(RPUSH, key, value))
}

func (r *Session) ZADDSend(key string, score int, value interface{}) (bool, error) {
	return redis.Bool(r.Do(ZADD, key, score, value))
}

func (r *Session) Lrange(key string, start, end int) ([]interface{}, error) {
	return redis.Values(r.Do(LRANGE, key, start, end))
}

func (r *Session) Brpop(key string, timeout int) ([]interface{}, error) {
	return redis.Values(r.Do(BRPOP, key, timeout))
}

func (r *Session) ScanCursor(key string) ([]interface{}, error) {
	values, err := redis.Values(r.Do(SCAN, 0, MATCH, key+"*"))
	if err != nil {
		return nil, err
	}

	var (
		cursor int64
		items  []interface{}
	)
	_, err = redis.Scan(values, &cursor, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
