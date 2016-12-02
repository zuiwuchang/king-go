package tcppool

import (
	"errors"
	"log"
	"time"
)

//創建 連接池
func NewPool(how PoolHow) *Pool {
	pool := &Pool{}
	pool.Init(how)
	return pool
}

const (
	POOL_MIN_PING_DURATION          = time.Hour
	POOL_MIN_PING_INTERVAL_DURATION = time.Minute * 30
	//POOL_MIN_PING_DURATION          = time.Second * 10
	//POOL_MIN_PING_INTERVAL_DURATION = time.Second * 10

	_POOL_CMD_OK    = 1
	_POOL_CMD_CLOSE = 2

	_POOL_CMD_MALLOC = 10

	_POOL_CMD_EXCEPTION = 20
	_POOL_CMD_PING      = 21
	_POOL_CMD_PING_OK   = 22

	_POOL_CMD_GET_SUM  = 100
	_POOL_CMD_GET_USE  = 101
	_POOL_CMD_GET_FREE = 102
)

type poolCmdParams struct {
	Cmd   int
	Ch    chan int //rs
	Conns map[*Conn]bool
}

//tcp 連接池
type Pool struct {
	//是否已經初始化
	isInit bool

	//連接池 設置
	initConns            int
	minConns             int
	maxConns             int
	pingDuration         time.Duration
	pingIntervalDuration time.Duration
	how                  PoolHow

	//哨兵 監聽 和服務器 連接狀態
	conn *Conn

	//正在ping的連接數量
	sumPing int
	//最後一次ping 開始時間
	lastPing time.Time
	//服務器 最後 異常時間
	lastClose time.Time

	//空閒 連接
	freeConns map[*Conn]bool

	//已使用 連接
	useConns map[*Conn]bool

	//channel
	chCmd    chan poolCmdParams //controller 控制指令
	chMalloc chan *Conn         //從池中 申請連接
	chFree   chan *Conn         //將連接放回 池中
}

//連接池 控制 goroutine
func (p *Pool) controller() {
	//創建 初始化 連接
	for i := 0; i < p.initConns; i++ {
		c := p.malloc()
		if c == nil {
			break
		}
		p.freeConns[c] = true
	}

	//創建 ping 通知 goroutine
	if p.pingDuration > 0 {
		go func(ch chan poolCmdParams, duration time.Duration) {
			defer func() {
				//send on ch close
				recover()
			}()

			for {
				time.Sleep(duration)
				ch <- poolCmdParams{Cmd: _POOL_CMD_PING}
			}
		}(p.chCmd, p.pingDuration)
	}

	//事件循環
	for {
		select {
		case c := <-p.chFree:
			//驗證是否是 池中 連接
			if _, ok := p.useConns[c]; ok {
				if c.IsOk {
					//將 連接 從 use 移動到 free
					delete(p.useConns, c)
					p.freeConns[c] = true

					//空閒容量 增加 縮容
					p.resizeLess()
				} else {
					//連接 已經 不可用 斷開
					delete(p.useConns, c)
					p.free(c)
				}
			}
		case params := <-p.chCmd:
			switch params.Cmd {
			case _POOL_CMD_PING:
				//執行 ping 驗證 連接 是否 可用
				if p.sumPing > 0 {
					//上次 ping 未結束 不需要 執行 新的 ping
					break
				}
				p.lastPing = time.Now()

				p.sumPing = len(p.freeConns)
				pingConns := p.freeConns
				p.freeConns = make(map[*Conn]bool)

				go func(ch chan poolCmdParams, conns map[*Conn]bool, duration time.Duration, how PoolHow) {
					for c, _ := range conns {
						last := c.GetLastActive()
						last.Add(duration)
						now := time.Now()
						if !now.After(last) {
							continue
						}

						//開始 ping
						if err := how.Ping(c); err != nil {
							c.Conn.Close()
							delete(conns, c)
						}
					}
					params := poolCmdParams{Cmd: _POOL_CMD_PING_OK}

					if len(conns) > 0 {
						params.Conns = conns
					}
					ch <- params
				}(p.chCmd, pingConns, p.pingIntervalDuration, p.how)
			case _POOL_CMD_PING_OK: //ping 結束 合併 ping 節點 到 空閒節點
				p.sumPing = 0
				//在ping結束前 服務器 沒有發送異常
				//則合併 否則 ping 的結果時無效的 直接 丟棄
				if !p.lastClose.IsZero() && p.lastClose.After(p.lastPing) {
					break
				}
				for c, _ := range params.Conns {
					p.freeConns[c] = true
				}
			case _POOL_CMD_EXCEPTION:
				//服務器 已斷開 關閉 所有連接
				p.lastClose = time.Now()
				if p.conn != nil {
					if p.conn.Conn != nil {
						p.conn.Conn.Close()
						p.conn.Conn = nil
					}
					p.conn.IsOk = false

				}
				for c, _ := range p.freeConns {
					c.Conn.Close()
					//delete(p.freeConns, c)
				}
				for c, _ := range p.useConns {
					c.Conn.Close()
					//delete(p.useConns, c)
				}
				p.freeConns = make(map[*Conn]bool)
				p.useConns = make(map[*Conn]bool)
			case _POOL_CMD_CLOSE:
				p.destory()
				params.Ch <- 0
				return
			case _POOL_CMD_GET_SUM:
				params.Ch <- len(p.useConns) + p.getFree()
			case _POOL_CMD_GET_USE:
				params.Ch <- len(p.useConns)
			case _POOL_CMD_GET_FREE:
				params.Ch <- p.getFree()
			case _POOL_CMD_MALLOC:
				needNew := true
				for c, _ := range p.freeConns {
					delete(p.freeConns, c)
					p.useConns[c] = true
					p.chMalloc <- c
					needNew = false

					if need, n := p.how.ResizeMore(len(p.useConns), p.getFree()); need {
						p.resizeMore(n)
					}
					break
				}
				if needNew {
					c := p.malloc()
					if c != nil {
						p.useConns[c] = true
					}

					p.chMalloc <- c

					//所有連接都被 佔用 擴容
					_, n := p.how.ResizeMore(len(p.useConns), 0)
					p.resizeMore(n)
				}
			}

		}
	}
}

