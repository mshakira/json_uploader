package loadfmt

// any Uploader should implement these
type Upload interface {
	UploadPayload([]byte) (err error)
}

// specific to data (implements the interface)
type DataFmt struct {
	//
}

func Init(init interface{}) (*DataFmt, error) {
	// no op
	return nil, nil
}

func Destroy(*DataFmt) error {
	// no op
	return nil
}

// Interface Functions
func (*DataFmt) UploadPayload([]byte) (err error) {
	return nil
}
