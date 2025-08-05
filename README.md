# go-cmd

A simplified console program framework for quickly building Go command-line applications.

## Features

- Easy to define commands (`run`, `start`, `stop`, `restart`, `status`, `version`, `help`)
- Built-in process management (PID file, start/stop/restart/status)
- Signal handling for graceful shutdown
- Customizable program info and hooks

## Installation

```sh
go get github.com/xyzj/gocmd
```

## Usage

Create a main program using the framework:

```go
package main

import (
    "github.com/xyzj/gocmd"
)

func main() {
    gocmd.DefaultProgram(&gocmd.Info{
        Title:    "a test program",
        Descript: "this is a console program",
        Ver:      "v0.0.1",
    }).Execute()
    // your code here...
}
```

## Built-in Commands

- `run`      - Run the program in the foreground.
- `start`    - Start the program in the background.
- `stop`     - Stop the running program.
- `restart`  - Restart the program.
- `status`   - Show process status.
- `version`  - Show version info.
- `help`     - Show help message.

## Custom Commands

You can add your own commands using:

```go
program := gocmd.NewProgram(&gocmd.Info{...})
program.AddCommand(&gocmd.Command{
    Name: "yourcmd",
    Descript: "your command description",
    RunWithExitCode: func(pinfo *gocmd.ProcInfo) int {
        // your logic
        return 0
    },
})
program.Execute()
```

## Process Management

- PID file is automatically managed.
- Supports graceful shutdown via signals (`SIGINT`, `SIGTERM`, `SIGQUIT`).

##