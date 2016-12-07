package keyStorages

// NilStorage - empty key storage
type NilStorage struct {}

// Get - always return not found
func (c *NilStorage) Get(name string) (value string, err error) {
	value = ""
	err = ErrStorageKeyNotFound
	return
}

// Set - save nothing but return nil error
func (c *NilStorage) Set(name, value string) (err error) {
	err = nil
	return
}

// Remove - remove nothing but return nil error
func (c *NilStorage) Remove(name string) (err error) {
	err = nil
	return
}