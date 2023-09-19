package kvdb

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const LogPrefix = ">>>>> "

type KVDB struct {
	mutex          sync.Mutex
	data           map[string]*entryOffset
	dataFileOffset uint64
	fd             *os.File
	closed         atomic.Bool
	options        *options
}

func Open(opts ...Option) *KVDB {
	db := &KVDB{}
	db.options = newOptions(opts...)
	db.Init()

	log.Println(LogPrefix, "KVDB is started now")
	return db
}

func (db *KVDB) isClosed() bool {
	return db.closed.Load()
}

func (db *KVDB) Close() (err error) {
	if db.isClosed() {
		log.Println(LogPrefix, "KVDB has already closed")
		return
	}

	defer db.fd.Close()
	defer log.Println(LogPrefix, "KVDB is closed now")

	db.closed.CompareAndSwap(false, true)
	db.data = nil

	return
}

func (db *KVDB) Init() {
	log.Println(LogPrefix, "KVDB.init() start")
	// 检查数据文件是否存在，如果不存在则创建新文件
	_, err := os.Stat(db.options.dbFileName)
	if err != nil {
		file, err := os.Create(db.options.dbFileName)
		if err != nil {
			return
		}
		file.Close()
	}

	db.fd, err = os.OpenFile(db.options.dbFileName, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	//defer db.fd.Close()

	// 初始化 kvdb.data
	if db.data == nil {
		db.data = make(map[string]*entryOffset)
	}

	curTimeUnix := time.Now().Unix()
	buf := make([]byte, entryMetaSize)
	var offset int64

	for {
		if _, err := db.fd.ReadAt(buf, offset); err != nil {
			if err == io.EOF {
				break
			}
		}
		//fmt.Printf("offset  is %d\n", offset)
		//fmt.Printf("buf is %d\n", buf)

		entry, err := decodeEntry(buf)
		if err != nil {
		}

		key := make([]byte, entry.kLen)
		if _, err = db.fd.ReadAt(key, offset+entryMetaSize); err != nil {
			if err == io.EOF {
				break
			}
		}

		if entry.mark == EntryDeleted || (entry.ttl > 0 && entry.ttl < curTimeUnix) {
			if _, ok := db.data[string(key)]; ok {
				delete(db.data, string(key))
			}
		} else {
			db.data[string(key)] = &entryOffset{
				offset: uint32(offset),
			}
		}

		//fmt.Printf("ety:%+v\n", entry)
		//fmt.Printf("key:%s\n", key)

		offset = offset + entryMetaSize + int64(entry.kLen) + int64(entry.vLen)
	}
	db.dataFileOffset = uint64(offset)

	log.Println(LogPrefix, "KBVD.init() okay")
}

func (db *KVDB) Printf() {
	log.Println(LogPrefix, "KVDB.Printf() start")
	buf := make([]byte, entryMetaSize)

	for k, v := range db.data {
		fmt.Printf("key:%s \n", k)
		//fmt.Printf("v is: %+v \n", v)
		_, err := db.fd.ReadAt(buf, int64(v.offset))
		if err != nil {
		}

		//fmt.Printf("buf %+v\n", buf)

		entry, err := decodeEntry(buf)
		val := make([]byte, entry.vLen)

		//fmt.Printf("entry %+v\n", entry)
		//fmt.Println("val offset ", int64(v.offset)+entryMetaSize+int64(entry.kLen))

		_, err = db.fd.ReadAt(val, int64(v.offset)+entryMetaSize+int64(entry.kLen))
		fmt.Printf("val:[%s] bytes:%+v \n", string(val), val)
	}

	log.Println(LogPrefix, "KBVD.Printf() okay")
}

func (db *KVDB) txBegin() (*Tx, error) {
	tx := &Tx{
		db: db,
	}

	return tx, nil
}

func (db *KVDB) Tx(fn func(tx *Tx)) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	tx, _ := db.txBegin()
	fn(tx)
}

func (db *KVDB) Put(key string, value []byte) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.put(key, value)
}

func (db *KVDB) put(key string, value []byte) error {
	entry := NewEntry(key, value)
	byte, err := entry.encode()
	if err != nil {
		return err
	}

	_, err = db.fd.WriteAt(byte, int64(db.dataFileOffset))
	if err != nil {
		return err
	}

	// 更新内存索引
	db.data[key] = &entryOffset{offset: uint32(atomic.LoadUint64(&db.dataFileOffset))}
	atomic.AddUint64(&db.dataFileOffset, uint64(entry.Size()))

	return nil
}

func (db *KVDB) Get(key string) ([]byte, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.get(key)
}

func (db *KVDB) get(key string) ([]byte, error) {
	entryOffset, ok := db.data[key]
	if !ok {
		return nil, errors.New("not found")
	}

	buf := make([]byte, entryMetaSize)
	db.fd.ReadAt(buf, int64(entryOffset.offset))
	entry, _ := decodeEntry(buf)

	val := make([]byte, entry.vLen)
	db.fd.ReadAt(val, int64(entryOffset.offset)+entryMetaSize+int64(entry.kLen))
	return val, nil
}

func (db *KVDB) Delete(key string) (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.delete(key)
}

func (db *KVDB) delete(key string) (err error) {
	// Search for key
	if _, ok := db.data[string(key)]; !ok {
		return
	}

	// Write to file
	e := NewDelEntry(key)
	byte, err := e.encode()
	if err != nil {
		return
	}
	_, err = db.fd.WriteAt(byte, int64(db.dataFileOffset))
	if err != nil {
		return
	}

	atomic.AddUint64(&db.dataFileOffset, uint64(e.Size()))

	delete(db.data, key)
	return nil
}
