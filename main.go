package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cirocosta/go-mod-dep-source-finder/homepage"
	"github.com/cirocosta/go-mod-dep-source-finder/parser"
	"github.com/cirocosta/go-mod-dep-source-finder/resolver"
	"github.com/cirocosta/go-mod-dep-source-finder/result"
	sa "github.com/cirocosta/go-mod-dep-source-finder/synchronousarray"
	"golang.org/x/sync/errgroup"
)

// [cc]: add flags for requests timeout
// [cc]: add flags for repository retrieval timeout
// [cc]: add flags for limiting concurrency

func failWithHelp(messageFormat string, args ...interface{}) {
	const usageFormat = `Usage: %s <line>\n`

	fmt.Fprintf(os.Stderr, messageFormat, args...)
	fmt.Fprintf(os.Stderr, usageFormat, os.Args[0])

	os.Exit(1)
}

func execute(ctx context.Context, text string, syncArray *sa.SynchronousArray) error {
	line, err := parser.ParseLine(text)
	if err != nil {
		return err
	}

	location, err := resolver.Resolve(
		context.Background(),
		"https://"+line.Dependency,
	)
	if err != nil {
		return err
	}

	url, unknownHost, err := homepage.Find(location.URL, line.Reference)
	if err != nil {
		return err
	}

	res := result.Result{
		Original:           text,
		DiscoveredHomePage: url,
	}

	if unknownHost {
		res.Problem = "unknwon host: " + location.URL
	}

	syncArray.Add(res)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		failWithHelp("error: not enough arguments.\n")
	}

	var reader io.Reader

	line := os.Args[1]
	if line == "-" {
		reader = os.Stdin
	} else {
		reader = bytes.NewBufferString(line)
	}

	scanner := bufio.NewScanner(reader)
	group, ctx := errgroup.WithContext(context.Background())

	syncArray := sa.SynchronousArray{}

	for scanner.Scan() {
		text := scanner.Text()

		group.Go(func() error {
			return execute(ctx, text, &syncArray)
		})
	}

	err := group.Wait()
	if err != nil {
		log.Panic(err)
	}

	jsonBytes, err := json.Marshal(syncArray.Results)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(jsonBytes))
}
