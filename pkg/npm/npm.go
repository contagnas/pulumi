package npm

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Pack runs `npm pack` in the given directory, packaging the Node.js app located there into a
// tarball and returning it as `[]byte`. `stdout` is ignored for the command, as it does not
// generate useful data.
func Pack(dir string, stderr io.Writer) ([]byte, error) {
	// TODO[pulumi/pulumi#1307]: move to the language plugins so we don't have to hard code here.
	command := "npm pack"
	var c *exec.Cmd
	// We pass `--loglevel=error` to prevent `npm` from printing warnings about missing
	// `description`, `repository`, and `license` fields in the package.json file.
	c = exec.Command("npm", "pack", "--loglevel=error")
	c.Dir = dir

	// Run the command. Note that `npm pack` doesn't have the ability to rename the resulting
	// filename, since it's meant to be uploaded directly to npm, which means that we have to get
	// that information by parsing the output of the command.
	var stdout bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return nil, errors.Wrapf(err, "could not publish policies because of error running %s", command)
	}
	packfile := strings.TrimSpace(stdout.String())
	defer os.Remove(packfile)

	packTarball, err := ioutil.ReadFile(packfile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not publish policies because of error running %s",
			command)
	}

	return packTarball, nil
}

// Install runs `npm install` in the given directory, installing the dependencies for the Node.js
// app located there.
func Install(dir string, stdout, stderr io.Writer) error {
	// TODO[pulumi/pulumi#1307]: move to the language plugins so we don't have to hard code here.
	var command string
	var c *exec.Cmd
	command = "npm install"
	// We pass `--loglevel=error` to prevent `npm` from printing warnings about missing
	// `description`, `repository`, and `license` fields in the package.json file.
	c = exec.Command("npm", "install", "--loglevel=error")
	c.Dir = dir

	// Run the command.
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "installing dependencies; rerun '%s' manually to try again, "+
			"then run 'pulumi up' to perform an initial deployment", command)
	}

	// Ensure the "node_modules" directory exists.
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		return errors.Errorf("installing dependencies; rerun '%s' manually to try again", command)
	}

	return nil
}
