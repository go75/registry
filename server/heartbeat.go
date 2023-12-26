package server

import (
	"net/http"
	"registry/service"
	"time"
)

func (s *RegistryService)heartBeat(attempCount int, attempDuration, checkDuration time.Duration) {
	go func() {
		for {
			unregistServices := make([]*service.ServiceInfo, 0)
			s.ServiceInfos.RLockRangeFunc(func(serviceName string, infos []*service.ServiceInfo) {
				for i := len(infos) - 1; i >= 0; i-- {
					srv := (infos)[i]
					for j := 0; j < attempCount; j++ {
						resp, err := http.Get("http://" + srv.Addr + "/heart-beat")
						if err == nil && resp.StatusCode == http.StatusOK {
							goto NEXT
						}
						time.Sleep(attempDuration)
					}
					
					unregistServices = append(unregistServices, srv)

					NEXT:
				}
			})

			for i := len(unregistServices) - 1; i >= 0; i-- {
				s.unregistService(unregistServices[i])
			}

			time.Sleep(checkDuration)
		}
	}()
}