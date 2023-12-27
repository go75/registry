package server

import (
	"encoding/json"
	"log"
	"net/http"
	"registry/config"
	"registry/packet"
	"registry/route"
	"registry/service"
	"time"
)

type RegistryService struct {
	ServiceInfos            *service.ServiceTable
	heartBeatAttempCount    int
	heartBeatAttempDuration time.Duration
	heartBeatCheckDuration  time.Duration
	routeTable route.RouteTable
}

func Default() *RegistryService {
	return New(3, time.Second, 60*time.Second)
}

func New(heartBeatAttempCount int, heartBeatAttempDuration, heartBeatCheckDuration time.Duration) *RegistryService {
	return &RegistryService{
		ServiceInfos:            service.NewServiceTable(),
		heartBeatAttempCount:    heartBeatAttempCount,
		heartBeatAttempDuration: heartBeatAttempDuration,
		heartBeatCheckDuration:  heartBeatCheckDuration,
		routeTable: make(route.RouteTable, 0),
	}
}

func (s *RegistryService) Regist(id uint32, fn func([]byte)) {
	s.routeTable.Regist(id, fn)
}

func (s *RegistryService) Run() error {
	connectRedis(config.REDIS_ADDR, config.REDIS_PASS, config.REDIS_DB)
	s.heartBeat(s.heartBeatAttempCount, s.heartBeatAttempDuration, s.heartBeatCheckDuration)
	s.monitor()
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serviceInfo, err := service.BuildServiceInfo(r.Body)
		if err != nil {
			log.Println("build service info err:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.registService(serviceInfo)
		if err != nil {
			log.Println("regist service err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serviceInfos := s.ServiceInfos.BuildRequiredServiceInfos(serviceInfo)
		data, err := json.Marshal(serviceInfos)
		if err != nil {
			log.Println("marshal srevice infos err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	return http.ListenAndServe(config.SERVICE_ADDR, nil)
}

func (s *RegistryService) registService(srv *service.ServiceInfo) error {
	s.ServiceInfos.Add(srv)
	data, err := packet.JsonMarshal(packet.ADD, srv)
	if err != nil {
		return err
	}

	err = rds.Publish(srv.Name, data).Err()

	return err
}

func (s *RegistryService) Publish(channel string, message interface{}) error {
	return rds.Publish(channel, message).Err()
}

func (s *RegistryService) unregistService(srv *service.ServiceInfo) error {
	s.ServiceInfos.Remove(srv)
	data, err := packet.JsonMarshal(packet.REMOVE, srv)
	if err != nil {
		return err
	}

	err = rds.Publish(srv.Name, data).Err()

	return err
}

func (s *RegistryService) Process(msg *packet.Message) {
	fn := s.routeTable[msg.ID]
	if fn != nil {
		fn(msg.Payload)
	}
}

func (s *RegistryService) monitor() {
	msgChan := rds.Subscribe(config.SERVICE_NAME).Channel()
	go func() {
		for msg := range msgChan {
			msg := packet.UnPack([]byte(msg.String()))
			s.Process(msg)
		}
	}()
}