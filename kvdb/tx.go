package kvdb

type Tx struct {
	Err   error
	db    *KVDB
	undos []*entry
}

func (this *Tx) Get(key string) ([]byte, error) {
	val, err := this.db.get(key)
	return val, err
}

func (this *Tx) Put(key string, value []byte) {
	oldValue, err := this.db.get(key)
	if err != nil {
		this.rollBack()
	}

	err = this.db.put(key, value)
	if err != nil {
		ety := &entry{
			mark:  EntryNormal,
			kLen:  uint32(len(key)),
			vLen:  uint32(len(oldValue)),
			key:   []byte(key),
			value: []byte(oldValue),
		}
		this.undos = append(this.undos, ety)
		this.rollBack()

		panic("db.put() 操作失败")
	}
}

func (this *Tx) Delete(key string) {
	oldValue, err := this.db.Get(key)
	if err != nil {
		this.rollBack()
	}

	err = this.db.delete(key)

	if err != nil {
		ety := &entry{
			mark:  EntryNormal,
			kLen:  uint32(len(key)),
			vLen:  uint32(len(oldValue)),
			key:   []byte(key),
			value: []byte(oldValue),
		}
		this.undos = append(this.undos, ety)
		this.rollBack()
	}
}

func (this *Tx) rollBack() (e error) {
	for _, ety := range this.undos {
		this.db.Put(string(ety.key), ety.value)
	}

	//this.db = nil
	this.undos = nil
	return nil
}

func (tx *Tx) commit() (e error) {
	return
}
