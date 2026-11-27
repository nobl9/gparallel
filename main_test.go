package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

const runAsGparallelEnv = "GPARALLEL_TEST_RUN_AS_GPARALLEL"

// TestMain lets the test binary stand in for the gparallel executable, so that
// tests can observe it as a separate process.
func TestMain(m *testing.M) {
	if os.Getenv(runAsGparallelEnv) != "" {
		main()
	}
	os.Exit(m.Run())
}

// TestManyShortCommandsWithoutATerminal verifies that gparallel does not panic
// and prints every output when many short commands finish while others still
// wait for a free slot. The failure is a race, so one run is not enough.
func TestManyShortCommandsWithoutATerminal(t *testing.T) {
	const runs = 20
	const commands = 300

	args := []string{"-P", "32", "echo", ":::"}
	for i := 1; i <= commands; i++ {
		args = append(args, strconv.Itoa(i))
	}
	want := strings.Join(args[4:], "\n") + "\n"

	for run := 1; run <= runs; run++ {
		var stderr strings.Builder
		cmd := exec.Command(os.Args[0], args...)
		cmd.Env = append(os.Environ(), runAsGparallelEnv+"=1")
		cmd.Stderr = &stderr

		got, err := cmd.Output()

		if err != nil {
			t.Fatalf("run %d: %v\n%s", run, err, stderr.String())
		}
		if string(got) != want {
			t.Fatalf("run %d: printed %d of %d lines", run, strings.Count(string(got), "\n"), commands)
		}
	}
}
