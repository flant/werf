package get_autogenerated_values

import (
	"fmt"

	"github.com/flant/logboek"
	helm_common "github.com/flant/werf/cmd/werf/helm/common"

	"github.com/flant/werf/cmd/werf/common"
	"github.com/flant/werf/pkg/deploy"
	"github.com/flant/werf/pkg/docker"
	"github.com/flant/werf/pkg/docker_registry"
	"github.com/flant/werf/pkg/lock"
	"github.com/flant/werf/pkg/ssh_agent"
	"github.com/flant/werf/pkg/true_git"
	"github.com/flant/werf/pkg/util"
	"github.com/flant/werf/pkg/werf"
	"github.com/spf13/cobra"
)

var CmdData struct {
}

var CommonCmdData common.CmdData

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-autogenerated-values",
		Short: "Get service values yaml generated by werf for helm chart during deploy",
		Long: common.GetLongCommandDescription(`Get service values generated by werf for helm chart during deploy.

These values includes project name, docker images ids and other`),
		DisableFlagsInUseLine: true,
		Annotations: map[string]string{
			common.CmdEnvAnno: common.EnvsDescription(common.WerfSecretKey),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetServiceValues()
		},
	}

	common.SetupDir(&CommonCmdData, cmd)
	common.SetupTmpDir(&CommonCmdData, cmd)
	common.SetupHomeDir(&CommonCmdData, cmd)
	common.SetupSSHKey(&CommonCmdData, cmd)

	common.SetupTag(&CommonCmdData, cmd)
	common.SetupEnvironment(&CommonCmdData, cmd)
	common.SetupNamespace(&CommonCmdData, cmd)

	common.SetupStagesStorage(&CommonCmdData, cmd)
	common.SetupImagesRepo(&CommonCmdData, cmd)
	common.SetupDockerConfig(&CommonCmdData, cmd, "Command needs granted permissions to read and pull images from the specified stages storage and images repo")
	common.SetupInsecureRepo(&CommonCmdData, cmd)

	return cmd
}

func runGetServiceValues() error {
	logboek.MuteOut()

	if err := werf.Init(*CommonCmdData.TmpDir, *CommonCmdData.HomeDir); err != nil {
		return fmt.Errorf("initialization error: %s", err)
	}

	if err := lock.Init(); err != nil {
		return err
	}

	if err := deploy.Init(deploy.InitOptions{WithoutHelm: true}); err != nil {
		return err
	}

	if err := true_git.Init(true_git.Options{Out: logboek.GetOutStream(), Err: logboek.GetErrStream()}); err != nil {
		return err
	}

	if err := docker_registry.Init(docker_registry.Options{AllowInsecureRepo: *CommonCmdData.InsecureRepo}); err != nil {
		return err
	}

	if err := docker.Init(*CommonCmdData.DockerConfig); err != nil {
		return err
	}

	projectDir, err := common.GetProjectDir(&CommonCmdData)
	if err != nil {
		return fmt.Errorf("getting project dir failed: %s", err)
	}

	werfConfig, err := common.GetWerfConfig(projectDir)
	if err != nil {
		return fmt.Errorf("bad config: %s", err)
	}

	imagesRepo := common.GetOptionalImagesRepo(werfConfig.Meta.Project, &CommonCmdData)
	withoutRepo := true
	if imagesRepo != "" {
		withoutRepo = false
	}

	imagesRepo = helm_common.GetImagesRepoOrStub(imagesRepo)

	environment := helm_common.GetEnvironmentOrStub(*CommonCmdData.Environment)

	namespace, err := common.GetKubernetesNamespace(*CommonCmdData.Namespace, environment, werfConfig)
	if err != nil {
		return err
	}

	tag, tagStrategy, err := helm_common.GetTagOrStub(&CommonCmdData)
	if err != nil {
		return err
	}

	if err := ssh_agent.Init(*CommonCmdData.SSHKeys); err != nil {
		return fmt.Errorf("cannot initialize ssh agent: %s", err)
	}
	defer func() {
		err := ssh_agent.Terminate()
		if err != nil {
			logboek.LogErrorF("WARNING: ssh agent termination failed: %s\n", err)
		}
	}()

	images := deploy.GetImagesInfoGetters(werfConfig.Images, werfConfig.ImagesFromDockerfile, imagesRepo, tag, withoutRepo)

	serviceValues, err := deploy.GetServiceValues(werfConfig.Meta.Project, imagesRepo, namespace, tag, tagStrategy, images, deploy.ServiceValuesOptions{Env: environment})
	if err != nil {
		return fmt.Errorf("error creating service values: %s", err)
	}

	fmt.Printf("%s", util.DumpYaml(serviceValues))

	return nil
}
