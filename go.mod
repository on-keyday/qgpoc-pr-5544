module github.com/on-keyday/qgpoc

go 1.25.5

// replace github.com/quic-go/quic-go => github.com/on-keyday/quic-go v0.0.0-20260118200636-5a7fa253b928

require (
	github.com/quic-go/quic-go v0.59.0
	golang.org/x/net v0.43.0
	golang.org/x/sys v0.35.0
)

require golang.org/x/crypto v0.41.0 // indirect
