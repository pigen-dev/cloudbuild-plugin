package cloudbuild

type StepTemplatesFile struct {
	Steps []CloudbuildStep `yaml:"steps"`
}

type CloudbuildStep struct {
	Id string `yaml:"id"`
	Name string `yaml:"name"`
	Entrypoint string `yaml:"entrypoint"`
	Args []string
}

