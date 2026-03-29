package app

type Config struct {
	ControlPlaneURL string
	NodeID          string
	ReservedRAM     uint
	ReservedCPU     uint
	ClientCertFile  string
	ClientKeyFile   string
	CACertFile      string
}
