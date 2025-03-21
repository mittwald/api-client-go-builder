# mittwald Golang-SDK utilities

> [!IMPORTANT]
> This repository contains tools for automatically generating the Golang SDK. If you only want to use the mittwald mStudio API in your Golang project, there should be no need for you to interact with this project; please use the [mittwald Go client](https://github.com/mittwald/api-client-go) in this case.

## Using the generator toolkit

### Generating the client locally

Install this builder, or invoke it directly from source:

```
$ go install github.com/mittwald/api-client-go-builder/cmd/mittwald-go-client-builder@latest

$ # alternatively:
$ git clone https://github.com/mittwald/api-client-go-builder
```

After cloning this repository, you can generate the client locally. The following commands assume that you have a local checkout of the `github.com/mittwald/api-client-go` package available in your local working directory:

```bash
$ mittwald-go-client-builder generate --url=https://api.mittwald.de/v2/openapi.json --target=./mittwaldv2/generated --pkg=mittwaldv2
$ # or from a local file instead:
$ mittwald-go-client-builder generate --path=to/your/local/openapi.json --target=./mittwaldv2/generated --pkg=mittwaldv2

$ # alternatively:
$ go run ./cmd/mittwald-go-client-builder/main.go generate --url=https://api.mittwald.de/v2/openapi.json --target=./mittwaldv2/generated --pkg=mittwaldv2
```

After generating, run the code formatting (not part of the `generate` command because it takes a long time) and the tests and commit the changes:

```bash
$ goimports -w ./mittwaldv2/generated
$ go test ./...
$ git add .
$ git commit -m "Update generated client"
```
