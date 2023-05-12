package ghutil

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/go-github/github"
)

// 缓存信息
type ghClientCache struct {
	UserInfos      map[string]*github.User                // 用户信息
	RepoInfos      map[string]*github.Repository          // 仓库信息
	UserRepos      map[string][]*github.Repository        // 每个用户或组织下的仓库信息
	RepoStargazers map[string][]*github.Stargazer         // 仓库的 Star 信息
	UserStarRepos  map[string][]*github.StarredRepository // 用户点赞仓库
}

func newGHClientCache() *ghClientCache {
	return &ghClientCache{
		UserInfos:      make(map[string]*github.User),
		RepoInfos:      make(map[string]*github.Repository),
		UserRepos:      make(map[string][]*github.Repository),
		RepoStargazers: make(map[string][]*github.Stargazer),
		UserStarRepos:  make(map[string][]*github.StarredRepository),
	}
}

func (p *ghClientCache) isEmpty() bool {
	if len(p.UserInfos) > 0 {
		return false
	}
	if len(p.UserRepos) > 0 {
		return false
	}
	if len(p.UserRepos) > 0 {
		return false
	}
	if len(p.RepoStargazers) > 0 {
		return false
	}
	return true
}

func (p *ghClientCache) load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var tmp ghClientCache
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	if len(tmp.UserInfos) > 0 {
		p.UserInfos = tmp.UserInfos
	}
	if len(tmp.RepoInfos) > 0 {
		p.RepoInfos = tmp.RepoInfos
	}
	if len(tmp.UserRepos) > 0 {
		p.UserRepos = tmp.UserRepos
	}
	if len(tmp.UserStarRepos) > 0 {
		p.UserStarRepos = tmp.UserStarRepos
	}
	if len(tmp.RepoStargazers) > 0 {
		p.RepoStargazers = tmp.RepoStargazers
	}

	return nil
}

func (p *ghClientCache) save(filename string) error {
	d, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}

	newName := filename + ".tmp"
	bakName := filename + ".bak." + time.Now().Format("20060102")
	err = os.WriteFile(newName, d, 0666)
	if err != nil {
		return err
	}

	os.Rename(filename, bakName)
	os.Rename(newName, filename)
	return nil
}
