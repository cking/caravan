package caravan

import (
	"context"
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// CloneGit will clone a git repo
func CloneGit(ctx context.Context, repositoryURL string, branch string) (*git.Worktree, error) {
	storer := memory.NewStorage()
	fs := memfs.New()

	branchRef := plumbing.NewBranchReferenceName(branch)
	repo, err := git.CloneContext(ctx, storer, fs, &git.CloneOptions{
		URL:           repositoryURL,
		SingleBranch:  true,
		ReferenceName: branchRef,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to clone repository %s for templating: %w", repositoryURL, err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("unable to get worktree: %w", err)
	}

	return w, nil
}
