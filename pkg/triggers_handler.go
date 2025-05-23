package cloudbuild

import (
	"context"
	"fmt"
	"log"

	cloudbuildv1 "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpbv1 "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"cloud.google.com/go/cloudbuild/apiv2"
	"github.com/google/uuid"
	"github.com/pigen-dev/cloudbuild-plugin/helpers"
	shared "github.com/pigen-dev/shared"
	"google.golang.org/api/iterator"
)

func (cb *Cloudbuild) CreateTrigger(pigenFile shared.PigenStepsFile) error {
	ctx := context.Background()
	err := cb.ParseConfig(pigenFile)
	if err != nil {
		return err
	}
	githubConfig, err := helpers.ParseGithubUrl(cb.GithubUrl)
	if err != nil {
		return err
	}

	trigger, err := cb.trigger_exist(ctx, githubConfig)
	if err != nil {
		return err
	}
	if trigger == nil {
		log.Printf("Creating trigger on: %v", githubConfig.Url)
		_ , err = cb.create_trigger(ctx, githubConfig)
		if err != nil {
			return err
		}
		return nil
	} else{
		log.Printf("Trigger already exist on: %v", githubConfig.Url)
	}
	return nil
}

func (cb Cloudbuild) trigger_exist(ctx context.Context, githubConfig helpers.GithubUrl) (*cloudbuildpbv1.BuildTrigger, error){
	parent := "projects/"+cb.Deployment.ProjectId+"/locations/" + cb.Deployment.ProjectRegion
	cv1, err := cloudbuildv1.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't create cloudbuild client to check if trigger exist: %v", err)
	}
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)

	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	if err != nil {
			return nil, fmt.Errorf("can't create cloudbuild client to check if trigger exist: %v", err)
	}
	listBuildTriggersRequest := &cloudbuildpbv1.ListBuildTriggersRequest{
		Parent: parent,
		ProjectId: cb.Deployment.ProjectId,
	}
	resp := cv1.ListBuildTriggers(ctx, listBuildTriggersRequest)
	for {
		repo, err := resp.Next()
		if err != nil {
			if err == iterator.Done {
				break
			} else{
				return nil, fmt.Errorf("error while iterating over triggers: %v", err)
			}
		}
		if repo.RepositoryEventConfig == nil {
			continue
		}
		if repo.RepositoryEventConfig.Repository == parent+"/connections/"+githubConfig.Parent+"/repositories/"+githubConfig.Repo && repo.RepositoryEventConfig.GetPush().GetBranch() == cb.TargetBranch {
			log.Printf("Trigger for repository %v on branch %v does exist", githubConfig.Url, cb.TargetBranch)
			return repo, nil
		}
	}
	return nil, nil
}

func (cb Cloudbuild) create_trigger(ctx context.Context, githubConfig helpers.GithubUrl) (*cloudbuildpbv1.BuildTrigger, error) {
	cv1, err := cloudbuildv1.NewClient(ctx)
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)
	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	if err != nil {
			return nil, fmt.Errorf("can't create cloudbuild client to create trigger: %v", err)
	}
	//a trigger name must not exceed 64 character that's why i can't use owner + repo name
	encodedName := uuid.New().String()
	parent := "projects/"+cb.Deployment.ProjectId+"/locations/" + cb.Deployment.ProjectRegion

	pushFilter_Branch := &cloudbuildpbv1.PushFilter_Branch{
		Branch: cb.TargetBranch,
	}

	pushFilter := &cloudbuildpbv1.PushFilter{
		GitRef: pushFilter_Branch,
	}

	buildTrigger_Filename := &cloudbuildpbv1.BuildTrigger_Filename{
		Filename: "cloudbuild.yaml",
	}
	repositoryEventConfig_Push := &cloudbuildpbv1.RepositoryEventConfig_Push{
		Push: pushFilter,
	}

	repositoryEventConfig:= &cloudbuildpbv1.RepositoryEventConfig{
		Repository: parent+"/connections/"+githubConfig.Parent+"/repositories/"+githubConfig.Repo,
		RepositoryType: 1,
		Filter:repositoryEventConfig_Push,
	}
	
	buildTrigger := &cloudbuildpbv1.BuildTrigger{
		Name:encodedName[:8],
		Description:"This is a trigger on " + githubConfig.Parent+"-"+githubConfig.Repo + "-" + cb.TargetBranch,
		BuildTemplate: buildTrigger_Filename,
		ServiceAccount: "projects/"+cb.Deployment.ProjectId+"/serviceAccounts/cloudbuild-sa@aidodev.iam.gserviceaccount.com",
		RepositoryEventConfig: repositoryEventConfig,
	}
	

	createBuildTriggerRequest := &cloudbuildpbv1.CreateBuildTriggerRequest{
		//Parent: "projects/"+cb.Deployment.Config.ProjectNumber+"/locations/" + cb.Deployment.Config.ProjectRegion,
		Parent: parent,
		ProjectId: cb.Deployment.ProjectId,
		Trigger: buildTrigger,
		
	}

	trigger, err := cv1.CreateBuildTrigger(ctx, createBuildTriggerRequest)
	if err != nil {
			return nil, fmt.Errorf("create build trigger failed: %+v", err)
	}
		
		return trigger, nil
}