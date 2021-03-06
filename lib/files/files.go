package files

import (
	"binlog-db-sync/lib/errors"
	"binlog-db-sync/lib/yml"
	
	"io"
	"os"
	"path/filepath"
	"reflect"
)

/**
 */
func AbsolutePath(relativePath string) (absolutePath string, err error) {
	dir, err := filepath.Abs(relativePath)
	if errors.CheckAndReturnIfError(err) {
		return
	}

	return dir, nil
}

/**
 */
func ReadTextFile(path string) (resultString string, err error) {
	resultString = ""

	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return
	}

	file, err := os.Open(filepath)
	if errors.CheckAndReturnIfError(err) {
		return
	}
	defer file.Close()

	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		resultString += string(data[:n])
	}

	err = nil
	return
}

/**
 */
func IsExists(path string) (bool, error) {
	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return false, err
	}

	_, err = os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/**
 */
func CreateDir(path string) error {
	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return err
	}

	return os.MkdirAll(filepath, 0755)
}

/* Copy the src file to dst. Any existing file will be overwritten and will not
   copy file attributes.
*/
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func readYAMLData(yamlFilePath string) (yamlData string, err error) {
	_, err = IsExists(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return "", err
	}

	absYamlFilePath, err := AbsolutePath(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return "", err
	}

	yamlData, err = ReadTextFile(absYamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return "", err
	}

	return
}

/* Reads  and parses 1 level yaml config file
 */
func ReadOneLevelYALM(yamlFilePath string) (yamlMap yml.OneLevelMap, err error) {
	yamlData, err := readYAMLData(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return yml.OneLevelMap{}, err
	}

	yamlMap, err = yml.ParseOneLevelYAML(yamlData)
	if errors.CheckAndReturnIfError(err) {
		return yml.OneLevelMap{}, err
	}

	return
}

/* Reads  and parses 2 level yaml config file
 */
func ReadTwoLevelYALM(yamlFilePath string) (yamlMap yml.TwoLevelMap, err error) {
	yamlData, err := readYAMLData(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return yml.TwoLevelMap{}, err
	}

	yamlMap, err = yml.ParseTwoLevelYAML(yamlData)
	if errors.CheckAndReturnIfError(err) {
		return yml.TwoLevelMap{}, err
	}

	return
}

/* Reads  and parses 3 level yaml config file
 */
func ReadThreeLevelYALM(yamlFilePath string) (yamlMap yml.ThreeLevelMap, err error) {
	yamlData, err := readYAMLData(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return yml.ThreeLevelMap{}, err
	}

	yamlMap, err = yml.ParseThreeLevelYAML(yamlData)
	if errors.CheckAndReturnIfError(err) {
		return yml.ThreeLevelMap{}, err
	}

	return
}

/* Reads  and parses 4 level yaml config file
 */
func ReadFourLevelYALM(yamlFilePath string) (yamlMap yml.FourLevelMap, err error) {
	yamlData, err := readYAMLData(yamlFilePath)
	if errors.CheckAndReturnIfError(err) {
		return yml.FourLevelMap{}, err
	}

	yamlMap, err = yml.ParseFourLevelYAML(yamlData)
	if errors.CheckAndReturnIfError(err) {
		return yml.FourLevelMap{}, err
	}

	return
}

/**
* Writes string to the of file
 */
func WriteAppendFile(filenamePath string, text string) (result int, err error) {
	absoluteFilenamePath, err := AbsolutePath(filenamePath)
	isFileExists, _ := IsExists(filenamePath)

	if !isFileExists {
		file, _ := os.Create(absoluteFilenamePath)
		defer file.Close()
	}

	file, _ := os.OpenFile(absoluteFilenamePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	result, err = file.WriteString(text + "\n")

	return
}

/**
* It adds spaces indent before string anf write this string to the end of file
 */
func WriteAppendFileWithIndent(filenamePath string, text string, indent int) (result int, err error) {
	return WriteAppendFile(filenamePath, indentInSpaces(indent)+text)
}

/**
* Writes tree stucture data (with 3 levels) to yaml file
 */
func WriteYaml(filenamePath string, data interface{}) (err error) {
	return WriteYamlBranch(filenamePath, data, 0)
}

/**
* Writes one branch of tree stucture data to yaml file
 */
func WriteYamlBranch(filenamePath string, data interface{}, depth int) (err error) {
	if data, ok := data.(map[string]map[string]map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": ", 2*depth)
			err = WriteYamlBranch(filenamePath, value, depth+1)
		}
	}

	if data, ok := data.(map[string]map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": ", 2*depth)
			err = WriteYamlBranch(filenamePath, value, depth+1)
		}
	}

	if data, ok := data.(map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": "+value, 2*depth)
		}
	}

	return
}

/**
* Returns indent with num spaces
 */
func indentInSpaces(indent int) (spacesIndent string) {
	spacesIndent = ""
	for i := 0; i < indent; i++ {
		spacesIndent += " "
	}
	return spacesIndent
}

type Pair struct {
	key   string
	value *interface{}
}

func (p Pair) writeAppendYAMLFile(filenamePath string, depth int) {
	v := reflect.ValueOf(*(p.value))
	switch v.Kind() {
	case reflect.String:
		WriteAppendFileWithIndent(filenamePath, p.key+": "+v.Elem().String(), depth)

	case reflect.Slice, reflect.Array:
		WriteAppendFileWithIndent(filenamePath, p.key+": ", depth)
		for i := 0; i < v.NumField(); i++ {
			item := v.Field(i)
			if reflect.ValueOf(item).Kind() == reflect.String {
				WriteAppendFileWithIndent(filenamePath, "- "+reflect.ValueOf(item).Elem().String(), depth+2)
			} else {
				WriteAppendFileWithIndent(filenamePath, " - the value couldn't be printed", depth+2)
			}
		}

	case reflect.Map:
		WriteAppendFileWithIndent(filenamePath, p.key+": ", depth)
		for _, key := range v.MapKeys() {
			// (Pair{key.Elem().String(), v.MapIndex(key).Pointer()}).writeAppendYAMLFile(filenamePath, depth + 2)
			// (Pair{key.Elem().String(), &"unknown value"}).writeAppendYAMLFile(filenamePath, depth + 2)
			WriteAppendFileWithIndent(filenamePath, key.Elem().String()+": unknown value", depth+2)
		}

	default:
		WriteAppendFileWithIndent(filenamePath, p.key+": the value couldn't be printed", depth)
	}
}

type YamlTree []Pair

func (t YamlTree) WriteYaml(filenamePath string) {
	for _, pair := range t {
		pair.writeAppendYAMLFile(filenamePath, 0)
	}
}
