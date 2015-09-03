package github

import (
	"fmt"
	"path/filepath"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"

	"github.com/go-gonzo/archive/tar"
	"github.com/go-gonzo/compress/gzip"
	"github.com/go-gonzo/util"
	"github.com/go-gonzo/web"
)

type Release struct {
	User  string
	Repo  string
	Tag   string
	Pluck []string
}

//Get returns a channel of gonzo.Files that match the provided patterns.
func Get(ctx context.Context, release Release) gonzo.Pipe {

	return web.Get(
		context.WithValue(ctx, "repo", fmt.Sprintf("%s/%s#%s", release.User, release.Repo, release.Tag)),
		fmt.Sprintf(
			"https://codeload.github.com/%s/%s/tar.gz/%s",
			release.User,
			release.Repo,
			release.Tag,
		),
	).Pipe(
		gzip.Uncompress(),
		tar.Untar(tar.Options{
			StripComponenets: 1,
			Pluck:            release.Pluck,
		}),
		util.Rename(func(old string) string {
			return filepath.Join(release.Repo, old)
		}),
	)
}
