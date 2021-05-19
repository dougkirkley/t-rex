package lambda

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"

	"github.com/dougkirkley/trex/config"
	"github.com/dougkirkley/trex/terraform"
)

const GO = "go"

func WorkingDir() string {
    wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
	}
	return wd
}

func notInConfig(configFile config.TrexConfig, file string) bool {
    for _, v := range configFile.Functions {
		if v.LambdaName == file {
			return true
		}
	}
	return false
}

func makeFuncDir(path string) error {
	wd := WorkingDir()
	filePath := filepath.Join(wd, fmt.Sprintf("/%s", path))
	err := os.Mkdir(filePath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func ChangeDir(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return nil
}

func GoBuild(path string) error {
	command := exec.Command(GO, "build", "-o", "bin/main")
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func BuildFuncs(cfg config.TrexConfig) error {
	rootPath := WorkingDir()
	for _, function := range cfg.Functions {
		ChangeDir(rootPath)
		var path = fmt.Sprintf("functions/%s", function.LambdaName)
		if !config.Exists(path) {
			err := makeFuncDir(path)
			if err != nil {
				log.Print("error dir already exists")
				return err
			}
			_ = makeFuncDir(path + "/dist")
			terraform.CreateGoFiles(path)
		}
		ChangeDir(rootPath + "/" + path)
		err := GoBuild(path)
		if err != nil {
			return err
		}
	}
	ChangeDir(rootPath)
	return nil
}

func RemoveOldFuncs(configFile config.TrexConfig) error {
	files, err := ioutil.ReadDir("functions")
	if err != nil {
		return err
	}
	for _, file := range files {
		if notInConfig(configFile, file.Name()) {
			os.Remove(file.Name())
		}
	}
	
	return nil
}
