package mongo

import (
	"log"
	"os"
	"os/exec"

	"github.com/elos/autonomous"
)

var (
	mongod exec.Cmd
)

type runner struct {
	autonomous.Life
	autonomous.Stopper
	autonomous.Managed

	mongod     *exec.Cmd
	ConfigFile string
}

var Runner = &runner{
	Life:    autonomous.NewLife(),
	Stopper: make(autonomous.Stopper),
	Managed: *new(autonomous.Managed),
}

func (r *runner) Start() {
	if r.ConfigFile != "" {
		r.mongod = exec.Command("mongod", "--config", r.ConfigFile)
	} else {
		r.mongod = exec.Command("mongod")
	}

	r.mongod.Stdout = os.Stdout
	r.mongod.Stderr = os.Stderr

	if err := r.mongod.Start(); err != nil {
		log.Print(err)
	} else {
		log.Print("Mongo successfully started")
	}

	r.Life.Begin()
	<-r.Stopper
	if err := r.mongod.Process.Signal(os.Interrupt); err != nil {
		log.Print(err)
	} else {
		log.Print("Mongo succesfully stopped")
	}
	r.Life.End()
}
