package main

import (
	"fmt"
	"log"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

func main() {

	opts := new(agent.AgentOptions)
	opts.ParseArgs()
	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}
	fs := file.NewFileStorage("mem.txt")
	storage := memory.NewStorage(40, fs, true)
	a := agent.NewAgent(storage, opts)

	fmt.Println("start app")
	a.Run()

}
