package main

import (
	"github.com/hashicorp/go-plugin"
	// "fmt"

	// "github.com/pigen-dev/cloudbuild-plugin/helpers"
	pkg "github.com/pigen-dev/cloudbuild-plugin/pkg"
	shared "github.com/pigen-dev/shared"
)



func main(){
	// pigenFileDummy := map[string]any{
	// 	"type": "cloudbuild",
	// 	"version": "v1.0.0",
  // 	"repo_url": "",
	// 	"config": map[string]any{
	// 		"deployment": map[string]any{
	// 			"target": "gcp",
	// 			"config": map[string]any{
	// 				"project_number": 775465758731,
	// 				"project_id":     "aidodev",
	// 				"project_region": "europe-west1",
	// 			},
	// 		},
	// 		"github_url":    "https://github.com/aidocons/java-complete-pipeline.git",
	// 		"target_branch": "^develop$",
	// 	},
	// 	"steps": []map[string]any{
	// 		{
	// 			"step":         "mvn-compile",
	// 			"placeholders": map[string]any{}, // empty placeholders as dummy
	// 		},
	// 	},
	// }
	// pigenFile := shared.PigenStepsFile{}
	// err := helpers.YamlConfigParser(pigenFileDummy, &pigenFile)
	// if err != nil {
	// 	panic(err)
	// }
	// cloudbuildPlugin := &pkg.Cloudbuild{}
	// action := cloudbuildPlugin.ConnectRepo(pigenFile)
	// fmt.Println(action)
	cloudbuildPlugin := &pkg.Cloudbuild{}
	pluginMap := map[string]plugin.Plugin{"cicdPlugin": &shared.CicdPlugin{Impl: cloudbuildPlugin}}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
	})
}