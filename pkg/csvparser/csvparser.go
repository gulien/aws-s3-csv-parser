package csvparser

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// Options gathers the parsing options.
type Options struct {
	// Comma is the separator to use while parsing.
	Comma rune
	// SkipFirstLine skips the first line if true.
	SkipFirstLine bool
}

// Callback is a function which is called while parsing a CSV file.
type Callback struct {
	// Every is the number of records to read before calling the callback.
	Every int

	// Do is the function called after Every number of records has been
	// reached. The records length is equal to Every.
	Do func(ctx context.Context, records [][]string) error
}

// Parse parses a CSV file and calls the given callback.
func Parse(ctx context.Context, path string, callback Callback, options Options) error {
	if callback.Do == nil {
		return errors.New("nil callback Do function")
	}

	if callback.Every <= 0 {
		return errors.New("negative or zero chunk of records to process in the callback")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open CSV file: %v", err)
	}

	r := csv.NewReader(f)
	r.Comma = options.Comma

	records := make([][]string, 0)
	firstLine := true

	for {
		// As the CSV file could be huge, we check if the provided context has
		// been canceled or if its deadline has exceeded.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		record, err := r.Read()

		if err == io.EOF {
			// End of file.

			if len(records) == 0 {
				// No records anymore, we're done.
				return nil
			}

			// Call the callback for the last records.
			errCallback := callback.Do(ctx, records)
			if errCallback != nil {
				return fmt.Errorf("call callback for last records: %v", errCallback)
			}
		}

		if err != nil {
			return fmt.Errorf("get record: %v", err)
		}

		if firstLine {
			firstLine = false

			if options.SkipFirstLine {
				continue
			}
		}

		// Add the records to the current chunk.
		records = append(records, record)

		if len(records) == callback.Every {
			// Chunk reached its limit, let's call the callback.
			errCallback := callback.Do(ctx, records)

			// Create a new chunk.
			records = make([][]string, 0)

			if errCallback != nil {
				return fmt.Errorf("call callback for records: %v", errCallback)
			}
		}
	}
}
