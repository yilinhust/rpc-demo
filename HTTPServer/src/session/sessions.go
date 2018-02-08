package session

import (
	"sync"
	"fmt"
	"io"
	"encoding/base64"
	"crypto/rand"
	"net/http"
	"net/url"
	"time"
)


type Provider interface {
	SessionInit(sid string) (Session, error)	//Session初始化
	SessionRead(sid string) (Session, error)	//返回sid所代表的Session，如果不存在则调用SessionInit函数创建一个新的
	SessionDestroy(sid string) error			//销毁sid对应的Session
	SessionGC(maxLifeTime int64)				//根据maxLifeTime删除过期的数据
}

type Session interface {
	Set(key, value interface{}) error				//设置值
	Get(key interface{}) interface{}				//读取值
	Delete(key interface{}) error					//删除值
	SessionID() string								//返回当前sessionID
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxlifetime int64
}

var provides = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provider
}


func NewManager(provideName string, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}


//创建全局唯一的sessionId
func (manager *Manager) createSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//session
func (manager *Manager) SessionStart(writer http.ResponseWriter, request *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := request.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.createSessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(writer, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

//session重置
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		// HttpOnly:禁止JS脚本访问
		// MaxAge小于0：临时性Cookie，不会被持久化，不会被写到Cookie文件中。
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

//session销毁
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })//定时器
}



