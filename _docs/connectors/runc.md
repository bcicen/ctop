# runc

`ctop` comes with a native connector for runc containers and can be enabled via the `--connector runc` option.

## Options

Default connector behavior can be changed by setting the below environment variables:

Var | Description
--- | ---
RUNC_ROOT | path to runc root (default: `/run/runc`)
RUNC_SYSTEMD_CGROUP | if set, enable systemd cgroups

