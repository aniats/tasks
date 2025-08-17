package main

import (
	"errors"
	"testing"
)

func TestQueue_Enqueue(t *testing.T) {
	queue := NewQueue()

	queue.Enqueue(1)
	if queue.IsEmpty() {
		t.Error("Очередь не должна быть пустой после добавления элемента")
	}

	if queue.Size() != 1 {
		t.Errorf("Размер очереди должен быть 1, получен: %d", queue.Size())
	}

	queue.Enqueue("hello")
	queue.Enqueue(3.14)
	queue.Enqueue(true)

	if queue.Size() != 4 {
		t.Errorf("Размер очереди должен быть 4, получен: %d", queue.Size())
	}
}

func TestQueue_Dequeue(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Dequeue()
	if err == nil {
		t.Error("Dequeue() на пустой очереди должен возвращать ошибку")
	}

	if !errors.Is(err, errorQueueEmpty) {
		t.Error("Dequeue() на пустой очереди должен возвращать errorQueueEmpty")
	}

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	item, err := queue.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() не должен возвращать ошибку: %v", err)
	}

	if item != 1 {
		t.Errorf("Dequeue() должен возвращать 1, получен: %v", item)
	}

	if queue.Size() != 2 {
		t.Errorf("Размер очереди после Dequeue() должен быть 2, получен: %d", queue.Size())
	}

	item, _ = queue.Dequeue()
	if item != 2 {
		t.Errorf("Второй Dequeue() должен возвращать 2, получен: %v", item)
	}

	item, _ = queue.Dequeue()
	if item != 3 {
		t.Errorf("Третий Dequeue() должен возвращать 3, получен: %v", item)
	}

	if !queue.IsEmpty() {
		t.Error("Очередь должна быть пустой после извлечения всех элементов")
	}
}

func TestQueue_Front(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Front()
	if err == nil {
		t.Error("Front() на пустой очереди должен возвращать ошибку")
	}

	if !errors.Is(err, errorQueueEmpty) {
		t.Error("Front() на пустой очереди должен возвращать errorQueueEmpty")
	}

	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	initialSize := queue.Size()

	item, err := queue.Front()
	if err != nil {
		t.Errorf("Front() не должен возвращать ошибку: %v", err)
	}

	if item != "first" {
		t.Errorf("Front() должен возвращать 'first', получен: %v", item)
	}

	if queue.Size() != initialSize {
		t.Errorf("Front() не должен изменять размер очереди. Был: %d, стал: %d",
			initialSize, queue.Size())
	}

	item2, _ := queue.Front()
	if item != item2 {
		t.Error("Повторный Front() должен возвращать тот же элемент")
	}
}

func TestQueue_Back(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Back()
	if err == nil {
		t.Error("Back() на пустой очереди должен возвращать ошибку")
	}

	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	item, err := queue.Back()
	if err != nil {
		t.Errorf("Back() не должен возвращать ошибку: %v", err)
	}

	if item != "third" {
		t.Errorf("Back() должен возвращать 'third', получен: %v", item)
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	queue := NewQueue()

	if !queue.IsEmpty() {
		t.Error("Новая очередь должна быть пустой")
	}

	queue.Enqueue(1)
	if queue.IsEmpty() {
		t.Error("Очередь не должна быть пустой после добавления элемента")
	}

	queue.Dequeue()
	if !queue.IsEmpty() {
		t.Error("Очередь должна быть пустой после извлечения единственного элемента")
	}
}

func TestQueue_Size(t *testing.T) {
	queue := NewQueue()

	if queue.Size() != 0 {
		t.Errorf("Размер пустой очереди должен быть 0, получен: %d", queue.Size())
	}

	for i := 1; i <= 5; i++ {
		queue.Enqueue(i)
		if queue.Size() != i {
			t.Errorf("Размер очереди должен быть %d, получен: %d", i, queue.Size())
		}
	}

	for i := 4; i >= 0; i-- {
		queue.Dequeue()
		if queue.Size() != i {
			t.Errorf("Размер очереди должен быть %d, получен: %d", i, queue.Size())
		}
	}
}

func TestQueue_Clear(t *testing.T) {
	queue := NewQueue()

	queue.Clear()
	if !queue.IsEmpty() {
		t.Error("Очередь должна быть пустой после очистки")
	}

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	queue.Clear()

	if !queue.IsEmpty() {
		t.Error("Очередь должна быть пустой после очистки")
	}

	if queue.Size() != 0 {
		t.Errorf("Размер очереди после очистки должен быть 0, получен: %d", queue.Size())
	}
}
