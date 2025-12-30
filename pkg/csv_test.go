package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateCSV(t *testing.T) {
	csv := CreateCSV("test.csv", "test_folder/", true)

	if csv == nil {
		t.Fatal("CreateCSV() returned nil")
	}

	if csv.name != "test.csv" {
		t.Errorf("Expected name 'test.csv', got '%s'", csv.name)
	}

	if csv.folder != "test_folder/" {
		t.Errorf("Expected folder 'test_folder/', got '%s'", csv.folder)
	}

	if csv.hasBOM != true {
		t.Errorf("Expected hasBOM true, got %v", csv.hasBOM)
	}
}

func TestCreateCSVWithoutBOM(t *testing.T) {
	csv := CreateCSV("test.csv", "test_folder/", false)

	if csv.hasBOM != false {
		t.Errorf("Expected hasBOM false, got %v", csv.hasBOM)
	}
}

func TestSetHeaders(t *testing.T) {
	csv := CreateCSV("test.csv", "", false)
	headers := []string{"Name", "Age", "City"}

	csv.SetHeaders(headers)

	if len(csv.headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(csv.headers))
	}

	if csv.headers[0] != "Name" {
		t.Errorf("Expected first header 'Name', got '%s'", csv.headers[0])
	}
}

func TestAddData(t *testing.T) {
	csv := CreateCSV("test.csv", "", false)
	data := [][]string{
		{"John", "25", "New York"},
		{"Jane", "30", "Los Angeles"},
	}

	csv.AddData(data)

	if len(csv.data) != 2 {
		t.Errorf("Expected 2 data rows, got %d", len(csv.data))
	}

	if csv.data[0][0] != "John" {
		t.Errorf("Expected first data item 'John', got '%s'", csv.data[0][0])
	}
}

func TestGenerateWithBOM(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"Name", "Age"})
	csv.AddData([][]string{{"John", "25"}, {"Jane", "30"}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	if !strings.HasPrefix(string(content), "\uFEFF") {
		t.Error("Generated file should start with BOM")
	}
}

func TestGenerateWithoutBOM(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), false)
	csv.SetHeaders([]string{"Name", "Age"})
	csv.AddData([][]string{{"John", "25"}, {"Jane", "30"}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	if strings.HasPrefix(string(content), "\uFEFF") {
		t.Error("Generated file should not start with BOM")
	}
}

func TestGenerateHeadersOnly(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"Name", "Age", "City"})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Name,Age,City") {
		t.Error("Generated file should contain headers")
	}
}

func TestGenerateDataOnly(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.AddData([][]string{{"John", "25"}, {"Jane", "30"}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "John,25") {
		t.Error("Generated file should contain data")
	}
}

func TestGenerateEmptyCSV(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	_, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}
}

func TestGenerateWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"姓名", "年龄", "城市"})
	csv.AddData([][]string{{"张三", "25", "北京"}, {"李四", "30", "上海"}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "姓名,年龄,城市") {
		t.Error("Generated file should contain Chinese headers")
	}

	if !strings.Contains(contentStr, "张三,25,北京") {
		t.Error("Generated file should contain Chinese data")
	}
}

func TestGenerateLargeData(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"ID", "Name", "Value"})

	var data [][]string
	for i := 0; i < 1000; i++ {
		data = append(data, []string{AnyToString(i), "Name" + AnyToString(i), "Value" + AnyToString(i)})
	}
	csv.AddData(data)

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	if len(content) == 0 {
		t.Error("Generated file should not be empty")
	}
}

func TestGenerateToInvalidPath(t *testing.T) {
	csv := CreateCSV("test.csv", "/nonexistent/path/", true)
	csv.SetHeaders([]string{"Name", "Age"})
	csv.AddData([][]string{{"John", "25"}})

	err := csv.Generate()
	if !HasError(err) {
		t.Error("Generate() expected error for invalid path, but got nil")
	}
}

func TestGenerateToReadOnlyDirectory(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test as running as root")
	}

	if IsWindows() {
		t.Skip("Skipping test on Windows due to different permission handling")
	}

	tempDir := t.TempDir()

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), false)
	csv.SetHeaders([]string{"Name", "Age"})
	csv.AddData([][]string{{"John", "25"}})

	chmodErr := os.Chmod(tempDir, 0o444)
	if chmodErr != nil {
		t.Fatalf("Failed to make directory read-only: %v", chmodErr)
	}
	defer os.Chmod(tempDir, 0o755)

	genErr := csv.Generate()

	if !HasError(genErr) {
		t.Error("Generate() expected error for read-only directory, but got nil")
	}
}

func TestGetLine(t *testing.T) {
	csv := CreateCSV("test.csv", "", false)
	line := csv.getLine()

	if IsWindows() {
		if line != "\r\n" {
			t.Errorf("Expected '\\r\\n' on Windows, got '%s'", line)
		}
	} else {
		if line != "\n" {
			t.Errorf("Expected '\\n' on non-Windows, got '%s'", line)
		}
	}
}

func TestGenerateWithEmptyData(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"Name", "Age"})
	csv.AddData([][]string{})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Name,Age") {
		t.Error("Generated file should contain headers even with empty data")
	}
}

func TestGenerateWithEmptyFields(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"Name", "Age", "City"})
	csv.AddData([][]string{{"John", "", "New York"}, {"", "25", ""}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "John,,New York") {
		t.Error("Generated file should handle empty fields correctly")
	}
}

func TestGenerateWithCommaInData(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")

	csv := CreateCSV("test.csv", tempDir+string(filepath.Separator), true)
	csv.SetHeaders([]string{"Name", "Description"})
	csv.AddData([][]string{{"Item1", "This is a test"}, {"Item2", "Another test"}})

	err := csv.Generate()
	if HasError(err) {
		t.Fatalf("Generate() failed: %v", err)
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("Failed to read generated file: %v", readErr)
	}

	if len(content) == 0 {
		t.Error("Generated file should not be empty")
	}
}
