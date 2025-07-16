package types

type Service struct {
	Name        string
	PodTemplate Pod
	Replicas    int
	Pods        []string //todo: change this to map[string]string represent node_id to pod/container id because service has to know in where node this pod/container deployed
	Traffic     int
	CPUUsage    float64
	//todo: add balancer here, every service has load balancer to distribute traffic across its pods
}
