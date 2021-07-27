module github.com/bcicen/ctop

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/c9s/goprocinfo v0.0.0-20170609001544-b34328d6e0cd
	github.com/fsouza/go-dockerclient v1.7.0
	github.com/gizak/termui v2.3.0+incompatible
	github.com/jgautheron/codename-generator v0.0.0-20150829203204-16d037c7cc3c
	github.com/mattn/go-runewidth v0.0.0-20170201023540-14207d285c6c
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/nsf/termbox-go v0.0.0-20180303152453-e2050e41c884
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/op/go-logging v0.0.0-20160211212156-b2cb9fa56473
	github.com/opencontainers/runc v1.0.0-rc95
	github.com/pkg/errors v0.9.1
)

replace github.com/gizak/termui => github.com/bcicen/termui v0.0.0-20180326052246-4eb80249d3f5

go 1.15
