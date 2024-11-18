package lib

import (
	"fmt"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

const (
	MAP_TYPE  = "map[interface {}]interface {}"
	LIST_TYPE = "[]interface {}"
)

type Merger struct {
	MergeFileList []string
	Voliations    []WrappedError
}

func (m *Merger) CountViolationByType() map[string]int {
	count := map[string]int{}

	for _, v := range m.Voliations {
		if _, ok := count[v.Violation]; ok {
			count[v.Violation] += 1
		} else {
			count[v.Violation] = 1
		}
	}
	return count

}

func (m *Merger) CountViolationByLevel() map[string]int {
	count := map[string]int{}

	for _, v := range m.Voliations {
		if _, ok := count[v.Level.String()]; ok {
			count[v.Level.String()] += 1
		} else {
			count[v.Level.String()] = 1
		}
	}
	return count

}

func (m *Merger) MergeFiles(fileList []string) map[interface{}]interface{} {

	m.MergeFileList = fileList

	var base map[interface{}]interface{}
	var override map[interface{}]interface{}

	baseFilePath := m.MergeFileList[0]
	base, err := m.ReadFilePathAsYaml(baseFilePath)

	for _, filePath := range m.MergeFileList[1:] {
		override, err = m.ReadFilePathAsYaml(filePath)
		m.Mergerecursive(baseFilePath, filePath, &base, &override, ".")
		baseFilePath = filePath
	}

	_, err = yaml.Marshal(base)
	if err != nil {
		panic(err)
	}

	return base
}

func (m *Merger) handleViolation(violation string, err error) error {
	wErr := NewWrappedError(violation, err)
	m.Voliations = append(m.Voliations, *wErr)
	return wErr
}

func (m *Merger) Mergerecursive(baseFilePath string, overrideFilePath string, base *map[interface{}]interface{}, override *map[interface{}]interface{}, objectPath string) {
	if *base == nil {
		m.handleViolation(VIOLATION_INVALID_INPUT_FILE, fmt.Errorf("basefile:%s is empty", baseFilePath))
		*base = *override
		return
	}

	if override == nil {
		m.handleViolation(VIOLATION_INVALID_INPUT_FILE, fmt.Errorf("overridefile:%s is empty", overrideFilePath))
		return
	}

	for k, overrideValue := range *override {
		currentObjectPath := fmt.Sprintf("%s%s", objectPath, k)
		baseValue, exists := (*base)[k]
		if exists {
			overrideValueType := reflect.TypeOf(overrideValue).String()
			switch overrideValueType {
			case MAP_TYPE:
				overrideNode := overrideValue.(map[interface{}]interface{})
				if reflect.TypeOf((*base)[k]).String() == MAP_TYPE {
					basenode := (*base)[k].(map[interface{}]interface{})
					m.Mergerecursive(baseFilePath, overrideFilePath, &basenode, &overrideNode, currentObjectPath)
				} else {
					m.handleViolation(VIOLATION_NON_MAP_MERGE, fmt.Errorf("%s contains map at path:%s, base is not a map", overrideFilePath, currentObjectPath))
				}

			case LIST_TYPE:
				existingList := (baseValue).([]interface{})
				overrideList := (overrideValue).([]interface{})
				joint := append(overrideList, existingList...)
				(*base)[k] = joint
				log.Debugf("append   %s %s\n", currentObjectPath, overrideFilePath)

			default:
				if overrideValue == baseValue {
					m.handleViolation(VIOLATION_DUPLICATE_OVERRIDE_VALUE, fmt.Errorf("%s duplicate value at path:%s", overrideFilePath, currentObjectPath))
				} else {
					(*base)[k] = overrideValue
				}
				log.Debugf("override %s %s\n", currentObjectPath, overrideFilePath)
			}

		} else {
			log.Debugf("add      %s %s\n", currentObjectPath, overrideFilePath)
			(*base)[k] = overrideValue
		}

	}
}

func (m *Merger) ReadFilePathAsYaml(path string) (map[interface{}]interface{}, error) {
	var asYaml map[interface{}]interface{}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		err = m.handleViolation(VIOLATION_MISSING_INPUT_FILE, err)
		return asYaml, err
	}
	if len(fileContent) == 0 {
		err = m.handleViolation(VIOLATION_INVALID_INPUT_FILE, err)
		return asYaml, err
	}

	if err := yaml.Unmarshal(fileContent, &asYaml); err != nil {
		err = m.handleViolation(VIOLATION_INVALID_INPUT_FILE, err)
		return asYaml, err
	}

	return asYaml, nil
}
