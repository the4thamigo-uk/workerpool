package workerpool_test

import (
	"github.com/stretchr/testify/require"
	"github.com/the4thamigo-uk/workerpool"
	"testing"
	"time"
)

func testWork(id int, d time.Duration, ids chan int) workerpool.Work {
	return func() {
		time.Sleep(d)
		ids <- id
	}
}

func TestWorkPool_ThreeJobs(t *testing.T) {
	p, err := workerpool.New(1, 3)
	require.NoError(t, err)
	require.NotNil(t, p)
	defer p.Close()

	res := make(chan int, 1)
	err = p.Add(testWork(1, time.Second, res))
	require.NoError(t, err)
	err = p.Add(testWork(2, time.Second, res))
	require.NoError(t, err)
	err = p.Add(testWork(3, time.Second, res))
	require.NoError(t, err)

	id1 := <-res
	require.Equal(t, 1, id1)
	id2 := <-res
	require.Equal(t, 2, id2)
	id3 := <-res
	require.Equal(t, 3, id3)
}

func TestWorkPool_AddAfterClose(t *testing.T) {
	p, err := workerpool.New(1, 1)
	require.NoError(t, err)
	require.NotNil(t, p)
	p.Close()

	err = p.Add(func() {})
	require.Error(t, err)
}

func TestWorkPool_Cancel(t *testing.T) {
	p, err := workerpool.New(1, 1)
	require.NoError(t, err)
	require.NotNil(t, p)

	res := make(chan int, 1)
	err = p.Add(testWork(1, 2*time.Second, res))
	require.NoError(t, err)
	time.Sleep(time.Second)

	p.Cancel()

	// after complete all the result should be sitting in the chan
	select {
	case id1 := <-res:
		require.Equal(t, 1, id1)
	default:
		t.Errorf("work 1 did not complete")
	}
}

func TestWorkPool_Complete(t *testing.T) {
	p, err := workerpool.New(1, 1)
	require.NoError(t, err)
	require.NotNil(t, p)

	res := make(chan int, 3)
	err = p.Add(testWork(1, 1*time.Second, res))
	require.NoError(t, err)
	err = p.Add(testWork(2, 1*time.Second, res))
	require.NoError(t, err)
	err = p.Add(testWork(3, 1*time.Second, res))
	require.NoError(t, err)

	p.Complete()

	// after complete all the results should be sitting in the chan
	select {
	case id1 := <-res:
		require.Equal(t, 1, id1)
	default:
		t.Errorf("work 1 did not complete")
	}
	select {
	case id2 := <-res:
		require.Equal(t, 2, id2)
	default:
		t.Errorf("work 2 did not complete")
	}
	select {
	case id3 := <-res:
		require.Equal(t, 3, id3)
	default:
		t.Errorf("work 3 did not complete")
	}
}
