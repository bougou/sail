package models

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	newexec "github.com/bougou/gopkg/exec"
)

//go:embed ansible.cfg
var defaultAnsibleCfg string

type RunningZone struct {
	zone     *Zone
	playbook string

	ansiblePlaybookArgs []string
}

func NewRunningZone(zone *Zone, playbook string) *RunningZone {
	rz := &RunningZone{
		zone:     zone,
		playbook: playbook,
	}

	rz.ansiblePlaybookArgs = []string{
		zone.PlaybookFile(playbook),
		"-i",
		zone.HostsFile,
		"-e",
		"@" + zone.VarsFile,
		"-e",
		"@" + zone.ComputedFile,
		"-e",
		"products_dir=" + zone.sailOption.ProductsDir,
		"-e",
		"packages_dir=" + zone.sailOption.PackagesDir,
		"-e",
		"targets_dir=" + zone.sailOption.TargetsDir,
		"-e",
		"target_dir=" + zone.TargetDir,
		"-e",
		"zone_dir=" + zone.ZoneDir,
		"-e",
		"target_name=" + zone.TargetName,
		"-e",
		"zone_name=" + zone.ZoneName,
	}

	// sudoParams := []string{
	// 	"--become",
	// 	"--become-user=root",
	// }
	// sudoParams = append(sudoParams, "--ask-become-pass")

	return rz
}
func (rz *RunningZone) Run(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(args) > 0 {
		rz.ansiblePlaybookArgs = append(rz.ansiblePlaybookArgs, args...)
	}

	if _, err := os.Stat(rz.zone.ansibleCfgFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("not found ansible.cfg, generate one")
			if err := os.WriteFile(rz.zone.ansibleCfgFile, []byte(defaultAnsibleCfg), 0644); err != nil {
				fmt.Println("write ansible.cfg file failed", err)
			}
		}
	}

	cmd := exec.CommandContext(ctx, "ansible-playbook", rz.ansiblePlaybookArgs...)
	env := []string{
		"ANSIBLE_FORCE_COLOR=true", // this env vars will make ansible-playbook always output color
		"ANSIBLE_CONFIG=" + rz.zone.ansibleCfgFile,
	}

	logFileName := "/tmp/log.txt"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("can not create log file: %s, exit", logFileName))
	}

	// thus, the cmd's output goes to terminal AND logfile
	cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
	cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
	// ansible-playbook checks if stdin is a tty device `os.isatty(0)`, then set column width accordingly
	// ref: https://github.com/ansible/ansible/blob/2cbfd1e350cbe1ca195d33306b5a9628667ddda8/lib/ansible/utils/display.py#L534
	// here, we specifically set to os.Stdin to simulate tty
	cmd.Stdin = os.Stdin

	cmdWrapper := newexec.NewCmdEnvWrapper(cmd, env...)
	fmt.Println("⛵ " + cmdWrapper.String())
	// cmdWrapper.SetDebug(true)
	return cmdWrapper.Run()

}

// Todo
func (rz *RunningZone) helmInstall() {
	helmArgs := []string{
		"install",
		rz.zone.ProductName,
		"the-helm-dir",
		"--kubeconfig",
		"/path/to/kubeconfig",
		"--values",
		"/path/to/vars.yml",
		"--debug",
	}
	fmt.Println(helmArgs)
}

// Todo
func (rz *RunningZone) helmUpgrade() {
	helmArgs := []string{
		"upgrade",
		rz.zone.ProductName,
		"the-helm-dir",
		"--kubeconfig",
		"/path/to/kubeconfig",
		"--values",
		"/path/to/vars.yml",
		"--force",
		"--debug",
	}
	fmt.Println(helmArgs)
}
