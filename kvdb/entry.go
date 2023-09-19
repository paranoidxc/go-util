package kvdb

import "encoding/binary"

type EntryMark int

const (
	EntryNormal EntryMark = iota
	EntryDeleted
)

const (
	entryMetaSize = 1 + 8 + 4 + 4
)

type entryOffset struct {
	offset uint32 // 数据在文件中的偏移量
}

func (e *entry) Size() uint32 {
	return entryMetaSize + e.kLen + e.vLen
}

type entry struct {
	mark  EntryMark // 区分是正常的还是删除的数据
	ttl   int64     // 过期时间
	kLen  uint32    // key的长度
	vLen  uint32    // value的长度
	key   []byte    // key本身
	value []byte    // value本身
}

func NewDelEntry(key string) *entry {
	return &entry{
		mark:  EntryDeleted,
		ttl:   0,
		kLen:  uint32(len(key)),
		vLen:  0,
		key:   []byte(key),
		value: nil,
	}
}

func NewEntry(key string, val []byte) *entry {
	return &entry{
		mark:  EntryNormal,
		ttl:   0,
		kLen:  uint32(len(key)),
		vLen:  uint32(len(val)),
		key:   []byte(key),
		value: []byte(val),
	}
}

func (e *entry) encode() ([]byte, error) {
	buf := make([]byte, e.Size())

	buf[0] = byte(e.mark)
	pos := 1

	binary.BigEndian.PutUint64(buf[pos:pos+8], uint64(e.ttl))
	pos += 8

	binary.BigEndian.PutUint32(buf[pos:pos+4], e.kLen)
	pos += 4

	binary.BigEndian.PutUint32(buf[pos:pos+4], e.vLen)
	pos += 4

	copy(buf[entryMetaSize:], e.key)
	copy(buf[entryMetaSize+e.kLen:], e.value)

	//fmt.Printf("entry buf is: %+v\n", buf)
	return buf, nil
}

func decodeEntry(buf []byte) (*entry, error) {
	pos := 1
	ttl := binary.BigEndian.Uint64(buf[pos : pos+8])
	pos += 8

	kLen := binary.BigEndian.Uint32(buf[pos : pos+4])
	pos += 4

	vLen := binary.BigEndian.Uint32(buf[pos:])
	pos += 4

	e := &entry{
		mark: EntryMark(buf[0]),
		ttl:  int64(ttl),
		kLen: kLen,
		vLen: vLen,
	}

	if len(buf) > entryMetaSize {
		key := make([]byte, kLen)
		val := make([]byte, vLen)
		copy(key, buf[entryMetaSize:entryMetaSize+kLen])
		copy(val, buf[entryMetaSize+kLen:])
		e.key = key
		e.value = val
	}

	return e, nil
}
