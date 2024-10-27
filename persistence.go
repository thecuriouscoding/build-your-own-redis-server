package main

import (
	"fmt"
	"os"
	"time"
)

func loadLastSnapshot() error {
	if _, err := os.Stat("snapshot.rdb"); os.IsNotExist(err) {
		return createSnapshotFile()
	} else {
		return loadFromSnapshot()
	}
}

func createSnapshotFile() error {
	file, err := os.Create("snapshot.rdb")
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func loadFromSnapshot() error {
	data, err := os.ReadFile("snapshot.rdb")
	if err != nil {
		return err
	}
	fmt.Println("File content:", string(data))
	// add this data to data store
	return nil
}

func createSnapshot() {
	// file, err := os.OpenFile("snapshot.rdb", os.O_APPEND|os.O_WRONLY, 0644)
	// if err != nil {
	// 	// fmt.Println("Error:", err)
	// 	return
	// }
	// defer file.Close()
	// // make data store and expirations such that it can be stored in this file
	// _, err = file.WriteString("\nAppending new line")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("Data appended successfully")

	//TODO- add a check to chekc whether last snapshot is in progress
	// Create a temporary file in the same directory as the target file
	tempFile, err := os.CreateTemp("", "temp-snapshot.txt")
	if err != nil {
		// return fmt.Errorf("failed to create temp file: %w", err)
		return
	}
	defer os.Remove(tempFile.Name())

	// Write new data to the temporary file
	_, err = tempFile.Write([]byte(fmt.Sprintf("new data added %v", time.Now())))
	if err != nil {
		tempFile.Close()
		// return fmt.Errorf("failed to write to temp file: %w", err)
		return
	}

	// Close the file to ensure all data is flushed
	err = tempFile.Close()
	if err != nil {
		// return fmt.Errorf("failed to close temp file: %w", err)
		return
	}

	// Rename temporary file to original file name
	err = os.Rename(tempFile.Name(), "snapshot.rdb")
	if err != nil {
		// return fmt.Errorf("failed to rename temp file: %w", err)
		return
	}

	// return
}
