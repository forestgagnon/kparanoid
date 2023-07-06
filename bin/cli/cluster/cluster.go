package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/forestgagnon/kparanoid/bin/cli/cfg"
	"github.com/forestgagnon/kparanoid/bin/cli/util"
	"gopkg.in/yaml.v2"
)

const (
	ClusterVariantGKE     = "gke"
	ClusterVariantVanilla = "vanilla"
)

const (
	RuntimeImageTagGKE     = "forestgagnon/kparanoid/local/gke-runtime:1"
	RuntimeImageTagVanilla = "forestgagnon/kparanoid/local/vanilla-runtime:1"
)

func ClusterRuntimeImageTag(c *ClusterConf) string {
	if c == nil {
		return RuntimeImageTagVanilla
	}
	switch c.Variant {
	case ClusterVariantGKE:
		return RuntimeImageTagGKE + "-kubectl-" + GetAptKubectlVersion(c)
	case ClusterVariantVanilla:
		return RuntimeImageTagVanilla + "-kubectl-" + GetAptKubectlVersion(c)
	}

	return ""
}

func clusterConfigDirSetup(name string, clusterConf *ClusterConf) error {
	if name == "" {
		return errors.New("cluster name cannot be blank")
	}

	clusterCfgPath := cfg.Cfg.ClusterPath(name)

	if err := os.MkdirAll(clusterCfgPath, 0777); err != nil {
		return err
	}

	j, err := json.MarshalIndent(clusterConf, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(clusterCfgPath, "cluster.json"), j, 0666)
	if err != nil {
		return err
	}

	if err := clusterConf.EnsureRequiredContainerEnvFiles(); err != nil {
		return err
	}

	return nil
}

func (c *ClusterConf) EnsureRequiredContainerEnvFiles() error {
	clusterCfgPath := cfg.Cfg.ClusterPath(c.Name)

	if err := os.MkdirAll(path.Join(clusterCfgPath, "/container-env"), 0777); err != nil {
		return err
	}

	if err := util.EnsureFile(path.Join(clusterCfgPath, "/container-env/.bash_history")); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(clusterCfgPath, "/container-env/.kube"), 0777); err != nil {
		return err
	}

	if err := util.EnsureFile(path.Join(clusterCfgPath, "/container-env/.kube/config")); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(clusterCfgPath, "/container-env/.config"), 0777); err != nil {
		return err
	}

	return nil
}

func AddGKECluster(name string) error {
	clusterConf := &ClusterConf{
		Variant:        ClusterVariantGKE,
		KubectlVersion: "1.27.3-00",
		Name:           name,
	}

	if err := clusterConfigDirSetup(name, clusterConf); err != nil {
		return err
	}

	if err := BuildGKERuntime(clusterConf); err != nil {
		return err
	}

	setupCmd := fmt.Sprintf(
		"gcloud auth login && %s",
		cfg.Cfg.Flags.AddClusterGKE.GcloudClusterCredentialsCommand,
	)
	if err := RunContainerInteractive(clusterConf, "-c", setupCmd); err != nil {
		return err
	}

	return nil
}

func AddVanillaCluster(name string) error {
	clusterConf := &ClusterConf{
		Variant:        ClusterVariantVanilla,
		KubectlVersion: "1.27.3-00",
		Name:           name,
	}

	if err := clusterConfigDirSetup(name, clusterConf); err != nil {
		return err
	}

	if err := BuildVanillaRuntime(clusterConf); err != nil {
		return err
	}

	msg := fmt.Sprintf("\n\nYou must manually place the kubeconfig file at exactly '%s'\n",
		path.Join(cfg.Cfg.ClusterPathHost(clusterConf.Name), "/container-env/.kube/config"),
	)
	if _, err := os.Stdout.WriteString(msg); err != nil {
		return err
	}

	return nil
}

func BuildGKERuntime(c *ClusterConf) error {
	if !cfg.Cfg.SquelchDockerBuildOutput {
		if _, err := os.Stderr.WriteString("Building runtime...\n"); err != nil {
			return err
		}
	}
	kubectlVersion := GetAptKubectlVersion(c)
	dockerBuildArgs := []string{
		"build", cfg.Cfg.EnvPath(),
		"-t", ClusterRuntimeImageTag(c),
		"--file", cfg.Cfg.ContainmentDockerfile("debian-common.Dockerfile"),
		"--build-arg", "DEBIAN_BASE_IMAGE=google/cloud-sdk:slim",
		"--build-arg", fmt.Sprintf("KUBECTL_VERSION=%s", kubectlVersion),
	}
	cmd := exec.Command("docker", dockerBuildArgs...)
	if !cfg.Cfg.SquelchDockerBuildOutput {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build GKE runtime: %w", err)
	}

	return nil
}

func BuildVanillaRuntime(c *ClusterConf) error {
	if !cfg.Cfg.SquelchDockerBuildOutput {
		if _, err := os.Stderr.WriteString("Building runtime...\n"); err != nil {
			return err
		}
	}
	kubectlVersion := GetAptKubectlVersion(c)
	dockerBuildArgs := []string{
		"build", cfg.Cfg.EnvPath(),
		"-t", ClusterRuntimeImageTag(c),
		"--file", cfg.Cfg.ContainmentDockerfile("debian-common.Dockerfile"),
		"--build-arg", fmt.Sprintf("KUBECTL_VERSION=%s", kubectlVersion),
	}
	cmd := exec.Command("docker", dockerBuildArgs...)
	if !cfg.Cfg.SquelchDockerBuildOutput {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build GKE runtime: %w", err)
	}

	return nil
}

