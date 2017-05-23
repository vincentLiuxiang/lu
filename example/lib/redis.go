package lib

import "github.com/garyburd/redigo/redis"

type Redis struct {
	redisConn redis.Conn
	lifeTime  int64
}

func (s *Redis) SetExpire(key, value interface{}, lifeTime int64) error {
	var err error = nil
	if lifeTime <= 0 {
		_, err = s.redisConn.Do("SET", key, value)
	} else {
		_, err = s.redisConn.Do("SET", key, value, "EX", lifeTime)
	}
	return err
}

func (s *Redis) Set(key, value interface{}) error {
	var err error = nil
	if s.lifeTime <= 0 {
		_, err = s.redisConn.Do("SET", key, value)
	} else {
		_, err = s.redisConn.Do("SET", key, value, "EX", s.lifeTime)
	}
	return err
}

func (s *Redis) Get(key interface{}) (interface{}, error) {
	return s.redisConn.Do("GET", key)
}

func (s *Redis) Delete(key interface{}) error {
	_, err := s.redisConn.Do("DEL", key)
	return err
}

func NewRedis(protocol, hostPort string, lifeTime int64) *Redis {
	c, err := redis.Dial(protocol, hostPort)
	if err != nil {
		panic(err)
		return nil
	}
	return &Redis{redisConn: c, lifeTime: lifeTime}
}
