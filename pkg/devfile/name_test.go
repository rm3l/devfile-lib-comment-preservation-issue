package devfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devfile/library/pkg/testingutil/filesystem"
	"github.com/google/go-cmp/cmp"
)

func TestUpdateName(t *testing.T) {
	fs := filesystem.DefaultFs{}
	tempDir, err := fs.TempDir("", "devfile")
	if err != nil {
		t.Errorf("could not create temporary file: %v\n", err)
		return
	}

	defer func() {
		_ = fs.RemoveAll(tempDir)
	}()

	type args struct {
		name string
	}
	for _, tt := range []struct {
		name           string
		devfileContent string
		args           args
		wantErr        bool
		want           string
	}{
		{
			name:           "without comments",
			devfileContent: devfileContentProviderWithoutComments("my-app-name"),
			args: args{
				name: "app-renamed",
			},
			want: devfileContentProviderWithoutComments("app-renamed"),
		},
		{
			name:           "with comments",
			devfileContent: devfileContentProviderWithComments("another-app-name"),
			args: args{
				name: "another-app-renamed",
			},
			want: devfileContentProviderWithComments("another-app-renamed"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(tempDir, "devfile.yaml")
			if err = fs.WriteFile(filename, []byte(tt.devfileContent), os.FileMode(0644)); err != nil {
				t.Errorf("could not write content to file %q: %v\n", filename, err)
				return
			}

			err = UpdateName(filename, tt.args.name)

			if tt.wantErr != (err != nil) {
				t.Errorf("wantErr=%v, but got err: %v", tt.wantErr, err)
				return
			}

			var b []byte
			b, err = fs.ReadFile(filename)
			if err != nil {
				t.Errorf("could not read content of Devfile at %q: %v\n", filename, err)
				return
			}

			if diff := cmp.Diff(tt.want, string(b)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func devfileContentProviderWithComments(name string) string {
	return fmt.Sprintf(`#
# A comment that should not get removed
#
commands:
- exec:
    # Comment explaining this build command
    commandLine: GOCACHE=${PROJECT_SOURCE}/.cache go build main.go
    component: runtime
    group:
      isDefault: true
      kind: build
    workingDir: ${PROJECT_SOURCE}
  id: build
- exec:
    # Comment explaining this run command
    commandLine: ./main
    component: runtime
    group:
      isDefault: true
      kind: run
    workingDir: ${PROJECT_SOURCE}
  id: run
components:
- container:
    endpoints:
    - name: http
      targetPort: 8080
    image: quay.io/devfile/golang:latest
    memoryLimit: 1024Mi
    mountSources: true
  name: runtime
metadata:
  description: Stack with the latest Go version
  displayName: Go Runtime
  icon: https://raw.githubusercontent.com/devfile-samples/devfile-stack-icons/main/golang.svg
  language: go
  name: %s
  projectType: go
  version: 1.0.0
schemaVersion: 2.1.0
`, name)
}

func devfileContentProviderWithoutComments(name string) string {
	return fmt.Sprintf(`commands:
- exec:
    commandLine: GOCACHE=${PROJECT_SOURCE}/.cache go build main.go
    component: runtime
    group:
      isDefault: true
      kind: build
    workingDir: ${PROJECT_SOURCE}
  id: build
- exec:
    commandLine: ./main
    component: runtime
    group:
      isDefault: true
      kind: run
    workingDir: ${PROJECT_SOURCE}
  id: run
components:
- container:
    endpoints:
    - name: http
      targetPort: 8080
    image: quay.io/devfile/golang:latest
    memoryLimit: 1024Mi
    mountSources: true
  name: runtime
metadata:
  description: Stack with the latest Go version
  displayName: Go Runtime
  icon: https://raw.githubusercontent.com/devfile-samples/devfile-stack-icons/main/golang.svg
  language: go
  name: %s
  projectType: go
  version: 1.0.0
schemaVersion: 2.1.0
`, name)
}
