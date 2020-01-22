package consul

import "time"

// ConsulRegister consul service register
type ConsulRegister struct {
    Address                        string
    Port                           int
    ConsulAddress                  string
    ConsulPort                     int
    Service                        string
    Tag                            []string
    DeregisterCriticalServiceAfter time.Duration
    Interval                       time.Duration
}
