// aws-s3-csv-parser A simple CLI which downloads a CSV file from S3 and
// inserts its data into a MySQL database.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gulien/aws-s3-csv-parser/pkg/csvparser"
	"github.com/gulien/aws-s3-csv-parser/pkg/datastorage"
	"github.com/gulien/aws-s3-csv-parser/pkg/filestorage"
	flag "github.com/spf13/pflag"
)

var version = "snapshot"

func main() {
	fs := flag.NewFlagSet("aws-s3-csv-parser", flag.ExitOnError)
	fs.String("region", "", "Set the AWS region")
	fs.String("bucket", "", "Set the AWS S3 bucket")
	fs.String("key", "", "Set the AWS S3 object key to download")
	fs.Bool("skip-download", false, "Skip download if true (i.e, file has already been downloaded)")
	fs.Int("timeout", 300, "Set the maximum duration in seconds before timing out")

	// Parses the flags...
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("[SYSTEM] version %s\n", version)

	// and gets their values.
	region, err := fs.GetString("region")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	bucket, err := fs.GetString("bucket")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	key, err := fs.GetString("key")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	skipDownload, err := fs.GetBool("skip-download")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	timeout, err := fs.GetInt("timeout")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	timeoutDuration := time.Duration(timeout) * time.Second

	inserter, err := datastorage.NewMySQLInserter(
		os.Getenv("MYSQL_DATABASE_URL"),
		"events",
		[]string{
			"id",
			"timestamp",
			"email",
			"country",
			"ip",
			"uri",
			"action",
			"tags",
		},
		timeoutDuration,
	)

	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	const tmpFilename = "tmp.csv"

	if !skipDownload {
		w, err := os.Create(tmpFilename)
		if err != nil {
			fmt.Printf("[FATAL] %s\n", err)
			os.Exit(1)
		}

		downloader := filestorage.NewS3PublicDownloader(region, bucket, key)

		fmt.Println("[INFO] Downloading file...")

		err = downloader.Download(ctx, w)
		if err != nil {
			fmt.Printf("[FATAL] %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("[INFO] Parsing CSV...")

	err = csvparser.Parse(
		ctx,
		tmpFilename,
		csvparser.Callback{
			Every: 1000,
			Do: func(ctx context.Context, records [][]string) error {
				fmt.Printf("[INFO] Processing %d records...\n", len(records))

				return inserter.Insert(ctx, records, datastorage.InsertOptions{
					Formatter: func(value string) string {
						return strings.ReplaceAll(value, "'", "\"")
					},
					IgnoreDuplicateError: true,
				})
			},
		},
		csvparser.Options{
			Comma:         ',',
			SkipFirstLine: true,
		},
	)

	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
