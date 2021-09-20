package ansible

// ActionHosts wrapps the Action and Hosts list.
// Action can be "add", "remove", "update".
// Hosts is a list of host addresses.
type ActionHosts struct {
	Action string
	Hosts  []string
}

func PatchAnsibleGroup(g *Group, ah *ActionHosts) {
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
