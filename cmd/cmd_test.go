package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/dontry/alfred-prompt-manager/service"
	"github.com/stretchr/testify/assert"
)

func setupTest() {
	// run before all tests in this file
	// run command to load environment variables
	os.Setenv("alfred_workflow_bundleid", "com.github.dongc.prompts")
	os.Setenv("alfred_workflow_cache", "cache")
	os.Setenv("alfred_workflow_data", "data")

	Init()
}

func cleanTest() {
	os.Remove("./custom_prompts.json")
	os.Remove("./awesome_prompts.json")
	os.RemoveAll("./cache")
	os.RemoveAll("./data")
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestDownload(t *testing.T) {
	// Create a temporary file for testing
	setupTest()
	t.Run("download successfully", func(t *testing.T) {
		_, err := executeCommand(rootCmd, "download")
		assert.Nil(t, err)

		// Verify that the prompt was added to the file
		data, err := ioutil.ReadFile("./awesome_prompts.json")
		assert.Nil(t, err)

		var prompts []service.Prompt
		err = json.Unmarshal(data, &prompts)
		assert.Nil(t, err)

		assert.Greater(t, len(prompts), 1)
		cleanTest()
	})

}

func TestAdd(t *testing.T) {
	setupTest()
	t.Run("add successfully", func(t *testing.T) {
		_, err := executeCommand(rootCmd, "add", "test:prompt")
		assert.Nil(t, err)

		data, err := ioutil.ReadFile("./custom_prompts.json")
		assert.Nil(t, err)

		var prompts []service.Prompt
		err = json.Unmarshal(data, &prompts)
		assert.Nil(t, err)

		// Verify prompts contains the new prompt
		found := false
		for _, p := range prompts {
			if p.Title == "test" && p.Subtitle == "prompt" {
				found = true
				break
			}
		}
		assert.True(t, found)
		cleanTest()
	})
}
