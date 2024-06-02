# http-go-server

## What this is

This repo contains a simple/basic HTTP server in Go, with a basic code organization.
We use:
* net/http package to start and serve HTTP server
* Gorilla mux to handle routes
* Swagger in lorder to serve a REST API compliant with OpenAPI specs

## Pre-requisites

Install Go in 1.13 version minimum.

## Build the app

`$ go build -o bin/http-go-server internal/main.go`

or

`$ make build`

## Run the app

`$ ./bin/http-go-server`

## Test the app

```
$ curl http://localhost:8080/healthz
OK

$ curl http://localhost:8080/hello/ForgeItUp

```



### Request & Response Example

Swagger doc: [http-go-server](https://github.com/scraly/http-go-server/doc/index.html)

|      URL      | Port  | HTTP Method | Operation                                         |
| :-----------: | :---: | :---------: | ------------------------------------------------- |
|   /healthz    | 8080  |     GET     | Test if the app is running                        |
| /hello/{name} | 8080  |     GET     | Returns message with {name} provided in the query |  |


`$ curl localhost:8080/hello/Forge`

## Generate swagger files

After editing `pkg/swagger/swagger.yml` file you need to generate swagger files again:

`$ make gen.swagger`

## Test swagger file validity

`$ make swagger.validate`

## Generate swagger documentation

`$ make swagger.doc`

# Making Changes

1. Make some changes to the source code or just this ```README.md``` file in a branch
2. create a PR
3. You will see checks appear in your PR shortly

The steps for this repo are:

```markdown

1. forge--pr-checks--build  
2. forge--pr-checks--test  
3. forge--pr-checks--security  
4. forge--pr-checks--sca

```

Get list of commands by just commenting on the Pull Request

```markdown
/forge help 
```  
