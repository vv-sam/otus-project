package repository

import (
	"fmt"

	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
)

var (
	agents  = make([]*agent.Info, 0)
	configs = make([]*configuration.Minecraft, 0)
	metric  = make([]*metrics.HostMetrics, 0)
	tasks   = make([]*task.Task, 0)
)

func PassStructs(structs []fmt.Stringer) {
	for _, s := range structs {
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
	}
}

func PrintValues() {
	for _, s := range agents {
		fmt.Println(s.String())
	}

	for _, s := range configs {
		fmt.Println(s.String())
	}

	for _, s := range metric {
		fmt.Println(s.String())
	}

	for _, s := range tasks {
		fmt.Println(s.String())
	}

	fmt.Println("\n\n\n")
}
