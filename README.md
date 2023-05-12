# 用 Go 获取 Github 常用信息

[hello.go](hello.go) 例子:

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/chai2010/ghutil"
)

var (
	flagUserName = flag.String("u", "chai2010", "set user name")
)

func main() {
	flag.Parse()

	c := ghutil.NewGHClient(&ghutil.GHOptions{
		Token: os.Getenv("GITHUB_ACCESS_TOKEN"),
	})
	repos, err := c.GetRepositories(context.Background(), *flagUserName)
	if err != nil {
		panic(err)
	}

	var starCount int
	var forkCount int
	for _, repo := range repos {
		starCount += *repo.StargazersCount
		forkCount += *repo.ForksCount
	}

	fmt.Printf("%s/{*}:StargazersCount: %d\n", *flagUserName, starCount)
	fmt.Printf("%s/{*}:ForksCount:      %d\n", *flagUserName, forkCount)
}
```

设置好 `GITHUB_ACCESS_TOKEN` 环境变量, 然后执行:

```
$ go run hello.go -u=chai2010
chai2010/{*}:StargazersCount: 25283
chai2010/{*}:ForksCount:      3978

$ go run hello.go -u=KusionStack
KusionStack/{*}:StargazersCount: 1206
KusionStack/{*}:ForksCount:      200
```
