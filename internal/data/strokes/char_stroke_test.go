package strokes

import "testing"

func TestCharStroke(t *testing.T) {
	tests := []struct {
		c    rune
		want int
	}{
		{'一', 1},
		{'乙', 1},
		{'人', 2},
		{'大', 3},
		{'王', 4},
		{'生', 5},
		{'安', 6},
		{'李', 7},
		{'明', 8},
		{'春', 9},
		{'高', 10},
		{'康', 11},
		{'博', 12},
		{'龘', 0}, // 不在数据库
	}
	for _, tt := range tests {
		if got := CharStroke(tt.c); got != tt.want {
			t.Errorf("CharStroke(%c) = %d, want %d", tt.c, got, tt.want)
		}
	}
}
