package ini

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

//解碼器 請使用 NewDecoder 創建
type Decoder struct {
}

//創建 編碼器
func NewDecoder() *Decoder {
	return &Decoder{}
}

type bytesReader struct {
	b []byte
}

func (b *bytesReader) Read(data []byte) (int, error) {
	if len(b.b) == 0 {
		return 0, io.EOF
	}
	n := copy(data, b.b)
	b.b = b.b[n:]
	return n, nil
}

//解碼 ini 到 struct
func (d *Decoder) Unmarshal(b []byte, v ...interface{}) (e error) {
	defer func() {
		if emsg := recover(); emsg != nil {
			//將異常包裝爲error
			e = fmt.Errorf("%v", emsg)
		}
	}()

	objs := make(map[string]interface{})
	for i := 0; i < len(v); i++ {
		t := reflect.TypeOf(v[i])
		if t.Kind() != reflect.Ptr {
			continue
		}
		tt := t.Elem()
		if tt.Kind() != reflect.Struct {
			continue
		}
		objs[getSectionName(v[i], t)] = v[i]
	}
	r := &bytesReader{b: b}
	e = d.unmarshalReader(bufio.NewReader(r), objs)
	return

}

//解碼 ini 到 struct
func (d *Decoder) UnmarshalReader(r io.Reader, v ...interface{}) (e error) {
	defer func() {
		if emsg := recover(); emsg != nil {
			//將異常包裝爲error
			e = fmt.Errorf("%v", emsg)
		}
	}()

	objs := make(map[string]interface{})
	for i := 0; i < len(v); i++ {
		t := reflect.TypeOf(v[i])
		if t.Kind() != reflect.Ptr {
			continue
		}
		tt := t.Elem()
		if tt.Kind() != reflect.Struct {
			continue
		}
		objs[getSectionName(v[i], t)] = v[i]
	}

	e = d.unmarshalReader(bufio.NewReader(r), objs)
	return
}

func (d *Decoder) unmarshalReader(r *bufio.Reader, objs map[string]interface{}) error {
	var section interface{}
	for {
		key, val, e := d.readNode(r, section != nil)
		//讀取完畢
		if e == io.EOF {
			break
		}
		//出錯
		if e != nil {
			return e
		}

		if key == "" {
			//尋找 新段名
			section = objs[val]
		} else {
			//設置 項目
			d.setItem(section, key, val)
		}
	}
	return nil
}

//設置 屬性
func (d *Decoder) setItem(v interface{}, key, val string) {
	t := reflect.TypeOf(v).Elem()
	v0 := reflect.ValueOf(v).Elem()

	//尋找 匹配 屬性
	if _, ok := t.FieldByName(key); ok {
		v1 := v0.FieldByName(key)
		d.setVal(&v1, val)
		return
	}

	//在 tag 中 尋找屬性
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		if name[0] >= 'A' && name[0] <= 'Z' {
			str := field.Tag.Get("ini")
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}
			if str == key {
				v1 := v0.FieldByName(name)
				d.setVal(&v1, val)
			}
		}
	}

}

