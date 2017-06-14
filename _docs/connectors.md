# connectors

`ctop` comes with the below native connectors, enabled via the `--connector` option.

Default connector behavior can be changed by setting the relevant environment variables.

## Docker

Default connector, configurable via standard [Docker commandline varaibles](https://docs.docker.com/engine/reference/commandline/cli/#environment-variables)

#### Options

Var | Description
--- | ---
DOCKER_HOST | Daemon socket to connect to (default: `unix://var/run/docker.sock`)

## RunC

Using this connector requires full privileges to the local runC root dir (default: `/run/runc`)

#### Options

Default connector behavior can be changed by setting the below environment variables:

Var | Description
--- | ---
RUNC_ROOT | path to runc root (default: `/run/runc`)
RUNC_SYSTEMD_CGROUP | if set, enable systemd cgroups
