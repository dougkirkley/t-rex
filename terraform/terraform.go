package terraform

import (
	"os/exec"
)

const Terraform = "terraform"

// ApplyArguments is args for apply
type ApplyArguments struct {
	BackupPath    *string           `flag:"-backup-path"`
	Lock          *bool             `flag:"-lock"`
	LockTimeout   *int              `flag:"-lock-timeout"`
	Input         *bool             `flag:"-input"`
	AutoApprove   *bool             `flag:"-auto-approve"`
	NoColor       *bool             `flag:"-no-color"`
	Parallelism   *int              `flag:"-parallelism"`
	Refresh       *bool             `flag:"-refresh"`
	State         *string           `flag:"-state"`
	StateOut      *string           `flag:"-state-out"`
	Targets       []string          `flag:"-target"`
	Variables     map[string]string `flag:"-var"`
	VariableFiles []string          `flag:"-var-file"`
	DirOrPlan     *string           `pos:"0"`
}

// PlanArguments arguments for plan command
type PlanArguments struct {
	Destroy          *bool    `flag:"-destroy"`
	DetailedExitCode *bool    `flag:"-detailed-exitcode"`
	Input            *bool    `flag:"-input"`
	Lock             *bool    `flag:"-lock"`
	LockTimeout      *int     `flag:"-lock-timeout"`
	NoColor          *bool    `flag:"-no-color"`
	Out              *string  `flag:"-out"`
	Parrallism       *int     `flag:"-parallelism"`
	Refresh          *bool    `flag:"-refresh"`
	State            *string  `flag:"-state"`
	Targets          []string `flag:"-target"`
	Variables        []string `flag:"-var"`
	VariableFiles    []string `flag:"-var-file"`
	Dir              *string  `pos:"0"`
}

// FmtArguments arguments for fmt
type FmtArguments struct {
	List  *bool   `flag:"list"`
	Write *bool   `flag:"write"`
	Diff  *bool   `flag:"diff"`
	Check *bool   `flag:"check"`
	Dir   *string `pos:"0"`
}

// InitArguments arguments for init command
type InitArguments struct {
	Dir         *string `pos:"0"`
	Input       *bool   `flag:"-input"`
	Lock        *bool   `flag:"-lock"`
	LockTimeout *int    `flag:"-lock-timeout"`
	NoColor     *bool   `flag:"-no-color"`
	Upgrade     *string `flag:"-upgrade"`
}

// OutputArguments arguments for output command
type OutputArguments struct {
	Name   *string `pos:"0"`
	JSON   *bool   `flag:"-json"`
	State  *string `flag:"-state"`
	Module *string `flag:"-module"`
}

// ImportArguments arguments for import command
type ImportArguments struct {
	Address string `pos:"0"`
	ID      string `pos:"1"`

	Backup        *string           `flag:"-backup"`
	Config        *string           `flag:"-config"`
	Lock          *bool             `flag:"-lock"`
	LockTimeout   *int              `flag:"-lock-timeout"`
	Input         *bool             `flag:"-input"`
	NoColor       *bool             `flag:"-no-color"`
	Parallelism   *int              `flag:"-parallelism"`
	Provider      *string           `flag:"-provider"`
	State         *string           `flag:"-state"`
	StateOut      *string           `flag:"-state-out"`
	Variables     map[string]string `flag:"-var"`
	VariableFiles []string          `flag:"-var-file"`
}

// ApplyCommmand builds a terraform apply command
func ApplyCommmand(arguments ApplyArguments) *exec.Cmd {
	args := append([]string{"apply"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// DestroyCommand creates a terraform destroy command
func DestroyCommand(arguments ApplyArguments) *exec.Cmd {
	args := append([]string{"destroy"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// PlanCommand creates a plan command
func PlanCommand(arguments PlanArguments) *exec.Cmd {
	args := append([]string{"plan"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// FmtCommand creates an fmt command
func FmtCommand(arguments FmtArguments) *exec.Cmd {
	args := append([]string{"fmt"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// InitCommand creates an init command
func InitCommand(arguments InitArguments) *exec.Cmd {
	args := append([]string{"init"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// OutputCommand creates an output command
func OutputCommand(arguments OutputArguments) *exec.Cmd {
	args := append([]string{"output"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}

// ImportCommand creates an import command
func ImportCommand(arguments ImportArguments) *exec.Cmd {
	args := append([]string{"import"}, BuildArgs(arguments)...)
	command := exec.Command(Terraform, args...)
	return command
}