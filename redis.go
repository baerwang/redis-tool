package tool

import (
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var session *Session

// Config redis 配置
type Config struct {
	Host        string
	Password    string
	Prefix      string
	Port        int
	DbName      int
	MaxIdle     int
	IdleTimeout int
}

type Session struct {
	pool    *redis.Pool
	prefix  string
	whether bool
}

func LoadRedisSession(c Config) error {
	session = &Session{}

	session.pool = &redis.Pool{
		MaxIdle:     c.MaxIdle,
		IdleTimeout: time.Duration(c.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", c.Host+":"+strconv.Itoa(c.Port),
				redis.DialPassword(c.Password), redis.DialDatabase(c.DbName))
		},
	}

	session.SetPrefix(c.Prefix)

	conn := session.pool.Get()
	defer conn.Close()

	return conn.Err()
}

func GetSession() *Session {
	if session == nil {
		session = &Session{}
	}
	return session
}

func (r *Session) SetPrefix(name string) {
	r.prefix = name + ":"
}

func (r *Session) GetPrefix() string {
	return r.prefix
}

func (r *Session) GetConn() redis.Conn {
	return r.pool.Get()
}

func (r Session) Whether(b bool) Session {
	r.whether = b
	return r
}

func (r *Session) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer Close(conn)
	if !r.whether {
		args[0] = r.GetPrefix() + args[0].(string)
	}
	return conn.Do(cmd, args...)
}

func Close(conn redis.Conn) {
	func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			log.Println("redis conn close failure")
		}
	}(conn)
}
