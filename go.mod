module github.com/jedisct1/piknik

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	golang.design/x/clipboard v0.5.3
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)

replace golang.design/x/clipboard v0.5.3 => github.com/andrewrech/clipboard v0.5.3
