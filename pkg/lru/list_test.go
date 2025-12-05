package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.len())
		require.Nil(t, l.front())
		require.Nil(t, l.back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.len())

		middle := l.front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.len())
		require.Equal(t, 80, l.front().Value)
		require.Equal(t, 70, l.back().Value)

		l.MoveToFront(l.front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.len())
		for i := l.front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestList2(t *testing.T) {
	doublyLL := NewList()

	li1 := doublyLL.PushFront(1)
	li2 := doublyLL.PushFront(2)
	li3 := doublyLL.PushFront(3)
	li4 := doublyLL.PushBack(4) // 3 2 1 4
	require.Equal(t, 4, doublyLL.length)

	liFront := doublyLL.front()
	require.Equal(t, 3, liFront.Value)

	liBack := doublyLL.back()
	require.Equal(t, 4, liBack.Value)

	doublyLL.Remove(li3) // 2 1 4
	require.Equal(t, 3, doublyLL.length)
	liFront = doublyLL.front()
	require.Equal(t, 2, liFront.Value)

	doublyLL.MoveToFront(li4) // 4 2 1
	liFront = doublyLL.front()
	require.Equal(t, 4, liFront.Value)

	// doublyLL.printList() // 4 <-> 2 <-> 1 <->

	doublyLL.DeleteLinkedList()
	require.Equal(t, 0, doublyLL.length)

	require.Nil(t, doublyLL.head)
	require.Nil(t, doublyLL.tail)

	require.Nil(t, li1.Next)
	require.Nil(t, li1.Prev)
	require.Nil(t, li2.Next)
	require.Nil(t, li2.Prev)
	require.Nil(t, li4.Next)
	require.Nil(t, li4.Prev)
}
