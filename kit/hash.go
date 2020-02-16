package kit

import (
	"crypto/md5"
	"io"
	"os"
)

// Hash sum and size of a file.
func Hashsum(filename string) (md5sum []byte, size int64, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	hash := md5.New()

	if size, err = io.Copy(hash, file); err != nil {
		return nil, 0, err
	}

	return hash.Sum(nil), size, nil
}
