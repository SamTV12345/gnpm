# gnpm — The package manager to rule them all

Gnpm is a versatile package manager that manages your package manager.
Ever bothered with switching between npm, yarn, and pnpm for different projects? 
Gnpm simplifies this by automatically detecting and using the appropriate package manager based on your project's configuration files.

## Description

Gnpm uses your package.json to auto-download, configure, and run the correct package manager for your project.
It supports npm, yarn, bun, and pnpm. Also, it supports downloading the specified runtime. So no matter if you run deno, node, or bun. 
If gnpm can determine with the help of your package.json which runtime you use, it will auto download and configure it for you.
So the only thing you need to remember is `gnpm install` and `gnpm run <script>`. Under the hood it will call the respective package manager with the correct runtime.
It also supports workspaces and monorepos.

## Features
- Auto-detects and uses the correct package manager (npm, yarn, pnpm, bun) based on your project's configuration files.
- Supports multiple runtimes (node, bun, deno) and auto-downloads the required.
- Makes starting in a new team easy because no matter how much experience you have you just need to download gnpm from a release and put it in your path.



## Installation
You can download the latest release from the [Releases](https://github.com/SamTV12345/gnpm/releases) page. After that you need to place the binary into your PATH.
Alternatively, you can build it from source using Go:

- Download and install Go from [golang.org](https://golang.org/dl/).
- Clone the repository: `git clone https://github.com/SamTV12345/gnpm.git`
- Navigate to the project directory: `cd gnpm`
- Build the binary: `go build -o gnpm .`
- Move the binary to a directory in your PATH, e.g., `/usr/local/bin`

## Usage

Once installed, don't remember anything of the above because it is unnecessary. Run your commands as you would with npm, yarn, or pnpm.
If you have, e.g., Node.js configured in your package.json, and you want to call node directly, you can do it like this:

```bash
gnpm node --version
```

The same applies also for bun and deno.

During performing actions like running your scripts the opened shell has the respective package manager and runtime in its PATH. 
So you don't need to do any code changes in your project. Other people can as well use their package manager instead of gnpm so you are not locked into gnpm as a team.
