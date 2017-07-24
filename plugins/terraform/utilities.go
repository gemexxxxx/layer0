package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/common/models"
)

func resolveTags(client client.Client, target, entityType string, params map[string]string) (string, error) {
	params["fuzz"] = target
	params["type"] = entityType

	taggedEntities, err := client.SelectByQuery(params)
	if err != nil {
		return "", err
	}

	if len(taggedEntities) == 0 {
		return "", fmt.Errorf("No entities of type %s found matching %s", entityType, target)
	}

	if len(taggedEntities) == 1 {
		return taggedEntities[0].EntityID, nil
	}

	entityIDs := []string{}
	for _, taggedEntity := range taggedEntities {
		entityIDs = append(entityIDs, taggedEntity.EntityID)

		for _, tag := range taggedEntity.Tags {
			if tag.Key == "name" && tag.Value == target {
				return taggedEntity.EntityID, nil
			}
		}
	}

	text := fmt.Sprintf("Multiple entities of type %s found matching %s: ", entityType, target)
	for _, id := range entityIDs {
		text += fmt.Sprintf("%s \n", id)
	}

	return "", fmt.Errorf(text)
}

func setResourceData(setter func(string, interface{}) error, values map[string]interface{}) error {
	for key, value := range values {
		if err := setter(key, value); err != nil {
			return err
		}
	}

	return nil
}

func suppressEquivalentDockerrunDiffs(k, old, new string, d *schema.ResourceData) bool {
	var oldDockerrun models.Dockerrun
	if err := json.Unmarshal([]byte(old), &oldDockerrun); err != nil {
		return false
	}

	var newDockerrun models.Dockerrun
	if err := json.Unmarshal([]byte(new), &newDockerrun); err != nil {
		return false
	}

	return reflect.DeepEqual(oldDockerrun, newDockerrun)
}