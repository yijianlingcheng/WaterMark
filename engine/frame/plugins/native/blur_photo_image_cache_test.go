package native

import (
	"sync"
	"testing"
	"time"
)

func TestInitBlurImageCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Init blur image cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initBlurImageCache()

			if blurImageFileListMaps == nil {
				t.Error("initBlurImageCache() should initialize blurImageFileListMaps")
			}

			if blurImageFileListMaps.blurImageMap == nil {
				t.Error("initBlurImageCache() should initialize blurImageMap")
			}

			if blurImageFileListMaps.blurImageWriteMap == nil {
				t.Error("initBlurImageCache() should initialize blurImageWriteMap")
			}
		})
	}
}

func TestAddBlurImageToWriteList(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Add blur image to write list",
			path: "test_blur_image_1.jpg",
		},
		{
			name: "Add another blur image to write list",
			path: "test_blur_image_2.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBlurImageToWriteList(tt.path)

			if !checkBlurImageInWriteList(tt.path) {
				t.Errorf("addBlurImageToWriteList() should add %s to write list", tt.path)
			}

			if checkBlurImageExist(tt.path) {
				t.Errorf("addBlurImageToWriteList() should not add %s to blurImageMap", tt.path)
			}
		})
	}
}

func TestMoveToBlurImageFileList(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Move blur image from write list to file list",
			path: "test_move_blur_image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBlurImageToWriteList(tt.path)

			if !checkBlurImageInWriteList(tt.path) {
				t.Errorf("Before move: %s should be in write list", tt.path)
			}

			moveToBlurImageFileList(tt.path)

			if checkBlurImageInWriteList(tt.path) {
				t.Errorf("moveToBlurImageFileList() should remove %s from write list", tt.path)
			}

			if !checkBlurImageExist(tt.path) {
				t.Errorf("moveToBlurImageFileList() should add %s to blurImageMap", tt.path)
			}
		})
	}
}

func TestCheckBlurImageExist(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		setup    func()
		expected bool
	}{
		{
			name:     "Check non-existent blur image",
			path:     "non_existent_blur.jpg",
			setup:    func() {},
			expected: false,
		},
		{
			name: "Check existing blur image",
			path: "existing_blur.jpg",
			setup: func() {
				addBlurImageToWriteList("existing_blur.jpg")
				moveToBlurImageFileList("existing_blur.jpg")
			},
			expected: true,
		},
		{
			name: "Check blur image in write list",
			path: "in_write_list_blur.jpg",
			setup: func() {
				addBlurImageToWriteList("in_write_list_blur.jpg")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got := checkBlurImageExist(tt.path)

			if got != tt.expected {
				t.Errorf("checkBlurImageExist() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckBlurImageInWriteList(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		setup    func()
		expected bool
	}{
		{
			name:     "Check non-existent blur image in write list",
			path:     "non_existent_write.jpg",
			setup:    func() {},
			expected: false,
		},
		{
			name: "Check blur image in write list",
			path: "in_write_list.jpg",
			setup: func() {
				addBlurImageToWriteList("in_write_list.jpg")
			},
			expected: true,
		},
		{
			name: "Check blur image moved to file list",
			path: "moved_to_file_list.jpg",
			setup: func() {
				addBlurImageToWriteList("moved_to_file_list.jpg")
				moveToBlurImageFileList("moved_to_file_list.jpg")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got := checkBlurImageInWriteList(tt.path)

			if got != tt.expected {
				t.Errorf("checkBlurImageInWriteList() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWaitBlurImageInList(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Wait for blur image to be added to list",
			path: "wait_for_blur.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(200 * time.Millisecond)
				addBlurImageToWriteList(tt.path)
				moveToBlurImageFileList(tt.path)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				waitBlurImageInList(tt.path)

				if !checkBlurImageExist(tt.path) {
					t.Errorf("waitBlurImageInList() should wait until %s exists", tt.path)
				}
			}()

			wg.Wait()
		})
	}
}

func TestBlurImageCache_ConcurrentOperations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Concurrent add and check operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			numGoroutines := 10
			pathsPerGoroutine := 100

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()
					for j := 0; j < pathsPerGoroutine; j++ {
						path := "concurrent_test_" + string(rune(goroutineID)) + "_" + string(rune(j)) + ".jpg"
						addBlurImageToWriteList(path)
						checkBlurImageInWriteList(path)
						moveToBlurImageFileList(path)
						checkBlurImageExist(path)
					}
				}(i)
			}

			wg.Wait()

			if blurImageFileListMaps == nil {
				t.Error("Concurrent operations should not cause blurImageFileListMaps to be nil")
			}
		})
	}
}

func TestBlurImageCache_Reset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Reset blur image cache state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBlurImageToWriteList("reset_test_1.jpg")
			addBlurImageToWriteList("reset_test_2.jpg")
			moveToBlurImageFileList("reset_test_1.jpg")

			blurImageFileListMaps.blurImageMtx.Lock()
			blurImageFileListMaps.blurImageMap = make(map[string]struct{})
			blurImageFileListMaps.blurImageWriteMap = make(map[string]struct{})
			blurImageFileListMaps.blurImageMtx.Unlock()

			if checkBlurImageExist("reset_test_1.jpg") {
				t.Error("After reset, blurImageMap should be empty")
			}

			if checkBlurImageInWriteList("reset_test_2.jpg") {
				t.Error("After reset, blurImageWriteMap should be empty")
			}
		})
	}
}
