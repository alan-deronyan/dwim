package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf"
	"github.com/deronyan-llc/dwim/internal/clients"
)

func main() {
	// Parse the RDF schema files
	// read all schema files in the schemas directory
	inputDir := os.Args[1]
	//outputDir := os.Args[2]
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//var parser parsers.Parser = parsers.RDFParser{}
	//var goSrcGenerator generators.SrcGenerator = golang.GoSrcGenerator{}
	flureeClient := clients.NewFlureeClient("http://localhost:58090")

	// convert RDF Turtle/TTL file into LD-JSON

	for _, file := range files {
		inputName := inputDir + "/" + file.Name()

		// get the file name and extension
		fileSplits := strings.Split(file.Name(), ".")
		fileSplitsLen := len(fileSplits)
		fileName := fileSplits[fileSplitsLen-2]
		fileExt := fileSplits[fileSplitsLen-1]

		// parse triples as either TTL or LD-JSON
		var triples []rdf.Triple
		reader, err := os.Open(inputName)
		if err != nil {
			fmt.Printf("Error opening file(%s): %v\n", file.Name(), err)
			continue
		}

		var buf []byte
		// read the TTL or LD-JSON file
		if fileExt == "ttl" {
			triples, err = rdf.DecodeTTL(reader)
			if err != nil {
				fmt.Printf("Error decoding TTL for file(%s): %v\n", file.Name(), err)
				continue
			}

			// encode the triples as LD-JSON
			out, err := os.CreateTemp("", "dwim-"+fileName+"-*.ld-json")
			if err != nil {
				fmt.Printf("Error creating temp file for JSON-LD for file(%s): %v\n", file.Name(), err)
				continue
			}
			tmpPath := out.Name()
			err = rdf.EncodeJSONLD(triples, out)
			out.Close()
			if err != nil {
				fmt.Printf("Error encoding JSON-LD for file(%s): %v\n", file.Name(), err)
				continue
			}

			in, err := os.Open(tmpPath)
			if err != nil {
				fmt.Printf("Error opening temp file(%s): %v\n", tmpPath, err)
				continue
			}
			inInfo, err := in.Stat()
			if err != nil {
				fmt.Printf("Error getting file info for temp file(%s): %v\n", tmpPath, err)
				continue
			}

			buf = make([]byte, inInfo.Size())
			_, err = in.Read(buf)
			if err != nil {
				fmt.Printf("Error creating ledger(%s): %v\n", file.Name(), err)
				continue
			}

		} else if fileExt == "json" || fileExt == "ld-json" || fileExt == "ldjson" || fileExt == "jsonld" || fileExt == "json-ld" {
			inInfo, err := reader.Stat()
			if err != nil {
				fmt.Printf("Error getting file info for temp file(%s): %v\n", inputName, err)
				continue
			}

			buf = make([]byte, inInfo.Size())
			_, err = reader.Read(buf)
			if err != nil {
				fmt.Printf("Error reading ld-json file(%s): %v\n", file.Name(), err)
				continue
			}
		} else {
			fmt.Printf("Error: Unsupported file extension(%s) for file(%s)\n", fileExt, file.Name())
			continue
		}

		// create a ledger in fluree from LD-JSON
		ledger := "schemas/" + fileName
		fmt.Printf("Loading file `%s`into fluree ledger `%s`... ", fileName, ledger)
		_, err = flureeClient.Create(ledger, nil, []string{string(buf)})
		if err != nil {
			fmt.Printf("Error creating ledger(%s): %v\n", file.Name(), err)
			continue
		}
		fmt.Printf("...done.\n")

		/*
			fmt.Printf("Parsing RDF schema for file(%s)\n", file.Name())
			context, err := parser.Parse(inputName)
			if err != nil {
				fmt.Printf("Error parsing RDF schema for file(%s): %v\n", file.Name(), err)
				continue
			}

			context.OutputPath = outputDir
			if err := goSrcGenerator.Generate(context); err != nil {
				fmt.Printf("Error generating GoLang code for file(%s): %v\n", file.Name(), err)
				continue
			}
		*/
	}
}
