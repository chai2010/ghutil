package ghutil

import (
	"context"
	"sort"

	"github.com/google/go-github/github"
)

// 链接选项
type GHOptions struct {
	Token         string // GITHUB_ACCESS_TOKEN
	CacheFilename string // 缓存文件
}

// Github 客户端
type GHClient struct {
	c *ghClient
}

func NewGHClient(opt *GHOptions) *GHClient {
	return &GHClient{c: newGHClient(opt)}
}

// 获取用户信息
func (p *GHClient) GetUserInfo(ctx context.Context, userName string) (user *github.User, err error) {
	return p.c.GetUserInfo(ctx, userName)
}

// 获取仓库信息
func (p *GHClient) GetRepoInfo(ctx context.Context, userName, repoName string) (repo *github.Repository, err error) {
	return p.c.GetRepoInfo(ctx, userName, repoName)
}

// 获取用户的仓库列表
func (p *GHClient) GetRepositories(ctx context.Context, userName string) (allRepos []*github.Repository, err error) {
	return p.c.GetRepositories(ctx, userName)
}

// 用户点赞仓库
func (p *GHClient) GetStarredRepos(ctx context.Context, userName string) (allRepos []*github.StarredRepository, err error) {
	return p.c.GetStarredRepos(ctx, userName)
}

// 获取仓库的关注信息
func (p *GHClient) GetRepoStargazers(ctx context.Context, userName, repoName string) (allStargazers []*github.Stargazer, err error) {
	return p.c.GetRepoStargazers(ctx, userName, repoName)
}

// 获取仓库编程语言类型
func (p *GHClient) GetRepoLanguages(repos []*github.Repository, max int) []string {
	if true {
		return []string{"TODO"}
	}
	var langMap = make(map[string]int)
	for _, repo := range repos {
		langMap[repo.GetLanguage()]++
	}

	var langInfos = []struct {
		Name  string
		Count int
	}{}

	for k, v := range langMap {
		langInfos = append(langInfos, struct {
			Name  string
			Count int
		}{
			Name:  k,
			Count: v,
		})
	}

	// 逆序排列
	sort.Slice(langInfos, func(i, j int) bool {
		return langInfos[i].Count > langInfos[j].Count
	})

	var langs []string
	for _, x := range langInfos {
		langs = append(langs, x.Name)
	}

	// 逆序排列
	sort.Slice(langs, func(i, j int) bool {
		return langs[i] > langs[j]
	})

	if max > 0 && len(langs) > max {
		langs = langs[:max]
	}

	return langs
}
