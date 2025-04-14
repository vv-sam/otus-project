package repository

import (
	"fmt"
	"sync"

	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
)

var (
	m sync.RWMutex

	agents  = make([]*agent.Info, 0)
	configs = make([]*configuration.Minecraft, 0)
	metric  = make([]*metrics.HostMetrics, 0)
	tasks   = make([]*task.Task, 0)

	agentsCount  = len(agents)
	configsCount = len(configs)
	metricCount  = len(metric)
	tasksCount   = len(tasks)
)

func PassStructs(ch chan fmt.Stringer) {
	for s := range ch {
		m.Lock()

		switch v := s.(type) {
		case *agent.Info:
			agents = append(agents, v)
		case *configuration.Minecraft:
			configs = append(configs, v)
		case *metrics.HostMetrics:
			metric = append(metric, v)
		case *task.Task:
			tasks = append(tasks, v)
		}

		m.Unlock()
	}
}

func CheckUpdates() {
	m.RLock()
	defer m.RUnlock()

	fmt.Println("Checking updates...")
	if len(agents) != agentsCount {
		for _, a := range agents[agentsCount:] {
			fmt.Printf("Update: %q\n", a.String())
		}
		agentsCount = len(agents)
	}

	if len(configs) != configsCount {
		for _, c := range configs[configsCount:] {
			fmt.Printf("Update: %q\n", c.String())
		}
		configsCount = len(configs)
	}

	if len(metric) != metricCount {
		for _, m := range metric[metricCount:] {
			fmt.Printf("Update: %q\n", m.String())
		}
		metricCount = len(metric)
	}

	if len(tasks) != tasksCount {
		for _, m := range tasks[tasksCount:] {
			fmt.Printf("Update: %q\n", m.String())
		}
		tasksCount = len(tasks)
	}
	fmt.Println("Checking done")
}
