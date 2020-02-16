package kit

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download file and save it to a local path. Creates directories if
// needed. Returns the file size and md5 sum of the file.
func Download(url, path string) (size int64, md5sum []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil, fmt.Errorf("getting %s: %s", url, http.StatusText(resp.StatusCode))
	}

	if dir, _ := filepath.Split(path); dir != "" {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return 0, nil, err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()

	hash := md5.New()
	writer := io.MultiWriter(file, hash)

	if size, err = io.Copy(writer, resp.Body); err != nil {
		return 0, nil, err
	}

	return size, hash.Sum(nil), nil
}
