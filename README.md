# GraphQL Schema Merger

This Go package provides tools for reading, parsing, and merging GraphQL schema files from a specified directory into a single file. It is designed to support developers in combining multiple GraphQL schema files effortlessly, facilitating easier management and integration of complex GraphQL schemas.

### Getting the Package

To install this package, execute the following command in your terminal:

`go get github.com/annibuliful-lab/merge-graphql-schema`

##### Example: Merging Schemas

Here's a simple example demonstrating how to use this package to merge all GraphQL schema files in a specified directory:

```go
package main

import (
    "log"
    merger "github.com/annibuliful-lab/merge-graphql-schema"

)

func main() {
    rootDir := "./path/to/schemas"
    pattern := ".graphql"
    outputFile := "merged_schema.graphql"

    // Merge all schema files matching the pattern
    if _, err := merger.MergeSchemas(rootDir, pattern, outputFile); err != nil {
        log.Fatalf("Failed to merge schemas: %v", err)
    }

    log.Printf("Schemas successfully merged into '%s'", outputFile)

}

```
