// Github 辅助函数
package ghutil

import (
	"context"
	"sort"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Github 客户端
type ghClient struct {
	opt *GHOptions
	c   *github.Client

	cache *ghClientCache
}

func newGHClient(opt *GHOptions) *ghClient {
	if opt == nil {
		opt = &GHOptions{}
	}

	client := github.NewClient(nil)
	if opt.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: opt.Token})
		client = github.NewClient(oauth2.NewClient(context.Background(), ts))
	}

	p := &ghClient{
		opt:   opt,
		c:     client,
		cache: newGHClientCache(),
	}

	if p.opt.CacheFilename != "" {
		p.cache.load(p.opt.CacheFilename)
	}

	return p
}

// 获取用户的仓库列表
func (p *ghClient) GetRepositories(ctx context.Context, userName string) (allRepos []*github.Repository, err error) {
	if repos, ok := p.cache.UserRepos[userName]; ok {
		return repos, nil
	}
	defer func() {
		if err == nil {
			p.cache.UserRepos[userName] = allRepos
			if p.opt.CacheFilename != "" {
				p.cache.save(p.opt.CacheFilename)
			}
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, resp, err := p.c.Repositories.List(ctx, userName, opt)
	if len(repos) > 0 {
		allRepos = append(allRepos, repos...)
	}
	if err != nil {
		return allRepos, err
	}

	if resp.NextPage != 0 && resp.LastPage > 2 {
		for page := 2; page <= resp.LastPage; page++ {
			opt := &github.RepositoryListOptions{
				ListOptions: github.ListOptions{PerPage: 100},
			}
			opt.Page = page

			repos, _, err := p.c.Repositories.List(ctx, userName, opt)
			if len(repos) > 0 {
				allRepos = append(allRepos, repos...)
			}
			if err != nil {
				return allRepos, err
			}
		}
	}

	sort.Slice(allRepos, func(i, j int) bool {
		var si, sj string
		if allRepos[i].FullName != nil {
			si = *allRepos[i].FullName
		}
		if allRepos[j].FullName != nil {
			sj = *allRepos[j].FullName
		}
		return si < sj
	})

	return allRepos, nil
}

// 用户点赞仓库
func (p *ghClient) GetStarredRepos(ctx context.Context, userName string) (allRepos []*github.StarredRepository, err error) {
	if repos, ok := p.cache.UserStarRepos[userName]; ok {
		return repos, nil
	}
	defer func() {
		if err == nil {
			p.cache.UserStarRepos[userName] = allRepos
			if p.opt.CacheFilename != "" {
				p.cache.save(p.opt.CacheFilename)
			}
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	opt := &github.ActivityListStarredOptions{}
	opt.ListOptions.PerPage = 100

	repos, resp, err := p.c.Activity.ListStarred(ctx, userName, opt)
	if len(repos) > 0 {
		allRepos = append(allRepos, repos...)
	}
	if err != nil {
		return allRepos, err
	}

	if resp.NextPage != 0 && resp.LastPage > 2 {
		for page := 2; page <= resp.LastPage; page++ {
			opt.ListOptions.Page = page

			repos, _, err := p.c.Activity.ListStarred(ctx, userName, opt)
			if len(repos) > 0 {
				allRepos = append(allRepos, repos...)
			}
			if err != nil {
				return allRepos, err
			}
		}
	}

	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].StarredAt.Before(allRepos[j].StarredAt.Time)
	})

	return allRepos, nil
}

// 获取仓库的关注信息
func (p *ghClient) GetRepoStargazers(ctx context.Context, userName, repoName string) (allStargazers []*github.Stargazer, err error) {
	if stargazer, ok := p.cache.RepoStargazers[userName+"/"+repoName]; ok {
		return stargazer, nil
	}
	defer func() {
		if err == nil {
			p.cache.RepoStargazers[userName+"/"+repoName] = allStargazers
			if p.opt.CacheFilename != "" {
				p.cache.save(p.opt.CacheFilename)
			}
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	opt := github.ListOptions{PerPage: 100}

	stargazers, resp, err := p.c.Activity.ListStargazers(ctx, userName, repoName, &opt)
	if len(stargazers) > 0 {
		allStargazers = append(allStargazers, stargazers...)
	}
	if err != nil {
		return allStargazers, err
	}

	if resp.NextPage != 0 && resp.LastPage > 2 {
		for page := 2; page <= resp.LastPage; page++ {
			opt.Page = page
			stargazers, _, err := p.c.Activity.ListStargazers(ctx, userName, repoName, &opt)
			if len(stargazers) > 0 {
				allStargazers = append(allStargazers, stargazers...)
			}
			if err != nil {
				return allStargazers, err
			}
		}
	}

	sort.Slice(allStargazers, func(i, j int) bool {
		var si, sj string
		if allStargazers[i].User.Name != nil {
			si = *allStargazers[i].User.Name
		}
		if allStargazers[j].User.Name != nil {
			sj = *allStargazers[j].User.Name
		}
		return si < sj
	})
	return allStargazers, nil
}

// 获取用户信息
func (p *ghClient) GetUserInfo(ctx context.Context, userName string) (user *github.User, err error) {
	if user, ok := p.cache.UserInfos[userName]; ok {
		return user, nil
	}
	defer func() {
		if err == nil {
			p.cache.UserInfos[userName] = user
			if p.opt.CacheFilename != "" {
				p.cache.save(p.opt.CacheFilename)
			}
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	user, _, err = p.c.Users.Get(ctx, userName)
	return user, err
}

// 获取仓库信息
func (p *ghClient) GetRepoInfo(ctx context.Context, userName, repoName string) (repo *github.Repository, err error) {
	if user, ok := p.cache.RepoInfos[userName+"/"+repoName]; ok {
		return user, nil
	}
	defer func() {
		if err == nil {
			p.cache.RepoInfos[userName+"/"+repoName] = repo
			if p.opt.CacheFilename != "" {
				p.cache.save(p.opt.CacheFilename)
			}
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	repo, _, err = p.c.Repositories.Get(ctx, userName, repoName)
	return repo, err
}
