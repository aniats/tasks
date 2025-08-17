package main

import (
	"errors"
	"testing"
)

func TestStack_Push(t *testing.T) {
	stack := NewStack()

	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("Стек не должен быть пустым после добавления элемента")
	}

	if stack.Size() != 1 {
		t.Errorf("Размер стека должен быть 1, получен: %d", stack.Size())
	}

	stack.Push("hello")
	stack.Push(3.14)
	stack.Push(true)

	if stack.Size() != 4 {
		t.Errorf("Размер стека должен быть 4, получен: %d", stack.Size())
	}
}

func TestStack_Pop(t *testing.T) {
	stack := NewStack()

	_, err := stack.Pop()
	if err == nil {
		t.Error("Pop() на пустом стеке должен возвращать ошибку")
	}

	if !errors.Is(err, errorStackEmpty) {
		t.Error("Pop() на пустом стеке должен возвращать errorStackEmpty")
	}

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	item, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() не должен возвращать ошибку: %v", err)
	}

	if item != 3 {
		t.Errorf("Pop() должен возвращать 3, получен: %v", item)
	}

	if stack.Size() != 2 {
		t.Errorf("Размер стека после Pop() должен быть 2, получен: %d", stack.Size())
	}

	item, _ = stack.Pop()
	if item != 2 {
		t.Errorf("Второй Pop() должен возвращать 2, получен: %v", item)
	}

	item, _ = stack.Pop()
	if item != 1 {
		t.Errorf("Третий Pop() должен возвращать 1, получен: %v", item)
	}

	if !stack.IsEmpty() {
		t.Error("Стек должен быть пустым после извлечения всех элементов")
	}
}

func TestStack_Peek(t *testing.T) {
	stack := NewStack()

	_, err := stack.Peek()
	if err == nil {
		t.Error("Peek() на пустом стеке должен возвращать ошибку")
	}

	if !errors.Is(err, errorStackEmpty) {
		t.Error("Peek() на пустом стеке должен возвращать errorStackEmpty")
	}

	stack.Push("first")
	stack.Push("second")
	stack.Push("third")

	initialSize := stack.Size()

	item, err := stack.Peek()
	if err != nil {
		t.Errorf("Peek() не должен возвращать ошибку: %v", err)
	}

	if item != "third" {
		t.Errorf("Peek() должен возвращать 'third', получен: %v", item)
	}

	if stack.Size() != initialSize {
		t.Errorf("Peek() не должен изменять размер стека. Был: %d, стал: %d",
			initialSize, stack.Size())
	}

	item2, _ := stack.Peek()
	if item != item2 {
		t.Error("Повторный Peek() должен возвращать тот же элемент")
	}
}

func TestStack_IsEmpty(t *testing.T) {
	stack := NewStack()

	if !stack.IsEmpty() {
		t.Error("Новый стек должен быть пустым")
	}

	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("Стек не должен быть пустым после добавления элемента")
	}

	stack.Pop()
	if !stack.IsEmpty() {
		t.Error("Стек должен быть пустым после извлечения единственного элемента")
	}
}

func TestStack_Size(t *testing.T) {
	stack := NewStack()

	if stack.Size() != 0 {
		t.Errorf("Размер пустого стека должен быть 0, получен: %d", stack.Size())
	}

	for i := 1; i <= 5; i++ {
		stack.Push(i)
		if stack.Size() != i {
			t.Errorf("Размер стека должен быть %d, получен: %d", i, stack.Size())
		}
	}

	for i := 4; i >= 0; i-- {
		stack.Pop()
		if stack.Size() != i {
			t.Errorf("Размер стека должен быть %d, получен: %d", i, stack.Size())
		}
	}
}

func TestStack_Clear(t *testing.T) {
	stack := NewStack()

	stack.Clear()
	if !stack.IsEmpty() {
		t.Error("Стек должен быть пустым после очистки")
	}

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	stack.Clear()

	if !stack.IsEmpty() {
		t.Error("Стек должен быть пустым после очистки")
	}

	if stack.Size() != 0 {
		t.Errorf("Размер стека после очистки должен быть 0, получен: %d", stack.Size())
	}
}

func TestStack_Order(t *testing.T) {
	stack := NewStack()

	expected := []int{1, 2, 3, 4, 5}
	for _, v := range expected {
		stack.Push(v)
	}

	for i := len(expected) - 1; i >= 0; i-- {
		item, err := stack.Pop()
		if err != nil {
			t.Fatalf("Ошибка при Pop(): %v", err)
		}

		if item != expected[i] {
			t.Errorf("Ожидался элемент %d, получен: %v", expected[i], item)
		}
	}
}
