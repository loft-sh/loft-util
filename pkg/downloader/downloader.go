package downloader

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/loft-sh/loft-util/pkg/downloader/commands"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Downloader interface {
	EnsureCommand(ctx context.Context) (string, error)
}

type downloader struct {
	httpGet getRequest
	command commands.Command
	log     logr.Logger
}

func NewDownloader(command commands.Command, log logr.Logger) Downloader {
	return &downloader{
		httpGet: http.Get,
		command: command,
		log:     log,
	}
}

func (d *downloader) EnsureCommand(ctx context.Context) (string, error) {
	command := d.command.Name()
	valid, err := d.command.IsValid(ctx, command)
	if err != nil {
		return "", err
	} else if valid {
		return command, nil
	}

	installPath, err := d.command.InstallPath()
	if err != nil {
		return "", err
	}

	valid, err = d.command.IsValid(ctx, installPath)
	if err != nil {
		return "", err
	} else if valid {
		return installPath, nil
	}

	return installPath, d.downloadExecutable(command, installPath, d.command.DownloadURL())
}

func (d *downloader) downloadExecutable(command, installPath, installFromURL string) error {
	err := os.MkdirAll(filepath.Dir(installPath), 0755)
	if err != nil {
		return err
	}

	err = d.downloadFile(command, installPath, installFromURL)
	if err != nil {
		return errors.Wrap(err, "download file")
	}

	err = os.Chmod(installPath, 0755)
	if err != nil {
		return errors.Wrap(err, "cannot make file executable")
	}

	return nil
}

type getRequest func(url string) (*http.Response, error)

func (d *downloader) downloadFile(command, installPath, installFromURL string) error {
	d.log.Info("Downloading " + command + "...")

	t, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(t)

	archiveFile := filepath.Join(t, "download")
	f, err := os.Create(archiveFile)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := d.httpGet(installFromURL)
	if err != nil {
		return errors.Wrap(err, "get url")
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "download file")
	}

	err = f.Close()
	if err != nil {
		return err
	}

	// install the file
	return d.command.Install(archiveFile)
}
