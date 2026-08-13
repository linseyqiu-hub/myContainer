package types

type Config struct {
	ID         string
	Child      int
	Hostname   string
	MemMax     string
	TargetCmd  string
	TargetArgs []string
}
