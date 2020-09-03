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

func TestWorkPool_WaitForCompletion(t *testing.T) {
	p, err := workerpool.New(1, 1)
	require.NoError(t, err)
	require.NotNil(t, p)

	res := make(chan int, 1)
	err = p.Add(testWork(1, 2*time.Second, res))
	require.NoError(t, err)
	time.Sleep(time.Second)
	p.Close()

	id1 := <-res
	require.Equal(t, 1, id1)
}
