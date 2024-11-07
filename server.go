package main

import (
	"fmt"
	"http-server/server"
)

func main() {

	serv := server.NewServer()

	serv.MapGet("/", func(request server.HttpRequest) server.HttpResponse {
		response := server.NewHttpResponse(200)
		return *response
	})

	serv.MapGet("/office/{office-id}/user/{user-id}", func(request server.HttpRequest) server.HttpResponse {

		officeId := request.RouteParameters["office-id"]
		userId := request.RouteParameters["user-id"]

		response := server.NewHttpResponse(200)

		response.SetTextBody(fmt.Sprintf("office-id=%s user-id=%s", officeId, userId))

		return *response
	})

	serv.MapGet("/echo/{echo-value}", func(request server.HttpRequest) server.HttpResponse {
		echoValue := request.RouteParameters["echo-value"]

		respone := server.NewHttpResponse(200)
		respone.SetTextBody(echoValue)

		return *respone
	})

	serv.MapGet("/user-agent", func(request server.HttpRequest) server.HttpResponse {
		userAgent, hasUserAgent := request.Headers.GetHeaderValue("User-Agent")

		if !hasUserAgent {
			return *server.NewNotFound()
		}

		response := server.NewHttpResponse(200)
		response.SetTextBody(userAgent)

		return *response
	})

	serv.MapGet("/user", func(request server.HttpRequest) server.HttpResponse {
		response := server.NewHttpResponse(200)

		response.SetJsonBody(User{
			Name: "John",
			Age:  40,
		})

		return *response
	})

	serv.RunOnPort("4221")
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
