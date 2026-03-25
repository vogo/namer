package version_test

import (
	"strings"
	"testing"

	"github.com/vogo/namer/integrations/helper"
)

func TestVersionFlag(t *testing.T) {
	stdout, _, err := helper.Run(t, "-v")
	if err != nil {
		t.Fatalf("namer -v failed: %v", err)
	}

	if !strings.Contains(stdout, "namer v") {
		t.Errorf("expected version output, got: %s", stdout)
	}
}

func TestVersionFlagLong(t *testing.T) {
	stdout, _, err := helper.Run(t, "--version")
	if err != nil {
		t.Fatalf("namer --version failed: %v", err)
	}

	if !strings.Contains(stdout, "namer v") {
		t.Errorf("expected version output, got: %s", stdout)
	}
}

func TestHelpFlag(t *testing.T) {
	stdout, _, err := helper.Run(t, "-h")
	if err != nil {
		t.Fatalf("namer -h failed: %v", err)
	}

	// Should contain usage info
	for _, keyword := range []string{"namer", "用法", "-xing", "-ming", "-web", "评分维度"} {
		if !strings.Contains(stdout, keyword) {
			t.Errorf("help output missing keyword %q", keyword)
		}
	}
}

func TestHelpFlagLong(t *testing.T) {
	stdout, _, err := helper.Run(t, "--help")
	if err != nil {
		t.Fatalf("namer --help failed: %v", err)
	}

	if !strings.Contains(stdout, "用法") {
		t.Errorf("expected help output, got: %s", stdout)
	}
}

func TestHelpCommand(t *testing.T) {
	stdout, _, err := helper.Run(t, "help")
	if err != nil {
		t.Fatalf("namer help failed: %v", err)
	}

	if !strings.Contains(stdout, "用法") {
		t.Errorf("expected help output, got: %s", stdout)
	}
}
