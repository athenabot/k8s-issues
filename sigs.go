package main

type Sig struct {
	name          string
	strongMatches []string
	weakMatches   []string
}

var sigClusterLifeCycle = Sig{
	name:          "cluster-lifecycle",
	strongMatches: []string{"kubeadm"},
	weakMatches:   []string{},
}

var sigNetwork = Sig{
	name:          "network",
	strongMatches: []string{"kube-dns", "kube dns", "kube-proxy", "kube proxy", "cni", "calico", "flannel", "istio", "linkerd"},
	weakMatches:   []string{"envoy", "network", "service", "ingress", "connection"},
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

var sigStorage = Sig{
	name:          "storage",
	strongMatches: []string{"persistentvolume"},
	weakMatches:   []string{"pv", "pvc"},
}

var allSigs = []Sig{
	sigClusterLifeCycle,
	sigNetwork,
	sigNode,
	sigScheduling,
	sigStorage,
}
