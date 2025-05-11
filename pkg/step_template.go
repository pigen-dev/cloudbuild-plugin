package cloudbuild

type CloudbuildStep struct {
	Id string `yaml:"id"`
	Name string `yaml:"name"`
	Entrypoint string `yaml:"entrypoint"`
	Args []string
}

