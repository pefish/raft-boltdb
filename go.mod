module github.com/hashicorp/raft-boltdb

go 1.12

require (
	github.com/hashicorp/go-msgpack v0.5.5
	github.com/pefish/bolt v1.3.1
	github.com/pefish/raft v1.1.0
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed // indirect
)

replace github.com/pefish/raft => ./raft

replace github.com/pefish/bolt => ./bolt
