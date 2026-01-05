package cmd

import "fmt"

func runStage(cmdOut interface{ Write([]byte) (int, error) }, stageID string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	proto, err := loadProtocol(activeProtocolName(cfg))
	if err != nil {
		return err
	}
	stage, ok := proto.StageByID(stageID)
	if !ok {
		return fmt.Errorf("unknown stage: %s", stageID)
	}

	s, err := loadState()
	if err != nil {
		return err
	}

	if !stageNoStrict {
		if err := requireStageDependencies(s, stage); err != nil {
			return err
		}
	}

	prompt, err := compilePrompt(stage, cfg, stageIDList(proto))
	if err != nil {
		return err
	}

	prompt = applyMaxChars(prompt, stageMaxChars)
	formatted, err := formatPrompt(stageFormat, stageID, prompt)
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
