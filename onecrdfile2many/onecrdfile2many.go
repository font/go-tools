package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO: Add flag to specify input file name as well as prefix name for each
// CRD file generated e.g. 'federationv2'.

func writeCRD(name string, buffer bytes.Buffer) error {
	filename := "federationv2-" + strings.ToLower(name) + ".crd.yaml"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(buffer.String())
	f.Sync()
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	validKind := regexp.MustCompile(`^    kind: ([a-zA-Z]+)$`)
	var crdName, line string
	var crdBuffer bytes.Buffer
	for scanner.Scan() {
		line = scanner.Text()
		if line == "---" { // Beginning of new YAML config
			// Write new CRD file
			err := writeCRD(crdName, crdBuffer)
			crdBuffer.Reset()
			if err != nil {
				fmt.Printf("Error writing file for CRD %v: %v", crdName, err)
				panic(err)
			}
			continue
		} else {
			// Process YAML config
			crdBuffer.WriteString(line + "\n")
			if validKind.MatchString(line) {
				strSlice := validKind.FindStringSubmatch(line)
				crdName = strSlice[1]
			}
		}
	}
	// Write the last CRD file
	err := writeCRD(crdName, crdBuffer)
	if err != nil {
		fmt.Printf("Error writing file for CRD %v: %v", crdName, err)
	}
}
