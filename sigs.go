package main

type Sig struct {
	name          string
	strongMatches []string
	weakMatches   []string
}

var sigNetwork = Sig{
	name:          "network",
	strongMatches: []string{"kube-proxy", "kube proxy"},
	weakMatches:   []string{"network", "service", "ingress", "connection"},
}

var sigNode = Sig{
	name:          "node",
	strongMatches: []string{"kubelet"},
	weakMatches:   []string{"node"},
}

var sigScheduling = Sig{
	name:          "scheduling",
	strongMatches: []string{"sheduler"},
	weakMatches:   []string{"schedule"},
}

var allSigs = []Sig{
	sigNetwork,
	sigNode,
	sigScheduling,
}
