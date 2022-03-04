package majorupgradechecklist

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePostChecklistFile(t *testing.T) {
	// remove json file
	removeFile()
	// Create json file
	cl, _ := NewPostChecklistManager("3")

	err := cl.CreatePostChecklistFile()
	// // return nil if json file is created successfully
	assert.NoError(t, err)
}

func TestReadPostChecklistByIdSuccess(t *testing.T) {
	cl, _ := NewPostChecklistManager("3")
	IsExecuted, err := cl.ReadPostChecklistById("pg_migrate")
	// return nil if checklist by id found
	assert.NoError(t, err)
	// check is executed is true or false
	assert.Equal(t, IsExecuted, false)
}

func TestReadPostChecklistSuccess(t *testing.T) {
	cl, _ := NewPostChecklistManager("3")
	result, err := cl.ReadPendingPostChecklistFile()
	assert.NoError(t, err)
	//get json data as result
	assert.NotEqual(t, result, []string{})
}

func removeFile() {
	IsExist := false
	_, err := os.Stat("/hab/svc/deployment-service/var/upgrade_metadata.json")
	if err == nil {
		IsExist = true
	}
	if os.IsNotExist(err) {
		if _, err := os.Stat("/hab/svc/deployment-service/var"); err == nil {
			c := exec.Command("mkdir", "-p", "/hab/svc/deployment-service/var")
			c.Stdin = os.Stdin
			c.Stdout = io.MultiWriter(os.Stdout)
			c.Stderr = io.MultiWriter(os.Stderr)
			err := c.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
		IsExist = false
	}
	if IsExist {
		e := os.Remove("/hab/svc/deployment-service/var/upgrade_metadata.json")
		if e != nil {
			log.Fatal(e)
		}
	}
}
