// Copyright Meshplay Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"fmt"

	meshkitutils "github.com/khulnasoft/meshkit/utils"
	meshkitkube "github.com/khulnasoft/meshkit/utils/kubernetes"
	"github.com/khulnasoft/meshplay/meshplayctl/internal/cli/root/config"
	"github.com/khulnasoft/meshplay/meshplayctl/internal/cli/root/constants"
	c "github.com/khulnasoft/meshplay/meshplayctl/pkg/constants"
	"github.com/khulnasoft/meshplay/meshplayctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkDocUpdate = map[string]string{
	"link":    "![update-usage](/assets/img/meshplayctl/update.png)",
	"caption": "Usage of meshplayctl system update",
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshplay images/manifest files.",
	Long:  `Pull new Meshplay container images and manifests from artifact repository.`,
	Example: `
// Pull new Meshplay images from Docker Hub. This does not update meshplayctl. This command may be executed while Meshplay is running.
meshplayctl system update

// Pull the latest manifest files alone
meshplayctl system update --skip-reset
	`,
	Annotations: linkDocUpdate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "update"))
		}

		// Get viper instance used for context
		mctlCfg, err := config.GetMeshplayCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		if tempContext != "" {
			err = mctlCfg.SetCurrentContext(tempContext)
			if err != nil {
				utils.Log.Error(ErrSetCurrentContext(err))
				return nil
			}
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		err = currCtx.ValidateVersion()
		if err != nil {
			return err
		}
		if currCtx.GetVersion() != "latest" {
			// ask confirmation if user has pinned the version in config
			log.Infof("You have pinned version: %s in your current context", currCtx.GetVersion())
			userResponse := false
			if utils.SilentFlag {
				userResponse = true
			} else {
				userResponse = utils.AskForConfirmation("Updating Meshplay container images/manifest files will supersede the version to latest. Are you sure you want to continue")
			}

			if !userResponse {
				log.Info("Update aborted.")
				return nil
			}
			currCtx.SetVersion("latest")
		}

		log.Info("Updating Meshplay...")

		switch currCtx.GetPlatform() {
		case "docker":
			log.Info("Updating Meshplay containers")
			err = utils.UpdateMeshplayContainers()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to update Meshplay containers"))
			}

			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())

			if err != nil {
				return err
			}

		case "kubernetes":
			// create a client
			kubeClient, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}
			meshplayImageVersion := currCtx.GetVersion()
			providerURL := viper.GetString(c.ProviderURLsENV)
			// If the user skips reset, then just restart the pods else fetch updated manifest files and apply them
			if !utils.SkipResetFlag {

				// Apply the latest helm chart along with the default image tag specified in the charts "stable-latest"
				if err = applyHelmCharts(kubeClient, currCtx, meshplayImageVersion, false, meshkitkube.UPGRADE, "", providerURL); err != nil {
					return errors.Wrap(err, "cannot update Meshplay")
				}
			}

			// run k8s checks to make sure if k8s cluster is running
			hcOptions := &HealthCheckOptions{
				PrintLogs:           false,
				IsPreRunE:           false,
				Subcommand:          "",
				RunKubernetesChecks: true,
			}
			hc, err := NewHealthChecker(hcOptions)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			// If k8s is available in case of platform docker than we deploy operator
			if err = hc.Run(); err != nil {
				return ErrHealthCheckFailed(err)
			}

			running, err := utils.AreMeshplayComponentsRunning(currCtx.GetPlatform())
			if err != nil {
				return err
			}
			if !running {
				// Meshplay is not running, run the start command
				if err := start(); err != nil {
					return ErrRestartMeshplay(err)
				}
			}

			currCtx.SetVersion("latest")
			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
			if err != nil {
				return err
			}

			log.Info("... updated Meshplay in the Kubernetes Cluster.")
		}

		log.Info("Meshplay is now up-to-date")
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		latestVersions, err := meshkitutils.GetLatestReleaseTagsSorted(c.GetMeshplayGitHubOrg(), c.GetMeshplayGitHubRepo())
		version := constants.GetMeshplayctlVersion()
		if err == nil {
			if len(latestVersions) == 0 {
				log.Warn("no versions found for Meshplay")
				return
			}
			latest := latestVersions[len(latestVersions)-1]
			if latest != version {
				log.Printf("A new release of meshplayctl is available: %s → %s", version, latest)
				log.Printf("https://github.com/meshplay/meshplay/releases/tag/%s", latest)
				log.Print("Check https://docs.meshplay.khulnasofy.com/installation/upgrades#upgrading-meshplay-cli for instructions on how to update meshplayctl\n")
			}
		}
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&utils.SkipResetFlag, "skip-reset", "", false, "(optional) skip checking for new Meshplay manifest files.")
}
