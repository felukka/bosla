# Installation

## With Nix

The repository root is a flake that exposes the `bosla` package for
`x86_64-linux` and `aarch64-linux`.

=== "From the repository"

    ```bash
    git clone git@github.com:felukka/bosla.git
    cd bosla
    nix build .#bosla
    ```

    The binary lands in `result/bin/bosla`.

=== "Anywhere (flake reference)"

    ```bash
    nix build github:felukka/bosla
    ```

You can also run it directly without installing:

```bash
nix run github:felukka/bosla -- --help
```

## From source

```bash
git clone git@github.com:felukka/bosla.git
cd bosla
go build -o bosla ./cmd/bosla
```

!!! note
    There are no tagged releases yet, so source builds track `main`.

## Development shell

```bash
nix develop
```

The shell includes the Go toolchain, `cobra-cli`, `sqlc`, `golangci-lint`,
`golines`, `goimports-reviser`, and a Python environment with
`mkdocs-material` for building the docs.

## Verify

```bash
bosla --help
```

Expected output starts with:

```
A serious, well-scoped versioning tool.
```
