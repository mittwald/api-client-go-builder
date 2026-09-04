# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A code generator that turns the mittwald mStudio OpenAPI v3 spec into the Go SDK published at
[github.com/mittwald/api-client-go](https://github.com/mittwald/api-client-go). This repo contains **only the
generator**; the generated client lives in the other repository and is fully overwritten on each run.

## ⚠️ Generated code does not compile until `goimports` has run

Emitted files unconditionally import `github.com/google/uuid` and `.../pkg/httperr`, and nothing else — the
generator does not track imports at all. It relies entirely on `goimports -w` being run over the output as a
separate step afterwards. This is deliberate: doing it in-process is too slow (see the commented-out block at the
end of `pkg/generator/type_store.go`) and the blind imports in `SchemaName.BuildRoot` exist precisely because
goimports will clean them up.

Consequences:

- Never judge a generation by compile errors before goimports has run — they are meaningless.
- Never "fix" a missing import by adding it to `BuildRoot`; check that goimports actually ran instead.
- Every generation recipe below ends with goimports. Do not skip it, even for a quick spot check.

## Commands

### Working on the generator itself

```bash
go build -v ./cmd/mittwald-go-client-builder      # build
go vet ./...                                      # what CI runs
go test ./...                                     # all tests
go test ./pkg/util -run TestConvertToTypename     # single test
go test ./pkg/util -run 'TestConvertToTypename/.*sftp'   # single subtest
```

CI (`.github/workflows/test.yml`) runs exactly build + vet + test on Go 1.23.

### One-time tool setup

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/mittwald/api-client-go-builder/cmd/mittwald-go-client-builder@latest   # optional; `go run` works too
```

### Full generation run (the real regression test)

This repo has almost no unit tests of its own, so **generating against the live spec and then compiling and
testing the result is the only meaningful way to validate a change.** Assuming a checkout of `api-client-go`
next to this repo:

```bash
# 1. Wipe the target. Required: EmitToFile (pkg/generator/emit.go) errors out
#    rather than overwrite an existing file.
rm -rf ../api-client-go/mittwaldv2/generated

# 2. Generate. Exactly one of --url / --path; --target and --pkg are required.
go run ./cmd/mittwald-go-client-builder generate \
  --target=../api-client-go/mittwaldv2/generated \
  --pkg=mittwaldv2 \
  --url=https://api.mittwald.de/v2/openapi.json

# 3. Fix up imports — mandatory, see above. Run it from the api-client-go root,
#    not from the generated/ subdirectory (imports resolve against the module).
cd ../api-client-go && goimports -w .

# 4. Compile and run the generated ginkgo suites.
go build ./... && go test ./...
```

Iterating against a local spec file avoids re-downloading and lets you shrink the input while debugging:

```bash
curl -o /tmp/openapi.json https://api.mittwald.de/v2/openapi.json
go run ./cmd/mittwald-go-client-builder generate \
  --target=/tmp/out --pkg=mittwaldv2 --path=/tmp/openapi.json
```

The generator logs at debug level and prints every type it observes and emits; grep that output for a struct
name to find out which phase mangled it.

`.github/workflows/test-generation.yml` performs the same four steps against a fresh `api-client-go` checkout on
every PR — a change is not done until that workflow's sequence passes locally.

The second CLI command, `next-version <semver>`, just bumps the patch level and is used by the release pipeline
of the client repo.

## Architecture

### Pipeline (`pkg/generator/generator.go`)

`Generator.Build` runs four phases against a shared `TypeStore`:

1. **Load** — `SpecLoader` (URL or file) parses the spec with `libopenapi` into a v3 document model.
2. **Construct** — every entry under `#/components/schemas` goes through `BuildTypeFromSchema`
   (`type_build.go`), which dispatches on the JSON-schema shape and returns one of the `*Type` implementations.
   Only the top level is built here; nested properties are *not* yet resolved.
3. **Expand** (`TypeStore.BuildSubtypes`) — every type implementing `TypeWithSubtypes` builds its children
   (object properties, array items, map values, oneOf alternatives, request params/bodies, response types),
   registering them back into the store. Since expanding a type *adds* types, this runs as a fixpoint loop until
   the visited count stops growing.
4. **Emit** (`TypeStore.EmitDeclarations`) — each type renders itself into a `gowrtr` statement tree, one file
   per type. Types implementing `TypeWithTestcases` additionally emit a `*_test.go`, and a `suite_test.go`
   ginkgo bootstrap is generated for each package that got test cases.

Phase ordering matters: `ObjectType.EmitDeclaration` panics if `BuildSubtypes` never ran for it.

### The Type abstraction (`pkg/generator/type.go`)

`Type` is the core interface (`Name`, `EmitDeclaration`, `EmitReference`, `IsPointerType`); `SchemaType` adds
`BuildExample`, `Schema`, `IsLightweight`. Behaviour is opted into via small marker interfaces, all checked with
type assertions at the call sites:

- `TypeWithSubtypes` — participates in phase 3
- `TypeWithTestcases` — emits ginkgo specs (see `type_object_testgen.go`, `type_oneof_testgen.go`)
- `TypeWithValidation` — contributes to the generated `Validate() error` method
- `TypeWithStringConversion` — can be rendered into a URL path/query parameter
- `UnpackableType` — wrapper types (`OptionalType`, `ReferenceType`) that delegate to an inner type

When adding a new type kind, add a `type_*.go` file, wire it into the `switch` in `BuildTypeFromSchema`, and
implement whichever marker interfaces apply.

`IsLightweight() == true` (currently `ReferenceType`) means "emits no file of its own" — such types are skipped
during emission but still resolve as references.

### Naming and file placement (`pkg/generator/schema_name.go`)

`SchemaName{PackageKey, PackagePath, StructName}` decides both the Go identifier and the output file for every
type. `ForSubtype(name)` derives child names (appends to the struct name, suffixes the filename), `ForTestcase()`
derives the `_test.go` / `_test` package variant. The mittwald strategy maps spec names like
`de.mittwald.v1.sshuser.SshUser` onto `schemas/sshuserv2/sshuser.go`, package `sshuserv2`, struct `SshUser`.

### Client generation (`generator_clients.go`, `client_set.go`, `client.go`, `client_operation_req.go`)

One `Client` per OpenAPI tag, all aggregated in a single `ClientSet` (`clients/clientset.go`). Each operation
yields a `<Op>Request` type (with `BuildRequest`) and, when the spec declares JSON responses, a `<Op>Response`
type — or a `OneOfType` over per-status alternatives when several 2xx/3xx responses exist. Generated methods
return `(*Response, *http.Response, error)`.

### Output layout

```
<target>/
  clients/clientset.go
  clients/<tag>client<ver>/{client.go,<op>_request.go,<op>_response.go,suite_test.go}
  schemas/<pkg><ver>/<type>.go
```

## Other things that will bite you

- **Naming hacks live in two hand-maintained lists**, and most recent commits to this repo are additions to
  them: `commonInitialisms` in `pkg/util/typename.go` (`Ssh` → `SSH`, `Api` → `API`, …) and `commonPrefixes` in
  `pkg/generator/client.go` (strips redundant prefixes off operation IDs, e.g. `ssh-user-create-ssh-user` →
  `CreateSSHUser`). Pin new initialism behaviour with a case in `pkg/util/typename_test.go`.
- **Explicit nullability.** A schema with exactly one `allOf` entry plus `nullable: true` is a mittwald-specific
  convention meaning "may be explicitly `null`" (as opposed to merely absent, for PATCH semantics) and maps to
  `ExplicitlyNullableType`, not to `OptionalType`.
- **Debugging comments.** `GeneratorContext.WithDebuggingComments` is hardcoded on, so every generated schema
  file embeds the originating JSON schema as a comment — the fastest way to see why a type came out wrong.
- The `mittwald-go-client-builder` binary at the repo root is a stale local build artifact, not a tracked file.
