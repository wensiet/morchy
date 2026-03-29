package app

type Config struct {
	Port              int
	DBConnString      string
	LeaseLifetimeSec  int
	EventListLimit    int
	StuckTimeoutSec   int
	SeedToken         string
	TLSCertFile       string
	TLSKeyFile        string
	TLSCAFile         string
	DevMode           bool
}
