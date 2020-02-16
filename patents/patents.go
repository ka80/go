// Package patents reads Google Patents data.
package patents

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Parse Google Patents CSV data and return US patent numbers.
// If nresult > 0, then only return nresult patents.
func USPatents(url string, nresult int) (patents []string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Google %d rejection: %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	scanner := bufio.NewScanner(resp.Body)

	ch := make(chan string) // channel for patent number

	nline := 0
	for scanner.Scan() {
		go func(line string) {
			col0 := strings.SplitN(line, ",", 2)[0]
			data := strings.SplitN(col0, "-", 3)
			if data[0] == "US" {
				ch <- patentID(data[1])
			} else {
				ch <- ""
			}
		}(scanner.Text())

		nline++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for i := 0; i < nline; i++ {
		id := <-ch
		if id != "" {
			patents = append(patents, id)
			nresult--
			if nresult == 0 {
				break
			}
		}
	}

	return patents, nil
}

// Patent ID number from CSV data (PG-PUBs are missing a digit).
func patentID(data string) string {
	if len(data) < 9 || len(data) == 11 {
		return data
	}

	year := data[:4]
	num, err := strconv.Atoi(data[4:])
	if err != nil {
		log.Printf("patentID unable to convert %s: %v", data, err)
		return ""
	}

	return fmt.Sprintf("%s%07d", year, num)
}
