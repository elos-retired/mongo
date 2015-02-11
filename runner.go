package mongo

import (
	// "github.com/elos/server/util"
	"os"
	"os/exec"
)

var (
	mongod exec.Cmd
)

func StartDatabaseServer() error {
	mongod := exec.Command("mongod", "--config", "./mongo.conf")
	mongod.Stdout = os.Stdout
	mongod.Stderr = os.Stderr

	if err := mongod.Start(); err != nil {
		return err
	}
	// Log("Mongo succesfully started") causes runtime panic?
	return nil
}

func StopDatabaseServer(sig os.Signal) error {
	if err := mongod.Process.Signal(sig); err != nil {
		return err
	}

	// Log("Mongo succesfully stopped") causes runtime panic?
	return nil
}
