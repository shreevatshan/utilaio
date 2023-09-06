package utilaio

import "sync"

type ShutdownHandler struct {
	listenerlist map[string]*ShutdownListener
	mu           sync.Mutex
}

type ShutdownListener struct {
	State bool
}

func InitShutdownHandler() *ShutdownHandler {
	shutdownhandler := &ShutdownHandler{
		listenerlist: make(map[string]*ShutdownListener),
	}
	return shutdownhandler
}

func (shutdownhandler *ShutdownHandler) AddListener(listener_name string, shutdownlistener *ShutdownListener) {
	shutdownhandler.mu.Lock()
	shutdownhandler.listenerlist[listener_name] = shutdownlistener
	shutdownhandler.mu.Unlock()
}

func (shutdownhandler *ShutdownHandler) RemoveListener(listener_name string) {
	shutdownhandler.mu.Lock()
	delete(shutdownhandler.listenerlist, listener_name)
	shutdownhandler.mu.Unlock()
}

func (shutdownhandler *ShutdownHandler) NotifyListener(listener_name string) {

	if _, exists := shutdownhandler.listenerlist[listener_name]; exists {
		shutdownlistener := shutdownhandler.listenerlist[listener_name]
		shutdownlistener.Shutdown()
	}
}

func (shutdownhandler *ShutdownHandler) NotifyAllListeners() {

	for _, shutdownlistener := range shutdownhandler.listenerlist {
		shutdownlistener.Shutdown()
	}
}

func (shutdownhandler *ShutdownHandler) GetListener(listener_name string) *ShutdownListener {
	shutdownlistener := &ShutdownListener{
		State: true,
	}
	shutdownhandler.AddListener(listener_name, shutdownlistener)
	return shutdownlistener
}

func (shutdownlistener *ShutdownListener) Shutdown() {
	shutdownlistener.State = false
}
