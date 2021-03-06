package xsql

import (
	"encoding/json"
	"fmt"
	"strings"
)

func PrintFieldType(ft FieldType) (result string) {
	switch t := ft.(type) {
	case *BasicType:
		result = t.Type.String()
	case *ArrayType:
		result = "array("
		if t.FieldType != nil {
			result += PrintFieldType(t.FieldType)
		} else {
			result += t.Type.String()
		}
		result += ")"
	case *RecType:
		result = "struct("
		isFirst := true
		for _, f := range t.StreamFields {
			if isFirst {
				isFirst = false
			} else {
				result += ", "
			}
			result = result + f.Name + " " + PrintFieldType(f.FieldType)
		}
		result += ")"
	}
	return
}

func PrintFieldTypeForJson(ft FieldType) (result interface{}) {
	r, q := doPrintFieldTypeForJson(ft)
	if q {
		return r
	} else {
		return json.RawMessage(r)
	}
}

func doPrintFieldTypeForJson(ft FieldType) (result string, isLiteral bool) {
	switch t := ft.(type) {
	case *BasicType:
		return t.Type.String(), true
	case *ArrayType:
		var (
			fieldType string
			q         bool
		)
		if t.FieldType != nil {
			fieldType, q = doPrintFieldTypeForJson(t.FieldType)
		} else {
			fieldType, q = t.Type.String(), true
		}
		if q {
			result = fmt.Sprintf(`{"Type":"array","ElementType":"%s"}`, fieldType)
		} else {
			result = fmt.Sprintf(`{"Type":"array","ElementType":%s}`, fieldType)
		}

	case *RecType:
		result = `{"Type":"struct","Fields":[`
		isFirst := true
		for _, f := range t.StreamFields {
			if isFirst {
				isFirst = false
			} else {
				result += ","
			}
			fieldType, q := doPrintFieldTypeForJson(f.FieldType)
			if q {
				result = fmt.Sprintf(`%s{"FieldType":"%s","Name":"%s"}`, result, fieldType, f.Name)
			} else {
				result = fmt.Sprintf(`%s{"FieldType":"%s","Name":"%s"}`, result, fieldType, f.Name)
			}
		}
		result += `]}`
	}
	return result, false
}

func GetStreams(stmt *SelectStatement) (result []string) {
	if stmt == nil {
		return nil
	}
	for _, source := range stmt.Sources {
		if s, ok := source.(*Table); ok {
			result = append(result, s.Name)
		}
	}

	for _, join := range stmt.Joins {
		result = append(result, join.Name)
	}
	return
}

func LowercaseKeyMap(m map[string]interface{}) map[string]interface{} {
	m1 := make(map[string]interface{})
	for k, v := range m {
		if m2, ok := v.(map[string]interface{}); ok {
			m1[strings.ToLower(k)] = LowercaseKeyMap(m2)
		} else {
			m1[strings.ToLower(k)] = v
		}
	}
	return m1
}
