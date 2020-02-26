package main

/* Go Module which allows you convert between HCL/JSON/YAML */

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "github.com/ghodss/yaml"
	hcl "github.com/hashicorp/hcl"
	hclPrinter "github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
)

// Format type
type Format int

const (
	HCL Format = iota
	JSON
	YAML
	AUTO
)

// RunConfig of the converter
type RunConfig struct {
	formatFrom Format
	formatTo Format
	translateIAMPolicy bool
}

func main() {
	cfg := parseArg()

	object, err := readIn(os.Stdin, cfg.formatFrom)

	fmt.Println("%+v", object)

	if err != nil {
		fmt.Printf("Unable to parse the input, %s", err)
		os.Exit(1)
	}

	printOut(cfg.formatTo, object, os.Stdout)
}

func printOut(format Format, obj interface{}, file *os.File) error {
	switch (format) {
	case HCL:
		return printHCL(file, obj)
	case JSON:
		return printJSON(file, obj)
	case YAML:
		return printYAML(file, obj)
	default:
		return fmt.Errorf("No valid output format")
	}
}

func printYAML(file *os.File, obj interface{}) error {
	yamlOut, err := yaml.Marshal(obj)

	if err != nil {
		return fmt.Errorf("Unable to Marshal yaml: %s", err)
	}

	fmt.Fprintln(file, string(yamlOut))

	return nil
}

func printJSON(file *os.File, obj interface{}) error {
	jsonOut, err := json.MarshalIndent(obj, "", "  ")

	if err != nil {
		return fmt.Errorf("Unable to Marshal json: %s", err)
	}

	fmt.Fprintln(file, string(jsonOut))

	return nil
}

func printHCL(file *os.File, obj interface{}) error {
	// TODO: Find a way to Marshal HCL directly
	jsonOut, err := json.Marshal(obj)

	if err != nil {
		return fmt.Errorf("Unable to Marshal json: %s", err)
	}

	hclAST, err := jsonParser.Parse(jsonOut)

	if err != nil {
		return fmt.Errorf("Unable to Parse JSON to hclAST: %s", err)
	}

	err = hclPrinter.Fprint(file, hclAST)
	// Extra newline to be consistent
	fmt.Println()

	if err != nil {
		return fmt.Errorf("Unable to print HCL: %s", err)
	}

	return nil
}

func readIn(file *os.File, format Format) (interface{}, error) {
	input, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file/stream: %s", err)
	}

	switch (format) {
	case HCL:
		return readHCL(input)
	case JSON:
		return readJSON(input)
	case YAML:
		return readYAML(input)
	case AUTO:
		fallthrough
	default:
		// Try parse, JSON -> YAML -> HCL
		if json, err := readJSON(input); err == nil {
			return json, nil
		} else if yaml, err := readYAML(input); err == nil {
			return yaml, nil
		} else {
			return readHCL(input)
		}
	}
}

func readHCL(input []byte) (interface{}, error) {
	var buffer interface{}

	err := hcl.Unmarshal(input, &buffer)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse HCL: %s", err)
	}

	return buffer, nil
}

func readYAML(input []byte) (interface{}, error) {
	jsonOut, err := yaml.YAMLToJSON(input)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse YAML: %s", err)
	}

	var buffer interface{}
	err = json.Unmarshal(jsonOut, &buffer)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json: %s", err)
	}

	return buffer, nil
}

func readJSON(input []byte) (interface{}, error) {
	var buffer interface{}
	err := json.Unmarshal(input, &buffer)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json: %s", err)
	}

	return buffer, nil

}

func stringToFormat(s string) (Format, bool) {
	switch (s) {
	case "JSON":
		return JSON, true
	case "YAML", "YML":
		return YAML, true
	case "HCL", "TF":
		return HCL, true
	case "AUTO":
		return AUTO, true
	default:
		return -1, false
	}
}


func parseArg() RunConfig {
	errored := false
	config := RunConfig{}
	fromFormatIn := flag.String("from", "AUTO", "From format, <JSON|YAML|HCL|AUTO>")
	toFormatIn := flag.String("to", "AUTO", "From format, <JSON|YAML|HCL>")
	translateIAMIn := flag.Bool("iam", false, "Whether or not to translate IAM json to HCL")

	flag.Parse()

	if fromFormat, ok := stringToFormat(*fromFormatIn); !ok {
		fmt.Println("From format must be one of <JSON|YAML|HCL|AUTO>")
		errored = true
	} else {
		config.formatFrom = fromFormat
	}

	if toFormat, ok := stringToFormat(*toFormatIn); !ok || toFormat == AUTO {
		fmt.Println("To format must be one of <JSON|YAML|HCL>")
		errored = true
	} else {
		config.formatTo = toFormat
	}

	config.translateIAMPolicy = *translateIAMIn

	if errored {
		os.Exit(1)
	}

	return config
}