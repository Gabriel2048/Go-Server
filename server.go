package main

import (
	"errors"
	"fmt"
	"http-server/server"
	"os"
	"slices"
	"strings"
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

	serv.MapGet("/files/{file-name}", func(request server.HttpRequest) server.HttpResponse {
		fileName := request.RouteParameters["file-name"]
		directory, err := getDirectory()

		if err != nil {
			return *server.NewNotFound()
		}

		filePath := directory + fileName
		fileContent, err := os.ReadFile(filePath)

		if err != nil {
			return *server.NewNotFound()
		}
		respone := server.NewHttpResponse(200)
		respone.SetOctetStreamBody(fileContent)

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

	serv.MapPost("/files/{file-name}", func(request server.HttpRequest) server.HttpResponse {
		fileName := request.RouteParameters["file-name"]
		directory, err := getDirectory()

		if err != nil {
			return *server.NewNotFound()
		}

		file, err := os.Create(directory + fileName)

		if err != nil {
			return *server.NewNotFound()
		}

		body, err := request.ReadBody()

		if err != nil {
			return *server.NewNotFound()
		}

		_, err = file.Write(body)

		if err != nil {
			return *server.NewNotFound()
		}

		respone := server.NewHttpResponse(201)

		return *respone
	})

	serv.RunOnPort("4221")
}

func getDirectory() (string, error) {
	args := os.Args
	if len(args) > 0 {
		directoryIndex := slices.IndexFunc(args, func(arg string) bool {
			return strings.HasPrefix(arg, "--directory")
		})
		if directoryIndex >= 0 {
			return args[directoryIndex+1], nil
		}
	}
	return "", errors.New("directory parameter not found")
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
