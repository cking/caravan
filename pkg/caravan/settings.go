package caravan

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"net/url"
	"os"
	"time"
)

type Settings struct {
	GitRepository string
	GitBranch     string
	GitPath       string

	Interval time.Duration
}

func NewSettings() *Settings {
	return &Settings{
		Interval:  time.Minute,
		GitPath:   "/",
		GitBranch: "main",
	}
}

func SettingsFromEnv() *Settings {
	settings := NewSettings()

	if key := os.Getenv("GIT_REPO"); key != "" {
		settings.GitRepository = key
	}
	if key := os.Getenv("GIT_BRANCH"); key != "" {
		settings.GitBranch = key
	}
	if key := os.Getenv("GIT_PATH"); key != "" {
		settings.GitPath = key
	}
	if key := os.Getenv("CARAVAN_INTERVAL"); key != "" {
		if dur, err := time.ParseDuration(key); err == nil {
			settings.Interval = dur
		}
	}

	return settings
}

func (s *Settings) Verify() (err error) {
	if s.Interval < 0 {
		err := multierr.Append(err, errors.New("negative time provided for interval"))
		if err != nil {
			return err
		}
	}

	if _, err := url.Parse(s.GitRepository); err != nil {
		err := multierr.Append(err, fmt.Errorf("failed to parse git url: %w", err))
		if err != nil {
			return err
		}
	}

	if len(s.GitBranch) < 0 {
		err := multierr.Append(err, errors.New("branch name empty"))
		if err != nil {
			return err
		}
	}

	return
}
