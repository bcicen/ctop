# ctop

Top-like interface for container metrics

## Install

Fetch the [latest release](https://github.com/bcicen/ctop/releases) for your platform:

#### Linux

```bash
wget https://github.com/bcicen/ctop/releases/download/v0.1/ctop-0.1-linux-amd64 -O ctop
sudo mv ctop /usr/local/bin/
sudo chmod +x /usr/local/bin/ctop
```

#### OS X

```bash
curl -Lo ctop https://github.com/bcicen/ctop/releases/download/v0.1/ctop-0.1-darwin-amd64
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
h | Open help dialog
s | Select container sort field
r | Reverse container sort order
q |  Quit ctop
