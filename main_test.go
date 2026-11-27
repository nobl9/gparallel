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
	assertEveryRunPrintsEveryCommand(t, 20, 32, 300)
}

// TestManyShortCommandsOneAtATime verifies that gparallel does not exit early
// when it may run only one command at once, so that every other command waits
// for the one before it. The failure is a race, so one run is not enough.
func TestManyShortCommandsOneAtATime(t *testing.T) {
	assertEveryRunPrintsEveryCommand(t, 20, 1, 100)
}

// assertEveryRunPrintsEveryCommand runs gparallel repeatedly without a terminal
// and fails the test on the first run that fails or loses output.
func assertEveryRunPrintsEveryCommand(t *testing.T, runs, maxConcurrent, commands int) {
	t.Helper()

	args := []string{"-P", strconv.Itoa(maxConcurrent), "echo", ":::"}
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
