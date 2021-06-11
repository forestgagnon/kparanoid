package cfg

import "path"

type Config struct {
	ConfigDirOnHost string
	Flags           struct {
		KubectlVersion string
		ConfigDir      string
		AddClusterGKE  struct {
			GcloudClusterCredentialsCommand string
		}
	}
	OnlyBuildRuntimeIfNeeded bool
	SquelchDockerBuildOutput bool
}

var Cfg Config

func (c *Config) ClustersPath() string {
	return path.Join(c.Flags.ConfigDir, "/clusters")
}

func (c *Config) ClustersPathHost() string {
	return path.Join(c.ConfigDirOnHost, "/clusters")
}

func (c *Config) ClusterPath(clusterName string) string {
	return path.Join(c.ClustersPath(), clusterName)
}

func (c *Config) ClusterPathHost(clusterName string) string {
	return path.Join(c.ClustersPathHost(), clusterName)
}

func (c *Config) ContainmentDockerfilesPath() string {
	return path.Join(c.EnvPath(), "/containment-dockerfiles")
}

func (c *Config) ContainmentDockerfile(name string) string {
	return path.Join(c.ContainmentDockerfilesPath(), name)
}

func (c *Config) EnvPath() string {
	return path.Join(c.Flags.ConfigDir, "/env")
}
