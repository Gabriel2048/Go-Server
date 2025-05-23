package main

import (
	"fmt"
	"http-server/server"
	"http-server/server/builder"
)

func main() {

	serv, _ := server.NewServer(
		builder.WithPort(4221),
		builder.WithHost("0.0.0.0"),
		builder.WithHttpsRedirect(),
	)

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

	serv.Run()
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
