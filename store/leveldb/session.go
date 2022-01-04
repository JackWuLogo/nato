package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

type Session struct {
	sync.Mutex
	file string
	db   *leveldb.DB
}

func (s *Session) Reload() error {
	s.Lock()
	defer s.Unlock()

	db, err := leveldb.OpenFile(s.file, nil)
	if err != nil {
		return err
	}

	if s.db != nil {
		_ = s.db.Close()
	}
	s.db = db

	return nil
}

func (s *Session) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
		s.db = nil
	}

	return nil
}

// DB 数据库
func (s *Session) DB() *leveldb.DB {
	return s.db
}

// Put 写入数据
func (s *Session) Put(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Get 读取数据. 需要判断返回数据是否为空
func (s *Session) Get(key []byte) ([]byte, error) {
	res, err := s.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

// Del 删除数据
func (s *Session) Del(key []byte) error {
	return s.db.Delete(key, nil)
}

func NewSession(file string) (*Session, error) {
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return nil, err
	}

	return &Session{
		file: file,
		db:   db,
	}, nil
}
