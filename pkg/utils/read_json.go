package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

)

func ReadJSON[T any](fileName string) (*T, error) {
    file, err := os.ReadFile(filepath.Join(".", fileName))
    if err != nil {
        return nil, err
    }

    var data T
    if err := json.Unmarshal(file, &data); err != nil {
        return nil, err
    }

    return &data, nil
}


func WriteJson[T any](fileName string, cfg T) (bool, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return false, err
	}

	err = os.WriteFile(
		filepath.Join(".", fmt.Sprintf("%s.json", fileName)),
		data,
		0644,
	)

	return err == nil, err
}
