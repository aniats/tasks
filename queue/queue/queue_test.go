package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	queue *Queue
}

func (suite *QueueTestSuite) SetupTest() {
	suite.queue = NewQueue()
}

func (suite *QueueTestSuite) TearDownTest() {
	suite.queue = nil
}

func (suite *QueueTestSuite) TestEnqueue() {
	suite.queue.Enqueue(1)
	assert.False(suite.T(), suite.queue.IsEmpty())
	assert.Equal(suite.T(), 1, suite.queue.Size())

	suite.queue.Enqueue("hello")
	suite.queue.Enqueue(3.14)
	suite.queue.Enqueue(true)

	assert.Equal(suite.T(), 4, suite.queue.Size())
}

func (suite *QueueTestSuite) TestDequeue() {
	_, err := suite.queue.Dequeue()
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, errorQueueEmpty)

	suite.queue.Enqueue(1)
	suite.queue.Enqueue(2)
	suite.queue.Enqueue(3)

	item, err := suite.queue.Dequeue()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, item)
	assert.Equal(suite.T(), 2, suite.queue.Size())

	item, _ = suite.queue.Dequeue()
	assert.Equal(suite.T(), 2, item)

	item, _ = suite.queue.Dequeue()
	assert.Equal(suite.T(), 3, item)

	assert.True(suite.T(), suite.queue.IsEmpty())
}

func (suite *QueueTestSuite) TestFront() {
	_, err := suite.queue.Front()
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, errorQueueEmpty)

	suite.queue.Enqueue("first")
	suite.queue.Enqueue("second")
	suite.queue.Enqueue("third")

	initialSize := suite.queue.Size()

	item, err := suite.queue.Front()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "first", item)
	assert.Equal(suite.T(), initialSize, suite.queue.Size())

	item2, _ := suite.queue.Front()
	assert.Equal(suite.T(), item, item2)
}

func (suite *QueueTestSuite) TestBack() {
	_, err := suite.queue.Back()
	assert.Error(suite.T(), err)

	suite.queue.Enqueue("first")
	suite.queue.Enqueue("second")
	suite.queue.Enqueue("third")

	item, err := suite.queue.Back()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "third", item)
}

func (suite *QueueTestSuite) TestIsEmpty() {
	assert.True(suite.T(), suite.queue.IsEmpty())

	suite.queue.Enqueue(1)
	assert.False(suite.T(), suite.queue.IsEmpty())

	suite.queue.Dequeue()
	assert.True(suite.T(), suite.queue.IsEmpty())
}

func (suite *QueueTestSuite) TestSize() {
	assert.Equal(suite.T(), 0, suite.queue.Size())

	for i := 1; i <= 5; i++ {
		suite.queue.Enqueue(i)
		assert.Equal(suite.T(), i, suite.queue.Size())
	}

	for i := 4; i >= 0; i-- {
		suite.queue.Dequeue()
		assert.Equal(suite.T(), i, suite.queue.Size())
	}
}

func (suite *QueueTestSuite) TestClear() {
	suite.queue.Clear()
	assert.True(suite.T(), suite.queue.IsEmpty())

	suite.queue.Enqueue(1)
	suite.queue.Enqueue(2)
	suite.queue.Enqueue(3)

	suite.queue.Clear()

	assert.True(suite.T(), suite.queue.IsEmpty())
	assert.Equal(suite.T(), 0, suite.queue.Size())
}

func (suite *QueueTestSuite) TestQueueOperationsSequence() {
	suite.queue.Enqueue("a")
	suite.queue.Enqueue("b")

	front, _ := suite.queue.Front()
	assert.Equal(suite.T(), "a", front)

	back, _ := suite.queue.Back()
	assert.Equal(suite.T(), "b", back)

	item, _ := suite.queue.Dequeue()
	assert.Equal(suite.T(), "a", item)

	front, _ = suite.queue.Front()
	assert.Equal(suite.T(), "b", front)

	back, _ = suite.queue.Back()
	assert.Equal(suite.T(), "b", back)
}

func (suite *QueueTestSuite) TestMixedTypes() {
	suite.queue.Enqueue(42)
	suite.queue.Enqueue("string")
	suite.queue.Enqueue(true)
	suite.queue.Enqueue(3.14)

	assert.Equal(suite.T(), 4, suite.queue.Size())

	item1, _ := suite.queue.Dequeue()
	assert.Equal(suite.T(), 42, item1)

	item2, _ := suite.queue.Dequeue()
	assert.Equal(suite.T(), "string", item2)

	item3, _ := suite.queue.Dequeue()
	assert.Equal(suite.T(), true, item3)

	item4, _ := suite.queue.Dequeue()
	assert.Equal(suite.T(), 3.14, item4)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
