package layout

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"WaterMark/internal"
	"WaterMark/pkg"
)

func setupTemplateTest(t *testing.T) string {
	tempDir := t.TempDir()

	testLayouts := FrameLayouts{
		List: []FrameLayout{
			{
				Name:            "test_layout_1",
				Type:            "frame",
				Layout:          "horizontal",
				BgColor:         "#FFFFFF",
				LogoRatio:       10,
				TextRatio:       20,
				LogoMarginRight: 5,
				TextOneContent:  "Test Text 1",
				TextTwoContent:  "Test Text 2",
			},
			{
				Name:            "test_layout_2",
				Type:            "frame",
				Layout:          "vertical",
				BgColor:         "#000000",
				LogoRatio:       15,
				TextRatio:       25,
				LogoMarginRight: 10,
				TextOneContent:  "Another Test",
			},
		},
	}

	layoutData, err := json.MarshalIndent(testLayouts, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test layout: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "layout.json"), layoutData, 0o644)
	if err != nil {
		t.Fatalf("Failed to write test layout file: %v", err)
	}

	return tempDir
}

func TestFindLayoutByName(t *testing.T) {
	frameLayouts = &FrameLayouts{
		List: []FrameLayout{
			{
				Name:            "test_layout_1",
				Type:            "frame",
				Layout:          "horizontal",
				BgColor:         "#FFFFFF",
				LogoRatio:       10,
				TextRatio:       20,
				LogoMarginRight: 5,
				TextOneContent:  "Test Text 1",
				TextTwoContent:  "Test Text 2",
			},
			{
				Name:            "test_layout_2",
				Type:            "frame",
				Layout:          "vertical",
				BgColor:         "#000000",
				LogoRatio:       15,
				TextRatio:       25,
				LogoMarginRight: 10,
				TextOneContent:  "Another Test",
			},
		},
	}

	tests := []struct {
		name        string
		layoutName  string
		wantFound   bool
		wantName    string
		wantLayout  string
		wantBgColor string
	}{
		{
			name:        "Find existing layout 1",
			layoutName:  "test_layout_1",
			wantFound:   true,
			wantName:    "test_layout_1",
			wantLayout:  "horizontal",
			wantBgColor: "#FFFFFF",
		},
		{
			name:        "Find existing layout 2",
			layoutName:  "test_layout_2",
			wantFound:   true,
			wantName:    "test_layout_2",
			wantLayout:  "vertical",
			wantBgColor: "#000000",
		},
		{
			name:        "Find non-existing layout",
			layoutName:  "non_existent",
			wantFound:   false,
			wantName:    "",
			wantLayout:  "",
			wantBgColor: "",
		},
		{
			name:        "Find with empty name",
			layoutName:  "",
			wantFound:   false,
			wantName:    "",
			wantLayout:  "",
			wantBgColor: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout, err := FindLayoutByName(tt.layoutName)

			if tt.wantFound {
				if pkg.HasError(err) {
					t.Errorf("FindLayoutByName() unexpected error: %v", err)
				}
				if layout.Name != tt.wantName {
					t.Errorf("FindLayoutByName() Name = %v, want %v", layout.Name, tt.wantName)
				}
				if layout.Layout != tt.wantLayout {
					t.Errorf("FindLayoutByName() Layout = %v, want %v", layout.Layout, tt.wantLayout)
				}
				if layout.BgColor != tt.wantBgColor {
					t.Errorf("FindLayoutByName() BgColor = %v, want %v", layout.BgColor, tt.wantBgColor)
				}
			} else {
				if !pkg.HasError(err) {
					t.Errorf("FindLayoutByName() expected error but got none")
				}
				if err.Code != pkg.LayoutNotFindError.Code {
					t.Errorf("FindLayoutByName() error code = %v, want %v", err.Code, pkg.LayoutNotFindError.Code)
				}
			}
		})
	}
}

func TestMustFindLayoutByName(t *testing.T) {
	frameLayouts = &FrameLayouts{
		List: []FrameLayout{
			{
				Name:            "test_layout_1",
				Type:            "frame",
				Layout:          "horizontal",
				BgColor:         "#FFFFFF",
				LogoRatio:       10,
				TextRatio:       20,
				LogoMarginRight: 5,
				TextOneContent:  "Test Text 1",
				TextTwoContent:  "Test Text 2",
			},
			{
				Name:            "test_layout_2",
				Type:            "frame",
				Layout:          "vertical",
				BgColor:         "#000000",
				LogoRatio:       15,
				TextRatio:       25,
				LogoMarginRight: 10,
				TextOneContent:  "Another Test",
			},
		},
	}

	tests := []struct {
		name       string
		layoutName string
		wantName   string
		wantLayout string
	}{
		{
			name:       "Find existing layout",
			layoutName: "test_layout_1",
			wantName:   "test_layout_1",
			wantLayout: "horizontal",
		},
		{
			name:       "Find another existing layout",
			layoutName: "test_layout_2",
			wantName:   "test_layout_2",
			wantLayout: "vertical",
		},
		{
			name:       "Find non-existing layout returns empty",
			layoutName: "non_existent",
			wantName:   "",
			wantLayout: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := MustFindLayoutByName(tt.layoutName)

			if layout.Name != tt.wantName {
				t.Errorf("MustFindLayoutByName() Name = %v, want %v", layout.Name, tt.wantName)
			}
			if layout.Layout != tt.wantLayout {
				t.Errorf("MustFindLayoutByName() Layout = %v, want %v", layout.Layout, tt.wantLayout)
			}
		})
	}
}

