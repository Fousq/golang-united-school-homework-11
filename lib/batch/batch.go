package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	userCh := make(chan user, n)
	defer close(userCh)

	var waitGroup sync.WaitGroup
	waitGroup.Add(int(pool))

	partitionCap := int64(n / pool)
	var partitionStartIndex int64 = 0
	var i int64

	for i = 0; i < pool; i++ {
		go func(start int64, end int64) {
			for i := start; i < end; i++ {
				userCh <- getOne(i)
			}
			waitGroup.Done()
		}(partitionStartIndex, partitionStartIndex+partitionCap)
		partitionStartIndex += partitionCap
	}

	waitGroup.Wait()

	var users []user
	for i = 0; i < n; i++ {
		users = append(users, <-userCh)
	}
	return users
}
