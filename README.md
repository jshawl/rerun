# 📺 rerun

a task rerunner

## Usage

Create a file `rerun.toml` with a list of steps:

```toml
steps = [
    "go mod tidy",
    "golangci-lint run",
    "go test -cover -v ./...",
    "go build"
]
```
