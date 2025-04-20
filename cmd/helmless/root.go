package helmless

import (
	"github.com/spf13/cobra"
)

// Knative object type manipulable by Helmless
type GCRType int

const (
	GCRService GCRType = iota // Service
	GCRJob                    // Job
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helmless",
		Short: "Helmless is a tool for managing Helm charts",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newMigrateCmd())
	return cmd
}

func Execute() error {
	return NewRootCmd().Execute()
}
