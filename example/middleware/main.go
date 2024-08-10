package main

import (
	"fmt"

	"github.com/cccaaannn/gohst"
)

func main() {
	server := gohst.CreateServer()

	// Authentication middleware
	server.Use(func(next gohst.HandlerFunc) gohst.HandlerFunc {
		return func(req *gohst.Request, res *gohst.Response) {

			var raw string = req.Headers["Authorization"]

			var bearer string = ""
			var token string = ""
			if len(raw) >= 7 {
				bearer = raw[:6]
				token = raw[7:]
			}

			if bearer != "Bearer" || token == "" {
				res.StatusCode = 401
				res.Body = `
				<body>
					<h1>401 Unauthorized</h1>
				</body>
				`
				return
			}

			req.Context["token"] = token

			next(req, res)
		}
	})

	// Authorization middleware
	server.Use(func(next gohst.HandlerFunc) gohst.HandlerFunc {
		return func(req *gohst.Request, res *gohst.Response) {

			var token string = req.Context["token"].(string)

			if token != "banana" {
				res.StatusCode = 403
				res.Body = `
				<body>
					<h1>403 Forbidden</h1>
				</body>
				`
				return
			}

			next(req, res)
		}
	})

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

	stop, err := server.ListenAndServe(":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
	defer close(stop)
	<-stop
}