func RunContainerInteractive(c *ClusterConf, args ...string) error {
	if err := c.EnsureRequiredContainerEnvFiles(); err != nil {
		return err
	}

	cmdArgs := []string{
		"run", "-it", "--rm",
		"--env", "KPARANOID_CLUSTER_NAME=" + c.Name,
		"--volume", fmt.Sprintf("%s:/root/.config/gcloud",
			path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.config/gcloud"),
		),
		"--volume", fmt.Sprintf("%s:/root/.kube",
			path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.kube"),
		),
		"--volume", fmt.Sprintf("%s:/root/.bash_history",
			path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.bash_history"), // TODO: make sure history works
		),
		ClusterRuntimeImageTag(c),
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func RunContainerExec(c *ClusterConf, execArgs []string) error {
	if err := c.EnsureRequiredContainerEnvFiles(); err != nil {
		return err
	}

	dockerCmdArgs := []string{
		"run", "--rm",
		"--entrypoint", execArgs[0],
		"--env", "KPARANOID_CLUSTER_NAME=" + c.Name,
		"--volume", fmt.Sprintf("%s:/root/.config/gcloud",
			path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.config/gcloud"),
		),
		"--volume", fmt.Sprintf("%s:/root/.kube",
			path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.kube"),
		),
		ClusterRuntimeImageTag(c),
	}
	dockerCmdArgs = append(dockerCmdArgs, execArgs[1:]...)
	cmd := exec.Command("docker", dockerCmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

type ClusterConf struct {
	Variant        string `json:"variant"`
	KubectlVersion string `json:"kubectlVersion"`
	Name           string `json:"-"`
}

func OpenClusterInteractiveSession(clusterName string) error {
	c, err := GetClusterConf(clusterName)
	if err != nil {
		return err
	}

	if err := PrepareRuntimeForCluster(c); err != nil {
		return err
	}
	return RunContainerInteractive(c)
}

func ExecClusterNonInteractive(clusterName string, execArgs []string) error {
	c, err := GetClusterConf(clusterName)
	if err != nil {
		return err
	}

	if err := PrepareRuntimeForCluster(c); err != nil {
		return err
	}
	return RunContainerExec(c, execArgs)
}

func PrepareRuntimeForCluster(c *ClusterConf) error {
	if err := ErrorIfMultipleKubectlContextsPossible(c); err != nil {
		return err
	}

	if cfg.Cfg.OnlyBuildRuntimeIfNeeded {
		err := exec.Command("docker", "image", "inspect", ClusterRuntimeImageTag(c)).Run()
		if err == nil {
			// Image exists
			return nil
		}
	}

	switch c.Variant {
	case ClusterVariantGKE:
		return BuildGKERuntime(c)
	case ClusterVariantVanilla:
		return BuildVanillaRuntime(c)
	default:
		return fmt.Errorf("unexpected cluster variant '%s' for cluster '%s'",
			c.Variant, c.Name,
		)
	}
}

func GetClusterConf(clusterName string) (*ClusterConf, error) {
	fileRaw, err := ioutil.ReadFile(path.Join(cfg.Cfg.ClusterPath(clusterName), "cluster.json"))
	if err != nil {
		return nil, err
	}
	conf := &ClusterConf{
		Name: clusterName,
	}
	if err := json.Unmarshal(fileRaw, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func PrintAvailableKubectlVersionsApt() error {
	if err := BuildVanillaRuntime(nil); err != nil {
		return err
	}

	cmd := exec.Command("docker", "run", RuntimeImageTagVanilla, "-c", "apt-get update > /dev/null && apt-cache madison kubectl")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

const defaultAptKubectlVersion = "1.27.3-00"

func GetAptKubectlVersion(c *ClusterConf) string {
	// Flags take precedence over the config file
	if v := cfg.Cfg.Flags.KubectlVersion; v != "" {
		return v
	}

	if c != nil && c.KubectlVersion != "" {
		return c.KubectlVersion
	}

	return defaultAptKubectlVersion
}

func ErrorIfMultipleKubectlContextsPossible(c *ClusterConf) error {
	type kubectlConfig struct {
		Contexts []interface{} `yaml:"contexts"`
	}

	kubeconfigPath := path.Join(cfg.Cfg.ClusterPath(c.Name), "/container-env/.kube/config")
	kubeconfigPathOnHost := path.Join(cfg.Cfg.ClusterPathHost(c.Name), "/container-env/.kube/config")

	b, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return err
	}
	actual := &kubectlConfig{}
	if err := yaml.Unmarshal(b, actual); err != nil {
		return err
	}

	if len(actual.Contexts) > 1 {
		return fmt.Errorf(
			"UNSAFE!!! kubeconfig for cluster '%s' contains multiple contexts. That is a recipe for disaster. Edit the kubeconfig at '%s' so there is only one context.",
			c.Name, kubeconfigPathOnHost,
		)
	}
	return nil
}

func ListClusters() ([]string, error) {
	items, err := os.ReadDir(cfg.Cfg.ClustersPath())
	if err != nil {
		return nil, err
	}

	var list []string
	for _, i := range items {
		if i.IsDir() {
			list = append(list, i.Name())
		}
	}
	return list, nil
}

func RemoveClusterConfig(clusterName string) error {
	return os.RemoveAll(cfg.Cfg.ClusterPath(clusterName))
}
