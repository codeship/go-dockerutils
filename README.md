[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/go-dockerutils)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/go-dockerutils/blob/master/LICENSE)

## Installation
```bash
go get -u github.com/peter-edge/go-dockerutils
```

## Import
```go
import (
    "github.com/peter-edge/go-dockerutils"
)
```


#### func  DockerBinds

```go
func DockerBinds(hostVolumeToVolume map[string]string) []string
```

#### func  DockerExposedPorts

```go
func DockerExposedPorts(expose []uint16) map[docker.Port]struct{}
```

#### func  DockerVolumes

```go
func DockerVolumes(volumes []string) map[string]struct{}
```

#### func  DockerVolumesFrom

```go
func DockerVolumesFrom(volumesFrom []string) string
```

#### func  IsDockerTls

```go
func IsDockerTls(dockerTlsVerify string) bool
```

#### func  NewDockerClient

```go
func NewDockerClient(dockerHost string, dockerTlsVerify bool, dockerCertPath string, apiVersion string) (*docker.Client, error)
```

#### func  NewDockerClientFromEnv

```go
func NewDockerClientFromEnv(apiVersion string) (*docker.Client, error)
```

#### type DockerEnvironment

```go
type DockerEnvironment struct {
	DockerHost      string
	DockerTlsVerify string
	DockerCertPath  string
}
```


#### func  GetDockerEnvironment

```go
func GetDockerEnvironment() (*DockerEnvironment, error)
```
TODO(pedge): we are assuming the DOCKER_CERT_PATH is a directory within the host
regardless of if we are running in a docker container or not, this should
actually still work but this is not great

#### func (*DockerEnvironment) EnvStrings

```go
func (this *DockerEnvironment) EnvStrings() []string
```

#### func (*DockerEnvironment) HostVolumeToVolume

```go
func (this *DockerEnvironment) HostVolumeToVolume() map[string]string
```

#### func (*DockerEnvironment) IsDockerTls

```go
func (this *DockerEnvironment) IsDockerTls() bool
```