//釋放 連接池
func (p *Pool) destory() {
	//p.isInit = false

	p.initConns = 0
	p.minConns = 0
	p.maxConns = 0
	p.how = nil

	for c, _ := range p.freeConns {
		c.Conn.Close()
	}
	p.freeConns = nil
	for c, _ := range p.useConns {
		c.Conn.Close()
	}
	p.useConns = nil

	close(p.chCmd)
	close(p.chMalloc)
	close(p.chFree)
	p.chCmd = nil
	p.chMalloc = nil
	p.chFree = nil

}

//擴容 連接池
func (p *Pool) resizeMore(need int) {
	use := len(p.useConns)
	free := p.getFree()
	sum := use + free
	if sum > p.maxConns {
		return
	}

	if need < 1 {
		need = use
	}

	if need < POOL_MIN_CONN {
		//當服務器 停止或異常 所有連接都會 斷開
		//故 len(p.useConns) 會爲0 需要重新 創建初始化 連接
		need = p.initConns
	}

	max := p.how.MaxAddStep()
	if max >= POOL_MIN_CONN && need > max {
		need = max
	}
	log.Printf("resizeMore sum=%v need=%v\n", sum, need)
	for i := 0; i < need; i++ {
		c := p.malloc()
		if c == nil {
			//無法 獲取 新 連接 直接 返回
			break
		} else {
			p.freeConns[c] = true
		}
	}
}

//縮容 連接池
func (p *Pool) resizeLess() {
	//計算 use free 比例
	use := len(p.useConns)
	if use < p.minConns {
		use = p.minConns
	}

	free := p.getFree()

	//不需要 縮容
	if free < use*5/4 {
		return
	}

	need := free - use*3/4
	log.Printf("resizeLess use=%v free=%v need=%v\n", use, free, need)
	if need > 0 {
		for c, _ := range p.freeConns {
			if need < 1 {
				break
			}
			delete(p.freeConns, c)
			p.free(c)

			need--
		}
	}
}

//返回 空閒 連接量
func (p *Pool) getFree() int {
	return len(p.freeConns) + p.sumPing
}

