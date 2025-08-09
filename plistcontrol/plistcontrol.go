package plistcontrol

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Control[T any] struct {
	filePath string
	mutex    sync.RWMutex
}

func NewPlistControl[T any](filePath string) *Control[T] {
	return &Control[T]{
		filePath: filePath,
	}
}

func (c *Control[T]) Read() (*T, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	data := new(T)
	content, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	if len(content) == 0 {
		return data, nil
	}

	// Create a wrapper type that matches PLIST structure
	wrapper := struct {
		XMLName xml.Name `xml:"plist"`
		Dict    struct {
			Entries []PlistEntry `xml:",any"`
		} `xml:"dict"`
	}{}

	if err := xml.Unmarshal(content, &wrapper); err != nil {
		return nil, err
	}

	// Convert entries to map with proper type assertions
	plistMap := make(map[string]interface{})
	key := ""
	for _, entry := range wrapper.Dict.Entries {
		if entry.XMLName.Local == "key" {
			if str, ok := entry.Value.(string); ok {
				key = str
			}
		} else if key != "" {
			plistMap[key] = parsePlistValue(entry)
			key = ""
		}
	}

	// Convert map to target type
	if err := mapToStruct(plistMap, data); err != nil {
		return nil, err
	}

	return data, nil
}

// parsePlistValue converts PLIST entry to proper Go type
func parsePlistValue(entry PlistEntry) interface{} {
	switch entry.XMLName.Local {
	case "string":
		if str, ok := entry.Value.(string); ok {
			return str
		}
	case "integer":
		if str, ok := entry.Value.(string); ok {
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				return i
			}
		}
	case "real":
		if str, ok := entry.Value.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
		}
	case "true":
		return true
	case "false":
		return false
	case "array":
		// Handle arrays if needed
	case "dict":
		// Handle nested dictionaries if needed
	}
	return entry.Value // fallback to raw value
}

func (c *Control[T]) Write(data *T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Convert struct to map
	plistMap, err := structToMap(data)
	if err != nil {
		return err
	}

	// Create PLIST entries
	var entries []PlistEntry
	for k, v := range plistMap {
		entries = append(entries, PlistEntry{XMLName: xml.Name{Local: "key"}, Value: k})
		entries = append(entries, valueToPlistEntry(v))
	}

	// Create XML structure
	wrapper := struct {
		XMLName xml.Name `xml:"plist"`
		Version string   `xml:"version,attr"`
		Dict    struct {
			Entries []PlistEntry `xml:",any"`
		} `xml:"dict"`
	}{
		Version: "1.0",
		Dict: struct {
			Entries []PlistEntry `xml:",any"`
		}(struct{ Entries []PlistEntry }{Entries: entries}),
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "\t")
	if err := encoder.Encode(wrapper); err != nil {
		return err
	}

	// Add DOCTYPE
	xmlContent := strings.Replace(buf.String(), "<plist", `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist`, 1)

	tempFile := c.filePath + ".tmp"
	if err := ioutil.WriteFile(tempFile, []byte(xmlContent), 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, c.filePath)
}

type PlistEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",innerxml"`
}

func valueToPlistEntry(value interface{}) PlistEntry {
	switch v := value.(type) {
	case string:
		return PlistEntry{XMLName: xml.Name{Local: "string"}, Value: v}
	case bool:
		if v {
			return PlistEntry{XMLName: xml.Name{Local: "true"}}
		}
		return PlistEntry{XMLName: xml.Name{Local: "false"}}
	case int, int8, int16, int32, int64:
		return PlistEntry{XMLName: xml.Name{Local: "integer"}, Value: fmt.Sprintf("%d", v)}
	case float32, float64:
		return PlistEntry{XMLName: xml.Name{Local: "real"}, Value: fmt.Sprintf("%f", v)}
	default:
		return PlistEntry{XMLName: xml.Name{Local: "string"}, Value: fmt.Sprintf("%v", v)}
	}
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		key := field.Tag.Get("plist")
		if key == "" {
			key = field.Name
		}

		result[key] = fieldVal.Interface()
	}

	return result, nil
}

func mapToStruct(m map[string]interface{}, out interface{}) error {
	val := reflect.ValueOf(out).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		key := field.Tag.Get("plist")
		if key == "" {
			key = field.Name
		}

		if value, ok := m[key]; ok {
			if fieldVal.CanSet() {
				rv := reflect.ValueOf(value)
				if rv.Type().ConvertibleTo(fieldVal.Type()) {
					fieldVal.Set(rv.Convert(fieldVal.Type()))
				}
			}
		}
	}

	return nil
}
