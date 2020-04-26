package git

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
)

// MirrorDirectory is the name of the directory inside
// the cache directory to which the git repo is cloned.
const MirrorDirectory = "mirror"

var mirror *Repository

func getMirror() *Repository {
	if mirror == nil {
		panic(errors.InternalService("git mirror not initialised"))
	}

	return mirror
}

// Init clones the repo if necessary and checks out the given revision
func Init(revision string) error {
	// Don't initialise twice
	if mirror != nil {
		return nil
	}

	mirror = &Repository{
		Dir:    filepath.Join(utils.CacheDir(), MirrorDirectory),
		Remote: config.Repository(),
	}

	op := output.Info("Cloning repo to %s", mirror.Dir)
	if err := mirror.clone(); err != nil {
		op.Failed()
		return errors.WithMessage(err, "failed to clone repo")
	}
	op.Complete()

	op = output.Info("Checking out %s", revision)
	if err := mirror.checkout(revision); err != nil {
		op.Failed()
		return errors.WithMessage(err, "failed to checkout %q", revision)
	}
	op.Complete()

	return nil
}

// ShortHash returns the short hash of the current git commit
func ShortHash() (string, error) {
	return getMirror().shortHash()
}

// Repository represents a git repository
type Repository struct {
	Dir    string
	Remote string
}

func (r *Repository) shortHash() (string, error) {
	result := r.revParse("--short", "HEAD")

	if result.Err != nil {
		return "", errors.WithMessage(result.Err, "failed to get short hash")
	}

	return result.Stdout, nil
}

func (r *Repository) clone() error {
	slog.NewStdoutLogger()
	if !path.IsAbs(r.Dir) {
		return errors.InternalService("specified repository directory is not absolute")
	}

	// Use MkdirAll instead of Mkdir because it
	// doesn't error if the directory already exists
	if err := os.MkdirAll(r.Dir, os.ModePerm); err != nil {
		return errors.WithMessage(err, "failed to create mirror directory")
	}

	// Return early if the repository already exists
	if exists, err := r.exists(); err != nil {
		return errors.WithMessage(err, "failed to check if repository exists")
	} else if exists {
		return nil
	}

	if err := r.exec("clone", r.Remote, r.Dir); err != nil {
		return errors.WithMessage(err, "failed to clone repository")
	}

	return nil
}

func (r *Repository) checkout(revision string) error {
	if err := r.exec("fetch", "--all", "--prune", "--prune-tags", "--force").Err; err != nil {
		return errors.WithMessage(err, "failed to fetch from remote", revision)
	}

	if err := r.exec("checkout", revision).Err; err != nil {
		return errors.WithMessage(err, "failed to checkout %q", revision)
	}

	if onBranch, err := r.isOnBranch(); err != nil {
		return errors.WithMessage(err, "failed to check if on branch")
	} else if !onBranch {
		return nil
	}

	if err := r.exec("reset", "--hard", "@{u}").Err; err != nil {
		return errors.WithMessage(err, "failed to reset branch")
	}

	return nil
}

func (r *Repository) exists() (bool, error) {
	result := r.revParse("--is-inside-work-tree")

	if strings.HasPrefix(result.Stderr, "fatal: not a git repository") {
		return false, nil
	}

	if result.Err != nil {
		return false, result.Err
	}

	if result.Stdout == "true" {
		return true, nil
	}

	return false, errors.InternalService("unexpected output: %v", result.Stdout)
}

func (r *Repository) revParse(args ...string) *exe.Result {
	args = append([]string{"rev-parse"}, args...)
	return r.exec(args...)
}

func (r *Repository) isOnBranch() (bool, error) {
	result := r.exec("branch", "--show-current")

	if result.Err != nil {
		return false, result.Err
	}

	return result.Stdout != "", nil
}

func (r *Repository) exec(args ...string) *exe.Result {
	return exe.Command("git", args...).Dir(r.Dir).Run()
}
