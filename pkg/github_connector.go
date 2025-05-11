package cloudbuild

import (
	"context"
	"fmt"
	"log"

	"github.com/pigen-dev/cloudbuild-plugin/helpers"
	shared "github.com/pigen-dev/shared"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv2/cloudbuildpb"
	"google.golang.org/api/iterator"
)



func (cb *Cloudbuild)create_github_connection(ctx context.Context, gitConfig helpers.GithubUrl) (shared.ActionRequired){
	
	log.Printf("Connecting git: %v", gitConfig.Parent)
	connection, actionUri, err := cb.create_connection(ctx, gitConfig)
	if err != nil {
		return shared.ActionRequired{
			ActionUrl: "",
			Error: err,
		}
	}
	if actionUri != "" {
		log.Printf("Connection %v is not complete, action required: %v", gitConfig.Parent, actionUri)
		return shared.ActionRequired{
			ActionUrl: actionUri,
			Error: nil,
		}
	}

	log.Printf("Connection %v does exist", gitConfig.Parent)

	repo, err := cb.repo_connected(ctx, gitConfig)
	if err != nil {
		return shared.ActionRequired{
			ActionUrl: "",
			Error: err,
		}
	}
	if repo == nil {
		log.Printf("Connecting repository: %v", gitConfig.Repo)
		_ , err = cb.connect_repository(ctx, gitConfig, connection)
		if err != nil {
			return shared.ActionRequired{
				ActionUrl: "",
				Error: err,
			}
		}
		return shared.ActionRequired{
			ActionUrl: "",
			Error: nil,
		}
	}
	log.Printf("Repository %v does exist", gitConfig.Repo)
	
	return shared.ActionRequired{
		ActionUrl: "",
		Error: nil,
	}
}

func (cb *Cloudbuild) create_connection(ctx context.Context, gitConfig helpers.GithubUrl) (connection *cloudbuildpb.Connection, actionUrl string, err error) {
	connectionExist, err := cb.connection_exist(ctx, gitConfig)
	if connectionExist != nil {
		state := connectionExist.GetInstallationState().GetStage()
		if state.String() == "COMPLETE" {
			return connectionExist, "", nil
		}
		return nil, connectionExist.GetInstallationState().GetActionUri(), nil
	}
	log.Printf("Connecting github %v: ", gitConfig.Url)
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)
	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	if err != nil {
			return nil, "", fmt.Errorf("can't create cloudbuild client: %v", err)
	}
	
	parent := "projects/" + cb.Deployment.ProjectNumber + "/locations/" + cb.Deployment.ProjectRegion
	
	connection = &cloudbuildpb.Connection{
		Name: parent + "/connections/" + gitConfig.Parent,
		ConnectionConfig: &cloudbuildpb.Connection_GithubConfig{},
	}
	req := &cloudbuildpb.CreateConnectionRequest{
		Parent: parent,
		Connection: connection,
		ConnectionId: gitConfig.Parent,

	}
	fmt.Println(connection.Name)
	createOp, err := c.CreateConnection(ctx, req)
	if err != nil {
		return nil, "", err
	}
	createOp.Wait(ctx)
	op_connection, err := createOp.Poll(ctx)
	if err != nil {
		return nil, "", err
	}

	return nil, op_connection.GetInstallationState().GetActionUri(), nil
}

func (cb *Cloudbuild) connect_repository(ctx context.Context, gitConfig helpers.GithubUrl, connection *cloudbuildpb.Connection) (*cloudbuildpb.Repository, error) {
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)
	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	if err != nil {
			return nil, fmt.Errorf("can't create cloudbuild client: %v", err)
	}

	repository := &cloudbuildpb.Repository{
		RemoteUri: gitConfig.Url,
	}
	createRepositoryRequest := &cloudbuildpb.CreateRepositoryRequest{
		Parent: connection.GetName(),
		Repository: repository,
		RepositoryId: gitConfig.Repo,
	}

	repo, err := c.CreateRepository(ctx, createRepositoryRequest)
	if err != nil {
		return nil, err
	}
	reposirtory, err := repo.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for repository creation: %v", err)
	}
	return reposirtory, nil
}

//Check whether git connection exist or not, if it exist return it else return nil

func (cb *Cloudbuild) connection_exist(ctx context.Context, gitConfig helpers.GithubUrl) (*cloudbuildpb.Connection, error) {
	parent := "projects/"+cb.Deployment.ProjectId+"/locations/" + cb.Deployment.ProjectRegion
	connectionName := parent + "/connections/"+gitConfig.Parent
	log.Printf("checking connection existing : %v", connectionName)
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)
	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	if err != nil {
			return nil, err
	}
	listConnectionsRequest := &cloudbuildpb.ListConnectionsRequest{
		Parent: parent,
	}
	
	connections := c.ListConnections(ctx, listConnectionsRequest)
	for {
		resp, err := connections.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if resp.Name == connectionName {
			return resp, nil
		}
	}
	return nil, nil
}


func (cb *Cloudbuild) repo_connected(ctx context.Context, gitConfig helpers.GithubUrl) (*cloudbuildpb.Repository, error) {
	parent := "projects/"+cb.Deployment.ProjectId+"/locations/" + cb.Deployment.ProjectRegion + "/connections/" + gitConfig.Parent
	repoName := parent + "/repositories/" + gitConfig.Repo
	log.Printf("checking connection existing : %v", repoName)
	c, err := cloudbuild.NewRepositoryManagerClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't create repository manager client")
	}
	defer func(){
		if err == nil {
			c.Close()
		}
	}()
	
	listRepositoriesRequest := &cloudbuildpb.ListRepositoriesRequest{
		Parent: parent,
	}
	
	repositories := c.ListRepositories(ctx, listRepositoriesRequest)

	for {
		resp, err := repositories.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("can't get next repository")
		}
		if resp.Name == repoName {
			return resp, nil
		}
	}
	return nil, nil
}