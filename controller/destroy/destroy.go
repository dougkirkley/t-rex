package destroy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dougkirkley/trex/config"
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

// Run
func (ctrl *Controller) Run() error {
	mainFile := filepath.Join(ctrl.cfg.WorkingDirectory(), "main.tf")
	if !fn.Exists(mainFile) {
		return fmt.Errorf("%s does not exist or is not accessible. Perhaps run trex init first", mainFile)
	}
	if err := ctrl.deleteStateBucket(); err != nil {
		return err
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

func (ctrl *Controller) deleteStateBucket() error {
	bucket := fmt.Sprintf("%s-tfstate", ctrl.cfg.APIName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucket)},
	)
	if err != nil {
		return err
	}

	// Create S3 service client
	svc := s3.New(sess)
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
