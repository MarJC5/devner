package cli

import (
	"context"
	"fmt"

	"github.com/devner/devner/internal/tui"
)

func runTUI() error {
	d, err := buildDeps()
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}
	defer d.Close()
	return tui.Run(context.Background(), d)
}
