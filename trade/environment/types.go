package environment

type SourceType int8

const (
	OpenSource SourceType = iota
	HighSource
	LowSource
	CloseSource
)
