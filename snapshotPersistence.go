package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type snapshotStorageFormat struct {
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

// Convert linked list in ListValue to an array of values for snapshot storage
func convertListToSnapshotFormat(list ListValue) []string {
	var values []string
	current := list.Start
	for current != nil {
		values = append(values, current.Data)
		current = current.Next
	}
	return values
}

// Convert array of values back into a linked list structure and return ListValue
func convertSnapshotFormatToList(values []string) ListValue {
	if len(values) == 0 {
		return ListValue{}
	}
	head := &Node{Data: values[0]}
	current := head
	for _, v := range values[1:] {
		newNode := &Node{Data: v}
		current.Next = newNode
		newNode.Prev = current
		current = newNode
	}
	return ListValue{Start: head, Tail: current, Length: len(values)}
}

func loadLastSnapshot() error {
	if _, err := os.Stat("snapshot.rdb"); os.IsNotExist(err) {
		return nil
	} else {
		return loadSnapshot("snapshot.rdb")
	}
}

func loadSnapshot(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var dataStoreUnmarshalled []snapshotStorageFormat
	err = json.Unmarshal(file, &dataStoreUnmarshalled)
	if err != nil {
		return fmt.Errorf("error while unmarshalling data store: %v", err)
	}

	for _, snapshot := range dataStoreUnmarshalled {
		if valueFormatData, ok := snapshot.Data.(map[string]interface{}); ok {
			valueType, _ := valueFormatData["value_type"].(string)
			if valueType == "list" {
				// Unmarshal Value field of list directly into []string
				valueBytes, _ := json.Marshal(valueFormatData["value"])
				var values []string
				if err := json.Unmarshal(valueBytes, &values); err != nil {
					return fmt.Errorf("error in converting list snapshot: %v", err)
				}
				dataStore[snapshot.Key] = valueFormat{
					ValueType: "list",
					Value:     convertSnapshotFormatToList(values),
				}
			} else {
				// For non-list types, unmarshal into valueFormat directly
				valueBytes, _ := json.Marshal(valueFormatData["value"])
				jsonString := fmt.Sprintf(`{"value_type": "%s", "value": %s}`, valueType, string(valueBytes))
				var nonListValue valueFormat
				if err := json.Unmarshal([]byte(jsonString), &nonListValue); err != nil {
					return fmt.Errorf("error in converting non-list data: %v", err)
				}
				dataStore[snapshot.Key] = nonListValue
			}
		} else {
			return fmt.Errorf("invalid snapshot format for key: %s", snapshot.Key)
		}
	}
	log.Println("Snapshot loaded")
	return nil
}

func createSnapshot() {
	var err error
	defer func() {
		if err != nil {
			log.Println("Error while creating snapshot: ", err.Error())
		}
	}()
	// Create a temporary file in the same directory as the target file
	tempFile, err := os.CreateTemp("", "temp-snapshot.rdb")
	if err != nil {
		return
	}
	defer os.Remove(tempFile.Name())
	var dataStoreSnapshot []snapshotStorageFormat
	for key, value := range dataStore {
		if value.ValueType == "list" {
			listValue := value.Value.(ListValue)
			snapshotList := convertListToSnapshotFormat(listValue)
			dataStoreSnapshot = append(dataStoreSnapshot, snapshotStorageFormat{
				Key: key,
				Data: valueFormat{
					ValueType: "list",
					Value:     snapshotList,
				},
			})
		} else {
			dataStoreSnapshot = append(dataStoreSnapshot, snapshotStorageFormat{
				Key: key,
				Data: valueFormat{
					ValueType: "string",
					Value:     value.Value,
				},
			})
		}
	}
	// Write new data to the temporary file
	dataStoreBuf, err := json.Marshal(dataStoreSnapshot)
	if err != nil {
		return
	} else {
		_, err = tempFile.Write(dataStoreBuf)
	}
	if err != nil {
		tempFile.Close()
		return
	}
	// Close the file to ensure all data is flushed
	err = tempFile.Close()
	if err != nil {
		return
	}
	// Rename temporary file to original file name
	err = os.Rename(tempFile.Name(), "snapshot.rdb")
	if err != nil {
		return
	}
	log.Println("Snapshot created")
}
