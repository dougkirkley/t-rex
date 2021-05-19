package init

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/dougkirkley/trex/config"
	//"github.com/dougkirkley/trex/controller"
	fn "github.com/dougkirkley/trex/function"
	"github.com/dougkirkley/trex/terraform"
)

// New creates a new init controller
func New(cfg *config.TrexConfig) *Controller {
	return &Controller{
		cfg: cfg,
	}
}

// Controller does the stuff
type Controller struct {
	cfg *config.TrexConfig
}

// Run runs initialization
func (ctrl *Controller) Run() error {
	if err := ctrl.createTrexYaml(); err != nil {
		return err
	}
	if err := ctrl.cfg.MakeVarJSON(); err != nil {
		return err
	}
	if err := ctrl.createMainTF(); err != nil {
		return err
	}
	if !config.Exists("functions") {
		err := os.Mkdir("functions", 0777)
		if err != nil {
			fmt.Print("error on creating functions dir")
			return err
		}
	}
	if err := fn.BuildFuncs(*ctrl.cfg); err != nil {
		return err
	}

	command := terraform.InitCommand(terraform.InitArguments{})
	log.Printf("Running %s", strings.Join(command.Args, " "))
	output, err := command.CombinedOutput()
	fmt.Print(string(output))
	if err != nil {
		return err
	}

	command = terraform.FmtCommand(terraform.FmtArguments{})
	output, _ = command.CombinedOutput()
	fmt.Print(string(output))

	return nil
}

// createMainTF creates the main.tf unless it already exists
func (ctrl *Controller) createMainTF() error {
	mainTF := filepath.Join(ctrl.cfg.WorkingDirectory(), "main.tf")
	if config.Exists(mainTF) {
		return nil
	}

	s, err := terraform.CreateMainFile(*ctrl.cfg)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(mainTF, s, 0755)
}

// createTrexYaml creates the trex.yaml unless it already exists
func (ctrl *Controller) createTrexYaml() error {
	yamlFile := filepath.Join(ctrl.cfg.WorkingDirectory(), "trex.yaml")
	if config.Exists(yamlFile) {
		return nil
	}

	err := ctrl.cfg.GenerateDefaultConfig()
	if err != nil {
		return err
	}
	return nil
}

func (ctrl *Controller) Config() *config.TrexConfig {
	return ctrl.cfg
}

func (ctrl *Controller) createStateBucket() error {
	bucket := fmt.Sprintf("%s-tfstate", ctrl.cfg.APIName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucket)},
	)
	if err != nil {
		return err
	}

	// Create S3 service client
	svc := s3.New(sess)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	// Wait until bucket is created before finishing
	log.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	log.Printf("Bucket %q successfully created\n", bucket)

	return nil
}
