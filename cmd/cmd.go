package cmd

import (
	"caravan/pkg/caravan"
	"caravan/pkg/genric"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	nomad "github.com/hashicorp/nomad-openapi/clients/go/v1"
	"github.com/joho/godotenv"
	"io"
	"os"
	"time"
)

// Execute executes the root command.
func Execute() error {
	log := slog.Make(sloghuman.Sink(os.Stdout))
	_ = godotenv.Load()

	err := godotenv.Load()
	if err != nil {
		log.Warn(context.Background(), "Error loading .env file", slog.Error(err))
	}

	settings, err := ParseSettings()
	if err != nil {
		log.Fatal(context.Background(), "Error parsing env", slog.Error(err))
		return err
	}

	// Create Nomad client
	client, err := caravan.NewClient()
	if err != nil {
		log.Fatal(context.Background(), "Failed to build nomad client", slog.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	worktree, err := caravan.CloneGit(ctx, settings.GitRepository, settings.GitBranch)
	defer cancel()

	// Reconcile Loop
	for {
		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			log.Error(context.Background(), "Failed to pull latest", slog.Error(err))
			return err
		}

		fs := worktree.Filesystem
		path := settings.GitPath
		files, err := fs.ReadDir(path)
		if err != nil {
			log.Error(context.Background(), "Failed to read files in git repo", slog.Error(err))
			return err
		}

		jobs, err := client.ListJobs()
		if err != nil {
			log.Error(context.Background(), "Failed to fetch nomad jobs", slog.Error(err))
			return err
		}

		// get caravan managed jobs
		jobs = genric.FilterMap(jobs, func(k string, j *nomad.Job) bool {
			if j.Meta == nil {
				return false
			}

			_, found := (*j.Meta)[caravan.MetaKeyName]
			return found
		})

		// Parse and apply all jobs from within the git repo
		for _, file := range files {
			filePath := fs.Join(path, file.Name())
			f, err := fs.Open(filePath)
			if err != nil {
				log.Error(context.Background(), "Failed to open job file", slog.F("file", filePath), slog.Error(err))
				return fmt.Errorf("failed to open %v: %w", filePath, err)
			}

			b, err := io.ReadAll(f)
			if err != nil {
				log.Error(context.Background(), "Failed to read job file", slog.F("file", filePath), slog.Error(err))
				return fmt.Errorf("failed to read %v: %w", filePath, err)
			}

			job, err := client.ParseJob(string(b))
			if err != nil {
				log.Warn(context.Background(), "Failed to parse job file, skipping", slog.F("file", filePath), slog.Error(err))
				continue
			}
			delete(jobs, job.GetName())

			fmt.Printf("Applying job [%s]\n", job.GetName())
			_, err = client.ApplyJob(job)
			if err != nil {
				log.Warn(context.Background(), "Failed to apply job", slog.F("file", filePath), slog.Error(err))
				return err
			}
		}

		for name, job := range jobs {
			meta := job.GetMeta()

			if _, isManaged := meta[caravan.MetaKeyName]; isManaged {
				err = client.DeleteJob(job)
				if err != nil {
					log.Warn(context.Background(), "Failed to parse job file, continueing", slog.F("job", name), slog.Error(err))
				}
			}
		}

		time.Sleep(settings.Interval)
	}
}

func ParseSettings() (*caravan.Settings, error) {
	var err error

	settings := caravan.SettingsFromEnv()

	err = settings.Verify()
	return settings, err
}
