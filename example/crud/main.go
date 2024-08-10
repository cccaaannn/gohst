package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cccaaannn/gohst"
)

type Result struct {
	Message string `json:"message,omitempty"`
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Define dummy users
var users = []User{
	{Id: 1, Name: "Can kurt", Age: 30},
	{Id: 2, Name: "Banana king", Age: 25},
}

func main() {
	// Create a new server
	server := gohst.CreateServer()

	// Set default headers
	headers := map[string]string{"Content-Type": "application/json"}
	server.SetHeaders(headers)

	server.AddHandler("GET /users", func(req *gohst.Request, res *gohst.Response) {
		s := req.Query["search"]
		filteredUsers := make([]User, 0)
		if s != "" {
			for _, u := range users {
				if strings.Contains(u.Name, s) {
					filteredUsers = append(filteredUsers, u)
				}
			}
		} else {
			filteredUsers = make([]User, len(users))
			copy(filteredUsers, users)
		}

		json, _ := json.Marshal(filteredUsers)
		res.Body = string(json)
	})

	server.AddHandler("GET /users/:id", func(req *gohst.Request, res *gohst.Response) {
		id := req.Params["id"]
		userId, err := strconv.Atoi(id)
		if err != nil {
			res.StatusCode = 400
			json, _ := json.Marshal(Result{Message: "Invalid user id"})
			res.Body = string(json)
			return
		}

		user := User{}
		for _, u := range users {
			if u.Id == userId {
				user = u
				break
			}
		}

		if user.Id == 0 {
			res.StatusCode = 404
			json, _ := json.Marshal(Result{Message: "User not found"})
			res.Body = string(json)
			return
		}

		json, _ := json.Marshal(user)
		res.Body = string(json)
	})

	server.AddHandler("POST /users", func(req *gohst.Request, res *gohst.Response) {
		user := User{}
		json.Unmarshal([]byte(req.Body), &user)
		user.Id = len(users) + 1
		users = append(users, user)

		json, _ := json.Marshal(user)
		res.Body = string(json)
		res.StatusCode = 201
	})

	server.AddHandler("PUT /users/:id", func(req *gohst.Request, res *gohst.Response) {
		id := req.Params["id"]
		userId, err := strconv.Atoi(id)
		if err != nil {
			res.StatusCode = 400
			json, _ := json.Marshal(Result{Message: "Invalid user id"})
			res.Body = string(json)
			return
		}

		user := User{}
		json.Unmarshal([]byte(req.Body), &user)
		user.Id = userId

		for i, u := range users {
			if u.Id == userId {
				users[i] = user
				break
			}
		}

		json, _ := json.Marshal(user)
		res.Body = string(json)
	})

	server.AddHandler("DELETE /users/:id", func(req *gohst.Request, res *gohst.Response) {
		id := req.Params["id"]
		userId, err := strconv.Atoi(id)
		if err != nil {
			res.StatusCode = 400
			json, _ := json.Marshal(Result{Message: "Invalid user id"})
			res.Body = string(json)
			return
		}

		for i, u := range users {
			if u.Id == userId {
				users = append(users[:i], users[i+1:]...)
				break
			}
		}

		json, _ := json.Marshal(Result{Message: "User deleted"})
		res.Body = string(json)
	})

	server.AddHandler("GET /about", func(req *gohst.Request, res *gohst.Response) {
		res.Headers["Content-Type"] = "text/html"
		res.Body = `
		<body>
			<h1>About</h1>
			<p>This is a simple web server written in Go</p>
		</body>
		`
	})

	server.AddHandler("/*", func(req *gohst.Request, res *gohst.Response) {
		json, _ := json.Marshal(Result{Message: "Page not found"})
		res.Body = string(json)
		res.StatusCode = 404
	})

	stop, err := server.ListenAndServe(":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
	defer close(stop)
	<-stop
}
