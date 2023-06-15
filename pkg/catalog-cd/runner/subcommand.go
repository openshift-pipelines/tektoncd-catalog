package runner

import (
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	"github.com/spf13/cobra"
)

// SubCommand wraps a cobra command with "business" logic, providing a sequence of actions to perform
// the intended workflow.
type SubCommand interface {
	// Cmd exposes the SubCommand's cobra command instance.
	Cmd() *cobra.Command

	// Complete should load the required information, arguments, kubernetes resources, etc.
	Complete(_ *config.Config, _ []string) error

	// Validate should validate before "run".
	Validate() error

	// Run performs the primary business logic.
	Run(_ *config.Config) error
}
