package main

import (
	"flag"
	"fmt"
	"highspot/data"
	"highspot/data/file"
	"highspot/data/http"
	"log"
	"os"
)

type CommandLine struct {
	InputUrl   string
	InputPath  string
	Changes    string
	OutputPath string
	Help       bool
}

var cmdline = CommandLine{}

// The main program.
func main() {
	// Parse the command line arguments.
	flag.Parse()

	if cmdline.Help {
		// Print usage and exit 0.
		flag.Usage()
		os.Exit(0)
	}

	//
	// Create an ingester and execute the take-home exercise
	//

	ingester := data.NewIngestor(getInputReader(), file.NewClient(cmdline.Changes), file.NewClient(cmdline.OutputPath))

	err := ingester.Execute()
	if err != nil {
		log.Fatalf("Error encountered. %v", err)
	}

	log.Printf("The output file %v was successfully created.", cmdline.OutputPath)
}

func getInputReader() data.Reader {
	if len(cmdline.InputPath) != 0 {
		return file.NewClient(cmdline.InputPath)
	} else {
		return http.NewClient((cmdline.InputUrl))
	}
}

// Initialize the command line arguments. Print usage highspot -h.
func init() {
	flag.StringVar(&cmdline.InputUrl, "u", "https://gist.githubusercontent.com/jmodjeska/0679cf6cd670f76f07f1874ce00daaeb/raw/a4ac53fa86452ac26d706df2e851fb7d02697b4b/mixtape-data.json", "The input file URL.")
	flag.StringVar(&cmdline.InputPath, "p", "", "The input file path.")
	flag.StringVar(&cmdline.OutputPath, "o", "output.json", "The output file path.")
	flag.StringVar(&cmdline.Changes, "c", "changes.json", "The changes file.")
	flag.BoolVar(&cmdline.Help, "h", false, "Print the help text.")
	flag.Usage = printUsage
}

func printUsage() {
	fmt.Print("The Highspot take-home coding exercise.\n\n")
	fmt.Print("Usage: highspot [arguments]\n\n")
	fmt.Print("The arguments are:\n\n")
	flag.PrintDefaults()
}
