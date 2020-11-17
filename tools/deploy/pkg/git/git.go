package git

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/libraries/cache"
)

// MirrorDirectory is the name of the directory inside
// the cache directory to which the git repo is cloned.
const MirrorDirectory = "mirror"

var mirror *Repository

func getMirror() *Repository {
	if mirror == nil {
		panic(oops.InternalService("git mirror not initialised"))
	}

	return mirror
}

// Commit represents a commit as told by git log
type Commit struct {
	ShortHash string
	TitleLine string
}

// Init clones the repo if necessary and checks out the given revision
func Init(revision string) error {
	// Don't initialise twice
	if mirror != nil {
		return nil
	}

	mirror = &Repository{
		Dir:    filepath.Join(cache.Dir(), MirrorDirectory),
		Remote: config.Get().Repository,
	}

	op := output.Info("Cloning repo to %s", mirror.Dir)
	if err := mirror.clone(); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to clone repo")
	}
	op.Success()

	op = output.Info("Checking out %s", revision)
	if err := mirror.checkout(revision); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to checkout %q", revision)
	}
	op.Success()

	return nil
}

// Dir returns the full path of the git repo
func Dir() string {
	return getMirror().Dir
}

// ShortHash returns the short hash of the given commit
func ShortHash(commit string) (string, error) {
	return getMirror().hash(commit, true)
}

// CurrentHash returns the hash of the current git commit
func CurrentHash() (long string, short string, err error) {
	long, err = getMirror().hash("HEAD", false)
	if err != nil {
		return
	}

	short, err = getMirror().hash("HEAD", true)
	if err != nil {
		return
	}

	return
}

// Log returns commits for the path between the range given
func Log(from, to, path string) ([]*Commit, error) {
	return getMirror().log(from, to, path)
}

// Repository represents a git repository
type Repository struct {
	Dir    string
	Remote string
}

func (r *Repository) hash(revision string, short bool) (string, error) {
	args := []string{revision}
	if short {
		args = append([]string{"--short"}, args...)
	}

	result := r.revParse(args...)

	if result.Err != nil {
		return "", oops.WithMessage(result.Err, "failed to get hash")
	}

	return result.Stdout, nil
}

func (r *Repository) clone() error {
	slog.NewStdoutLogger()
	if !path.IsAbs(r.Dir) {
		return oops.InternalService("specified repository directory is not absolute")
	}

	// Use MkdirAll instead of Mkdir because it
	// doesn't error if the directory already exists
	if err := os.MkdirAll(r.Dir, os.ModePerm); err != nil {
		return oops.WithMessage(err, "failed to create mirror directory")
	}

	// Return early if the repository already exists
	if exists, err := r.exists(); err != nil {
		return oops.WithMessage(err, "failed to check if repository exists")
	} else if exists {
		return nil
	}

	if err := r.exec("clone", "--recurse-submodules", r.Remote, r.Dir).Err; err != nil {
		return oops.WithMessage(err, "failed to clone repository")
	}

	return nil
}

func (r *Repository) checkout(revision string) error {
	if err := r.exec("fetch", "--all", "--prune", "--prune-tags", "--recurse-submodules", "--force").Err; err != nil {
		return oops.WithMessage(err, "failed to fetch revision %s from remote", revision)
	}

	if err := r.exec("checkout", "--recurse-submodules", "--force", revision).Err; err != nil {
		return oops.WithMessage(err, "failed to checkout %q", revision)
	}

	if onBranch, err := r.isOnBranch(); err != nil {
		return oops.WithMessage(err, "failed to check if on branch")
	} else if !onBranch {
		return nil
	}

	if err := r.exec("reset", "--hard", "--recurse-submodules", "@{u}").Err; err != nil {
		return oops.WithMessage(err, "failed to reset branch")
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

	return false, oops.InternalService("unexpected output: %v", result.Stdout)
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

func (r *Repository) log(from, to, path string) ([]*Commit, error) {
	result := r.exec("log", "--oneline", from+"..."+to, path)

	if result.Err != nil {
		return nil, result.Err
	}

	var commits []*Commit
	lines := strings.Split(result.Stdout, "\n")

	for _, str := range lines {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		parts := strings.SplitN(str, " ", 2)
		if len(parts) != 2 {
			return nil, oops.InternalService("unexpected log line: %s", str)
		}

		commits = append(commits, &Commit{
			ShortHash: parts[0],
			TitleLine: parts[1],
		})
	}

	return commits, nil
}

func (r *Repository) exec(args ...string) *exe.Result {
	return exe.Command("git", args...).Dir(r.Dir).Run()
}
