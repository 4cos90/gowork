package cacheServer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type option func(*cache)

type cacheValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Time  time.Time   `json:"time"`
}

type cache struct {
	Name               string
	Cache              []cacheValue
	LocalIp            string
	Timeout            time.Duration
	RedisClient        RedisServer
	Lock               sync.Mutex
	ClearAfterGet      bool
	SyncFromOtherCache bool
}

type CacheServer interface {
	Get(key string) ([]cacheValue, error)
	Set(key string, kvalue interface{}) error
	Clear(key string) error
	GetAll() ([]cacheValue, error)
}

func (s *cache) clearOutTimeCache() {
	for startindex := 0; startindex < len(s.Cache); startindex++ {
		if time.Now().Sub(s.Cache[startindex].Time) < s.Timeout {
			s.Cache = s.Cache[startindex:]
			break
		} else if startindex == (len(s.Cache) - 1) {
			s.Cache = make([]cacheValue, 0)
		}
	}
}

func (s *cache) GetAll() ([]cacheValue, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.getLastWriterFromRedis()
	s.clearOutTimeCache()
	return s.Cache, nil
}

func (s *cache) Get(key string) ([]cacheValue, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.getLastWriterFromRedis()
	s.clearOutTimeCache()
	rlt := make([]cacheValue, 0)
	if s.ClearAfterGet {
		newCache := make([]cacheValue, 0)
		for i := 0; i < len(s.Cache); i++ {
			if s.Cache[i].Key == key {
				rlt = append(rlt, s.Cache[i])
			} else {
				newCache = append(newCache, s.Cache[i])
			}
		}
		s.Cache = newCache
	} else {
		for i := 0; i < len(s.Cache); i++ {
			if s.Cache[i].Key == key {
				rlt = append(rlt, s.Cache[i])
			}
		}
	}
	s.setLastWriterToRedis()
	return rlt, nil
}

func (s *cache) Clear(key string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.getLastWriterFromRedis()
	s.clearOutTimeCache()
	newCache := make([]cacheValue, 0)
	for i := 0; i < len(s.Cache); i++ {
		if s.Cache[i].Key != key {
			newCache = append(newCache, s.Cache[i])
		}
	}
	s.Cache = newCache
	s.setLastWriterToRedis()
	return nil
}

func (s *cache) Set(key string, kvalue interface{}) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.getLastWriterFromRedis()
	s.clearOutTimeCache()
	value := cacheValue{
		Key:   key,
		Value: kvalue,
		Time:  time.Now(),
	}
	s.Cache = append(s.Cache, value)
	s.setLastWriterToRedis()
	return nil
}

func (s *cache) getLastWriterFromRedis() {
	if s.SyncFromOtherCache {
		IP, err := s.RedisClient.Get(s.Name)
		if err == nil {
			if IP != s.LocalIp {
				resp, err := http.Get(IP + "/syncCache")
				if err != nil {
					return
				}
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				var res []cacheValue
				json.Unmarshal([]byte(body), &res)
				s.Cache = res
			}
		}
	}
}

func (s *cache) setLastWriterToRedis() error {
	if s.SyncFromOtherCache {
		return s.RedisClient.Set(s.Name, s.LocalIp)
	} else {
		return nil
	}
}

func DialClearAfterGet(sign bool) option {
	return func(c *cache) {
		c.ClearAfterGet = sign
	}
}

func DialSyncFromOtherCache(sign bool) option {
	return func(c *cache) {
		c.SyncFromOtherCache = sign
	}
}

func NewCache(name string, timeout time.Duration, localIp string, RedisClient RedisServer, options ...option) CacheServer {
	newCache := cache{
		Name:               name,
		Cache:              make([]cacheValue, 0),
		Timeout:            timeout,
		RedisClient:        RedisClient,
		LocalIp:            localIp,
		Lock:               sync.Mutex{},
		ClearAfterGet:      false,
		SyncFromOtherCache: false,
	}
	for _, option := range options {
		option(&newCache)
	}
	return &newCache
}
