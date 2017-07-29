package ini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

//編碼器 請使用 NewEncoder 創建
type Encoder struct {
	//換行符 默認 \n
	Line []byte
}

//創建 編碼器
func NewEncoder() *Encoder {
	return &Encoder{Line: []byte(lineFlag)}
}

//將 struct 自動 創建爲 ini 數據 並寫入到  io.Writer
func (encoder *Encoder) MarshalWriter(w io.Writer, v ...interface{}) (e error) {
	defer func() {
		if emsg := recover(); emsg != nil {
			//將異常包裝爲error
			e = fmt.Errorf("%v", emsg)
		}
	}()

	if encoder.Line == nil {
		encoder.Line = []byte(lineFlag)
	}

	for i := 0; i < len(v); i++ {
		e = encoder.marshalNode(w, v[i])

		if e != nil {
			return
		}
	}
	return
}

//將 struct 自動 創建爲 ini 數據
func (encoder *Encoder) Marshal(v ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	e := encoder.MarshalWriter(&buf, v...)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}

func (encoder *Encoder) marshalNode(w io.Writer, v interface{}) error {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("marshal need struct")
	}

	//write [section]
	e := writeStringLine(w, encoder.Line, "[", getSectionName(v, t), "]")
	if e != nil {
		return e
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		if name[0] >= 'A' && name[0] <= 'Z' {
			str := field.Tag.Get("ini")
			str = strings.TrimSpace(str)
			if str != "" {
				name = str
			}

			dec := field.Tag.Get("dec")
			dec = strings.TrimSpace(dec)
			if dec != "" {
				e := writeStringLine(w, encoder.Line, ";", dec)
				if e != nil {
					return e
				}
			}
			fieldVal := val.Field(i)
			if fieldVal.Type().Kind() == reflect.Ptr {
				fieldVal = fieldVal.Elem()
			}
			kind := fieldVal.Type().Kind()

			if kind != reflect.Int &&
				kind != reflect.Int8 &&
				kind != reflect.Int16 &&
				kind != reflect.Int32 &&
				kind != reflect.Int64 &&

				kind != reflect.Uint &&
				kind != reflect.Uint8 &&
				kind != reflect.Uint16 &&
				kind != reflect.Uint32 &&
				kind != reflect.Uint64 &&

				kind != reflect.Float32 &&
				kind != reflect.Float64 &&

				kind != reflect.String &&
				kind != reflect.Bool {
				continue
			}
			e := writeStringLine(w, encoder.Line, name, "=", fmt.Sprint(fieldVal))
			if e != nil {
				fmt.Println(e)
				return e
			}
		}
	}
	return nil
}

//將 struct 自動 創建爲 ini 數據 並寫入到  io.Writer
func MarshalWriter(w io.Writer, v ...interface{}) error {
	return NewEncoder().MarshalWriter(w, v...)
}

//將 struct 自動 創建爲 ini 數據
func Marshal(v ...interface{}) ([]byte, error) {
	return NewEncoder().Marshal(v...)
}
