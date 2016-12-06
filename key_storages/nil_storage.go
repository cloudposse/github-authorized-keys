package key_storages

type NilStorage struct {}

func (c *NilStorage) Get(name string) (value string, err error) {
	value = ""
	err = ErrStorageKeyNotFound
	return
}

func (c *NilStorage) Set(name, value string) (err error) {
	err = nil
	return
}

func (c *NilStorage) Remove(name string) (err error) {
	err = nil
	return
}