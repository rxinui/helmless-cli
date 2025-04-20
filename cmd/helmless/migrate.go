package helmless

import (
	"fmt"
	"os"

	"regexp"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const MIGRATE_SERVICE_REGEXP = `service(s)?/`
const MIGRATE_JOB_REGEXP = `job(s)?/`
const MIGRATE_ARG_REGEXP = `^(` + MIGRATE_SERVICE_REGEXP + `|` + MIGRATE_JOB_REGEXP + `)[\w-]+`

func buildMigrateCmdExample() string {
	return `
  # Migrate using the interactive CLI
  migrate

  # Migrate a specific service named my-service
  migrate services/my-service

  # Migrate a specific job name my-job
  migrate jobs/my-job

  # Migrate several service and/or jobs 
  migrate services/svc-1 jobs/job-1 services/svc-2 ...
`
}

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:        "migrate [...{services, jobs}/NAME] ",
		Aliases:    []string{},
		SuggestFor: []string{},
		Short:      "Migrate a Knative Service or Job into Helmless template",
		GroupID:    "",
		Example:    buildMigrateCmdExample(),
		ValidArgs:  []string{"services/", "jobs/"},
		Args:       cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				cmd.Help()
				os.Exit(1)
			}
			fmt.Println("migrate command is working")
			opts := newMigrateOptions(args[0])
			log.Info().Msgf("result: %v", opts)
		},
		SuggestionsMinimumDistance: 0,
	}
}

type migrateOptions struct {
	objectName string // TODO convert into []string to handle multiple objects
	fs         afero.Fs
	objectType GCRType
}

func newMigrateOptions(objectName string) migrateOptions {
	var subs []string
	reService := regexp.MustCompile(MIGRATE_SERVICE_REGEXP)
	reJob := regexp.MustCompile(MIGRATE_JOB_REGEXP)
	re := regexp.MustCompile(MIGRATE_ARG_REGEXP)
	if subs = re.FindStringSubmatch(objectName); len(subs) == 0 {
		// TODO handle error gracefully
		panic(fmt.Errorf("invalid object name '%s'", objectName))
	}
	switch {
	case reService.MatchString(subs[0]):
		log.Debug().Msgf("service matched: %s", subs[0])
		return migrateOptions{objectName: objectName, objectType: GCRService, fs: afero.NewOsFs()}
	case reJob.MatchString(subs[0]):
		log.Debug().Msgf("job matched: %s", subs[0])
		return migrateOptions{objectName: objectName, objectType: GCRJob, fs: afero.NewOsFs()}
	default:
		panic("invalid migrate")
	}
}
