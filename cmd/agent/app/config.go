package app

type Config struct {
	ControlPlaneURL string
	NodeID          string
	ReservedRAM     uint
	ReservedCPU     uint
}
