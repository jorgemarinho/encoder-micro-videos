package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {

	json := `{
			"id": "5435345345355",
			"file_path": "sega.mp4",
			"status": "pending"
		}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `wes`
	err = utils.IsJson(json)
	require.Error(t, err)
}
