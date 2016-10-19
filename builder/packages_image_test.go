package builder

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hpcloud/fissile/model"
	"github.com/hpcloud/termui"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDockerfile(t *testing.T) {
	assert := assert.New(t)

	ui := termui.New(
		os.Stdin,
		ioutil.Discard,
		nil,
	)

	workDir, err := os.Getwd()
	assert.NoError(err)

	compiledPackagesDir := filepath.Join(workDir, "../test-assets/tor-boshrelease-fake-compiled")
	targetPath, err := ioutil.TempDir("", "fissile-test")
	assert.NoError(err)
	defer os.RemoveAll(targetPath)

	packagesImageBuilder, err := NewPackagesImageBuilder("foo", compiledPackagesDir, targetPath, "3.14.15", ui)
	assert.NoError(err)

	baseImage := "scratch:latest"

	dockerfile := bytes.Buffer{}
	err = packagesImageBuilder.generateDockerfile(baseImage, &dockerfile)
	assert.NoError(err)

	var lines []string
	for _, line := range strings.Split(dockerfile.String(), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	assert.Equal([]string{
		"FROM scratch:latest",
		"ADD specs /opt/hcf/specs",
		"ADD packages-src /var/vcap/packages-src/",
	}, lines, "Unexpected dockerfile contents found")
}

func TestCreatePackagesDockerStream(t *testing.T) {
	assert := assert.New(t)

	ui := termui.New(
		os.Stdin,
		ioutil.Discard,
		nil,
	)

	workDir, err := os.Getwd()
	assert.NoError(err)

	releasePath := filepath.Join(workDir, "../test-assets/tor-boshrelease")
	releasePathCache := filepath.Join(releasePath, "bosh-cache")

	compiledPackagesDir := filepath.Join(workDir, "../test-assets/tor-boshrelease-fake-compiled")
	targetPath, err := ioutil.TempDir("", "fissile-test")
	assert.NoError(err)
	defer os.RemoveAll(targetPath)

	release, err := model.NewDevRelease(releasePath, "", "", releasePathCache)
	assert.NoError(err)

	roleManifestPath := filepath.Join(workDir, "../test-assets/role-manifests/tor-good.yml")
	rolesManifest, err := model.LoadRoleManifest(roleManifestPath, []*model.Release{release})
	assert.NoError(err)

	packagesImageBuilder, err := NewPackagesImageBuilder("foo", compiledPackagesDir, targetPath, "3.14.15", ui)
	assert.NoError(err)

	tarStream, errors, err := packagesImageBuilder.CreatePackagesDockerStream(
		rolesManifest,
		filepath.Join(workDir, "../test-assets/test-opinions/opinions.yml"),
		filepath.Join(workDir, "../test-assets/test-opinions/dark-opinions.yml"),
	)
	assert.NoError(err)
	defer tarStream.Close()

	pkg := getPackage(rolesManifest.Roles, "myrole", "tor", "tor")
	if !assert.NotNil(pkg) {
		return
	}

	expectedContents := map[string]string{
		"Dockerfile": `
			FROM foo-role-base:3.14.15
			ADD specs /opt/hcf/specs
			ADD packages-src /var/vcap/packages-src/`,
		"specs/foorole/tor.json": `{
				"job": {
					"name": "foorole",
					"templates": [{"name": "tor"}]
				},
				"networks": {
					"default": {}
				},
				"parameters": {},
				"properties": {
					"tor": {
						"client_keys": null,
						"hashed_control_password": null,
						"hostname": "localhost",
						"private_key": null
					}
				}
			}`,
		// The next file is empty
		"packages-src/b9973278a447dfb5e8e67661deaa5fe7001ad742/foo": ``,
	}
	tarReader := tar.NewReader(tarStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if !assert.NoError(err) {
			break
		}
		if expected, ok := expectedContents[header.Name]; ok {
			actual, err := ioutil.ReadAll(tarReader)
			assert.NoError(err)
			if strings.HasSuffix(header.Name, ".json") {
				assert.JSONEq(expected, string(actual))
			} else {
				var expectedLines, actualLines []string
				for _, line := range strings.Split(expected, "\n") {
					line = strings.TrimSpace(line)
					if line != "" {
						expectedLines = append(expectedLines, line)
					}
				}
				for _, line := range strings.Split(string(actual), "\n") {
					line = strings.TrimSpace(line)
					if line != "" {
						actualLines = append(actualLines, line)
					}
				}
				assert.Equal(expectedLines, actualLines)
			}
		}
	}

	for err := range errors {
		assert.NoError(err)
	}
}
