package main

type Sig struct {
	name          string
	strongMatches []string
	weakMatches   []string
}

var sigApps = Sig{
	name:          "apps",
	strongMatches: []string{"cronjob"},
	weakMatches:   []string{"deployment"},
}

var sigAutoscaling = Sig{
	name:          "autoscaling",
	strongMatches: []string{"hpa", "autoscaler"},
	weakMatches:   []string{},
}

var sigAws = Sig{
	name:          "aws",
	strongMatches: []string{"aws", "eks", "cloud-provider-aws", "aws-alb-ingress-controller", "aws-iam-authenticator", "aws-encryption-provider", "aws-ebs-csi-driver", "aws alb ingress controller", "aws iam authenticator", "aws encryption provider", "aws ebs csi driver"},
	weakMatches:   []string{"iam", "efs", "ebs", "alb ingress", "heptio authenticator"},
}

var sigAzure = Sig{
	name:          "azure",
	strongMatches: []string{"azure"},
	weakMatches:   []string{},
}

var sigCli = Sig{
	name:          "cli",
	strongMatches: []string{},
	weakMatches:   []string{"kubectl"},
}

var sigClusterLifeCycle = Sig{
	name:          "cluster-lifecycle",
	strongMatches: []string{"kubeadm"},
	weakMatches:   []string{},
}

var sigGcp = Sig{
	name:          "gcp",
	strongMatches: []string{"gcp", "glb", "compute engine"},
	weakMatches:   []string{"google"},
}

var sigMulticluster = Sig{
	name:          "multicluster",
	strongMatches: []string{"federation", "cluster registry"},
	weakMatches:   []string{},
}

var sigNetwork = Sig{
	name:          "network",
	strongMatches: []string{"ipv6", "ipvs", "ingress", "kube-dns", "kube dns", "kube-proxy", "kube proxy", "cni"},
	weakMatches:   []string{"envoy", "network", "service", "connection", "calico", "flannel", "istio", "linkerd"},
}

var sigNode = Sig{
	name:          "node",
	strongMatches: []string{},
	weakMatches:   []string{"node", "kubelet"},
}

var sigScheduling = Sig{
	name:          "scheduling",
	strongMatches: []string{"sheduler"},
	weakMatches:   []string{"schedule"},
}

var sigStorage = Sig{
	name:          "storage",
	strongMatches: []string{"persistentvolume"},
	weakMatches:   []string{"pv", "pvc", "efs", "ebs"},
}

var sigVmWare = Sig{
	name:          "vmware",
	strongMatches: []string{"vmware", "vsphere"},
	weakMatches:   []string{},
}

var sigWindows = Sig{
	name:          "windows",
	strongMatches: []string{"windows"},
	weakMatches:   []string{},
}

var allSigs = []Sig{
	sigApps,
	sigAutoscaling,
	sigAws,
	sigAzure,
	sigCli,
	sigClusterLifeCycle,
	sigGcp,
	sigMulticluster,
	sigNetwork,
	sigNode,
	sigScheduling,
	sigStorage,
	sigVmWare,
	sigWindows,
}
