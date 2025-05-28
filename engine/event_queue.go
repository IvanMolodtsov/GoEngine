package engine

type EventQueue struct {
	messages chan Event
	// quit         chan any
	MessageCount uint32
}

// func NewEventQueue() *EventQueue {
// 	var queue EventQueue
// 	queue.messages = make(chan Event)
// 	queue.MessageCount = 0
// 	return &queue
// }

// func (queue *EventQueue) Push(event Event) {
// 	select {
// 	case queue.messages <- event:
// 		queue.MessageCount++
// 		return
// 	// case <-queue.quit:
// 	// 	return
// 	default:
// 		return
// 	}
// }

// func (queue *EventQueue) Pop() Event {
// 	select {
// 	case event := <-queue.messages:
// 		queue.MessageCount--
// 		return event
// 	// case <-queue.quit:
// 	// 	return nil
// 	default:
// 		return nil
// 	}
// }
