package git

import (
	gogit "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	stacktrace "github.com/palantir/stacktrace"
)

type backend struct {
	repository *gogit.Repository
}

// Open returns a repo struct that exposes an api for getting information
// about project's git repository
func Open() (*backend, error) {
	r, err := gogit.PlainOpen(".")
	if err != nil {
		err = stacktrace.Propagate(err, "could not instantiate git repo")
		return nil, err
	}
	if r == nil {
		err = stacktrace.NewError("go-git returned repository as a nil struct")
		return nil, err
	}
	result := &backend{
		repository: r,
	}
	return result, nil
}

// Commit returns current branch name
func (r *backend) Branch() (string, error) {
	branchRefs, err := r.repository.Branches()
	if err != nil {
		return "", err
	}
	headRef, err := r.repository.Head()
	if err != nil {
		return "", err
	}
	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash().String() == headRef.Hash().String() {
			currentBranchName = branchRef.Name().Short()
			return nil
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return currentBranchName, nil
}

// Commit returns current commit hash
func (r *backend) Commit() (string, error) {
	headRef, err := r.repository.Head()
	if err != nil {
		return "", err
	}
	hash := headRef.Hash().String()
	return hash, nil
}

// Tag returns tag that matches the current commit hash ( if any )
func (r *backend) Tag() (string, error) {
	headHash, err := r.Commit()
	if err != nil {
		return "", err
	}
	tagRefs, err := r.repository.Tags()
	if err != nil {
		return "", err
	}
	var tag string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := r.repository.ResolveRevision(revision)
		if err != nil {
			return err
		}
		if tagCommitHash.String() == headHash {
			tag = tagRef.Name().String()
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return tag, nil
}
