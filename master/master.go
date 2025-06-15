// mini-k8s: master/master.go
package master

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	types "dude/types"
)

var nodes = []types.Node{
	{Name: "node1", Address: "http://node1:9001"},
	{Name: "node2", Address: "http://node2:9002"},
}

var services = make(map[string]*types.Service)
var mu sync.Mutex

func StartMaster() {
	loadServices()
	go scheduleInitialPods()
	select {} // keep running
}

func loadServices() {
	data, _ := ioutil.ReadFile("config/services.json")
	var svcList []*types.Service
	_ = json.Unmarshal(data, &svcList)
	for _, svc := range svcList {
		services[svc.Name] = svc
	}
}

func scheduleInitialPods() {
	mu.Lock()
	defer mu.Unlock()

	for _, svc := range services {
		for i := 0; i < svc.Replicas; i++ {
			containerID := deployPod(svc.PodTemplate, i)
			svc.Pods = append(svc.Pods, containerID)
		}
	}
}

func deployPod(pod types.Pod, index int) string {
	node := nodes[index%len(nodes)]
	data, _ := json.Marshal(pod)
	resp, err := http.Post(node.Address+"/run", "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("Deployment error:", err)
		return ""
	}
	defer resp.Body.Close()
	id, _ := ioutil.ReadAll(resp.Body)
	return string(id)
}

func scaleService(svc *types.Service, newReplicas int) {
	diff := newReplicas - svc.Replicas
	if diff > 0 {
		for i := 0; i < diff; i++ {
			id := deployPod(svc.PodTemplate, i)
			svc.Pods = append(svc.Pods, id)
		}
	} else {
		for i := 0; i < -diff; i++ {
			last := svc.Pods[len(svc.Pods)-1]
			terminatePod(last)
			svc.Pods = svc.Pods[:len(svc.Pods)-1]
		}
	}
	svc.Replicas = newReplicas
}

func terminatePod(containerID string) {
	for _, node := range nodes {
		http.Post(node.Address+"/terminate", "application/json", bytes.NewBufferString(containerID))
	}
}

func StartAutoScaler() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			mu.Lock()
			for _, svc := range services {
				if svc.Traffic > 50 || svc.CPUUsage > 70 {
					fmt.Println("[AutoScaler] Scaling up", svc.Name)
					scaleService(svc, svc.Replicas+1)
				} else if svc.Traffic < 10 && svc.CPUUsage < 20 && svc.Replicas > 1 {
					fmt.Println("[AutoScaler] Scaling down", svc.Name)
					scaleService(svc, svc.Replicas-1)
				}
			}
			mu.Unlock()
		}
	}()
}
