<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="ctop"/></p>
#

Top-like interface for container metrics

`ctop` provides a concise and condensed overview of real-time metrics for multiple containers:
<p align="center"><img src="_docs/img/grid.gif" alt="ctop"/></p>

as well as an [expanded view][expanded_view] for inspecting a specific container.

`ctop` currently comes with built-in support for Docker; connectors for other container and cluster systems are planned for future releases.

## Install

Fetch the [latest release](https://github.com/bcicen/ctop/releases) for your platform:

#### Linux

```bash
wget https://github.com/bcicen/ctop/releases/download/v0.4.1/ctop-0.4.1-linux-amd64 -O ctop
sudo mv ctop /usr/local/bin/
sudo chmod +x /usr/local/bin/ctop
```

#### OS X

```bash
curl -Lo ctop https://github.com/bcicen/ctop/releases/download/v0.4.1/ctop-0.4.1-darwin-amd64
sudo mv ctop /usr/local/bin/
sudo chmod +x /usr/local/bin/ctop
```

or run via Docker:
```bash
docker run -ti -v /var/run/docker.sock:/var/run/docker.sock quay.io/vektorlab/ctop:latest
```

`ctop` is also available for Arch in the [AUR](https://aur.archlinux.org/packages/ctop/)

#### Build

```
go get -u github.com/kardianos/govendor
go get -u github.com/bcicen/ctop
govendor sync
go build -o ctop
```

## Usage

`ctop` requires no arguments and will configure itself using the `DOCKER_HOST` environment variable
```bash
export DOCKER_HOST=tcp://127.0.0.1:4243
ctop
```

### Keybindings

Key | Action
--- | ---
a | Toggle display of all (running and non-running) containers
f | Filter displayed containers
H | Toggle ctop header
h | Open help dialog
s | Select container sort field
r | Reverse container sort order
q | Quit ctop

[expanded_view]: _docs/expanded.md
