package command

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mundanelizard/envi/pkg/cli"
)

func Push() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Action: handlePush,
	}
}

func handlePush(values *cli.ActionArgs, args []string) {
	_, err := srv.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	comDir, err := compressEnvironment(string(repo))
	if err != nil {
		logger.Fatal(err)
		return
	}

	encDir, err := encryptCompressedEnvironment(comDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = srv.PushRepo(string(repo), encDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Println("Environment has been pushed to server")
}

func encryptCompressedEnvironment(zip string) (string, error) {

	return "", nil
}

func compressEnvironment(repo string) (string, error) {
	zipDir := path.Join(os.TempDir(), path.Base(repo)+".envi.temp.zip")

	zipFile, err := os.Create(zipDir)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	dirname := path.Join(wd, ".envi")
	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file to add to the ZIP archive
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a new file header for the file in the ZIP archive
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the file in the ZIP archive
		header.Name = path[len(dirname)+1:]

		// Add the file header to the ZIP archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the contents of the file to the ZIP archive
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	return zipDir, err
}