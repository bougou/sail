package models

import "github.com/bougou/sail/pkg/ansible"

type ActionHosts struct {
	Action string
	Hosts  []string
}

func PatchAnsibleGroup(g *ansible.Group, ah *ActionHosts) {
	switch ah.Action {
	case "add":
		g.AddHosts(ah.Hosts...)
	case "remove":
		g.RemoveHosts(ah.Hosts...)
	case "update":
		m := make(map[string]map[string]interface{})
		g.Hosts = &m
		g.AddHosts(ah.Hosts...)
	}
}
