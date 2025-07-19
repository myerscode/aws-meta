package cmd

import (
	"fmt"

	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/internal/github"
	"github.com/myerscode/aws-meta/internal/util"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate AWS metadata files.",
	Run: func(cmd *cobra.Command, args []string) {

		botoRepo := github.Repo{
			Config: github.Config{
				Owner:    "boto",
				RepoName: "botocore",
				Branch:   "main",
			},
			Client: github.NewGitHubClient(""),
		}

		botocore := aws.Botocore{Repo: botoRepo}

		tags, err := botocore.Repo.FetchTags(1)

		if err != nil {
			util.PrintErrorAndExit(err)
		}

		for _, tag := range tags {
			util.LogInfo(fmt.Sprintf("Generating Partition List for Tag: %s", tag.Name))
			botocore.GeneratePartitionList(tag)
			util.LogInfo(fmt.Sprintf("Generating Service List for Tag: %s", tag.Name))
			botocore.GenerateServiceList(tag)
			util.LogInfo(fmt.Sprintf("Generating Region Service List for Tag: %s", tag.Name))
			botocore.GenerateRegionServicesList(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
