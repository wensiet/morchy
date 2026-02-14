package app

type Config struct {
	Port              int
	DBConnString      string
	LeaseLifetimeSec  int
	EventListLimit    int
	StuckTimeoutSec   int
}
