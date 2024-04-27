package main

import (
	"log"

	"github.com/Community-Sourced-Minecraft/Gate-Proxy/internal/hosting"
	"github.com/Community-Sourced-Minecraft/Gate-Proxy/plugins/core"
	"github.com/Community-Sourced-Minecraft/Gate-Proxy/plugins/motd"
	"github.com/Community-Sourced-Minecraft/Gate-Proxy/plugins/permissions"
	"github.com/Community-Sourced-Minecraft/Gate-Proxy/plugins/tab"
	"github.com/Community-Sourced-Minecraft/Gate-Proxy/plugins/whitelist"

	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type PluginCreator = func(h *hosting.Hosting) (proxy.Plugin, error)

func main() {
	permissionsFile, err := permissions.ReadFile("permissions.json")
	if err != nil {
		log.Fatal(err)
	}

	h, err := hosting.Init()
	if err != nil {
		log.Fatal(err)
	}

	var plugins = []PluginCreator{
		func(h *hosting.Hosting) (proxy.Plugin, error) {
			return core.New(h)
		},
		func(_ *hosting.Hosting) (proxy.Plugin, error) {
			return permissions.New(permissionsFile)
		},
		func(_ *hosting.Hosting) (proxy.Plugin, error) {
			return whitelist.New(h, permissionsFile)
		},
		motd.New,
	}

	proxy.Plugins = append(proxy.Plugins,
		tab.Plugin,
	)

	for _, create := range plugins {
		p, err := create(h)
		if err != nil {
			log.Fatal(err)
		}
		proxy.Plugins = append(proxy.Plugins, p)
	}

	gate.Execute()
}
