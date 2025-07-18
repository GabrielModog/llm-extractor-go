package main

import (
	"flag"
	"fmt"

	ocr "github.com/GabrielModog/llm-extractor-go/internal"
)

type ExtractConfig struct {
	data   string
	output string
}

func main() {
	var extractConfig ExtractConfig

	flag.StringVar(&extractConfig.data, "data", "assets/", "Data directory where is the data files to extract")
	flag.StringVar(&extractConfig.output, "output", "output/", "Output directory")

	flag.Parse()

	fmt.Printf("%+v\n", extractConfig)

	ocr.ExtractFromPDF(extractConfig.data, extractConfig.output)
}
