package runtime

type ObjectMeta struct {
	Version string `json:"version"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourceLimits struct {
	CPU uint `json:"cpu"`
	RAM uint `json:"ram"`
}

type Container struct {
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Command   []string       `json:"command"`
	Env       []EnvVar       `json:"env"`
	Resources ResourceLimits `json:"resources"`
}
