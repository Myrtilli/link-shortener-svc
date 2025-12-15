package data

type MasterQ interface {
	New() MasterQ

	URL() URLdb
}