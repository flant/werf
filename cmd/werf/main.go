package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flant/werf/cmd/werf/build"
	"github.com/flant/werf/cmd/werf/build_and_publish"
	"github.com/flant/werf/cmd/werf/cleanup"
	"github.com/flant/werf/cmd/werf/deploy"
	"github.com/flant/werf/cmd/werf/dismiss"
	"github.com/flant/werf/cmd/werf/publish"
	"github.com/flant/werf/cmd/werf/purge"
	"github.com/flant/werf/cmd/werf/run"

	helm_secret_decrypt "github.com/flant/werf/cmd/werf/helm/secret/decrypt"
	helm_secret_encrypt "github.com/flant/werf/cmd/werf/helm/secret/encrypt"
	helm_secret_file_decrypt "github.com/flant/werf/cmd/werf/helm/secret/file/decrypt"
	helm_secret_file_edit "github.com/flant/werf/cmd/werf/helm/secret/file/edit"
	helm_secret_file_encrypt "github.com/flant/werf/cmd/werf/helm/secret/file/encrypt"
	helm_secret_generate_secret_key "github.com/flant/werf/cmd/werf/helm/secret/generate_secret_key"
	helm_secret_rotate_secret_key "github.com/flant/werf/cmd/werf/helm/secret/rotate_secret_key"
	helm_secret_values_decrypt "github.com/flant/werf/cmd/werf/helm/secret/values/decrypt"
	helm_secret_values_edit "github.com/flant/werf/cmd/werf/helm/secret/values/edit"
	helm_secret_values_encrypt "github.com/flant/werf/cmd/werf/helm/secret/values/encrypt"

	"github.com/flant/werf/cmd/werf/ci_env"
	"github.com/flant/werf/cmd/werf/slugify"

	images_managed_add "github.com/flant/werf/cmd/werf/images/managed/add"
	images_managed_ls "github.com/flant/werf/cmd/werf/images/managed/ls"
	images_managed_rm "github.com/flant/werf/cmd/werf/images/managed/rm"

	images_cleanup "github.com/flant/werf/cmd/werf/images/cleanup"
	images_publish "github.com/flant/werf/cmd/werf/images/publish"
	images_purge "github.com/flant/werf/cmd/werf/images/purge"

	stages_build "github.com/flant/werf/cmd/werf/stages/build"
	stages_cleanup "github.com/flant/werf/cmd/werf/stages/cleanup"
	stages_purge "github.com/flant/werf/cmd/werf/stages/purge"

	stage_image "github.com/flant/werf/cmd/werf/stage/image"

	host_cleanup "github.com/flant/werf/cmd/werf/host/cleanup"
	host_project_list "github.com/flant/werf/cmd/werf/host/project/list"
	host_project_purge "github.com/flant/werf/cmd/werf/host/project/purge"
	host_purge "github.com/flant/werf/cmd/werf/host/purge"

	helm_delete "github.com/flant/werf/cmd/werf/helm/delete"
	helm_dependency "github.com/flant/werf/cmd/werf/helm/dependency"
	helm_deploy_chart "github.com/flant/werf/cmd/werf/helm/deploy_chart"
	helm_get "github.com/flant/werf/cmd/werf/helm/get"
	helm_get_autogenerated_values "github.com/flant/werf/cmd/werf/helm/get_autogenerated_values"
	helm_get_namespace "github.com/flant/werf/cmd/werf/helm/get_namespace"
	helm_get_release "github.com/flant/werf/cmd/werf/helm/get_release"
	helm_history "github.com/flant/werf/cmd/werf/helm/history"
	helm_lint "github.com/flant/werf/cmd/werf/helm/lint"
	helm_list "github.com/flant/werf/cmd/werf/helm/list"
	helm_render "github.com/flant/werf/cmd/werf/helm/render"
	helm_repo "github.com/flant/werf/cmd/werf/helm/repo"
	helm_rollback "github.com/flant/werf/cmd/werf/helm/rollback"

	config_list "github.com/flant/werf/cmd/werf/config/list"
	config_render "github.com/flant/werf/cmd/werf/config/render"

	"github.com/flant/werf/cmd/werf/completion"
	"github.com/flant/werf/cmd/werf/docs"
	"github.com/flant/werf/cmd/werf/version"

	"github.com/flant/werf/cmd/werf/common"
	"github.com/flant/werf/cmd/werf/common/templates"
	"github.com/flant/werf/pkg/logging"
	"github.com/flant/werf/pkg/process_exterminator"
)

