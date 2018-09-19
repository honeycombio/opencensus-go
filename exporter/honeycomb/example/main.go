package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/honeycombio/opencensus-exporter/honeycomb"
	"go.opencensus.io/trace"
)

func main() {
	// YOUR-HONEYCOMB-API-KEY is your Honeycomb API key (also known as your write key)
	// YOUR-HONEYCOMB-DATASET is the name of your Honeycomb dataset to send trace events to
	//
	// Don't have a Honeycomb account? Sign up at https://ui.honeycomb.io/signup
	exporter := honeycomb.NewExporter("YOUR-HONEYCOMB-API-KEY", "YOUR-HONEYCOMB-DATASET")
	defer exporter.Close()

	trace.RegisterExporter(exporter)

	br := bufio.NewReader(os.Stdin)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1.0)})

	// repl is the read, evaluate, print, loop
	for {
		if err := readEvaluateProcess(br); err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
	}
}

func readEvaluateProcess(br *bufio.Reader) error {
	ctx, span := trace.StartSpan(context.Background(), "repl")
	defer span.End()

	fmt.Printf("> ")

	line, err := readLine(ctx, br)
	if err != nil {
		return err
	}

	out, err := processLine(ctx, line)
	if err != nil {
		return err
	}
	fmt.Printf("< %s\n\n", out)
	return nil
}

func readLine(ctx context.Context, br *bufio.Reader) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "readLine")
	defer span.End()

	line, _, err := br.ReadLine()
	if err != nil {
		return nil, err
	}

	return line, err
}

// processLine takes in a line of text and
// transforms it. Currently it just capitalizes it.
func processLine(ctx context.Context, in []byte) (out []byte, err error) {
	_, span := trace.StartSpan(ctx, "processLine")
	defer span.End()

	return bytes.ToUpper(in), nil
}
