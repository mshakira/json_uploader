package store

// any Uploader should implement these
type GenericStore interface {
	UploadToStore(string, string, []byte) (interface{}, error)
}

// specific to data (implements the interface)
type GenStore struct {
	//
}

func Init(init interface{}) (*GenStore, error) {
	// no op
	return nil, nil
}

func Destroy(*GenStore) error {
	// no op
	return nil
}

// Interface Functions
func (*GenStore) UploadToStore(string, string, []byte) (interface{}, error) {
	return nil, nil
}
