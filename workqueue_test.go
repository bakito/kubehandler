package kubehandler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/extensions/v1beta1"
)

func TestShouldEnqueueIntoTheUnderlyingWorkQueue(t *testing.T) {
	workQueue := NewWorkQueue("WorkqueueTest")
	workQueue.EnqueueAdd("someKind", &appsv1.Deployment{})
	timeCompleted := make(chan string, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeCompleted <- "done"
	}()
	select {
	case <-timeCompleted:
		assert.Equal(t, 1, workQueue.Length())
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Nothing in the work queue after timeout")
	}
}

func TestShouldCallRegisteredAddFuncWhenAddEventIsReceived(t *testing.T) {
	workQueue := NewWorkQueue("WorkqueueTest2")
	kind := "Foo"
	addHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterAddHandler(kind, func(namespace, name string) error {
		addHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueAdd(kind, &appsv1.Deployment{})
	go workQueue.Run(1, stopChan)
	assert.True(t, <-addHandlerCalled)
	close(stopChan)
}

func TestShouldCallRegisteredUpdateFuncWhenUpdateEventIsReceived(t *testing.T) {
	workQueue := NewWorkQueue("WorkqueueTest3")
	kind := "Foo"
	updateHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterUpdateHandler(kind, func(namespace, name string) error {
		updateHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueUpdate(kind, &appsv1.Deployment{})
	go workQueue.Run(1, stopChan)
	assert.True(t, <-updateHandlerCalled)
	close(stopChan)
}

func TestShouldCallRegisteredDeleteFuncWhenDeleteEventIsReceived(t *testing.T) {
	workQueue := NewWorkQueue("WorkqueueTest4")
	kind := "Foo"
	deleteHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterDeleteHandler(kind, func(namespace, name string) error {
		deleteHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueDelete(kind, &appsv1.Deployment{})
	go workQueue.Run(1, stopChan)
	assert.True(t, <-deleteHandlerCalled)
	close(stopChan)
}
