package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func gitClone(repo string) (*git.Repository, error) {
	fs := memfs.New()

	rep, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})

	if err != nil {
		return nil, err
	}

	return rep, nil
}

func gitLog(file string, rep *git.Repository) (object.CommitIter, error) {
	cIter, err := rep.Log(&git.LogOptions{
		FileName: &file,
	})

	if err != nil {
		return nil, err
	}

	return cIter, nil
}

func GetCommit(repoLink string, fileHash string, file string) error {
	var (
		r       *git.Repository
		err     error
		cIter   object.CommitIter
		w       *git.Worktree
		f       *object.File
		content string
		found   bool
	)

	// Clones  the repoLink
	if r, err = gitClone(repoLink); err != nil {
		return err
	}

	// Gets the log of that repo
	if cIter, err = gitLog(file, r); err != nil {
		return err
	}

	found = false
	// For every commit , grab the commit hash from the given file
	// then read the file content in that commit branch and output
	// the commit location
	if err = cIter.ForEach(func(c *object.Commit) error {
		if w, err = r.Worktree(); err != nil {
			return err
		}

		// Checkout to the current commit (c.Hash)
		if err = w.Checkout(&git.CheckoutOptions{Hash: c.Hash}); err != nil {
			return err
		}

		// Get the target file in the specific commit
		if f, err = c.File(file); err != nil {
			return err
		}

		// Read the file content
		if content, err = f.Contents(); err != nil {
			return err
		}

		// Convert the file content to md5sum
		hash := md5.Sum([]byte(content))
		hashToString := hex.EncodeToString(hash[:])

		// Verify if the lookup file corresponds to the target file
		if hashToString == fileHash {
			commitLink := fmt.Sprintf("%s/commit/%s", strings.Split(repoLink, ".git")[0], c.Hash)
			fmt.Printf("\n\nCommit link: %s\n\n", commitLink)
			fmt.Printf(c.String())
			found = true
		}
		return nil
	}); err != nil {
		return err
	}

	if !found {
		log.Printf("Hash %s or file %s not found\n", fileHash, file)
	}

	return nil
}
