/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package install

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/spf13/cobra"
	"github.com/walle/targz"
)

type Asset struct {
	Name string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	Url string `json:"html_url"`
	Tag string `json:"tag_name"`
	Assets []Asset `json:"assets"`
}

// installCmd represents the install command
var CliCmd = &cobra.Command{
	Use:   "cli [version]",
	Short: "Install(update) NoCloud CLI",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var version string
		if len(args) == 0 {
			version = "latest"
		} else {
			version = args[0]
		}
		fmt.Println("Looking for tag: ", version)

		client := github.NewClient(nil)
		repo, _, err := client.Repositories.Get(context.Background(), "slntopp", "nocloud-cli")
		if err != nil {
			return err
		}

		release := repo.GetReleasesURL()
		release = strings.Replace(release, "{/id}", "/" + version, 1)
		resp, err := http.Get(release)
		if err != nil {
			return err
		}
		if resp.StatusCode == 404 {
			return fmt.Errorf("Release tag '%s' wasn't found", version)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var r Release
		json.Unmarshal(body, &r)

		fmt.Println("Found tag: ", r.Tag)
		asset_base := fmt.Sprintf("nocloud-%s-%s-%s", r.Tag, runtime.GOOS, runtime.GOARCH)
		asset_name := asset_base + ".tar.gz"
		var asset_url string
		for _, asset := range r.Assets {
			if asset.Name == asset_name {
				fmt.Println("Found asset: ", asset.Name)
				asset_url = asset.DownloadUrl
				goto asset_found
			}
		}
		return errors.New("required asset not found")
		
		asset_found:
		fmt.Println("Downloading: ", asset_url)
		err = DownloadFile(asset_name, asset_url)
		if err != nil {
			return err
		}

		fmt.Println("Decompressing: ", asset_name)
		targz.Extract(asset_name, asset_base)

		var path string
		pathB, err := exec.Command("which", "nocloud").Output()
		if err == nil {
			path = string(pathB)
			fmt.Println("Found nocloud executable at '" + path + "'. Replacing.")
		} else {
			fmt.Println("nocloud executable going to be put into '/usr/local/bin'")
			path = "/usr/local/bin"
		}

		debug, _ := cmd.Flags().GetBool("debug")

		err = os.Rename(asset_base + "/nocloud", path + "/nocloud")
		
		if !debug {
			os.RemoveAll(asset_base)
			os.Remove(asset_base + ".tar.gz")
		}
			
		os.Chmod(path + "/nocloud", os.FileMode(0755))

		fmt.Println("Installation Finished")
		return err
	},
}

func init() {
	CliCmd.Flags().Bool("debug", false, "Keeps assets and tarball after unpacking and moving")
}