func TestGetAllLayout(t *testing.T) {
	frameLayouts = &FrameLayouts{
		List: []FrameLayout{
			{
				Name:   "test_layout_1",
				Type:   "frame",
				Layout: "horizontal",
			},
			{
				Name:   "test_layout_2",
				Type:   "frame",
				Layout: "vertical",
			},
		},
	}

	layouts := GetAllLayout()

	if len(layouts) != 2 {
		t.Errorf("GetAllLayout() returned %d layouts, want 2", len(layouts))
	}

	layoutNames := make([]string, len(layouts))
	for i, layout := range layouts {
		layoutNames[i] = layout.Name
	}

	expectedNames := []string{"test_layout_1", "test_layout_2"}
	for i, expected := range expectedNames {
		if layoutNames[i] != expected {
			t.Errorf("GetAllLayout() layout[%d].Name = %v, want %v", i, layoutNames[i], expected)
		}
	}
}

func TestLoadandInitLayout(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantError bool
		errorCode int
	}{
		{
			name: "Load valid layout file",
			setupFunc: func() string {
				tempDir := t.TempDir()

				testLayouts := FrameLayouts{
					List: []FrameLayout{
						{
							Name:   "test_layout",
							Type:   "frame",
							Layout: "horizontal",
						},
					},
				}

				layoutData, err := json.MarshalIndent(testLayouts, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal test layout: %v", err)
				}

				err = os.WriteFile(filepath.Join(tempDir, "layout.json"), layoutData, 0o644)
				if err != nil {
					t.Fatalf("Failed to write test layout file: %v", err)
				}

				return tempDir
			},
			wantError: false,
		},
		{
			name: "Load non-existent file",
			setupFunc: func() string {
				tempDir := t.TempDir()
				return tempDir
			},
			wantError: true,
			errorCode: pkg.FILE_NOT_EXIST_ERROR,
		},
		{
			name: "Load invalid JSON file",
			setupFunc: func() string {
				tempDir := t.TempDir()

				err := os.WriteFile(filepath.Join(tempDir, "layout.json"), []byte("invalid json content"), 0o644)
				if err != nil {
					t.Fatalf("Failed to write test layout file: %v", err)
				}

				return tempDir
			},
			wantError: true,
			errorCode: pkg.FILE_NOT_READ_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := tt.setupFunc()

			frameLayouts = &FrameLayouts{}
			layoutPath := filepath.Join(tempDir, "layout.json")

			if !tt.wantError {
				layoutStr, err := os.ReadFile(layoutPath)
				if err != nil {
					t.Fatalf("Failed to read layout file: %v", err)
				}

				err = json.Unmarshal(layoutStr, &frameLayouts)
				if err != nil {
					t.Fatalf("Failed to unmarshal layout: %v", err)
				}
			}

			if tt.wantError {
				if !internal.PathExists(layoutPath) {
					return
				}

				layoutStr, err := os.ReadFile(layoutPath)
				if err != nil {
					return
				}

				err = json.Unmarshal(layoutStr, &frameLayouts)
				if err != nil {
					return
				}
			}
		})
	}
}

func TestReloadandInitLayout(t *testing.T) {
	tempDir := t.TempDir()

	testLayouts := FrameLayouts{
		List: []FrameLayout{
			{
				Name:   "initial_layout",
				Type:   "frame",
				Layout: "horizontal",
			},
		},
	}

	layoutData, err := json.MarshalIndent(testLayouts, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test layout: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "layout.json"), layoutData, 0o644)
	if err != nil {
		t.Fatalf("Failed to write test layout file: %v", err)
	}

	frameLayouts = &FrameLayouts{
		List: []FrameLayout{
			{
				Name: "old_layout",
			},
		},
	}

	layoutPath := filepath.Join(tempDir, "layout.json")
	layoutStr, err := os.ReadFile(layoutPath)
	if err != nil {
		t.Fatalf("Failed to read layout file: %v", err)
	}

	err = json.Unmarshal(layoutStr, &frameLayouts)
	if err != nil {
		t.Fatalf("Failed to unmarshal layout: %v", err)
	}

	if len(frameLayouts.List) != 1 {
		t.Errorf("ReloadandInitLayout() expected 1 layout, got %d", len(frameLayouts.List))
	}

	if frameLayouts.List[0].Name != "initial_layout" {
		t.Errorf("ReloadandInitLayout() layout name = %v, want initial_layout", frameLayouts.List[0].Name)
	}
}
