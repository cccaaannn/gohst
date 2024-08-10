package main

import (
	"fmt"

	"github.com/cccaaannn/gohst"
)

func main() {
	server := gohst.CreateServer()

	server.AddHandler("GET /hi", func(req *gohst.Request, res *gohst.Response) {
		res.Body = `
		<body>
			<h1>Hello world!</h1>
		</body>
		`
	})

	server.AddHandler("/*", func(req *gohst.Request, res *gohst.Response) {
		res.StatusCode = 404
		res.Body = `
		<body>
			<h1>404</h1>
		</body>
		`
	})

	stop, err := server.ListenAndServeTLS(":8080", "../../test/cert/localhost.crt", "../../test/cert/localhost.key")
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
	defer close(stop)
	<-stop
}
