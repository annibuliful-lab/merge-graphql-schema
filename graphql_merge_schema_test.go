package mergegraphqlschema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadAndParseSchema tests the readAndParseSchema function for both successful and error scenarios.
func TestReadAndParseSchema(t *testing.T) {
	// Create a temporary file with GraphQL schema content
	tmpFile, err := os.CreateTemp("", "schema-*.graphql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `type Query { hello: String }`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test successful read and parse
	doc, err := readAndParseSchema(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}
	if len(doc.Definitions) == 0 {
		t.Errorf("Expected definitions, got none")
	}

	// Test with non-existent file
	_, err = readAndParseSchema("nonexistent.graphql")
	if err == nil {
		t.Error("Expected error for non-existent file, got none")
	}
}

// TestWalkMatch tests the walkMatch function to ensure it correctly finds files matching a specific pattern.
func TestWalkMatch(t *testing.T) {
	// Setup a temporary directory with nested structure
	tmpDir := t.TempDir()

	nestedDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy file matching the pattern
	file := filepath.Join(nestedDir, "match.graphql")
	if err := os.WriteFile(file, []byte("type Query { hello: String }"), 0644); err != nil {
		t.Fatal(err)
	}

	// Execute walkMatch
	matches, err := walkMatch(tmpDir, ".graphql")
	if err != nil {
		t.Fatalf("Error walking the directory: %v", err)
	}
	if len(matches) == 0 {
		t.Errorf("Expected to find matches, but found none")
	}
}

// TestMergeSchemas tests the MergeSchemas function to ensure it merges multiple schemas correctly.
func TestMergeSchemas(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some schema files
	for i := 0; i < 2; i++ {
		file := filepath.Join(tmpDir, fmt.Sprintf("schema%d.graphql", i))
		content := fmt.Sprintf("type Query%d { hello: String }", i)
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	outFile := filepath.Join(tmpDir, "merged.graphql")
	mergedSchema, err := MergeSchemas(tmpDir, ".graphql", outFile)
	if err != nil {
		t.Fatalf("Failed to merge schemas: %v", err)
	}

	expectedSubstrings := []string{"type Query0 {", "type Query1 {"}
	for _, substr := range expectedSubstrings {
		if !strings.Contains(mergedSchema, substr) {
			t.Errorf("Merged schema does not contain expected content '%s'", substr)
		}
	}
}

// Main test entry point to run all tests
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
