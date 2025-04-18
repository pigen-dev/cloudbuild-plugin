package cloudbuild

import (
	"context"

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

func (cb *Cloudbuild) ConnectRepo(pigenFile shared.PigenStepsFile) (shared.ActionRequired) {
	// err := helpers.OpenBrowser("https://www.google.com")
	// if err != nil {
	// 	return err
	// }
	ctx := context.Background()
	err := cb.ParseConfig(pigenFile)
	if err != nil {
		return shared.ActionRequired{
			ActionUrl: "",
			Error: err,
		}
	}
	githubConfig, err := helpers.ParseGithubUrl(cb.GithubUrl)
	if err != nil {
		return shared.ActionRequired{
			ActionUrl: "",
			Error: err,
		}
	}
	actionResponse := cb.create_github_connection(ctx, githubConfig)
	return actionResponse
}

func (cb *Cloudbuild) ParseConfig(pigenFile shared.PigenStepsFile) error {
	err := helpers.YamlConfigParser(pigenFile.Config, cb)
	if err != nil {
		return err
	}
	return nil
}

//projects/aidodev/locations/europe-west1/triggers/0317d4ad-717e-471a-98fa-5acaa4f4787f

