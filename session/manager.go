// manager
package session

import (
	"fmt"
	"gosws/config"
)

var (
	sessions map[string]SessionFactory = map[string]SessionFactory{}
)

func Register(name string, sf SessionFactory) error {
	sessions[name] = sf
	return nil
}

func Get(name string) SessionFactory {
	sf, ok := sessions[name]
	if !ok {
		panic(fmt.Sprintf("not found %s SessionFactory", name))
	}
	return sf
}

func Init() error {
	confSession := config.Session()
	if confSession.Sessname == "" {
		return nil
	}

	secretKey := config.SecretKey()
	name := confSession.Sessname
	lifetime := confSession.Lifetime

	for _, c := range confSession.Config {
		if c.Name == name {
			return Get(name).Init(lifetime, secretKey, c.Jsonconfig)
		}
	}

	panic(fmt.Sprintf("not found session type:%s, config, please check conf file", name))

	return nil
}
