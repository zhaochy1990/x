package deque

import "fmt"

// minCapacity is the smallest capacity that deque may have. Must be power of 2
// for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

type Deque[T any] struct {
	buf    []T
	head   int
	tail   int
	count  int
	minCap int
}

// New creates a new Deque, optionally setting the current and minimum capacity
// when non-zero values are given for these. The Deque instance returns
// operates on items of the type specified by the type argument. For example,
// to create a Deque that contains strings,
//
//	stringDeque := deque.New[string]()
//
// To create a Deque with capacity to store 2048 ints without resizing, and
// that will not resize below space for 32 items when removing items:
//
//	d := deque.New[int](2048, 32)
//
// To create a Deque that has not yet allocated memory, but after it does will
// never resize to have space for less than 64 items:
//
//	d := deque.New[int](0, 64)
//
// Any size values supplied here are rounded up to the nearest power of 2.
func New[T any](size ...int) *Deque[T] {
	var capacity, minimum int
	if len(size) >= 1 {
		capacity = size[0]
		if len(size) >= 2 {
			minimum = size[1]
		}
	}

	minCap := minCapacity
	for minCap < minimum {
		minCap <<= 1
	}

	var buf []T
	if capacity > 0 {
		bufSize := minCap
		for bufSize < capacity {
			bufSize <<= 1
		}
		buf = make([]T, bufSize)
	}

	return &Deque[T]{
		buf:    buf,
		minCap: minCap,
	}
}

// Cap returns the current capacity of the Deque. If q is nil, q.Cap() is zero.
func (q *Deque[T]) Cap() int {
	if q == nil {
		return 0
	}
	return len(q.buf)
}

// Front returns the element at the front of the queue. This is the element
// that would be returned by PopFront(). This call panics if the queue is
// empty.
func (q *Deque[T]) Front() T {
	if q.count <= 0 {
		panic("deque: Front() called when empty")
	}
	return q.buf[q.head]
}

// Back returns the element at the back of the queue. This is the element that
// would be returned by PopBack(). This call panics if the queue is empty.
func (q *Deque[T]) Back() T {
	if q.count <= 0 {
		panic("deque: Back() called when empty")
	}
	return q.buf[q.prev(q.tail)]
}

func (q *Deque[T]) PopBack() T {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}
	q.tail = q.prev(q.tail)

	back := q.buf[q.tail]
	var zero T
	q.buf[q.tail] = zero

	q.count--

	q.shrinkIfExcess()
	return back
}

func (q *Deque[T]) PopFront() T {
	if q.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	front := q.buf[q.head]
	var zero T
	q.buf[q.head] = zero
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return front
}

// PushBack appends an element to the back of the queue. Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *Deque[T]) PushBack(elem T) {
	q.growIfFull()
	q.buf[q.tail] = elem
	// Calculate new tail position.
	q.tail = q.next(q.tail)
	q.count++
}

// func (q *Deque[T]) PushFront(elem T) {
//
// }
func (q *Deque[T]) Len() int {
	return len(q.buf)
}

// growIfFull resizes up if the buffer is full.
func (q *Deque[T]) growIfFull() {
	if q.count != len(q.buf) {
		return
	}
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]T, q.minCap)
		return
	}
	// q.resize()
}

// next returns the next buffer position wrapping around buffer.
func (q *Deque[T]) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

// prev returns the prev buffer position wrapping around buffer.
func (q *Deque[T]) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1) // bitwise modulus
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *Deque[T]) shrinkIfExcess() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the deque to fit exactly twice its current contents. This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (q *Deque[T]) resize() {
	newBuf := make([]T, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *Deque[T]) Debug() {
	fmt.Println("count: ", q.count)
	fmt.Println("head: ", q.head)
	fmt.Println("tail: ", q.tail)
}
