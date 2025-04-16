package cloudbuild

import (
	"context"
	"fmt"

	"github.com/pigen-dev/cloudbuild-plugin/helpers"
	shared "github.com/pigen-dev/shared"
)

type Cloudbuild struct {
	Deployment Deployment `yaml:"deployment"`
	GithubUrl string `yaml:"github_url"`
	TargetBranch string `yaml:"target_branch"`
}

type Deployment struct {
	Target string `yaml:"target"`
	Config Config `yaml:"config"`
}

type Config struct {
	ProjectNumber string `yaml:"project_number"`
	ProjectId string `yaml:"project_id"`
	ProjectRegion string `yaml:"project_region"`
}

func (cb *Cloudbuild) ConnectRepo(pigenFile map[string] any) error {
	ctx := context.Background()
	err := cb.ParseConfig(pigenFile)
	if err != nil {
		return err
	}
	githubConfig, err := helpers.ParseGithubUrl(cb.GithubUrl)
	if err != nil {
		return err
	}
	err = cb.create_github_connection(ctx, githubConfig)
	if err != nil {
		return fmt.Errorf("error creating github connection")
	}
	return nil
}

func (cb *Cloudbuild) ParseConfig(pigenFile map[string] any) error {
	pigen := shared.PigenSteps{}
	err := helpers.YamlConfigParser(pigenFile, &pigen)
	if err != nil {
		return err
	}
	err = helpers.YamlConfigParser(pigenFile, cb)
	if err != nil {
		return err
	}
	return nil
}

//projects/aidodev/locations/europe-west1/triggers/0317d4ad-717e-471a-98fa-5acaa4f4787f