func main() {
	common.EnableTerminationSignalsTrap()

	if err := logging.Init(); err != nil {
		common.TerminateWithError(fmt.Sprintf("logger initialization failed: %s", err), 1)
	}

	if err := process_exterminator.Init(); err != nil {
		common.TerminateWithError(fmt.Sprintf("process exterminator initialization failed: %s", err), 1)
	}

	rootCmd := &cobra.Command{
		Use:   "werf",
		Short: "werf helps to implement and support Continuous Integration and Continuous Delivery",
		Long: common.GetLongCommandDescription(`werf helps to implement and support Continuous Integration and Continuous Delivery.

Find more information at https://werf.io`),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	groups := templates.CommandGroups{
		{
			Message: "Main Commands:",
			Commands: []*cobra.Command{
				build.NewCmd(),
				publish.NewCmd(),
				build_and_publish.NewCmd(),
				run.NewCmd(),
				deploy.NewCmd(),
				dismiss.NewCmd(),
				cleanup.NewCmd(),
				purge.NewCmd(),
			},
		},
		{
			Message: "Toolbox:",
			Commands: []*cobra.Command{
				slugify.NewCmd(),
				ci_env.NewCmd(),
			},
		},
		{
			Message: "Lowlevel Management Commands:",
			Commands: []*cobra.Command{
				configCmd(),
				stagesCmd(),
				imagesCmd(),
				helmCmd(),
				hostCmd(),
			},
		},
	}
	groups.Add(rootCmd)

	templates.ActsAsRootCommand(rootCmd, groups...)

	rootCmd.AddCommand(
		completion.NewCmd(rootCmd),
		version.NewCmd(),
		docs.NewCmd(),
		stageCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		common.TerminateWithError(err.Error(), 1)
	}
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Work with werf.yaml",
	}
	cmd.AddCommand(
		config_render.NewCmd(),
		config_list.NewCmd(),
	)

	return cmd
}

func managedImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "managed",
		Short: "Work with managed images which will be preserved during cleanup procedure",
	}
	cmd.AddCommand(
		images_managed_add.NewCmd(),
		images_managed_ls.NewCmd(),
		images_managed_rm.NewCmd(),
	)

	return cmd
}

func imagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Work with images",
	}
	cmd.AddCommand(
		images_publish.NewCmd(),
		images_cleanup.NewCmd(),
		images_purge.NewCmd(),
		managedImagesCmd(),
	)

	return cmd
}

func stagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stages",
		Short: "Work with stages, which are cache for images",
	}
	cmd.AddCommand(
		stages_build.NewCmd(),
		stages_cleanup.NewCmd(),
		stages_purge.NewCmd(),
	)

	return cmd
}

func stageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "stage",
		Hidden: true,
	}
	cmd.AddCommand(
		stage_image.NewCmd(),
	)

	return cmd
}

func helmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helm",
		Short: "Manage application deployment with helm",
	}
	cmd.AddCommand(
		helm_get_namespace.NewCmd(),
		helm_get_release.NewCmd(),
		helm_get_autogenerated_values.NewCmd(),
		helm_deploy_chart.NewCmd(),
		helm_lint.NewCmd(),
		helm_render.NewCmd(),
		helm_list.NewCmd(),
		helm_delete.NewCmd(),
		helm_rollback.NewCmd(),
		helm_get.NewCmd(),
		helm_history.NewCmd(),
		secretCmd(),
		helm_repo.NewRepoCmd(),
		helm_dependency.NewDependencyCmd(),
	)

	return cmd
}

func hostCmd() *cobra.Command {
	hostCmd := &cobra.Command{
		Use:   "host",
		Short: "Work with werf cache and data of all projects on the host machine",
	}

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Work with projects",
	}

	projectCmd.AddCommand(
		host_project_list.NewCmd(),
		host_project_purge.NewCmd(),
	)

	hostCmd.AddCommand(
		host_cleanup.NewCmd(),
		host_purge.NewCmd(),
		projectCmd,
	)

	return hostCmd
}

func secretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Work with secrets",
	}

	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "Work with secret files",
	}

	fileCmd.AddCommand(
		helm_secret_file_encrypt.NewCmd(),
		helm_secret_file_decrypt.NewCmd(),
		helm_secret_file_edit.NewCmd(),
	)

	valuesCmd := &cobra.Command{
		Use:   "values",
		Short: "Work with secret values files",
	}

	valuesCmd.AddCommand(
		helm_secret_values_encrypt.NewCmd(),
		helm_secret_values_decrypt.NewCmd(),
		helm_secret_values_edit.NewCmd(),
	)

	cmd.AddCommand(
		fileCmd,
		valuesCmd,
		helm_secret_generate_secret_key.NewCmd(),
		helm_secret_encrypt.NewCmd(),
		helm_secret_decrypt.NewCmd(),
		helm_secret_rotate_secret_key.NewCmd(),
	)

	return cmd
}