//創建一個新的 連接
func (p *Pool) malloc() *Conn {
	if p.conn == nil {
		p.conn = &Conn{}
	}
	if !p.conn.IsOk {
		//創建 哨兵 連接
		c, e := p.how.Conect()
		if e != nil {
			log.Println("Conect", e)
			return nil
		}
		p.conn.Conn = c
		p.conn.lastActive = time.Now()
		p.conn.IsOk = true
		//啓動 哨兵 監控服務器 狀態
		ch := p.chCmd
		go func(c *Conn, ch chan poolCmdParams) {
			defer func() {
				recover() //close on chan<-
			}()
			b := make([]byte, 10, 10)
			for {
				if _, err := c.Read(b); err != nil {
					break
				}
			}
			//通知 控制器 服務器 已斷開
			log.Println("server exception")
			ch <- poolCmdParams{Cmd: _POOL_CMD_EXCEPTION}
		}(p.conn, ch)
	}

	c, e := p.how.Conect()
	if e != nil {
		log.Println("Conect", e)
		return nil
	}
	return &Conn{Conn: c, lastActive: time.Now(), IsOk: true}
}

//關閉 連接
func (p *Pool) free(c *Conn) {
	if c.Conn != nil {
		c.Conn.Close()
	}
}

//初始化 連接池
func (p *Pool) Init(how PoolHow) {
	if p.isInit {
		return
	}

	//格式化 參數
	min := how.MinConn()
	if min < POOL_MIN_CONN || min > POOL_MAX_CONN {
		min = POOL_MIN_CONN
	}
	max := how.MaxConn()
	if max < POOL_MAX_CONN {
		max = POOL_MAX_CONN
	} else if max < min {
		max = min
	}
	init := how.MinConn()
	if init < min || init > max {
		init = min
	}

	//保存 連接池 設置
	p.initConns = init
	p.minConns = min
	p.maxConns = max
	p.how = how
	p.pingDuration = how.PingDuration()
	if p.pingDuration < POOL_MIN_PING_DURATION {
		p.pingDuration = 0
	}
	p.pingIntervalDuration = how.PingIntervalDuration()
	if p.pingIntervalDuration < POOL_MIN_PING_INTERVAL_DURATION {
		p.pingIntervalDuration = POOL_MIN_PING_INTERVAL_DURATION
	}

	//創建 map
	p.freeConns = make(map[*Conn]bool)
	p.useConns = make(map[*Conn]bool)

	//創建 channel
	p.chCmd = make(chan poolCmdParams)
	p.chMalloc = make(chan *Conn)
	p.chFree = make(chan *Conn)

	//啓動 控制 goroutine
	go p.controller()

	//初始化 完成
	p.isInit = true

}

//申請一個 連接
func (p *Pool) Malloc() (*Conn, error) {
	p.chCmd <- poolCmdParams{Cmd: _POOL_CMD_MALLOC}

	c := <-p.chMalloc
	if c != nil {
		return c, nil
	}
	return nil, errors.New("cann't create a new conn")
}

//將一個 連接 釋放回 連接池
//如果 需要 關閉 此連接 而非 僅僅時放回 連接池
//在調用 Free 前 設置 c.IsOk = false
func (p *Pool) Free(c *Conn) {
	if c != nil {
		p.chFree <- c
	}
}

//返回 連接量
func (p *Pool) GetSum() int {
	params := poolCmdParams{Cmd: _POOL_CMD_GET_SUM,
		Ch: make(chan int)}
	p.chCmd <- params

	return <-params.Ch
}

//返回 已使用 連接量
func (p *Pool) GetUse() int {
	params := poolCmdParams{Cmd: _POOL_CMD_GET_USE,
		Ch: make(chan int)}
	p.chCmd <- params

	return <-params.Ch
}

//返回 空閒 連接量
func (p *Pool) GetFree() int {
	params := poolCmdParams{Cmd: _POOL_CMD_GET_FREE,
		Ch: make(chan int)}
	p.chCmd <- params

	return <-params.Ch
}

//關閉 連接池 釋放 所有資源
//一旦 Close 後 將不能再使用 此 pool
//需要 重新 使用 pool 使用 NewPool 創建 新的 pool
func (p *Pool) Close() {
	params := poolCmdParams{Cmd: _POOL_CMD_CLOSE,
		Ch: make(chan int)}
	p.chCmd <- params

	<-params.Ch
}
