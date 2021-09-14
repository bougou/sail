package models

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	newexec "github.com/bougou/gopkg/exec"
	"github.com/bougou/sail/pkg/ansible"
)

//go:embed ansible.cfg
var defaultAnsibleCfg string

type RunningZone struct {
	zone     *Zone
	playbook string

	serverComponents []string
	podComponents    []string

	startAtPlay         string
	ansiblePlaybookTags []string

	ansiblePlaybookArgs []string
}

func (rz *RunningZone) WithServerComponents(serverComponents map[string]string) {
	for componentName := range serverComponents {
		rz.serverComponents = append(rz.serverComponents, componentName)
	}
}

func (rz *RunningZone) WithPodComponents(podComponents map[string]string) {
	for componentName := range podComponents {
		rz.podComponents = append(rz.podComponents, componentName)
	}
}

func (rz *RunningZone) WithStartAtPlay(startAtPlay string) {
	rz.startAtPlay = startAtPlay
}

func (rz *RunningZone) WithAnsiblePlaybookTags(ansiblePlaybookTags []string) {
	rz.ansiblePlaybookTags = ansiblePlaybookTags
}

func NewRunningZone(zone *Zone, playbookName string) *RunningZone {
	rz := &RunningZone{
		zone:     zone,
		playbook: playbookName,
	}

	rz.ansiblePlaybookArgs = []string{
		zone.PlaybookFile(playbookName),
		"-i",
		zone.HostsFile,
		"-e",
		"@" + zone.VarsFile,
		"-e",
		"@" + zone.ComputedFile,
		"-e",
		"sail_products_dir=" + zone.sailOption.ProductsDir,
		"-e",
		"sail_packages_dir=" + zone.sailOption.PackagesDir,
		"-e",
		"sail_targets_dir=" + zone.sailOption.TargetsDir,
		"-e",
		"sail_target_dir=" + zone.TargetDir,
		"-e",
		"sail_zone_dir=" + zone.ZoneDir,
		"-e",
		"sail_target_name=" + zone.TargetName,
		"-e",
		"sail_zone_name=" + zone.ZoneName,
	}

	// sudoParams := []string{
	// 	"--become",
	// 	"--become-user=root",
	// }
	// sudoParams = append(sudoParams, "--ask-become-pass")

	return rz
}

func (rz *RunningZone) Run(args []string) error {
	// If not specify any components, it means all components.
	// So run ansible-playbook, then helm.
	if len(rz.serverComponents) == 0 && len(rz.podComponents) == 0 {
		if err := rz.RunAnsiblePlaybook(args); err != nil {
			return err
		}

		if err := rz.RunHelm(args); err != nil {
			return err
		}
	}

	if len(rz.serverComponents) > 0 {
		if err := rz.RunAnsiblePlaybook(args); err != nil {
			return err
		}
	}

	if len(rz.podComponents) > 0 {
		if err := rz.RunHelm(args); err != nil {
			return err
		}
	}

	return nil
}

func (rz *RunningZone) RunAnsiblePlaybook(args []string) error {
	// ansible-playbook tags set by sail commands.
	if len(rz.ansiblePlaybookTags) != 0 {
		rz.ansiblePlaybookArgs = append(rz.ansiblePlaybookArgs, "--tags", strings.Join(rz.ansiblePlaybookTags, ","))
	}

	// parse ansible playbook file to get tags for startAtPlay
	playbookFile := rz.zone.PlaybookFile(rz.playbook)
	playbook, err := ansible.NewPlaybookFromFile(playbookFile)
	if err != nil {
		return err
	}
	if rz.startAtPlay != "" {
		playbookTags := playbook.PlaysTagsStartAt(rz.startAtPlay)
		if len(playbookTags) != 0 {
			rz.ansiblePlaybookArgs = append(rz.ansiblePlaybookArgs, "--tags", strings.Join(playbookTags, ","))
		}
	}

	if len(args) > 0 {
		rz.ansiblePlaybookArgs = append(rz.ansiblePlaybookArgs, args...)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		"ANSIBLE_FORCE_COLOR=true", // this env var will make ansible-playbook always output color
		"ANSIBLE_CONFIG=" + rz.zone.ansibleCfgFile,
	}

	logFileName := "/tmp/sail.log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("can not create log file: %s, exit", logFileName)
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

func (rz *RunningZone) RunHelm(args []string) error {
	switch rz.zone.SailHelmMode {
	case SailHelmModeComponent:
		for componentName, component := range rz.zone.Product.Components {
			if component.Form == "pod" {
				// fmt.Printf("component: (%s) form: (%s)\n", componentName, component.Form)
				helmRelease := fmt.Sprintf("%s-%s", rz.zone.SailProduct, componentName)
				helmChartDir := rz.zone.HelmChartDirOfComponent(componentName)
				k8s := rz.zone.GetK8SForComponent(componentName)
				if err := helmCmd(helmRelease, helmChartDir, k8s); err != nil {
					return err
				}
			}
		}

	case SailHelmModeProduct:
		helmRelease := rz.zone.SailProduct
		helmChartDir := rz.zone.HelmChartDirOfProduct()
		k8s := rz.zone.GetK8SForProduct()
		return helmCmd(helmRelease, helmChartDir, k8s)

	default:
		return fmt.Errorf("not supported helm mode: (%s)", rz.zone.SailHelmMode)
	}

	return nil
}

func helmCmd(release string, chartDir string, k8s *K8S) error {
	args := []string{
		"upgrade",
		release,
		chartDir,
		"--install",
	}

	if k8s != nil {
		if k8s.KubeContext != "" {
			args = append(args, "--kube-context", k8s.KubeContext)
		}
		if k8s.KuebConfig != "" {
			args = append(args, "--kubeconfig", k8s.KuebConfig)
		}
		if k8s.Namespace != "" {
			args = append(args, "--namespace", k8s.Namespace)
		}
	}
	args = append(args, "--debug")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "helm", args...)

	logFileName := "/tmp/sail.log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("can not create log file: %s, exit", logFileName)
	}

	// thus, the cmd's output goes to terminal AND logfile
	cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
	cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
	// ansible-playbook checks if stdin is a tty device `os.isatty(0)`, then set column width accordingly
	// ref: https://github.com/ansible/ansible/blob/2cbfd1e350cbe1ca195d33306b5a9628667ddda8/lib/ansible/utils/display.py#L534
	// here, we specifically set to os.Stdin to simulate tty
	cmd.Stdin = os.Stdin

	cmdWrapper := newexec.NewCmdEnvWrapper(cmd)
	fmt.Println("⛵ " + cmdWrapper.String())
	// cmdWrapper.SetDebug(true)
	return cmdWrapper.Run()

}
