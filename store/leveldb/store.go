package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Store struct {
	file string
}

// 写入数据
func (s *Store) Put(key []byte, value []byte) error {
	// open leveldb
	db, err := leveldb.OpenFile(s.file, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Put(key, value, nil); err != nil {
		return err
	}

	return nil
}

// 读取数据. 需要判断返回数据是否为空
func (s *Store) Get(key []byte) ([]byte, error) {
	// open leveldb
	db, err := leveldb.OpenFile(s.file, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	val, err := db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return val, nil
}

// 删除数据
func (s *Store) Del(key []byte) error {
	// open leveldb
	db, err := leveldb.OpenFile(s.file, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(key, nil); err != nil {
		return err
	}

	return nil
}

func NewStore(file string) *Store {
	return &Store{
		file: file,
	}
}
