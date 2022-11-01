# Termite

[![codecov](https://codecov.io/gh/twelvelabs/termite/branch/main/graph/badge.svg?token=7BSJPVRDPZ)](https://codecov.io/gh/twelvelabs/termite)

Termite is collection of utilities for building CLI tools in Go.

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
    "github.com/twelvelabs/termite/ioutil"
    "github.com/twelvelabs/termite/ui"
)

func main() {
    ios := ioutil.System()
    messenger := ui.NewMessenger(ios)

    messenger.Info("Working...")
    ios.StartProgressIndicator()
    ios.StopProgressIndicator()
    messenger.Success("Done")
}
```

Output:

```text
• Working...
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
