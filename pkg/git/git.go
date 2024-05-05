package git

import (
	"github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Repository struct {
	token string
	fs    billy.Filesystem
	r     *gogit.Repository
	w     *gogit.Worktree
}

func New(token, url string) (*Repository, error) {
	fs := memfs.New()

	r, err := gogit.Clone(memory.NewStorage(), fs, &gogit.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "0815", // can be anything except an empty string
			Password: token,
		},
	})
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	return &Repository{
		token: token,
		fs:    fs,
		r:     r,
		w:     w,
	}, nil
}

func (r *Repository) AddFile(filename string, content []byte) error {
	file, err := r.fs.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	_, err = r.w.Add(filename)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Commit(msg string, author *object.Signature) error {
	_, err := r.w.Commit(msg, &gogit.CommitOptions{
		Author: author,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Push() error {
	return r.r.Push(&gogit.PushOptions{
		Auth: &http.BasicAuth{
			Username: "abc123", // can be anything except an empty string
			Password: r.token,
		},
	})
}

// type Client struct {
// 	token string
// }

// func New(token string) *Client {
// 	return &Client{
// 		token: token,
// 	}
// }

// type PushFileInput struct {
// 	URL      string
// 	Filename string
// 	Content  []byte
// 	Message  string
// 	Author   *object.Signature
// }

// func (g *Client) PushFile(input *PushFileInput) error {
// 	fs := memfs.New()

// 	r, err := gogit.Clone(memory.NewStorage(), fs, &gogit.CloneOptions{
// 		URL: input.URL,
// 		Auth: &http.BasicAuth{
// 			Username: "0815", // can be anything except an empty string
// 			Password: g.token,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	file, err := fs.Create(input.Filename)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	_, err = file.Write(input.Content)
// 	if err != nil {
// 		return err
// 	}

// 	w, err := r.Worktree()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Add(input.Filename)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Commit(input.Message, &gogit.CommitOptions{
// 		Author: input.Author,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	err = r.Push(&gogit.PushOptions{
// 		Auth: &http.BasicAuth{
// 			Username: "abc123", // can be anything except an empty string
// 			Password: g.token,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
