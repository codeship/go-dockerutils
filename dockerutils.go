package dockerutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codeship/go-dockerclient"
)

type DockerEnvironment struct {
	DockerHost      string
	DockerTlsVerify string
	DockerCertPath  string
}

func (d *DockerEnvironment) IsDockerTls() bool {
	return IsDockerTls(d.DockerTlsVerify)
}

func (d *DockerEnvironment) EnvStrings() []string {
	envStrings := []string{
		fmt.Sprintf("DOCKER_HOST=%s", d.DockerHost),
	}
	if d.IsDockerTls() {
		envStrings = append(
			envStrings,
			fmt.Sprintf("DOCKER_TLS_VERIFY=%s", d.DockerTlsVerify),
			fmt.Sprintf("DOCKER_CERT_PATH=%s", d.DockerCertPath),
		)
	}
	return envStrings
}

func (d *DockerEnvironment) HostVolumeToVolume() map[string]string {
	hostVolumeToVolume := make(map[string]string)
	if d.IsDockerTls() && d.DockerCertPath != "" {
		hostVolumeToVolume[d.DockerCertPath] = d.DockerCertPath
	}
	unixPrefix := "unix://"
	if strings.HasPrefix(d.DockerHost, unixPrefix) {
		host := strings.TrimPrefix(d.DockerHost, unixPrefix)
		hostVolumeToVolume[host] = host
	}
	if len(hostVolumeToVolume) == 0 {
		return nil
	}
	return hostVolumeToVolume
}

// TODO(pedge): we are assuming the DOCKER_CERT_PATH is a directory within the host
// regardless of if we are running in a docker container or not, d should actually
// still work but this is not great
func GetDockerEnvironment() (*DockerEnvironment, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}
	dockerTlsVerify := os.Getenv("DOCKER_TLS_VERIFY")
	var dockerCertPath string
	var err error
	if IsDockerTls(dockerTlsVerify) {
		dockerCertPath = os.Getenv("DOCKER_CERT_PATH")
		if dockerCertPath == "" {
			home := os.Getenv("HOME")
			if home == "" {
				return nil, errors.New("dockerutils: environment variable HOME must be set if DOCKER_CERT_PATH is not set")
			}
			dockerCertPath = filepath.Join(home, ".docker")
			dockerCertPath, err = filepath.Abs(dockerCertPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return &DockerEnvironment{
		DockerHost:      dockerHost,
		DockerTlsVerify: dockerTlsVerify,
		DockerCertPath:  dockerCertPath,
	}, nil
}

func NewDockerClientFromEnv(apiVersion string) (*docker.Client, error) {
	dockerEnvironment, err := GetDockerEnvironment()
	if err != nil {
		return nil, err
	}
	return NewDockerClient(dockerEnvironment, apiVersion)
}

func NewDockerClient(dockerEnvironment *DockerEnvironment, apiVersion string) (*docker.Client, error) {
	dockerHost := dockerEnvironment.DockerHost
	if dockerEnvironment.IsDockerTls() {
		parts := strings.Split(dockerHost, "://")
		if len(parts) != 2 {
			return nil, fmt.Errorf("dockerutils: could not split %s into two parts by ://", dockerHost)
		}
		if strings.HasPrefix(parts[0], "http") {
			dockerHost = fmt.Sprintf("tcp://%s", parts[1])
		}
		dockerCertPath := dockerEnvironment.DockerCertPath
		cert := filepath.Join(dockerCertPath, "cert.pem")
		key := filepath.Join(dockerCertPath, "key.pem")
		ca := filepath.Join(dockerCertPath, "ca.pem")
		if err := checkFileExists(dockerCertPath); err != nil {
			return nil, err
		}
		if err := checkFileExists(cert); err != nil {
			return nil, err
		}
		if err := checkFileExists(key); err != nil {
			return nil, err
		}
		if err := checkFileExists(ca); err != nil {
			return nil, err
		}
		return docker.NewVersionedTLSClient(dockerHost, cert, key, ca, apiVersion)
	}
	return docker.NewVersionedClient(dockerHost, apiVersion)
}

func IsDockerTls(dockerTlsVerify string) bool {
	return dockerTlsVerify != ""
}

func DockerPorts(expose []uint16, ports []string) (map[docker.Port]struct{}, map[docker.Port][]docker.PortBinding, error) {
	if (expose == nil || len(expose) == 0) && (ports == nil || len(ports) == 0) {
		return nil, nil, nil
	}
	m := make(map[docker.Port]struct{})
	n := make(map[docker.Port][]docker.PortBinding)
	for _, port := range expose {
		m[docker.Port(fmt.Sprintf("%v/tcp", port))] = emptyStruct()
	}
	for _, port := range ports {
		var err error
		var hostPort string
		split := strings.Split(port, ":")
		containerPort := split[len(split)-1]

		if len(split) == 2 && split[0] != "" {
			hostPort = split[0]
			_, err = strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				return nil, nil, err
			}
		}

		_, err = strconv.ParseInt(containerPort, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		dockerPort := docker.Port(fmt.Sprintf("%s/tcp", containerPort))
		m[dockerPort] = emptyStruct()
		n[dockerPort] = []docker.PortBinding{
			docker.PortBinding{
				HostPort: fmt.Sprintf("%s", hostPort),
			},
		}
	}
	return m, n, nil
}

func DockerVolumes(volumes []string) map[string]struct{} {
	if volumes == nil || len(volumes) == 0 {
		return nil
	}
	m := make(map[string]struct{})
	for _, volume := range volumes {
		m[volume] = emptyStruct()
	}
	return m
}

func DockerVolumesFrom(volumesFrom []string) string {
	return strings.Join(volumesFrom, ",")
}

func DockerBinds(hostVolumeToVolume map[string]string) []string {
	if hostVolumeToVolume == nil || len(hostVolumeToVolume) == 0 {
		return nil
	}
	binds := make([]string, len(hostVolumeToVolume))
	i := 0
	for hostVolume, volume := range hostVolumeToVolume {
		binds[i] = fmt.Sprintf("%s:%s:rw", hostVolume, volume)
		i++
	}
	return binds
}

// ***** PRIVATE *****

func checkFileExists(path string) error {
	exists, err := isFileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("docker: file %s does not exist", path)
	}
	return nil
}

func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// TODO(pedge): what?
func emptyStruct() struct{} {
	var str struct{}
	return str
}
