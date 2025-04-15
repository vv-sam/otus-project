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
	am sync.Mutex
	cm sync.Mutex
	mm sync.Mutex
	tm sync.Mutex

	agents  = make([]*agent.Info, 0)
	configs = make([]*configuration.Minecraft, 0)
	metric  = make([]*metrics.HostMetrics, 0)
	tasks   = make([]*task.Task, 0)

	agentsCount  = len(agents)
	configsCount = len(configs)
	metricCount  = len(metric)
	tasksCount   = len(tasks)
)

func PassStruct(s fmt.Stringer) {
	switch v := s.(type) {
	case *agent.Info:
		am.Lock()
		agents = append(agents, v)
		am.Unlock()
	case *configuration.Minecraft:
		cm.Lock()
		configs = append(configs, v)
		cm.Unlock()
	case *metrics.HostMetrics:
		mm.Lock()
		metric = append(metric, v)
		mm.Unlock()
	case *task.Task:
		tm.Lock()
		tasks = append(tasks, v)
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
