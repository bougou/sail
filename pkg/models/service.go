package models

import (
	"errors"
	"fmt"
	"strings"
)

// Service represents the service exposed by component.
type Service struct {
	ComponentName string `yaml:"-"`
	Name          string `yaml:"-"`

	Scheme string `yaml:"scheme"`

	// Host, IPv4, IPv6 does not means where the component installed.
	// Host can be domain name or ip address
	Host string `yaml:"host"`
	IPv4 string `yaml:"ipv4"`
	IPv6 string `yaml:"ipv6"`

	// The port should be set to the actual port on which the process listend,
	// thus ansible playbook can use this variable to render configuration files.
	Port int `yaml:"port"`

	Addr string `yaml:"addr"`

	Endpoints []string `yaml:"endpoints"`
	URLs      []string `yaml:"urls"`

	PubPort int `yaml:"pubPort"`
	LBPort  int `yaml:"lbPort"`
}

// NewService returns a service
func NewService(componentName string, serviceName string) *Service {
	s := &Service{
		Name:          serviceName,
		ComponentName: componentName,
		Endpoints:     make([]string, 0),
		URLs:          make([]string, 0),
	}

	return s
}

func (s *Service) Check() error {
	errs := []error{}

	if s.Scheme == "" {
		msg := fmt.Sprintf("the scheme of service (%s) can not be empty", s.Name)
		errs = append(errs, errors.New(msg))
	}

	if s.Port == 0 {
		msg := fmt.Sprintf("the port of service (%s) can not be 0", s.Name)
		errs = append(errs, errors.New(msg))
	}
	if len(errs) == 0 {
		return nil
	}

	errmsgs := []string{}
	for _, err := range errs {
		errmsgs = append(errmsgs, err.Error())
	}

	msg := fmt.Sprintf("check service (%s) failed, err: %s", s.Name, strings.Join(errmsgs, "; "))
	return errors.New(msg)
}

// ServiceComputed represents the computed configuration for a service.
type ServiceComputed struct {
	Scheme    string   `yaml:"scheme"`
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	Addr      string   `yaml:"addr"`
	Endpoints []string `yaml:"endpoints"`
	URLs      []string `yaml:"urls"`
}

// NewServiceComputed returns a computed service.
func NewServiceComputed() *ServiceComputed {
	sc := &ServiceComputed{
		Endpoints: make([]string, 0),
		URLs:      make([]string, 0),
	}

	return sc
}

// Compute return a ServiceComputed for this service.
func (s *Service) Compute(external bool, cmdb *CMDB) (*ServiceComputed, error) {
	if external {
		return s.computeExternal()
	}

	return s.computeNonExternal(cmdb)
}

func (s *Service) computeNonExternal(cmdb *CMDB) (*ServiceComputed, error) {
	svcComputed := NewServiceComputed()

	scheme := ""
	if s.Scheme != "" {
		scheme = s.Scheme
	} else {
		scheme = "tcp"
	}
	svcComputed.Scheme = scheme

	host := ""
	if s.Host != "" {
		host = s.Host
	} else {
		if cmdb.Inventory.HasGroup(s.ComponentName) {
			g, err := cmdb.Inventory.GetGroup(s.ComponentName)
			if err != nil {
				return nil, fmt.Errorf("get component (%s) from inventory failed, err: %s", s.ComponentName, err)
			}
			hosts := g.HostsList()
			if len(hosts) > 0 {
				host = hosts[0]
			}
		}
	}
	if host == "" {
		host = "127.0.0.1"
	}
	svcComputed.Host = host

	var port int
	if s.PubPort != 0 {
		port = s.PubPort
	} else if s.LBPort != 0 {
		port = s.LBPort
	} else {
		port = s.Port
	}
	svcComputed.Port = port

	addr := ""
	if s.Addr != "" {
		addr = s.Addr
	} else {
		addr = fmt.Sprintf("%s:%d", svcComputed.Host, svcComputed.Port)
	}
	svcComputed.Addr = addr

	endpoints := []string{}
	if len(s.Endpoints) != 0 {
		endpoints = append(endpoints, s.Endpoints...)
	} else {
		if cmdb.Inventory.HasGroup(s.ComponentName) {
			g, err := cmdb.Inventory.GetGroup(s.ComponentName)
			if err != nil {
				return nil, fmt.Errorf("get component (%s) from inventory failed, err: %s", s.ComponentName, err)
			}
			hosts := g.HostsList()
			for _, host := range hosts {
				endpoint := fmt.Sprintf("%s:%d", host, svcComputed.Port)
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	if len(endpoints) == 0 {
		endpoints = append(endpoints, svcComputed.Addr)
	}
	svcComputed.Endpoints = endpoints

	urls := []string{}
	if len(s.URLs) != 0 {
		urls = append(urls, s.URLs...)
	} else {
		for _, v := range svcComputed.Endpoints {
			url := fmt.Sprintf("%s://%s", svcComputed.Scheme, v)
			urls = append(urls, url)
		}
	}
	svcComputed.URLs = urls

	return svcComputed, nil
}

// computeExternal returns ServiceComputed
func (s *Service) computeExternal() (*ServiceComputed, error) {
	svcComputed := NewServiceComputed()

	var scheme string
	if s.Scheme != "" {
		scheme = s.Scheme
	} else {
		scheme = "tcp"
	}
	svcComputed.Scheme = scheme

	var host string
	if s.Host != "" {
		host = s.Host
	} else {
		host = "127.0.0.1"
	}
	svcComputed.Host = host

	svcComputed.Port = s.Port

	var addr string
	if s.Addr != "" {
		addr = s.Addr
	} else {
		addr = fmt.Sprintf("%s:%d", svcComputed.Host, svcComputed.Port)
	}
	svcComputed.Addr = addr

	var endpoints []string
	if len(s.Endpoints) != 0 {
		endpoints = append(endpoints, s.Endpoints...)
	} else {
		endpoints = append(endpoints, svcComputed.Addr)
	}
	svcComputed.Endpoints = endpoints

	urls := []string{}
	if len(s.URLs) != 0 {
		urls = append(urls, s.URLs...)
	} else {
		for _, v := range svcComputed.Endpoints {
			url := fmt.Sprintf("%s://%s", svcComputed.Scheme, v)
			urls = append(urls, url)
		}
	}
	svcComputed.URLs = urls

	return svcComputed, nil
}
