package main

import (
	"github.com/hashicorp/go-plugin"
	pkg"github.com/pigen-dev/cloudbuild-plugin/pkg"
	shared "github.com/pigen-dev/shared"
)



func main(){
	cloudbuildPlugin := &pkg.Cloudbuild{}
	pluginMap := map[string]plugin.Plugin{"cicdPlugin": &shared.CicdPlugin{Impl: cloudbuildPlugin}}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
	})
}