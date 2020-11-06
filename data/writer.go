package data

type Writer interface {
	Write(data []byte) error
}
