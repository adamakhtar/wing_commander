package filesnippet

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractLines_ExtractsCorrectRangeOfLinesAroundCenterLine(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snippet, err := ExtractLines(filePath, 5, 2)
	if err != nil {
		t.Fatalf("ExtractLines failed: %v", err)
	}

	if snippet.FilePath != filePath {
		t.Errorf("expected FilePath %s, got %s", filePath, snippet.FilePath)
	}

	expectedLines := []Line{
		{Number: 3, Content: "line3", IsCenter: false},
		{Number: 4, Content: "line4", IsCenter: false},
		{Number: 5, Content: "line5", IsCenter: true},
		{Number: 6, Content: "line6", IsCenter: false},
		{Number: 7, Content: "line7", IsCenter: false},
	}

	if len(snippet.Lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(snippet.Lines))
	}

	for i, expected := range expectedLines {
		if snippet.Lines[i].Number != expected.Number {
			t.Errorf("line %d: expected Number %d, got %d", i, expected.Number, snippet.Lines[i].Number)
		}
		if snippet.Lines[i].Content != expected.Content {
			t.Errorf("line %d: expected Content %q, got %q", i, expected.Content, snippet.Lines[i].Content)
		}
		if snippet.Lines[i].IsCenter != expected.IsCenter {
			t.Errorf("line %d: expected IsCenter %v, got %v", i, expected.IsCenter, snippet.Lines[i].IsCenter)
		}
	}
}

func TestExtractLines_ClipsToStartWhenRangeExtendsBeforeFirstLine(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snippet, err := ExtractLines(filePath, 1, 2)
	if err != nil {
		t.Fatalf("ExtractLines failed: %v", err)
	}

	expectedLines := []Line{
		{Number: 1, Content: "line1", IsCenter: true},
		{Number: 2, Content: "line2", IsCenter: false},
		{Number: 3, Content: "line3", IsCenter: false},
	}

	if len(snippet.Lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(snippet.Lines))
	}

	for i, expected := range expectedLines {
		if snippet.Lines[i].Number != expected.Number || snippet.Lines[i].Content != expected.Content || snippet.Lines[i].IsCenter != expected.IsCenter {
			t.Errorf("line %d: expected %+v, got %+v", i, expected, snippet.Lines[i])
		}
	}
}

func TestExtractLines_ClipsToEndWhenRangeExtendsAfterLastLine(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snippet, err := ExtractLines(filePath, 5, 2)
	if err != nil {
		t.Fatalf("ExtractLines failed: %v", err)
	}

	expectedLines := []Line{
		{Number: 3, Content: "line3", IsCenter: false},
		{Number: 4, Content: "line4", IsCenter: false},
		{Number: 5, Content: "line5", IsCenter: true},
	}

	if len(snippet.Lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(snippet.Lines))
	}

	for i, expected := range expectedLines {
		if snippet.Lines[i].Number != expected.Number || snippet.Lines[i].Content != expected.Content || snippet.Lines[i].IsCenter != expected.IsCenter {
			t.Errorf("line %d: expected %+v, got %+v", i, expected, snippet.Lines[i])
		}
	}
}

func TestExtractLines_ReturnsErrorWhenFileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.txt")

	_, err := ExtractLines(filePath, 1, 3)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("expected IsNotExist error, got: %v", err)
	}
}

func TestExtractLines_ReturnsErrorWhenCenterLineIsLessThanOne(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ExtractLines(filePath, 0, 3)
	if err == nil {
		t.Fatal("expected error for line number < 1")
	}
}

func TestExtractLines_ReturnsErrorWhenCenterLineExceedsFileLength(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ExtractLines(filePath, 10, 3)
	if err == nil {
		t.Fatal("expected error for line number > file length")
	}

	emptyFilePath := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(emptyFilePath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("failed to create empty test file: %v", err)
	}

	_, err = ExtractLines(emptyFilePath, 1, 3)
	if err == nil {
		t.Fatal("expected error for empty file with centerLine >= 1")
	}
}

func TestExtractLines_ReturnsErrorWhenSizeIsNegative(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ExtractLines(filePath, 2, -1)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestExtractLines_PreservesEmptyLinesAndCountsThemWhenDeterminingFileLength(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\n\nline3\n\nline5"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snippet, err := ExtractLines(filePath, 3, 1)
	if err != nil {
		t.Fatalf("ExtractLines failed: %v", err)
	}

	expectedLines := []Line{
		{Number: 2, Content: "", IsCenter: false},
		{Number: 3, Content: "line3", IsCenter: true},
		{Number: 4, Content: "", IsCenter: false},
	}

	if len(snippet.Lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(snippet.Lines))
	}

	for i, expected := range expectedLines {
		if snippet.Lines[i].Number != expected.Number || snippet.Lines[i].Content != expected.Content || snippet.Lines[i].IsCenter != expected.IsCenter {
			t.Errorf("line %d: expected %+v, got %+v", i, expected, snippet.Lines[i])
		}
	}

	_, err = ExtractLines(filePath, 6, 1)
	if err == nil {
		t.Fatal("expected error: empty lines should be counted, so line 6 exceeds file length of 5")
	}
}

func TestExtractLines_ReturnsOnlyCenterLineWhenSizeIsZero(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snippet, err := ExtractLines(filePath, 3, 0)
	if err != nil {
		t.Fatalf("ExtractLines failed: %v", err)
	}

	if len(snippet.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(snippet.Lines))
	}

	if snippet.Lines[0].Number != 3 || snippet.Lines[0].Content != "line3" || !snippet.Lines[0].IsCenter {
		t.Errorf("expected {Number: 3, Content: \"line3\", IsCenter: true}, got %+v", snippet.Lines[0])
	}
}
