package native

import (
	"testing"

	"WaterMark/layout"
)

func TestSimpleBorderFactory_CreateBorder(t *testing.T) {
	layout.GetAllLayout()

	tests := []struct {
		name       string
		f          *SimpleBorderFactory
		layoutName string
		wantNil    bool
		wantType   any
		checkProps func(any) bool
	}{
		{
			name:       "Create fixed bottom logo text left layout",
			f:          &SimpleBorderFactory{},
			layoutName: "固定布局左下logo模板",
			wantNil:    false,
			wantType:   &fixedBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*fixedBottomLogoTextLayoutBorder)
				return !border.IsRight
			},
		},
		{
			name:       "Create fixed bottom logo text right layout",
			f:          &SimpleBorderFactory{},
			layoutName: "固定布局右下logo模板",
			wantNil:    false,
			wantType:   &fixedBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*fixedBottomLogoTextLayoutBorder)
				return border.IsRight
			},
		},
		{
			name:       "Create auto bottom logo text left no separator layout",
			f:          &SimpleBorderFactory{},
			layoutName: "经典-左logo",
			wantNil:    false,
			wantType:   &autoBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*autoBottomLogoTextLayoutBorder)
				return !border.IsRight && !border.HasSeparator
			},
		},
		{
			name:       "Create auto bottom logo text right no separator layout",
			f:          &SimpleBorderFactory{},
			layoutName: "经典-右logo",
			wantNil:    false,
			wantType:   &autoBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*autoBottomLogoTextLayoutBorder)
				return border.IsRight && !border.HasSeparator
			},
		},
		{
			name:       "Create auto bottom logo text left layout",
			f:          &SimpleBorderFactory{},
			layoutName: "经典-左logo-2",
			wantNil:    false,
			wantType:   &autoBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*autoBottomLogoTextLayoutBorder)
				return !border.IsRight && border.HasSeparator
			},
		},
		{
			name:       "Create auto bottom logo text right layout",
			f:          &SimpleBorderFactory{},
			layoutName: "经典-右logo-1",
			wantNil:    false,
			wantType:   &autoBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*autoBottomLogoTextLayoutBorder)
				return border.IsRight && border.HasSeparator
			},
		},
		{
			name:       "Create auto bottom logo text average layout",
			f:          &SimpleBorderFactory{},
			layoutName: "经典-右logo-对比",
			wantNil:    false,
			wantType:   &autoBottomLogoTextAverageLayoutBorder{},
			checkProps: func(v any) bool {
				_, ok := v.(*autoBottomLogoTextAverageLayoutBorder)
				return ok
			},
		},
		{
			name:       "Create simple bottom text center layout",
			f:          &SimpleBorderFactory{},
			layoutName: "简约-居中-无logo",
			wantNil:    false,
			wantType:   &simpleBottomLogoTextCenterBorder{},
			checkProps: func(v any) bool {
				border := v.(*simpleBottomLogoTextCenterBorder)
				return !border.HasLogo
			},
		},
		{
			name:       "Create simple bottom logo text center layout",
			f:          &SimpleBorderFactory{},
			layoutName: "简约-居中",
			wantNil:    false,
			wantType:   &simpleBottomLogoTextCenterBorder{},
			checkProps: func(v any) bool {
				border := v.(*simpleBottomLogoTextCenterBorder)
				return border.HasLogo
			},
		},
		{
			name:       "Create blur bottom text center layout",
			f:          &SimpleBorderFactory{},
			layoutName: "高斯模糊-居中",
			wantNil:    false,
			wantType:   &blurBottomTextCenterLayout{},
			checkProps: func(v any) bool {
				_, ok := v.(*blurBottomTextCenterLayout)
				return ok
			},
		},
		{
			name:       "Create unknown layout (should return default)",
			f:          &SimpleBorderFactory{},
			layoutName: "unknown_layout",
			wantNil:    false,
			wantType:   &fixedBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*fixedBottomLogoTextLayoutBorder)
				return !border.IsRight
			},
		},
		{
			name:       "Create empty layout (should return default)",
			f:          &SimpleBorderFactory{},
			layoutName: "",
			wantNil:    false,
			wantType:   &fixedBottomLogoTextLayoutBorder{},
			checkProps: func(v any) bool {
				border := v.(*fixedBottomLogoTextLayoutBorder)
				return !border.IsRight
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			got := tt.f.createBorder(tt.layoutName)

			if tt.wantNil && got != nil {
				t.Errorf("SimpleBorderFactory.createBorder() expected nil but got %v", got)
			}

			if !tt.wantNil && got == nil {
				t.Error("SimpleBorderFactory.createBorder() returned nil")
			}

			if !tt.wantNil && got != nil {
				switch tt.wantType.(type) {
				case *fixedBottomLogoTextLayoutBorder:
					border, ok := got.(*fixedBottomLogoTextLayoutBorder)
					if !ok {
						t.Errorf("Expected *fixedBottomLogoTextLayoutBorder, got %T", got)
					} else if tt.checkProps != nil && !tt.checkProps(border) {
						t.Errorf("Properties check failed for *fixedBottomLogoTextLayoutBorder")
					}
				case *autoBottomLogoTextLayoutBorder:
					border, ok := got.(*autoBottomLogoTextLayoutBorder)
					if !ok {
						t.Errorf("Expected *autoBottomLogoTextLayoutBorder, got %T", got)
					} else if tt.checkProps != nil && !tt.checkProps(border) {
						t.Errorf("Properties check failed for *autoBottomLogoTextLayoutBorder")
					}
				case *autoBottomLogoTextAverageLayoutBorder:
					border, ok := got.(*autoBottomLogoTextAverageLayoutBorder)
					if !ok {
						t.Errorf("Expected *autoBottomLogoTextAverageLayoutBorder, got %T", got)
					} else if tt.checkProps != nil && !tt.checkProps(border) {
						t.Errorf("Properties check failed for *autoBottomLogoTextAverageLayoutBorder")
					}
				case *simpleBottomLogoTextCenterBorder:
					border, ok := got.(*simpleBottomLogoTextCenterBorder)
					if !ok {
						t.Errorf("Expected *simpleBottomLogoTextCenterBorder, got %T", got)
					} else if tt.checkProps != nil && !tt.checkProps(border) {
						t.Errorf("Properties check failed for *simpleBottomLogoTextCenterBorder")
					}
				case *blurBottomTextCenterLayout:
					border, ok := got.(*blurBottomTextCenterLayout)
					if !ok {
						t.Errorf("Expected *blurBottomTextCenterLayout, got %T", got)
					} else if tt.checkProps != nil && !tt.checkProps(border) {
						t.Errorf("Properties check failed for *blurBottomTextCenterLayout")
					}
				}
			}
		})
	}
}

func TestFixedBottomLogoTextLayoutBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *fixedBottomLogoTextLayoutBorder
	}{
		{
			name: "Fixed bottom logo text layout border (left)",
			border: &fixedBottomLogoTextLayoutBorder{
				IsRight: false,
			},
		},
		{
			name: "Fixed bottom logo text layout border (right)",
			border: &fixedBottomLogoTextLayoutBorder{
				IsRight: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}

			if tt.border != nil {
				if tt.border.IsRight && tt.name == "Fixed bottom logo text layout border (left)" {
					t.Error("FixedBottomLogoTextLayoutBorder IsRight mismatch")
				}
				if !tt.border.IsRight && tt.name == "Fixed bottom logo text layout border (right)" {
					t.Error("FixedBottomLogoTextLayoutBorder IsRight mismatch")
				}
			}
		})
	}
}

func TestAutoBottomLogoTextLayoutBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Auto bottom logo text layout border (left, no separator)",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Auto bottom logo text layout border (right, no separator)",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
		},
		{
			name: "Auto bottom logo text layout border (left, with separator)",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
		},
		{
			name: "Auto bottom logo text layout border (right, with separator)",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}
		})
	}
}

func TestAutoBottomLogoTextAverageLayoutBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextAverageLayoutBorder
	}{
		{
			name:   "Auto bottom logo text average layout border",
			border: &autoBottomLogoTextAverageLayoutBorder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}
		})
	}
}

func TestSimpleBottomLogoTextCenterBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *simpleBottomLogoTextCenterBorder
	}{
		{
			name: "Simple bottom text center layout (no logo)",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: false,
			},
		},
		{
			name: "Simple bottom logo text center layout (with logo)",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}
		})
	}
}

func TestBlurBottomTextCenterLayout(t *testing.T) {
	tests := []struct {
		name   string
		border *blurBottomTextCenterLayout
	}{
		{
			name:   "Blur bottom text center layout",
			border: &blurBottomTextCenterLayout{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}
		})
	}
}

func TestBaseBottomLogoTextLayoutBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *baseBottomLogoTextLayoutBorder
	}{
		{
			name: "Base bottom logo text layout border (left)",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Base bottom logo text layout border (right)",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
		},
		{
			name: "Base bottom logo text layout border (left with separator)",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
		},
		{
			name: "Base bottom logo text layout border (right with separator)",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.border == nil {
				t.Error("border is nil")
			}

			if tt.border != nil {
				if tt.border.IsRight && tt.name == "Base bottom logo text layout border (left)" {
					t.Error("BaseBottomLogoTextLayoutBorder IsRight mismatch")
				}
				if !tt.border.IsRight && tt.name == "Base bottom logo text layout border (right)" {
					t.Error("BaseBottomLogoTextLayoutBorder IsRight mismatch")
				}
				if tt.border.HasSeparator && tt.name == "Base bottom logo text layout border (left)" {
					t.Error("BaseBottomLogoTextLayoutBorder HasSeparator mismatch")
				}
			}
		})
	}
}

func TestAutoBottomLogoShowInfo(t *testing.T) {
	tests := []struct {
		name string
		info *autoBottomLogoShowInfo
	}{
		{
			name: "Auto bottom logo show info",
			info: &autoBottomLogoShowInfo{
				r1:         map[string]int{"width": 100, "height": 50},
				r2:         map[string]int{"width": 80, "height": 40},
				diffWidth:  20,
				diffHeight: 10,
			},
		},
		{
			name: "Auto bottom logo show info with zero diff",
			info: &autoBottomLogoShowInfo{
				r1:         map[string]int{"width": 100, "height": 50},
				r2:         map[string]int{"width": 100, "height": 50},
				diffWidth:  0,
				diffHeight: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.info == nil {
				t.Error("info is nil")
			}

			if tt.info != nil {
				if tt.info.r1 == nil || tt.info.r2 == nil {
					t.Error("autoBottomLogoShowInfo r1 or r2 is nil")
				}
			}
		})
	}
}
