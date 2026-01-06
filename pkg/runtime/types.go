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

type NetConfig struct {
	ContainerPort int    `json:"container_port"`
	HostPort      int    `json:"host_port"`
	Protocol      string `json:"protocol"`
}

type Container struct {
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Command   []string          `json:"command"`
	Env       []EnvVar          `json:"env"`
	Resources ResourceLimits    `json:"resources"`
	Labels    map[string]string `json:"labels"`
	NetConfig *NetConfig        `json:"net_config"`
}

type ContainerBrief struct {
	Name   string            `json:"name"`
	Image  string            `json:"image"`
	Labels map[string]string `json:"labels"`
}

type ContainerFilters struct {
	Labels map[string]string `json:"labels"`
}
