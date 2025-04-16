package cloudbuild

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"text/template"

	"github.com/pigen-dev/cloudbuild-plugin/helpers"
	shared "github.com/pigen-dev/shared"

	"gopkg.in/yaml.v2"
)


type CloudbuildFile struct {
	Steps []CloudbuildStep `yaml:"steps"`
}

func (cb *Cloudbuild) GeneratScript(pigenFile map[string] any)(error){
	pigen := shared.PigenSteps{}
	err := helpers.YamlConfigParser(pigenFile, &pigen)
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	destFile := filepath.Join(currentDir, "step-templates", "steps.yaml")
	b, err := os.ReadFile(destFile)
	if err != nil {
		return err
	}
	cloudbuildStepTemplates := StepTemplatesFile{}
	cloudbuildScript := CloudbuildFile{}
	//Get step templates and parse it into the cicd tool script file struct
	err = yaml.Unmarshal(b, &cloudbuildStepTemplates)
	if err != nil {
		return err
	}
	//Go through all steps in pi-gen.yaml and find it in the templates file
	for _,pigenStep := range pigen.Steps {
		for _,cloudbuildStepTemplate := range cloudbuildStepTemplates.Steps {
			//Find the desired step in the templates
			if pigenStep.Step == cloudbuildStepTemplate.Id {
				cloudbuildStep, err := ReplacePlaceholders(pigenStep, cloudbuildStepTemplate)
				if err != nil {
					return err
				}
				cloudbuildScript.Steps = append(cloudbuildScript.Steps, cloudbuildStep)
				
			}
		}
	}
	b, err = yaml.Marshal(cloudbuildScript)
	if err != nil {
		return err
	}
	return os.WriteFile("cloudbuild.yaml",b,0666)
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