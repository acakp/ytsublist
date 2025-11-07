package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
)

type Channel struct {
	ID   string
	Name string
}

func AddChannel(csvPath, url string) error {
	file, err := os.OpenFile(csvPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error while opening file:")
		return fmt.Errorf("open %s for append: %w", csvPath, err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if err := w.Write([]string{extractID(url), "test"}); err != nil {
		return fmt.Errorf("write record: %w", err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("flush csv: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("fsync %s: %w", csvPath, err)
	}
	return nil
}

func ReadCsv(filePath string) ([]Channel, error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("open file %s: %w", filePath, err)
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	channels, err := readAll(r)
	return channels, nil
}

func readAll(r *csv.Reader) (channels []Channel, err error) {
	for {
		ch, err := r.Read()
		if err == io.EOF {
			return channels, nil
		}
		if err != nil {
			return nil, err
		}
		channels = append(channels, Channel{ID: ch[0], Name: ch[1]})
	}
}

func extractID(link string) string {
	reUsername := regexp.MustCompile(`(?:^|\/)(@[A-Za-z0-9._-]+)`)
	reID := regexp.MustCompile(`(?:^|/)channel/([A-Za-z0-9_-]+)`)
	if match := reUsername.FindStringSubmatch(link); len(match) > 1 {
		return match[1]
	} else if match = reID.FindStringSubmatch(link); len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}
