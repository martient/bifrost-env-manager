package environmentmanager

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Config struct {
	Version              string `json:"version"`
	SoftwareTarget       string `json:"software_target"`
	Filename             string `json:"filename"`
	StaticVariables      []Variable
	RandomValueVariables []Variable
	CustomValueVariables []Variable
	EnvVariables         []Variable
	ReadOnlyVariables    []Variable
	ExistingVariables    []Variable
}
