module github.com/bcicen/ctop

require (
	github.com/BurntSushi/toml v0.3.0
	github.com/c9s/goprocinfo v0.0.0-20170609001544-b34328d6e0cd
	github.com/checkpoint-restore/go-criu v0.0.0-20190109184317-bdb7599cd87b // indirect
	github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50 // indirect
	github.com/coreos/go-systemd v0.0.0-20151104194251-b4a58d95188d // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/fsouza/go-dockerclient v1.4.1
	github.com/gizak/termui v2.3.0+incompatible
	github.com/godbus/dbus v0.0.0-20151105175453-c7fdd8b5cd55 // indirect
	github.com/jgautheron/codename-generator v0.0.0-20150829203204-16d037c7cc3c
	github.com/mattn/go-runewidth v0.0.0-20170201023540-14207d285c6c // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/mrunalp/fileutils v0.0.0-20171103030105-7d4729fb3618 // indirect
	github.com/nsf/termbox-go v0.0.0-20180303152453-e2050e41c884
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/op/go-logging v0.0.0-20160211212156-b2cb9fa56473
	github.com/opencontainers/runc v1.0.0-rc8
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/opencontainers/selinux v1.2.2 // indirect
	github.com/pkg/errors v0.8.1
	github.com/seccomp/libseccomp-golang v0.0.0-20150813023252-1b506fc7c24e // indirect
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/vishvananda/netlink v0.0.0-20150820014904-1e2e08e8a2dc // indirect
	github.com/vishvananda/netns v0.0.0-20180720170159-13995c7128cc // indirect
)

replace github.com/gizak/termui => github.com/bcicen/termui v0.0.0-20180326052246-4eb80249d3f5

go 1.13
