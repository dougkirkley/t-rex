package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"gitlab.com/dougkirkley/trex/config"
	//"github.com/dougkirkley/trex/controller"
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

// Run
func (ctrl *Controller) Run() error {
	mainFile := filepath.Join(ctrl.cfg.WorkingDirectory(), "/main.tf")
	if !config.Exists(mainFile) {
		return fmt.Errorf("%s does not exist or is not accessible. Perhaps run trex init first", mainFile)
	}
	command := terraform.ApplyCommmand(terraform.ApplyArguments{
		AutoApprove: aws.Bool(true),
		NoColor:     aws.Bool(true),
	})

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func (ctrl *Controller) Config() *config.TrexConfig {
	return ctrl.cfg
}
