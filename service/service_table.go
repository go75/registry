package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
)

type ServiceTable struct {
	serviceInfos map[string][]*ServiceInfo
	lock *sync.RWMutex
}

func NewServiceTable() *ServiceTable {
	return &ServiceTable{
		serviceInfos: make(map[string][]*ServiceInfo),
		lock: new(sync.RWMutex),
	}
}

func (t *ServiceTable) ParseServiceInfos(reader io.ReadCloser) (err error){
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	defer func() {
		err = reader.Close()
	}()
	t.lock.Lock()
	defer t.lock.Unlock()
	err = json.Unmarshal(data, &t.serviceInfos)
	return
}

func (t *ServiceTable) BuildRequiredServiceInfos(service *ServiceInfo) map[string][]*ServiceInfo {
	m := make(map[string][]*ServiceInfo, len(service.RequiredServices))
	t.lock.RLock()
	defer t.lock.RUnlock()
	
	for _, serviceName := range service.RequiredServices {
		m[serviceName] = t.serviceInfos[serviceName]
	}

	return m
}

func (t *ServiceTable) Add(service *ServiceInfo) {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Printf("Service table add %s with address %s\n", service.Name, service.Addr)
	t.serviceInfos[service.Name] = append(t.serviceInfos[service.Name], service)
}

func (t *ServiceTable) Remove(service *ServiceInfo) {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Printf("Service table remove %s with address %s\n", service.Name, service.Addr)
	services := t.serviceInfos[service.Name]
	for i := len(services) - 1; i >= 0; i-- {
		if services[i].Addr == service.Addr {
			t.serviceInfos[service.Name] = append(services[:i], services[i+1:]...)
		}
	}
}

func (t *ServiceTable) Get(serviceName string) *ServiceInfo {
	t.lock.RLock()
	defer t.lock.RUnlock()
	services, ok := t.serviceInfos[serviceName]
	if !ok || len(services) < 1 {
		return nil
	}
	idx := rand.Intn(len(services))
	return services[idx]
}

func (t *ServiceTable) Dump() {
	t.lock.RLock()
	defer t.lock.RUnlock()
	fmt.Println("==========Dump Service Table Start==========")
	for k, v := range t.serviceInfos {
		fmt.Print("Service " + k + ": [ ")
		for i := 0; i < len(v); i++ {
			fmt.Print(v[i].Addr + " ")
		}
		fmt.Println("]")
	}
	fmt.Println("==========Dump Service Table End==========")
}

func (t *ServiceTable) RLockFunc(fn func()) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	fn()
}

func (t *ServiceTable) LockFunc(fn func()) {
	t.lock.Lock()
	defer t.lock.Unlock()

	fn()
}

func (t *ServiceTable) RLockRangeFunc(fn func(string, []*ServiceInfo)) {
	t.RLockFunc(func() {
		for k, v := range t.serviceInfos {
			fn(k, v)
		}
	})
}

func (t *ServiceTable) LockRangeFunc(fn func(string, []*ServiceInfo)) {
	t.LockFunc(func() {
		for k, v := range t.serviceInfos {
			fn(k, v)
		}
	})
}