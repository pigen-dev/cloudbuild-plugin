package helpers

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func GetStepTemplate(stepID string)([]byte, error){
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			// TODO: Handle error.
	}
	bkt := client.Bucket("pigen-templates")
	stepTemplateFile := fmt.Sprintf("%s.yaml", stepID)
	objectPath := fmt.Sprintf("step-templates/%s", stepTemplateFile)
	obj := bkt.Object(objectPath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var content bytes.Buffer
	_, err = r.WriteTo(&content)
	if err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}