//設置 ptr 值
func (d *Decoder) setPtrVal(kind string, v *reflect.Value, val string) {
	switch kind {
	/*int*/
	case "int":
		x, _ := strconv.ParseInt(val, 10, 64)
		x1 := int(x)
		v.Set(reflect.ValueOf(&x1))
	case "int8":
		x, _ := strconv.ParseInt(val, 10, 64)
		x1 := int8(x)
		v.Set(reflect.ValueOf(&x1))
	case "int16":
		x, _ := strconv.ParseInt(val, 10, 64)
		x1 := int16(x)
		v.Set(reflect.ValueOf(&x1))
	case "int32":
		x, _ := strconv.ParseInt(val, 10, 64)
		x1 := int32(x)
		v.Set(reflect.ValueOf(&x1))
	case "int64":
		x, _ := strconv.ParseInt(val, 10, 64)
		v.Set(reflect.ValueOf(&x))
	/*uint*/
	case "uint":
		x, _ := strconv.ParseUint(val, 10, 64)
		x1 := uint(x)
		v.Set(reflect.ValueOf(&x1))
	case "uint8":
		x, _ := strconv.ParseUint(val, 10, 64)
		x1 := uint8(x)
		v.Set(reflect.ValueOf(&x1))
	case "uint16":
		x, _ := strconv.ParseUint(val, 10, 64)
		x1 := uint16(x)
		v.Set(reflect.ValueOf(&x1))
	case "uint32":
		x, _ := strconv.ParseUint(val, 10, 64)
		x1 := uint32(x)
		v.Set(reflect.ValueOf(&x1))
	case "uint64":
		x, _ := strconv.ParseUint(val, 10, 64)
		v.Set(reflect.ValueOf(&x))
	/*float*/
	case "float32":
		x, _ := strconv.ParseFloat(val, 64)
		x1 := float32(x)
		v.Set(reflect.ValueOf(&x1))
	case "float64":
		x, _ := strconv.ParseFloat(val, 64)
		v.Set(reflect.ValueOf(&x))
	case "string":
		v.Set(reflect.ValueOf(&val))
	case "bool":
		val = strings.ToLower(val)
		ok := true
		if val == "false" || val == "0" || val == "" {
			v.SetBool(false)
		}
		v.Set(reflect.ValueOf(&ok))
	}
}

//設置值
func (d *Decoder) setVal(v *reflect.Value, val string) {
	kind := v.Kind()
	if kind == reflect.Ptr {
		str := v.String()[2:]
		str = str[:strings.Index(str, " ")]
		d.setPtrVal(str, v, val)
		return
	}

	if kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 {
		x, _ := strconv.ParseInt(val, 10, 64)
		v.SetInt(x)
	} else if kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 {
		x, _ := strconv.ParseUint(val, 10, 64)
		v.SetUint(x)

	} else if kind == reflect.Float32 ||
		kind == reflect.Float64 {
		x, _ := strconv.ParseFloat(val, 64)
		v.SetFloat(x)
	} else if kind == reflect.String {
		v.SetString(val)
	} else if kind == reflect.Bool {
		val = strings.ToLower(val)
		if val == "false" || val == "0" || val == "" {
			v.SetBool(false)
		} else {
			v.SetBool(true)
		}
	}
}

//讀取一個數據行 section 是否已經獲取到段名
func (d *Decoder) readNode(r *bufio.Reader, section bool) ( /*key*/ string /*val or section*/, string, error) {
	for {
		b, _, e := r.ReadLine()
		if e != nil {
			return "", "", e
		}
		str := string(b)
		str = strings.TrimSpace(str)
		//爲空 讀取下行
		if str == "" {
			continue
		}
		//直接跳過
		if str[0] == ';' || //註釋
			str[0] == '=' /*錯誤數據*/ {
			continue
		}

		//無效數據 直接跳過
		n := len(str)
		if n < 3 {
			continue
		}

		//返回段
		if str[0] == '[' && str[n-1] == ']' {
			return "", strings.TrimSpace(str[1 : n-1]), nil
		}

		if !section { //不需要 項
			continue
		}
		//返回 項
		find := strings.Index(str, "=")
		if find == -1 { //無效數據
			continue
		}

		key := strings.TrimSpace(str[:find])
		if key == "" { //無效數據
			continue
		}
		return key, strings.TrimSpace(str[find+1:]), nil
	}
	return "", "", nil
}

//解碼 ini 到 struct
func UnmarshalReader(r io.Reader, v ...interface{}) (e error) {
	return NewDecoder().UnmarshalReader(r, v...)
}

//解碼 ini 到 struct
func Unmarshal(b []byte, v ...interface{}) (e error) {
	return NewDecoder().Unmarshal(b, v...)
}
