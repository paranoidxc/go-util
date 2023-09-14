package kvdb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

const (
	entryHeaderSize = 9
)

func encodeEntry(e *entry) ([]byte, error) {
	buf := make([]byte, e.Size())

	buf[0] = byte(e.mark)
	binary.BigEndian.PutUint32(buf[1:5], e.kLen)
	binary.BigEndian.PutUint32(buf[5:9], e.vLen)
	copy(buf[entryHeaderSize:], e.key)
	copy(buf[entryHeaderSize+e.kLen:], e.value)

	//fmt.Printf("entry buf is: %+v\n", buf)
	return buf, nil
}

func decodeEntry(buf []byte) (*entry, error) {
	kLen := binary.BigEndian.Uint32(buf[1:5])
	vLen := binary.BigEndian.Uint32(buf[5:])

	e := &entry{
		mark: EntryMark(buf[0]),
		kLen: kLen,
		vLen: vLen,
	}

	if len(buf) > entryHeaderSize {
		key := make([]byte, kLen)
		val := make([]byte, vLen)
		copy(key, buf[entryHeaderSize:entryHeaderSize+kLen])
		copy(val, buf[entryHeaderSize+kLen:])
		e.key = key
		e.value = val
	}

	return e, nil
}

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

	log.Println("DB is started now")
	return db
}

func (db *KVDB) isClosed() bool {
	return db.closed.Load()
}

func (db *KVDB) Close() (err error) {
	if db.isClosed() {
		log.Println("DB has already closed")
		return
	}

	defer db.fd.Close()
	defer log.Println("DB closed")

	db.closed.CompareAndSwap(false, true)
	db.data = nil

	return
}

func (db *KVDB) Init() {
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

	buf := make([]byte, entryHeaderSize)
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
		if _, err = db.fd.ReadAt(key, offset+entryHeaderSize); err != nil {
			if err == io.EOF {
				break
			}
		}

		if entry.mark == EntryDeleted {
			if _, ok := db.data[string(key)]; ok {
				delete(db.data, string(key))
			}
		} else {
			db.data[string(key)] = &entryOffset{
				offset: uint32(offset),
			}
		}

		//fmt.Printf("entry %+v\n", entry)
		//fmt.Printf("key %s\n", key)

		offset = offset + entryHeaderSize + int64(entry.kLen) + int64(entry.vLen)
	}
	db.dataFileOffset = uint64(offset)
}

func (db *KVDB) Printf() {
	fmt.Printf("====== db Printf start ======\n")
	buf := make([]byte, entryHeaderSize)

	for k, v := range db.data {
		fmt.Printf("key is:%s \n", k)
		//fmt.Printf("v is: %+v \n", v)
		_, err := db.fd.ReadAt(buf, int64(v.offset))
		if err != nil {
		}

		//fmt.Printf("buf %+v\n", buf)

		entry, err := decodeEntry(buf)
		val := make([]byte, entry.vLen)

		//fmt.Printf("entry %+v\n", entry)
		//fmt.Println("val offset ", int64(v.offset)+entryHeaderSize+int64(entry.kLen))

		_, err = db.fd.ReadAt(val, int64(v.offset)+entryHeaderSize+int64(entry.kLen))
		fmt.Printf("val is:[%s] %+v \n", string(val), val)
	}
	fmt.Printf("====== db Printf end ======\n")
}

type entryOffset struct {
	offset uint32 // 数据在文件中的偏移量
}

type entry struct {
	mark  EntryMark // 区分是正常的还是删除的数据
	kLen  uint32    // key的长度
	vLen  uint32    // value的长度
	key   []byte    // key本身
	value []byte    // value本身
}

func (e *entry) Size() uint32 {
	return entryHeaderSize + e.kLen + e.vLen
}

type EntryMark int

const (
	EntryNormal EntryMark = iota
	EntryDeleted
)

func (db *KVDB) Put(key string, value []byte) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	entry := &entry{
		mark:  EntryNormal,
		kLen:  uint32(len(key)),
		vLen:  uint32(len(value)),
		key:   []byte(key),
		value: []byte(value),
	}

	byte, err := encodeEntry(entry)
	if err != nil {
		return
	}

	_, err = db.fd.WriteAt(byte, int64(db.dataFileOffset))
	if err != nil {
		return
	}

	// 更新内存索引
	db.data[key] = &entryOffset{offset: uint32(atomic.LoadUint64(&db.dataFileOffset))}
	atomic.AddUint64(&db.dataFileOffset, uint64(entry.Size()))
}

func (db *KVDB) Get(key string) ([]byte, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	entryOffset, ok := db.data[key]
	if !ok {
		return nil, errors.New("not found")
	}

	buf := make([]byte, entryHeaderSize)
	db.fd.ReadAt(buf, int64(entryOffset.offset))
	entry, _ := decodeEntry(buf)

	val := make([]byte, entry.vLen)
	db.fd.ReadAt(val, int64(entryOffset.offset)+entryHeaderSize+int64(entry.kLen))
	return val, nil
}

func (db *KVDB) Delete(key string) (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Search for key
	if _, ok := db.data[string(key)]; !ok {
		return
	}

	// Write to file
	e := &entry{
		mark:  EntryDeleted,
		kLen:  uint32(len(key)),
		vLen:  0,
		key:   []byte(key),
		value: nil,
	}

	byte, err := encodeEntry(e)
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
