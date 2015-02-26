package mongo

import (
	"log"
	"os"
	"os/exec"
	"sync"

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
	*log.Logger
	sync.Mutex
}

var Runner = &runner{
	Life:    autonomous.NewLife(),
	Stopper: make(autonomous.Stopper),
	Managed: *new(autonomous.Managed),
	Logger:  DefaultLogger,
}

func (r *runner) SetLogger(l *log.Logger) {
	r.Lock()
	defer r.Unlock()

	r.Logger = l
}

func (r *runner) SetConfigFile(s string) {
	r.Lock()
	defer r.Unlock()

	r.ConfigFile = s
}

func (r *runner) Start() error {
	// Lock the runner
	r.Lock()
	defer r.Unlock()

	if r.ConfigFile != "" {
		r.mongod = exec.Command("mongod", "--config", r.ConfigFile)
	} else {
		r.mongod = exec.Command("mongod")
	}

	r.mongod.Stdout = os.Stdout
	r.mongod.Stderr = os.Stderr

	if err := r.mongod.Start(); err != nil {
		r.Print(err)
		return err
	}

	r.Life.Begin()
	r.Print("Mongo successfully started")

	<-r.Stopper

	if err := r.mongod.Process.Signal(os.Interrupt); err != nil {
		r.Print(err)
		return err
	}

	r.Life.End()
	r.Print("Mongo succesfully stopped")
	return nil
}

func testify(r *runner) {
	Runner.SetConfigFile("./test.conf")
	Runner.SetLogger(NullLogger)
	go Runner.Start()
	Runner.WaitStart()
}

func detestify(r *runner) {

	go r.Stop()
	r.WaitStop()

	r.SetConfigFile("")
	r.SetLogger(DefaultLogger)
}
