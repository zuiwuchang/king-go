/*	靈活的爲 xorm 執行的sql 中的 where 創建 ? 參數 	*/
package params

import (
	"bytes"
	"github.com/go-xorm/xorm"
)

type Params struct {
	buf    bytes.Buffer
	params []interface{}
}

//創建 一個 參數集
//cap	? 參數 數量的 參考值
func NewParams(cap int) Params {
	return Params{params: make([]interface{}, 0, cap)}
}

//增加 where 條件
func (p *Params) WriteWhere(where string) {
	p.buf.WriteString(where)
}

//增加 ?
func (p *Params) WriteParam(param interface{}) {
	p.params = append(p.params, param)
}

//執行 Where
func (p *Params) Where(session *xorm.Session) {
	switch len(p.params) {
	case 0:
		session.Where(p.buf.String())
	case 1:
		session.Where(p.buf.String(), p.params[0])
	case 2:
		session.Where(p.buf.String(), p.params[0], p.params[1])
	case 3:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2])
	case 4:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3])
	case 5:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4])
	case 6:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5])
	case 7:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6])
	case 8:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7])
	case 9:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8])
	case 10:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9])
	case 11:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10])
	case 12:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11])
	case 13:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12])
	case 14:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13])
	case 15:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14])
	case 16:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14], p.params[15])
	case 17:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14], p.params[15], p.params[16])
	case 18:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14], p.params[15], p.params[16], p.params[17])
	case 19:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14], p.params[15], p.params[16], p.params[17], p.params[18])
	case 20:
		session.Where(p.buf.String(), p.params[0], p.params[1], p.params[2], p.params[3], p.params[4], p.params[5], p.params[6], p.params[7], p.params[8], p.params[9], p.params[10], p.params[11], p.params[12], p.params[13], p.params[14], p.params[15], p.params[16], p.params[17], p.params[18], p.params[19])
	}
}
