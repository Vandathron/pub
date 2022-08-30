package pub

type subQueue struct {
	subscribers []subFunc
}

func newQueue(cf subFunc) *subQueue {
	if cf == nil {
		return &subQueue{subscribers: []subFunc{}}
	}
	return &subQueue{subscribers: []subFunc{cf}}
}

func(q *subQueue) pushFunc(chainFunc subFunc) {
	q.subscribers = append(q.subscribers,chainFunc)
}

func(q *subQueue) popFunc() (subFunc, bool) {
	if len(q.subscribers) == 0 {
		return nil,  false
	}
	var pf subFunc
	pf, q.subscribers = q.subscribers[0], q.subscribers[1:]
	return pf, true
}

func(q *subQueue) popFuncAt(idx int) {
	if idx >= len(q.subscribers) {
		return
	}
	q.subscribers = append(q.subscribers[:idx], q.subscribers[idx+1:]...)
}

func(q *subQueue) popAll() {
	q.subscribers = []subFunc{}
}

func(q *subQueue) getAllSubs()[]subFunc {
	var subs []subFunc

	for i := range q.subscribers {
		var f = q.subscribers[i]
		subs = append(subs, f)
	}
	return subs
}