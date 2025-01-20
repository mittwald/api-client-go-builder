# mittwald Golang-SDK utilities

> [!IMPORTANT]
> This repository contains tools for automatically generating the Golang SDK. If you only want to use the mittwald mStudio API in your Golang project, there should be no need for you to interact with this project; please use the [mittwald Go client](https://github.com/mittwald/api-client-go) in this case.

## Using the generator toolkit

### Automatic client generation

This repository is configured to automatically build and publish the PHP client using the `generate` GitHub Action. This action is triggered by a daily schedule, but can also be triggered manually:

```
$ gh workflow run generate.yml
```

### Generating the client locally

Install this builder, or invoke it directly from source:

```
$ go install github.com/mittwald/api-client-go-builder/cmd/mittwald-go-client-builder@latest

$ # alternatively:
$ git clone https://github.com/mittwald/api-client-go-builder
```

After cloning this repository, you can generate the client locally. The following commands assume that you
have a local checkout of the `github.com/mittwald/api-client-go` package available at `path/to/api-client-go`:

```bash
$ mittwald-go-client-builder https://api.mittwald.de/v2/openapi.json path/to/api-client-go/mittwaldv2/generated mittwaldv2

$ # alternatively:
$ go run ./cmd/mittwald-go-client-builder/main.go https://api.mittwald.de/v2/openapi.json path/to/api-client-go/mittwaldv2/generated mittwaldv2
```

After generating, you should switch into the client directory, run the code formatting (not part of the `generate` command because it takes a long time) and the tests and commit the changes:

```bash
$ cd path/to/api-client-go
$ goimports -w mittwaldv2/generated
$ go test ./...
$ git add .
$ git commit -m "Update generated client"
```
