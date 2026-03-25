package helper

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

var (
	buildOnce sync.Once
	binPath   string
	buildErr  error
)

// NamerBinary compiles the namer binary to a temp directory (once) and returns its path.
func NamerBinary(t *testing.T) string {
	t.Helper()

	buildOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "namer-integration-*")
		if err != nil {
			buildErr = err
			return
		}

		bin := filepath.Join(tmpDir, "namer")
		if runtime.GOOS == "windows" {
			bin += ".exe"
		}

		cmd := exec.Command("go", "build", "-o", bin, ".")
		cmd.Dir = projectRoot()
		out, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = &BuildError{Output: string(out), Err: err}
			return
		}

		binPath = bin
	})

	if buildErr != nil {
		t.Fatalf("failed to build namer binary: %v", buildErr)
	}

	return binPath
}

// BuildError wraps a build failure with compiler output.
type BuildError struct {
	Output string
	Err    error
}

func (e *BuildError) Error() string {
	return e.Err.Error() + "\n" + e.Output
}

// Run executes the namer binary with the given args and returns stdout, stderr, and error.
func Run(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	bin := NamerBinary(t)
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir()) // isolate from user config

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// RunWithStdin executes the namer binary with the given args and stdin input.
func RunWithStdin(t *testing.T, stdin string, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	bin := NamerBinary(t)
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	cmd.Stdin = bytes.NewBufferString(stdin)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// RunWithEnv executes the namer binary with extra environment variables.
func RunWithEnv(t *testing.T, env []string, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	bin := NamerBinary(t)
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// WriteConfigFile creates a config file with the given JSON content and returns the path.
func WriteConfigFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "namer.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return path
}

func projectRoot() string {
	// helper.go is at integrations/helper/helper.go
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..")
}
