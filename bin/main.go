package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/housepower/clickhouse_sinker/creator"
	"github.com/housepower/clickhouse_sinker/task"
	// "github.com/housepower/clickhouse_sinker/util"

	_ "github.com/kshvakov/clickhouse"
	"github.com/wswz/go_commons/app"
)

var (
	config     string
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func init() {
	flag.StringVar(&config, "conf", "", "config dir")

	flag.Parse()
}

func main() {

	var cfg creator.Config
	var runner *Sinker

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app.Run("clickhouse_sinker", func() error {
		cfg = *creator.InitConfig(config)
		runner = NewSinker(cfg)
		return runner.Init()
	}, func() error {
		runner.Run()
		return nil
	}, func() error {
		// util.MemStat()
		runner.Close()
		return nil
	})
}

type Sinker struct {
	tasks   []*task.TaskService
	config  creator.Config
	stopped chan struct{}
}

func NewSinker(config creator.Config) *Sinker {
	s := &Sinker{config: config, stopped: make(chan struct{})}
	return s
}

func (s *Sinker) Init() error {
	s.tasks = s.config.GenTasks()
	for _, t := range s.tasks {
		if err := t.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sinker) Run() {
	for i, _ := range s.tasks {
		go s.tasks[i].Run()
	}
	<-s.stopped
}

func (s *Sinker) Close() {
	for i, _ := range s.tasks {
		s.tasks[i].Stop()
	}
	close(s.stopped)
}
