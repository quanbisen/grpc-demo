package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
)

type Server struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Version string `json:"version"`
	Weight  int    `json:"weight"`
}

func BuildPrefix(server Server) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.Name)
	}
	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
}

func BuildRegisterPath(server Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(server), server.Address)
}

func ParseValue(value []byte) (Server, error) {
	s := Server{}
	if err := json.Unmarshal(value, &s); err != nil {
		return s, err
	}
	return s, nil
}

func SplitPath(path string) (Server, error) {
	s := Server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return s, errors.New("invalid path")
	}
	s.Address = strs[len(strs)-1]
	return s, nil
}

func Exist(l []resolver.Address, address resolver.Address) bool {
	for i := range l {
		if l[i].Addr == address.Addr {
			return true
		}
	}
	return false
}

// Remove helper function
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
