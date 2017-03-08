<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="cTop"/></p>
#

Top-like interface for container metrics

ctop provides a concise and condensed overview of real-time metrics for multiple containers:
![compact][compact]

as well as an expanded view for examining a specific container:
<p align="center"><img width="80%" src="_docs/img/grid.gif" alt="cTop"/></p>

## Install

Fetch the [latest release](https://github.com/bcicen/ctop/releases) for your platform:

#### Linux

```bash
wget https://github.com/bcicen/ctop/releases/download/v0.4/ctop-0.4-linux-amd64 -O ctop
sudo mv ctop /usr/local/bin/
sudo chmod +x /usr/local/bin/ctop
```

#### OS X

```bash
curl -Lo ctop https://github.com/bcicen/ctop/releases/download/v0.4/ctop-0.4-darwin-amd64
sudo mv ctop /usr/local/bin/
sudo chmod +x /usr/local/bin/ctop
```

## Usage

cTop requires no arguments and will configure itself using the `DOCKER_HOST` environment variable
```bash
export DOCKER_HOST=tcp://127.0.0.1:4243
ctop
```

### Keybindings

Key | Action
--- | ---
a | Toggle display of all (running and non-running) containers
f | Filter displayed containers
H | Toggle cTop header
h | Open help dialog
s | Select container sort field
r | Reverse container sort order
q | Quit cTop

[compact]: http://i.imgur.com/uDUq33N.gif "cTop"
