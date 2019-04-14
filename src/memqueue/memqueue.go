package memqueue

import (
	"fmt"
	"sync"
)

//Queue is a in memory message queue for URLs
type Queue struct {
	URLs []URL
	Qlen int
}

//URL is a single URL object, composed by URL domain and cleaned URL
type URL struct {
	URL    string
	Domain string
	Clean  string
}

//Push will add an element Url to the queue
func (q *Queue) Push(newURL URL) {
	var locker sync.RWMutex
	locker.Lock()
	q.URLs = append(q.URLs, newURL)
	q.checkLen()
	locker.Unlock()
}

func (q *Queue) checkLen() {
	if len(q.URLs) > q.Qlen {
		q.URLs = q.URLs[:q.Qlen]
	}
}

//Peek will get an element Url from the queue without deleting it
func (q Queue) Peek() (URL, error) {
	if ok := q.URLs[0]; ok.Domain != "" && ok.Clean != "" {
		return q.URLs[0], nil
	}
	return URL{}, fmt.Errorf("No elements in queue")
}

//PeekN will get an array of Url from the queue without deleting it
func (q Queue) PeekN(qty int) ([]URL, error) {
	var retSlice []URL
	for index := 0; index < qty; {
		if ok := q.URLs[index]; ok.Domain != "" && ok.Clean != "" {
			retSlice = append(retSlice, q.URLs[index])
		}
	}
	if len(retSlice) > 0 {
		return retSlice, nil
	}
	return retSlice, fmt.Errorf("No elements in queue")
}

//Pop will get and delete an element Url from the queue
func (q *Queue) Pop() (URL, error) {
	var locker sync.RWMutex
	locker.Lock()
	if len(q.URLs) > 0 {
		ret := q.URLs[0]
		q.URLs = q.URLs[1:]
		locker.Unlock()
		return ret, nil
	}
	locker.Unlock()
	return URL{}, fmt.Errorf("No elements in queue")
}

//Trim will  delete n elements from the start of the queue
func (q *Queue) Trim(qty int) error {
	var locker sync.RWMutex
	locker.Lock()
	for index := 0; index < qty; {
		index++
		_, err := q.Pop()
		if err != nil {
			break
		}
	}
	locker.Unlock()
	return fmt.Errorf("No elements in queue")
}

//Delete will delete given URL
func (q *Queue) Delete(given string) {
	var locker sync.RWMutex
	var found int = -1
	for k, URL := range q.URLs {
		if URL.URL == given {
			found = k
		}
	}
	locker.Lock()
	if found != -1 {
		q.URLs[len(q.URLs)-1], q.URLs[found] = q.URLs[found], q.URLs[len(q.URLs)-1]
		q.URLs = q.URLs[:len(q.URLs)-1]
	}
	locker.Unlock()
}
