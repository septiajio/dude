// mini-k8s: node/worker.go
package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	types "dude/types"
)

var containerMap = make(map[string]*types.Pod)
var trafficMap = make(map[string]int)
var cpuMap = make(map[string]float64)
var mu sync.Mutex

func RunServer(port string) {
	http.HandleFunc("/run", handleRun)
	http.HandleFunc("/terminate", handleTerminate)
	http.HandleFunc("/simulate-traffic", handleTraffic)
	http.HandleFunc("/metrics", handleMetrics)

	fmt.Println("Worker running on port", port)
	http.ListenAndServe(":"+port, nil)
}

// todo: this handle not yet deploying/build/start real image or container
func handleRun(w http.ResponseWriter, r *http.Request) {
	var pod types.Pod
	_ = json.NewDecoder(r.Body).Decode(&pod)
	id := fmt.Sprintf("container-%d", rand.Intn(99999))
	containerMap[id] = &pod
	trafficMap[id] = 0
	cpuMap[id] = 0
	go simulateCPU(id)
	fmt.Println("deploy simulated")
	w.Write([]byte(id))
}

func handleTerminate(w http.ResponseWriter, r *http.Request) {
	id, _ := ioutil.ReadAll(r.Body)
	mu.Lock()
	delete(containerMap, string(id))
	delete(trafficMap, string(id))
	delete(cpuMap, string(id))
	mu.Unlock()
}

func handleTraffic(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mu.Lock()
	trafficMap[id] += rand.Intn(10)
	mu.Unlock()
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	totalTraffic := 0
	totalCPU := 0.0
	mu.Lock()
	for id := range trafficMap {
		totalTraffic += trafficMap[id]
		totalCPU += cpuMap[id]
	}
	mu.Unlock()
	resp := map[string]interface{}{
		"traffic": totalTraffic,
		"cpu":     totalCPU,
	}
	json.NewEncoder(w).Encode(resp)
}

func simulateCPU(id string) {
	for {
		time.Sleep(5 * time.Second)
		mu.Lock()
		cpuMap[id] = float64(rand.Intn(100))
		mu.Unlock()
	}
}
