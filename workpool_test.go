package workerpool_test

import (
	"github.com/stretchr/testify/require"
	"github.com/the4thamigo-uk/workerpool"
	"testing"
	"time"
)

func testWork(id int, ids chan int) workerpool.Work {
	return func() {
		time.Sleep(time.Second)
		ids <- id
	}
}

func TestWorkPool_ThreeJobs(t *testing.T) {
	p, err := workerpool.New(1, 3)
	require.NoError(t, err)
	require.NotNil(t, p)
	res := make(chan int, 1)
	p.Add(testWork(1, res))
	p.Add(testWork(2, res))
	p.Add(testWork(3, res))

	id1 := <-res
	require.Equal(t, 1, id1)
	id2 := <-res
	require.Equal(t, 2, id2)
	id3 := <-res
	require.Equal(t, 3, id3)
	p.Close()
}
