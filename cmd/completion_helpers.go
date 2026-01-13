package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"specfirst/internal/app"
	"specfirst/internal/repository"
	"specfirst/internal/starter"
)

func stageIDCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return filterPrefix(loadStageIDs(), toComplete), cobra.ShellCompDirectiveNoFileComp
}

func stageIDThenFilesCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveDefault
	}
	return filterPrefix(loadStageIDs(), toComplete), cobra.ShellCompDirectiveNoFileComp
}

func protocolNameCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return filterPrefix(loadProtocolNames(), toComplete), cobra.ShellCompDirectiveNoFileComp
}

func starterNameCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return filterPrefix(loadStarterNames(), toComplete), cobra.ShellCompDirectiveNoFileComp
}

func archiveVersionCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return filterPrefix(loadArchiveVersions(), toComplete), cobra.ShellCompDirectiveNoFileComp
}

func attestRoleCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	application, err := app.Load(protocolFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	stageID := args[0]
	roles := make(map[string]struct{})
	for _, approval := range application.Protocol.Approvals {
		if approval.Stage == stageID {
			roles[approval.Role] = struct{}{}
		}
	}
	if len(roles) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	values := make([]string, 0, len(roles))
	for role := range roles {
		values = append(values, role)
	}
	sort.Strings(values)
	return filterPrefix(values, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func loadStageIDs() []string {
	application, err := app.Load(protocolFlag)
	if err != nil {
		return nil
	}
	stageIDs := make([]string, 0, len(application.Protocol.Stages))
	for _, stage := range application.Protocol.Stages {
		stageIDs = append(stageIDs, stage.ID)
	}
	return stageIDs
}

func loadProtocolNames() []string {
	entries, err := os.ReadDir(repository.ProtocolsPath())
	if err != nil {
		return nil
	}
	var protocols []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			protocols = append(protocols, strings.TrimSuffix(name, ".yaml"))
		}
	}
	sort.Strings(protocols)
	return protocols
}

func loadStarterNames() []string {
	starters, err := starter.List()
	if err != nil {
		return nil
	}
	values := make([]string, 0, len(starters))
	for _, s := range starters {
		values = append(values, s.Name)
	}
	sort.Strings(values)
	return values
}

func loadArchiveVersions() []string {
	repo := repository.NewSnapshotRepository(repository.ArchivesPath())
	versions, err := repo.List()
	if err != nil {
		return nil
	}
	sort.Strings(versions)
	return versions
}

func filterPrefix(values []string, prefix string) []string {
	if prefix == "" {
		return values
	}
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		if strings.HasPrefix(value, prefix) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}
