package types

type Service struct {
	Name        string
	PodTemplate Pod
	Replicas    int
	Pods        []string
	Traffic     int
	CPUUsage    float64
}
