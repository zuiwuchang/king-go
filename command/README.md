# command
command 包 提供了一個 ICommander 接口 此接口 提供了 Done(command interface{}) (e error) 方法

Done 可以 依據 傳入的 command 型別 自動 調用對應的 處理函數

```go
// ICommander 命令執行器
type ICommander interface {
	// Done 執行一個 命令
	// 如果 command 未註冊 返回 errCommandUnknow
	Done(command interface{}) (e error)

	// HanderType 返回 commandType 的處理 Hander 如果不存在 返回 nil
	HanderType(commandType reflect.Type) (f CommanderHander)
	// Hander HanderType(reflect.TypeOf(command)) 的語法糖
	Hander(command interface{}) (f CommanderHander)

	// Register 註冊 一個 命令
	// 如果 f == nil 則 註銷此 處理器
	RegisterType(commandType reflect.Type, f CommanderHander)

	// Register RegisterType(reflect.TypeOf(command),f) 的語法糖
	Register(command interface{}, f CommanderHander)

	// NumHander 返回 Hander 數量
	NumHander() int
}
```

# func RegisterCommander
RegisterCommander 是一個 輔助函數 可以自動 將一個 struct 或 *struct 中 指定簽名 的函數 註冊到 ICommander 中

注入 如果 有多個 型別的 處理函數 最後一個被查找到的 函數將 作為最終處理函數

# ICommanderSignal
go 提供了 方便的 chan 然 同一個 chan 中 傳輸 多種型別的一種方法是 使用 chan interface{}

ICommanderSignal 結合了 ICommander 和 chan interface{} 功能

```go
// ICommanderSignal 一個 帶有 ICommander 的 chan
type ICommanderSignal interface {
	// Commander 返回 綁定的 ICommander
	Commander() ICommander
	// Done 如果 command 未註冊 返回 errCommandUnknow 否則執行 ch <- command
	Done(command interface{}) (e error)

	// Close close(ch) 此後 不能在 執行 任何 Done 操作
	Close()

	// Run 運行 commnad<-ch Done(e) 直到 Done 返回 錯誤 或 close(ch)
	// 如果 Run 返回 錯誤 於是你調用了 Close 此時應該 繼續調用 Run 或 RunNull 以便 在使用 帶緩存的 chan 時 緩存 goroutine 能夠正常退出
	Run() (e error)

	// RunNull 只執行 commnad<-ch 而不執行 Done(e)
	// 通常只在 使用了 帶緩存的 chan 時 Close 後 讓 chan 的 goroutine 退出
	RunNull()
}
```

# Example ICommander
```go
package main

import (
	"fmt"
	"github.com/zuiwuchang/king-go/command"
)

// Hander 定義一個 struct 用於 自動 註冊 處理 函數
type Hander struct {
	val   uint8
	val16 uint16
}

// DoneUint8 uint8 處理 函數
func (h *Hander) DoneUint8(v uint8) (e error) {
	h.val += v
	return
}

// DoneUint16 uint16 處理 函數
func (h *Hander) DoneUint16(v uint16) (e error) {
	h.val16 += v
	return
}

func main() {
	// 創建 命令 集合
	commander := command.New()
	// 如果要 修改 Hander 狀態 顯然 需要 指針
	h := &Hander{}
	// 自動 註冊 處理 函數
	command.RegisterCommander(commander,
		h,
		"Done", // 名稱 使用 此前綴的 導出函數 才會被註冊 如果傳入空字符串 不驗證名稱前綴
	)

	// 執行 命令
	commander.Done(uint8(8))   //nil
	commander.Done(uint16(16)) //nil

	// command not registered : int
	fmt.Println(commander.Done(8))

	// 8 16
	fmt.Println(h)
}
```

# Example ICommanderSignal
```go
package main

import (
	"fmt"
	"github.com/zuiwuchang/king-go/command"
	"time"
)

// Hander 定義一個 struct 用於 自動 註冊 處理 函數
type Hander struct {
	val   uint8
	val16 uint16
}

// DoneUint8 uint8 處理 函數
func (h *Hander) DoneUint8(v uint8) (e error) {
	fmt.Println("DoneUint8", v)
	h.val += v
	return
}

// DoneUint16 uint16 處理 函數
func (h *Hander) DoneUint16(v uint16) (e error) {
	fmt.Println("DoneUint16", v)
	h.val16 += v
	return
}

func main() {
	// 創建 signal
	signal := command.NewSignal(
		nil, //nil 自動 創建 chan interface{}
		nil, //nil 自動創建 ICommander
	)
	// 如果要 修改 Hander 狀態 顯然 需要 指針
	h := &Hander{}
	// 自動 註冊 處理 函數
	command.RegisterCommander(signal.Commander(),
		h,
		"Done", // 名稱 使用 此前綴的 導出函數 才會被註冊 如果傳入空字符串 不驗證名稱前綴
	)

	go func() {
		// chan <- uint8
		signal.Done(uint8(8))

		// 沒註冊的 command Done 會立刻 返回 錯誤
		// command not registered : int
		e := signal.Done(64)
		fmt.Println(e)
	}()
	go func() {
		time.Sleep(time.Second)
		// chan <- uint16
		signal.Done(uint16(16))
	}()
	go func() {
		time.Sleep(time.Second)
		// 關閉 chan 之後 不可在 Done Close
		signal.Close()
	}()

	// main goroutine deal
	for {
		e := signal.Run()
		if e == nil {
			// nil 說明 chan 已經 close

			signal.RunNull() // 只有 使用 帶緩存 的 chan 才需要 調用此函數

			break
		} else {
			if command.IsUnknow(e) {
				// 除非 運行中 動態 移除了 某個 command 或 直接 操作 ICommander.Done
				// 否則 不會出現 unknow 因為 Done 時 接直接返回錯誤了
				fmt.Println("unkonw commnad")
			}
			fmt.Println(e)
		}
	}

	// 8 16
	fmt.Println(h)
}
```
