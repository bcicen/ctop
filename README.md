<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="ctop"/></p>

#

![release][release] ![homebrew][homebrew]

Top-like interface for container metrics

`ctop` provides a concise and condensed overview of real-time metrics for multiple containers:
<p align="center"><img src="_docs/img/grid.gif" alt="ctop"/></p>

as well as an [single container view][single_view] for inspecting a specific container.

`ctop` comes with built-in support for Docker and runC; connectors for other container and cluster systems are planned for future releases.

## Install

Fetch the [latest release](https://github.com/bcicen/ctop/releases) for your platform:

#### Linux

```bash
sudo wget https://github.com/bcicen/ctop/releases/download/v0.7/ctop-0.7-linux-amd64 -O /usr/local/bin/ctop
sudo chmod +x /usr/local/bin/ctop
```

#### OS X

```bash
brew install ctop
```
or
```bash
sudo curl -Lo /usr/local/bin/ctop https://github.com/bcicen/ctop/releases/download/v0.7/ctop-0.7-darwin-amd64
sudo chmod +x /usr/local/bin/ctop
```

#### Docker

```bash
docker run --rm -ti \
  --name=ctop \
  -v /var/run/docker.sock:/var/run/docker.sock \
  quay.io/vektorlab/ctop:latest
```

`ctop` is also available for Arch in the [AUR](https://aur.archlinux.org/packages/ctop-bin/)

## Building

Build steps can be found [here][build].

## Usage

`ctop` requires no arguments and uses Docker host variables by default. See [connectors][connectors] for further configuration options.

### Config file

While running, use `S` to save the current filters, sort field, and other options to a default config path. These settings will be loaded and applied the next time `ctop` is started.

### Options

Option | Description
--- | ---
-a	| show active containers only
-f \<string\> | set an initial filter string
-h	| display help dialog
-i  | invert default colors
-r	| reverse container sort order
-s  | select initial container sort field
-scale-cpu	| show cpu as % of system total
-v	| output version information and exit

### Keybindings

Key | Action
--- | ---
\<enter\> | Open container menu
a | Toggle display of all (running and non-running) containers
f | Filter displayed containers (`esc` to clear when open)
H | Toggle ctop header
h | Open help dialog
s | Select container sort field
r | Reverse container sort order
o | Open single view
l | View container logs (`t` to toggle timestamp when open)
S | Save current configuration to file
q | Quit ctop

[build]: _docs/build.md
[connectors]: _docs/connectors.md
[single_view]: _docs/single.md
[release]: https://img.shields.io/github/release/bcicen/ctop.svg "ctop"
[homebrew]: https://img.shields.io/homebrew/v/ctop.svg "ctop"
