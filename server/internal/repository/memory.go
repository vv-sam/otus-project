package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
)

var (
	agentPath  string
	configPath string
	metricPath string
	taskPath   string

	am sync.Mutex
	cm sync.Mutex
	mm sync.Mutex
	tm sync.Mutex

	agents  = make([]*agent.Info, 0)
	configs = make([]*configuration.Minecraft, 0)
	metric  = make([]*metrics.HostMetrics, 0)
	tasks   = make([]*task.Task, 0)

	agentsCount  int
	configsCount int
	metricCount  int
	tasksCount   int
)

func Initialize(dir string) error {
	if dir == "" {
		return fmt.Errorf("directory has to be set")
	}

	agentPath = path.Join(dir, "agents.json")
	configPath = path.Join(dir, "configs.json")
	metricPath = path.Join(dir, "metrics.json")
	taskPath = path.Join(dir, "tasks.json")

	if err := loadAgents(); err != nil {
		return err
	}

	if err := loadConfigs(); err != nil {
		return err
	}

	if err := loadMetrics(); err != nil {
		return err
	}

	if err := loadTasks(); err != nil {
		return err
	}

	return nil
}

func PassStruct(s fmt.Stringer) {
	switch v := s.(type) {
	case *agent.Info:
		am.Lock()
		agents = append(agents, v)
		updateAgents()
		am.Unlock()
	case *configuration.Minecraft:
		cm.Lock()
		configs = append(configs, v)
		updateConfigs()
		cm.Unlock()
	case *metrics.HostMetrics:
		mm.Lock()
		metric = append(metric, v)
		updateMetrics()
		mm.Unlock()
	case *task.Task:
		tm.Lock()
		tasks = append(tasks, v)
		updateTasks()
		tm.Unlock()
	}
}

func CheckUpdates() {
	fmt.Println("Checking updates...")

	am.Lock()
	if len(agents) != agentsCount {
		for _, a := range agents[agentsCount:] {
			fmt.Printf("Update: %q\n", a.String())
		}
		agentsCount = len(agents)
	}
	am.Unlock()

	cm.Lock()
	if len(configs) != configsCount {
		for _, c := range configs[configsCount:] {
			fmt.Printf("Update: %q\n", c.String())
		}
		configsCount = len(configs)
	}
	cm.Unlock()

	mm.Lock()
	if len(metric) != metricCount {
		for _, m := range metric[metricCount:] {
			fmt.Printf("Update: %q\n", m.String())
		}
		metricCount = len(metric)
	}
	mm.Unlock()

	tm.Lock()
	if len(tasks) != tasksCount {
		for _, m := range tasks[tasksCount:] {
			fmt.Printf("Update: %q\n", m.String())
		}
		tasksCount = len(tasks)
	}
	tm.Unlock()

	fmt.Println("Checking done")
}

func loadRawData(p string) ([]byte, error) {
	file, err := os.OpenFile(p, os.O_RDONLY, 0400)
	defer file.Close()

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return []byte{}, err
	}

	return io.ReadAll(file)
}

func saveRawData(p string, data []byte) error {
	file, err := os.OpenFile(p, os.O_WRONLY, 0200)
	defer file.Close()

	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		file, err = os.Create(p)
	}

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func loadAgents() error {
	data, err := loadRawData(agentPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &agents)
	agentsCount = len(agents)
	return err
}

func loadConfigs() error {
	data, err := loadRawData(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &configs)
	configsCount = len(configs)
	return err
}

func loadMetrics() error {
	data, err := loadRawData(metricPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &metric)
	metricCount = len(metric)
	return err
}

func loadTasks() error {
	data, err := loadRawData(taskPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &taskPath)
	tasksCount = len(tasks)
	return err
}

func updateAgents() error {
	data, err := json.Marshal(agents)
	if err != nil {
		return err
	}
	return saveRawData(agentPath, data)
}

func updateConfigs() error {
	data, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	return saveRawData(configPath, data)
}

func updateMetrics() error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	return saveRawData(metricPath, data)
}

func updateTasks() error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return saveRawData(taskPath, data)
}
