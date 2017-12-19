# rbtree
go的map雖然 很方便 快速 然 hash 是沒有順序的 如果需要保存的 元素 有 順序 紅黑樹 就顯得比較有用了 rbtree 就是 孤實現的一個 紅黑樹

rbtree 只提供了 幾個 對外 接口 相當 易用 你只需要為 key 實現 IKey 接口 即可 孤預定義了 幾個 IKeyInt IKeyString .... 為 基本型別 實現了 IKey 接口

```Go
//key 定義
type IKey interface {
	// -1(<) 0(==) 1(>)
	Compare(k IKey) int
}

//value 定義
type IValue interface{}

//節點定義
type IElement interface {
	//返回 key
	Key() IKey
	//返回 value
	Value() IValue
	//設置 value
	SetValue(v interface{})
}
```

# Example
```Go
package main

import (
	"fmt"
	"king-go/container/rbtree"
)

func main() {
	//創建 樹
	tree := rbtree.New()

	//插入
	for i := 0; i < 10; i++ {
		tree.Insert(
			rbtree.IKeyInt(i),
			i+100,
		)
	}

	//查找 元素
	ele := tree.Get(rbtree.IKeyInt(2))
	if ele != nil {
		fmt.Println(ele.Value().(int) == 2+100)

		//設置
		ele.SetValue(-102)
	}
	//設置
	tree.Insert(rbtree.IKeyInt(1), -101)

	//range
	tree.Do(func(ele rbtree.IElement) bool {
		fmt.Println(ele.Value())

		//返回 false 將 停止 遍歷
		return true
	})
	//逆向 range
	tree.DoReverse(func(ele rbtree.IElement) bool {
		fmt.Println(ele.Value())

		//返回 false 將 停止 遍歷
		return true
	})

	//刪除
	tree.Erase(tree.Get(rbtree.IKeyInt(100)))
	tree.Erase(tree.Get(rbtree.IKeyInt(9)))
	tree.EraseByKey(rbtree.IKeyInt(100))
	tree.EraseByKey(rbtree.IKeyInt(0))

	//返回 最大 最小 key
	fmt.Println("min =", tree.Min().Key())
	fmt.Println("max =", tree.Max().Key())

}
```
