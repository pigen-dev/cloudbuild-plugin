package cloudbuild

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"text/template"
	shared "github.com/pigen-dev/shared"

	"gopkg.in/yaml.v2"
)


type CloudbuildFile struct {
	Steps []CloudbuildStep `yaml:"steps"`
	Options map[string]any `yaml:"options"`
}

func (cb *Cloudbuild) GeneratScript(pigenFile shared.PigenStepsFile)(shared.CICDFile){
	currentDir, err := os.Getwd()
	if err != nil {
		return shared.CICDFile{
			Error: err,
			FileScript: nil,
		}
	}
	cloudbuildScript := CloudbuildFile{}

	
	for _,pigenStep := range pigenFile.Steps {
		//TODO: Get step template file from bucket
		stepTemplateFile := fmt.Sprintf("%s.yaml", pigenStep.Step)
		destFile := filepath.Join(currentDir, "step-templates", stepTemplateFile)
		b, err := os.ReadFile(destFile)
		if err != nil {
			return shared.CICDFile{
				Error: err,
				FileScript: nil,
			}
		}
		cloudbuildPigenTemplate := CloudbuildStep{}
		err = yaml.Unmarshal(b, &cloudbuildPigenTemplate)
		if err != nil {
			return shared.CICDFile{
				Error: err,
				FileScript: nil,
			}
		}
		cloudbuildStep, err := ReplacePlaceholders(pigenStep, cloudbuildPigenTemplate)
		if err != nil {
			return shared.CICDFile{
				Error: err,
				FileScript: nil,
			}
		}
		cloudbuildScript.Steps = append(cloudbuildScript.Steps, cloudbuildStep)
			
	}
	//TODO: Change hard coded options field
	cloudbuildScript.Options = map[string]any{
		"logging": "CLOUD_LOGGING_ONLY",
	}
	b, err := yaml.Marshal(cloudbuildScript)
	if err != nil {
		return shared.CICDFile{
			Error: err,
			FileScript: nil,
		}
	}
	return shared.CICDFile{
		Error: nil,
		FileScript: b,
	}
}



func ReplacePlaceholders (pigen_step shared.Step ,cloudbuildStepTemplate CloudbuildStep) (cloudbuildStep CloudbuildStep, err error) {
	//can't return nil instead of a struct
	cloudbuildStep = CloudbuildStep{}
	b_cloudbuildStepTemplate, err := yaml.Marshal(cloudbuildStepTemplate)
	if err != nil {
		return cloudbuildStep, err
	}
	//Generate a template from the extracted step template
	tmpl, err := template.New("step_template").Parse(string(b_cloudbuildStepTemplate))
	if err != nil {
		return cloudbuildStep, err
	}
	//to Execute the template i need a var of type io.Writer 
	var b bytes.Buffer
	err = tmpl.Execute(&b, pigen_step)
	if err != nil {
		return cloudbuildStep, err
	}
	//Umarshal the result bytes of the template exection into a new var
	//The output cloudbuildStep will have always the same type as the cloudbuildStepTemplate
	
	err = yaml.Unmarshal(b.Bytes(), &cloudbuildStep)
	if err != nil {
		return cloudbuildStep, fmt.Errorf("error marshaling cloudbuildStep %v", err)
	}
	return cloudbuildStep, nil
}