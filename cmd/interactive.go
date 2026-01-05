package cmd

import "specfirst/internal/assets"

type interactiveData struct {
	ProjectName  string
	ProtocolName string
	Stages       []stageSummary
	Language     string
	Framework    string
	CustomVars   map[string]string
	Constraints  map[string]string
}

type stageSummary struct {
	ID      string
	Name    string
	Intent  string
	Outputs []string
}

func runInteractive(cmdOut interface{ Write([]byte) (int, error) }) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	proto, err := loadProtocol(activeProtocolName(cfg))
	if err != nil {
		return err
	}

	stages := make([]stageSummary, 0, len(proto.Stages))
	for _, stage := range proto.Stages {
		stages = append(stages, stageSummary{
			ID:      stage.ID,
			Name:    stage.Name,
			Intent:  stage.Intent,
			Outputs: stage.Outputs,
		})
	}

	data := interactiveData{
		ProjectName:  cfg.ProjectName,
		ProtocolName: proto.Name,
		Stages:       stages,
		Language:     cfg.Language,
		Framework:    cfg.Framework,
		CustomVars:   cfg.CustomVars,
		Constraints:  cfg.Constraints,
	}

	prompt, err := renderInlineTemplate(assets.InteractiveTemplate, data)
	if err != nil {
		return err
	}

	prompt = applyMaxChars(prompt, stageMaxChars)
	formatted, err := formatPrompt(stageFormat, "interactive", prompt)
	if err != nil {
		return err
	}

	if stageOut != "" {
		if err := writeOutput(stageOut, formatted); err != nil {
			return err
		}
	}
	if _, err := cmdOut.Write([]byte(formatted)); err != nil {
		return err
	}
	return nil
}
