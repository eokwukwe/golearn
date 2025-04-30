package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

var DeleteFileDefinition = ToolDefinition{
	Name:        "delete_file",
	Description: "Delete a file of a given relative file path. Use this when you want to delete a file. Do not use this with directory names.",
	InputSchema: DeleteFileInputSchema,
	Function:    DeleteFile,
}

type DeleteFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory to delete"`
}

var DeleteFileInputSchema = GenerateSchema[DeleteFileInput]()

func DeleteFile(input json.RawMessage) (string, error) {
	deleteFileInput := DeleteFileInput{}

	err := json.Unmarshal(input, &deleteFileInput)
	if err != nil {
		panic(err)
	}

	// Delete the file
	err = os.Remove(deleteFileInput.Path)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully deleted file %s", deleteFileInput.Path), nil

}
