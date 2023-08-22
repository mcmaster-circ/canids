package websocket

import "sync"

// Housing items being deleted
type Deleted struct {
	m sync.Mutex
	d map[string]bool
}

func (del *Deleted) update(assetID string, isDeleted bool) {
	del.m.Lock()
	del.d[assetID] = isDeleted
	del.m.Unlock()
}
func (del *Deleted) delete(assetID string) {
	del.m.Lock()
	delete(del.d, assetID)
	del.m.Unlock()
}
func (del *Deleted) getItem(assetID string) bool {
	del.m.Lock()
	defer del.m.Unlock()
	return del.d[assetID]
}
func (del *Deleted) getIDs() []string {
	del.m.Lock()
	defer del.m.Unlock()
	var list []string
	for item := range del.d {
		list = append(list, item)
	}

	return list
}

// Housing items waiting for approval
type Waiting struct {
	m sync.Mutex
	w map[string]Authorization
}

func (waiting *Waiting) update(assetID string, auth Authorization) {
	waiting.m.Lock()
	waiting.w[assetID] = auth
	waiting.m.Unlock()
}
func (waiting *Waiting) approve(assetID string) {
	waiting.m.Lock()
	curr := waiting.w[assetID]
	curr.Approved = true
	waiting.w[assetID] = curr
	waiting.m.Unlock()
}
func (waiting *Waiting) delete(assetID string) {
	waiting.m.Lock()
	delete(waiting.w, assetID)
	waiting.m.Unlock()
}
func (waiting *Waiting) getItem(assetID string) Authorization {
	waiting.m.Lock()
	defer waiting.m.Unlock()
	return waiting.w[assetID]
}
func (waiting *Waiting) getAllItems() []string {
	waiting.m.Lock()
	defer waiting.m.Unlock()
	var list []string
	for item := range waiting.w {
		list = append(list, item)
	}

	return list
}

// Housing active clients
type Active struct {
	m sync.Mutex
	a []string
}

func (active *Active) append(assetID string) {
	active.m.Lock()
	active.a = append(active.a, assetID)
	active.m.Unlock()
}
func (active *Active) delete(assetID string) {
	active.m.Lock()
	i := 0
	for _, item := range active.a {
		if item == assetID {
			break
		}
		i += 1
	}
	active.a[i] = active.a[len(active.a)-1]
	active.a = active.a[:len(active.a)-1]
	active.m.Unlock()
}
func (active *Active) exists(assetID string) bool {
	active.m.Lock()
	defer active.m.Unlock()
	for _, item := range active.a {
		if item == assetID {
			return true
		}
	}
	return false
}
