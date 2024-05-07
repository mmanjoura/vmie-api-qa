// Package main is the main entry point for the MEQA (Mocked Endpoints Quality Assurance) tool.
// It provides functionality for generating test plans based on Swagger API specifications.
// The tool supports various algorithms for generating test plans, including simple, object-based, and path-based.
// It can load Swagger API specifications from a JSON or YAML file, and generate test plans in YAML format.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gbatanov/meqa/mqplan"
	"github.com/gbatanov/meqa/mqswag"
	"github.com/gbatanov/meqa/mqutil"
)

// Constants for the algorithm types
const (
	meqaDataDir = "meqa_data"
	algoSimple  = "simple"
	algoObject  = "object"
	algoPath    = "path"
	algoAll     = "all"
)

// List of available algorithms
var algoList []string = []string{algoSimple, algoObject, algoPath}

func main() {
	// Set up logger
	mqutil.Logger = mqutil.NewStdLogger()

	// Default file paths
	swaggerJSONFile := filepath.Join(meqaDataDir, "swagger.yml")

	// Define command-line flags
	meqaPath := flag.String("d", meqaDataDir, "the directory where we put the generated files")
	swaggerFile := flag.String("s", swaggerJSONFile, "the swagger.yml file location")
	algorithm := flag.String("a", "all", "the algorithm - simple, object, path, all")
	verbose := flag.Bool("v", false, "turn on verbose mode")
	whitelistFile := flag.String("w", "", "the whitelist.txt file location")

	// Parse command-line flags
	flag.Parse()

	// Run the program with the provided options
	run(meqaPath, swaggerFile, algorithm, verbose, whitelistFile)
}

// Function to run the program with the provided options
func run(meqaPath *string, swaggerFile *string, algorithm *string, verbose *bool, whitelistFile *string) {
	// Set verbose mode
	mqutil.Verbose = *verbose

	// Validate swagger file path
	swaggerJsonPath := *swaggerFile
	if fi, err := os.Stat(swaggerJsonPath); os.IsNotExist(err) || fi.Mode().IsDir() {
		fmt.Printf("Can't load swagger file at the following location %s", swaggerJsonPath)
		os.Exit(1)
	}

	// Validate whitelist file path
	whitelistPath := *whitelistFile
	var whitelist map[string]bool
	if len(whitelistPath) > 0 {
		if fi, err := os.Stat(whitelistPath); os.IsNotExist(err) || fi.Mode().IsDir() {
			fmt.Printf("Can't load whitelist file at the following location %s", whitelistPath)
			os.Exit(1)
		}
		wl, err := mqswag.GetWhitelistSuites(whitelistPath)
		whitelist = wl
		if err != nil {
			mqutil.Logger.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	}

	// Validate test plan directory path
	testPlanPath := *meqaPath
	if fi, err := os.Stat(testPlanPath); os.IsNotExist(err) {
		err = os.Mkdir(testPlanPath, 0755)
		if err != nil {
			fmt.Printf("Can't create the directory at %s\n", testPlanPath)
			os.Exit(1)
		}
	} else if !fi.Mode().IsDir() {
		fmt.Printf("The specified location is not a directory: %s\n", testPlanPath)
		os.Exit(1)
	}

	// Load swagger.json
	swagger, err := mqswag.CreateSwaggerFromURL(swaggerJsonPath, *meqaPath)
	if err != nil {
		mqutil.Logger.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	// Create and populate DAG (Directed Acyclic Graph)
	dag := mqswag.NewDAG()
	err = swagger.AddToDAG(dag)
	if err != nil {
		mqutil.Logger.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	// Sort and check weight of DAG
	dag.Sort()
	dag.CheckWeight()

	// Generate test plans based on selected algorithms
	var plansToGenerate []string
	if *algorithm == algoAll {
		plansToGenerate = algoList
	} else {
		plansToGenerate = append(plansToGenerate, *algorithm)
	}

	for _, algo := range plansToGenerate {
		var testPlan *mqplan.TestPlan
		switch algo {
		case algoPath:
			testPlan, err = mqplan.GeneratePathTestPlan(swagger, dag, whitelist)
		case algoObject:
			testPlan, err = mqplan.GenerateTestPlan(swagger, dag)
		default:
			testPlan, err = mqplan.GenerateSimpleTestPlan(swagger, dag)
		}
		if err != nil {
			mqutil.Logger.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
		testPlanFile := filepath.Join(testPlanPath, algo+".yml")
		err = testPlan.DumpToFile(testPlanFile)
		if err != nil {
			mqutil.Logger.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println("Test plans generated at:", testPlanFile)
	}
}

