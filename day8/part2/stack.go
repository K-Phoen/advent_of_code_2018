package main

type Item struct {
	value    interface{}
	previous *Item
}

type Stack struct {
	top    *Item
	length int
}

func NewStack() *Stack {
	return &Stack{top: nil, length: 0}
}

func (stack *Stack) Len() int {
	return stack.length
}

func (stack *Stack) Peek() interface{} {
	if stack.length == 0 {
		panic("Can not call Peek() on an empty stack")
	}

	return stack.top.value
}

func (stack *Stack) Pop() interface{} {
	if stack.length == 0 {
		panic("Can not call Pop() on an empty stack")
	}

	top := stack.top

	stack.top = top.previous
	stack.length--

	return top.value
}

func (stack *Stack) Push(item interface{}) {
	newItem := &Item{
		value:    item,
		previous: stack.top,
	}

	stack.top = newItem
	stack.length++
}
