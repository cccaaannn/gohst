# gohst

### Go HTTP Server Tool
A simple http server

[![Go Reference](https://pkg.go.dev/badge/github.com/cccaaannn/gohst.svg)](https://pkg.go.dev/github.com/cccaaannn/gohst) [![codecov](https://codecov.io/github/cccaaannn/gohst/graph/badge.svg?token=CM770U3IB4)](https://codecov.io/github/cccaaannn/gohst) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cccaaannn/gohst) ![GitHub top language](https://img.shields.io/github/languages/top/cccaaannn/gohst?color=008B8B&style=flat-square) ![GitHub repo size](https://img.shields.io/github/repo-size/cccaaannn/gohst?color=FF6347&style=flat-square) [![GitHub](https://img.shields.io/github/license/cccaaannn/gohst?color=green&style=flat-square)](https://github.com/cccaaannn/gohst/blob/master/LICENSE)

---

## Usage

### Install package
```shell
go get github.com/cccaaannn/gohst
```

### Minimal example
```go
package main

import (
	"fmt"

	"github.com/cccaaannn/gohst"
)

func main() {
	server := gohst.CreateServer()

	server.AddHandler("GET /hi/:name", func(req *gohst.Request, res *gohst.Response) {
		name := req.Params["name"]
		res.Body = fmt.Sprintf(`
		<body>
			<h1>Hello %s!</h1>
		</body>
		`, name)
	})

	server.AddHandler("/*", func(req *gohst.Request, res *gohst.Response) {
		res.StatusCode = 404
		res.Body = `
		<body>
			<h1>404</h1>
		</body>
		`
	})

	stop, _ := server.ListenAndServe(":8080")
	<-stop
}

```



## Development

### Test
```shell
go test -v
```
### Coverage
```shell
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```
