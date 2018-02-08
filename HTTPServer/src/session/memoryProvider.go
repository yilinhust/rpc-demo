package session

import (
	"container/list"
	"sync"
	"time"
)


var gmp = &MemoryProvider{list: list.New()}

func init() {
	gmp.sessions = make(map[string]*list.Element, 0)
	Register("memory", gmp)
}


//Provider实现
type MemoryProvider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (mp *MemoryProvider) SessionInit(sid string) (Session, error) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := mp.list.PushBack(newsess)
	mp.sessions[sid] = element
	return newsess, nil
}

func (mp *MemoryProvider) SessionRead(sid string) (Session, error) {
	if element, ok := mp.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := mp.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (mp *MemoryProvider) SessionDestroy(sid string) error {
	if element, ok := mp.sessions[sid]; ok {
		delete(mp.sessions, sid)
		mp.list.Remove(element)
		return nil
	}
	return nil
}

func (mp *MemoryProvider) SessionGC(maxlifetime int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	for {
		element := mp.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			mp.list.Remove(element)
			delete(mp.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

//刷新最后访问时间
func (mp *MemoryProvider) SessionUpdate(sid string) error {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if element, ok := mp.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		mp.list.MoveToFront(element)
		return nil
	}
	return nil
}




//Session实现
type SessionStore struct {
	sid          string                      //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	gmp.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	gmp.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	gmp.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}







