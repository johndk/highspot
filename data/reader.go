package data

type Reader interface {
	Read() ([]byte, error)
}
