package main

import (
	"github.com/hashicorp/go-plugin"
	//"os"
	//"github.com/pigen-dev/cloudbuild-plugin/helpers"
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
	// 				"project_number": "775465758731",
	// 				"project_id":     "aidodev",
	// 				"project_region": "europe-west1",
	// 		},
	// 		"options": map[string]any{
	// 			"logging": "CLOUD_LOGGING_ONLY",
	// 		},
	// 		"github_url":    "https://github.com/aidocons/java-complete-pipeline.git",
	// 		"target_branch": "^develop$",
	// 	},
	// 	"steps": []map[string]any{
	// 		{
	// 			"step":         "mvn-compile",
	// 			"placeholders": map[string]any{
	// 				"flags": "",
	// 			}, // empty placeholders as dummy
	// 		},
	// 		{
	// 			"step":         "Docker-build-push",
	// 			"placeholders": map[string]any{
	// 				"image": "testing/add:latest",
	// 			}, // empty placeholders as dummy
	// 		},
	// 	},
	// }
	// pigenFile := shared.PigenStepsFile{}
	// err := helpers.YamlConfigParser(pigenFileDummy, &pigenFile)
	// if err != nil {
	// 	panic(err)
	// }
	// cloudbuildPlugin := &pkg.Cloudbuild{}
	// cicdFile := cloudbuildPlugin.GeneratScript(pigenFile)
	// if cicdFile.Error != nil {
	// 	panic(cicdFile.Error)
	// }
	// os.WriteFile("cloudbuild.yaml", cicdFile.FileScript, 0644)
	cloudbuildPlugin := &pkg.Cloudbuild{}
	pluginMap := map[string]plugin.Plugin{"cicdPlugin": &shared.CicdPlugin{Impl: cloudbuildPlugin}}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
	})
}