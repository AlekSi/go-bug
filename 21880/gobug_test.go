package gobug

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sleep", "5")

	err := cmd.Run()
	expected := "signal: killed"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err)
	}
}

func TestCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	cancel()

	err := cmd.Wait()
	expected := "signal: killed"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err)
	}
}
