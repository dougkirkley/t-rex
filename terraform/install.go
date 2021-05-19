package terraform


import (
    "archive/zip"
	"path/filepath"
	"path"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

const terraformURL = "https://releases.hashicorp.com/terraform/0.13.0/terraform_0.13.0_%s_%s.zip"


func getTerraformURL() string {
	return fmt.Sprintf(terraformURL, runtime.GOOS, runtime.GOARCH)
}

func downloadFile(url string, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if out != nil {
			out.Close()
		}
	}()

	_, err = io.Copy(out, response.Body)

	return err
}

func unzip(src, dest string) ([]string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) (string, error) {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return "", err
			}
		}
		return path, nil
	}

	filePaths := []string{}
	for _, f := range r.File {
		p, err := extractAndWriteFile(f)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, p)
	}

	return filePaths, nil
}

// InstallTerraform installs Terraform 0.13
func InstallTerraform() error {
	terraformURL := getTerraformURL()
	fileName := path.Base(terraformURL)
	zipPath := fmt.Sprintf("/tmp/%s", fileName)

	log.Printf("Download %s to %s", terraformURL, zipPath)

	if err := downloadFile(terraformURL, zipPath); err != nil {
		return err
	}

	log.Printf("Unzipping %s", zipPath)

	results, err := unzip(zipPath, "/usr/local/bin/")
	if err != nil {
		return err
	}

	for _, result := range results {
		log.Printf("Extracted to %s", result)
	}

	log.Println("Successfully installed Terraform")
	return nil
}

// GetTerraformVersion gets the version of terraform running
func GetTerraformVersion() (string, error) {
	cmd := exec.Command("terraform", "version")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// WhichTerraform gets the path of the Terraform
func WhichTerraform() (string, error) {
	cmd := exec.Command("which", "terraform")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// IsTerraformInstalled returns true if Terraform is installed
func IsTerraformInstalled() bool {
	var err error
	_, err = WhichTerraform()
	if err != nil {
		return false
	}

	_, err = GetTerraformVersion()

	if err != nil {
		return false
	}

	return true
}