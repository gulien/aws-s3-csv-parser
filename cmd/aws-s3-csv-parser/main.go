package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws-s3-csv-parser/pkg/storage"
	flag "github.com/spf13/pflag"
)

var version = "snapshot"

func main() {
	fs := flag.NewFlagSet("aws-s3-csv-parser", flag.ExitOnError)
	fs.String("region", "", "Set the AWS region")
	fs.String("bucket", "", "Set the AWS S3 bucket")
	fs.String("key", "", "Set the AWS S3 object key to download")
	fs.Int("timeout", 300, "Set the maximum duration in seconds before timing out the download")

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

	timeout, err := fs.GetInt("timeout")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	w, err := os.Create("tmp.csv")
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	downloader := storage.NewS3PublicDownloader(region, bucket, key)

	err = downloader.Download(ctx, w)
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
