/*將golang struct 和 ini 檔案 進行 映射

每個 struct 對應一個 ini 段

struct 的屬性 可以是 intXXX uintXXX string floatXXX bool
*/
package ini

import (
	"io"
	"reflect"
	"strings"
)

const (
	lineFlag = "\n"
)

//實現了此接口的 struct 將使用 SectionName()返回值作爲 段名 而非 型別名
type ISection interface {
	//返回 段的 名稱
	SectionName() string
}

func writeStringLine(w io.Writer, line []byte, str ...string) error {
	for i := 0; i < len(str); i++ {
		_, e := w.Write([]byte(str[i]))
		if e != nil {
			return e
		}
	}
	_, e := w.Write(line)
	if e != nil {
		return e
	}
	return nil
}

//返回 段命 名稱
func getSectionName(v interface{}, t reflect.Type) string {
	if section, ok := v.(ISection); ok {
		return section.SectionName()
	}

	str := t.String()
	f := strings.LastIndex(str, ".")
	return str[f+1:]
}
