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

	Addr      string   `yaml:"addr"`
	Endpoints []string `yaml:"endpoints"`
	URLs      []string `yaml:"urls"`

	PubPort int `yaml:"pub_port"`
	LBPort  int `yaml:"lb_port"`
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

func (s *Service) Compute(external bool, cmdb *CMDB) (*ServiceComputed, error) {
	if external {
		return s.computeExternal()
	}

	return s.computeNonExternal(cmdb)
}

func (s *Service) computeNonExternal(cmdb *CMDB) (*ServiceComputed, error) {
	svcComputed := NewServiceComputed()
	svcComputed.Scheme = s.Scheme
	if svcComputed.Scheme == "" {
		svcComputed.Scheme = "tcp"
	}

	var host string
	if s.Host != "" {
		host = s.Host
	} else if cmdb.Inventory.HasGroup(s.ComponentName) {
		g, _ := cmdb.Inventory.GetGroup(s.ComponentName)
		hosts := g.HostsList()
		if len(hosts) > 0 {
			host = hosts[0]
		} else {
			host = "127.0.0.1"
		}
	} else {
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

	if s.Addr != "" {
		svcComputed.Addr = s.Addr
	} else {
		svcComputed.Addr = fmt.Sprintf("%s:%d", svcComputed.Host, svcComputed.Port)
	}

	if len(s.Endpoints) != 0 {
		for _, v := range s.Endpoints {
			svcComputed.Endpoints = append(s.Endpoints, v)
		}
	} else {
		svcComputed.Endpoints = append(svcComputed.Endpoints, svcComputed.Addr)
	}

	if len(s.URLs) != 0 {
		svcComputed.URLs = append(svcComputed.URLs, s.URLs...)
	} else {
		for _, v := range svcComputed.Endpoints {
			url := fmt.Sprintf("%s://%s", svcComputed.Scheme, v)
			svcComputed.URLs = append(svcComputed.URLs, url)
		}
	}

	return svcComputed, nil
}

func (s *Service) computeExternal() (*ServiceComputed, error) {
	svcComputed := NewServiceComputed()

	svcComputed.Scheme = s.Scheme

	var host string
	if s.Host != "" {
		host = s.Host
	} else {
		host = "127.0.0.1"
	}
	svcComputed.Host = host

	svcComputed.Port = s.Port

	if s.Addr != "" {
		svcComputed.Addr = s.Addr
	} else {
		svcComputed.Addr = fmt.Sprintf("%s:%d", svcComputed.Host, svcComputed.Port)
	}

	if len(s.Endpoints) != 0 {
		for _, v := range s.Endpoints {
			svcComputed.Endpoints = append(s.Endpoints, v)
		}
	} else {
		svcComputed.Endpoints = append(svcComputed.Endpoints, svcComputed.Addr)
	}

	if len(s.URLs) != 0 {
		svcComputed.URLs = append(svcComputed.URLs, s.URLs...)
	} else {
		for _, v := range svcComputed.Endpoints {
			url := fmt.Sprintf("%s://%s", svcComputed.Scheme, v)
			svcComputed.URLs = append(svcComputed.URLs, url)
		}
	}

	return svcComputed, nil
}
