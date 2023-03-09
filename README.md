# Termite

[![docs](https://pkg.go.dev/badge/github.com/twelvelabs/termite.svg)](https://pkg.go.dev/github.com/twelvelabs/termite)
[![build](https://github.com/twelvelabs/termite/actions/workflows/build.yml/badge.svg)](https://github.com/twelvelabs/termite/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/termite/branch/main/graph/badge.svg?token=7BSJPVRDPZ)](https://codecov.io/gh/twelvelabs/termite)

Termite is a collection of utilities for building CLI tools in Go.

A few things to note:

- **This is still a WIP**. I'm in the process of extracting things here when I find reusable patterns that I'm repeating often. This probably won't be stable or useful for others for a while, so... use at your own risk :grimacing: (at least until a 1.0 release).
- Some of these packages were inspired or copied from the [GitHub CLI](https://github.com/cli/cli) project. I'm a big fan :heart: and I've done my best to give proper attribution per the license. If anyone from that team feels I've misappropriated anything, please let me know and I'll try to make amends.

## Install

```text
go get github.com/twelvelabs/termite
```

## Usage

```go
package main

import (
    "github.com/twelvelabs/termite/api"
    "github.com/twelvelabs/termite/conf"
    "github.com/twelvelabs/termite/ui"
)

func main() {
    type Config struct {
        BaseURL string `default:"https://0.0.0.0/api/v1" validate:"required,url"`
        Debug   bool   `default:"true"`
    }

    // Loads and validates config values from ~/.config/my-app/config.yaml
    config, _ := conf.NewLoader(&Config{}, ConfigFile("my-app")).Load()
    client := api.NewRESTClient(&api.ClientOptions{
        BaseURL: config.BaseURL,
    })

    type APIResponse struct {
        Status string `json:"statusMessage"`
    }

    u := ui.NewUserInterface(ui.NewIOStreams())
    ok, _ := u.Confirm("Proceed?", true, "Some help text...")
    if ok {
        u.StartProgressIndicator("Requesting")
        resp := &APIResponse{}
        err := client.Get("/some/endpoint", resp)
        if err != nil {
            u.Out(u.FailureIcon() + " API failure: %v\n", err)
        } else {
            u.Out(u.InfoIcon() + " API success: %v\n", resp.Status)
        }
        u.StopProgressIndicator()
    }
    u.Out(u.SuccessIcon() + " Done\n")
}
```

Output:

```text
? Proceed? Yes
• API success: some status message
✓ Done
```

## Development

```bash
git clone git@github.com:twelvelabs/termite.git
cd termite

make setup
make build
make test

# Show full usage
make
```
