package caravan

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"net/url"
	"time"
)

type Settings struct {
	GitRepository string `emp:"GIT_REPO"`
	GitBranch     string `emp:"GIT_BRANCH"`
	GitPath       string `emp:"GIT_PATH"`

	Interval time.Duration `emp:"CARAVAN_INTERVAL"`
}

func NewSettings() *Settings {
	return &Settings{
		Interval:  time.Minute,
		GitPath:   "/",
		GitBranch: "main",
	}
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
