package queue

type Queue []string

func (q *Queue) Enqueue(v string) {
	*q = append(*q, v)
}

func (q *Queue) Dequeue() (string, bool) {
	if len(*q) == 0 {
		return "", false
	}
	elementRemoved := (*q)[0]
	*q = (*q)[1:]
	return elementRemoved, true
}
