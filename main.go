package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
)

// readAndParseSchema reads a GraphQL schema file and returns its AST document.
func readAndParseSchema(filename string) (*ast.Document, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parser.Parse(parser.ParseParams{Source: string(content)})
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, pattern) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// mergeSchemas takes a glob pattern and returns a merged schema.
func MergeSchemas(root, pattern, outfile string) (string, error) {
	var mergedAST *ast.Document

	// Find all files matching the glob pattern
	files, err := walkMatch(root, pattern)

	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files matched the pattern")
	}

	for _, file := range files {
		doc, err := readAndParseSchema(file)
		if err != nil {
			return "", fmt.Errorf("error parsing file %s: %v", file, err)
		}

		if mergedAST == nil {
			// Initialize mergedAST with the first document
			mergedAST = doc
		} else {
			// Append definitions from each parsed document
			mergedAST.Definitions = append(mergedAST.Definitions, doc.Definitions...)
		}
	}

	// Print the merged AST back to a schema string
	mergedSchema := printer.Print(mergedAST).(string)

	err = os.WriteFile(outfile, []byte(mergedSchema), 0777)
	if err != nil {
		panic(err)
	}

	return mergedSchema, nil
}
