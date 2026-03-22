package supports

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testKey = "support_test"
)

func TestArgonHash(t *testing.T) {
	t.Parallel()

	t.Run("Hash", func(t *testing.T) {
		require.NotPanics(t, func() {
			hash := ArgonHash("str")
			require.NotEmpty(t, hash)
		})
	})

	t.Run("Check Ok", func(t *testing.T) {
		hash := ArgonHash("str")
		is, err := IsStringArgonHash("str", hash)
		require.Nil(t, err)
		require.True(t, is)
	})

	t.Run("Check Error", func(t *testing.T) {
		is, err := IsStringArgonHash("str", "wrong hash")
		require.NotNil(t, err)
		require.False(t, is)
	})

	t.Run("Check False", func(t *testing.T) {
		hash := ArgonHash("str")
		is, err := IsStringArgonHash("another str", hash)
		require.Nil(t, err)
		require.False(t, is)
	})
}

func TestConcat(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		res := Concat("Power ", "of ", "Winx!")
		require.Equal(t, "Power of Winx!", res)
	})

}

func TestReadSecretFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, testKey)

	err := os.WriteFile(filePath, []byte(testKey), 0644)
	require.NoError(t, err)

	t.Run("Ok", func(t *testing.T) {
		v, err := ReadSecretFile(filePath)
		require.Nil(t, err)
		require.Equal(t, testKey, v)
	})

	t.Run("error", func(t *testing.T) {
		v, err := ReadSecretFile("unexisted")
		require.NotNil(t, err)
		require.Equal(t, "", v)
	})
}

func TestIsInContainer(t *testing.T) {
	t.Parallel()

	require.Equal(t, IsInContainer(), false)
}

func TestMakeKVMessagesJSON(t *testing.T) {
	t.Parallel()

	b, _ := MakeKVMessagesJSON("key1", "val1", "key2", "val2", "KeyNoVal")
	data := map[string]string{}
	err := json.Unmarshal(b, &data)
	require.Nil(t, err)

	require.Equal(t, data["key1"], "val1")
	require.Equal(t, data["key2"], "val2")
	_, exists := data["KeyNoVal"]
	require.False(t, exists)
}

func TestFNV1Hash(t *testing.T) {
	t.Parallel()

	v := FNV1Hash([]byte("str"))
	require.True(t, v != "")
}

func TestTextPatches(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		text1 := "todo"
		text2 := "in progress"
		textPatch := MakePatchFromTexts(text1, text2)
		patchedText, err := ApplyPatchToText(text1, textPatch)
		require.Nil(t, err)
		require.Equal(t, text2, patchedText)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		text1 := "todo"
		text2 := "in progress"
		textPatch := MakePatchFromTexts(text1, text2)

		editedText1 := "something wrong"

		_, err := ApplyPatchToText(editedText1, textPatch)
		require.NotNil(t, err)
	})
}
