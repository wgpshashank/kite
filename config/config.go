package config

import (
	"net/url"
	"os"

	"github.com/koding/kite/kitekey"
)

// Options is passed to kite.New when creating new instance.
type Config struct {
	// Options for Kite
	Username              string
	Environment           string
	Region                string
	KiteKey               string
	DisableAuthentication bool

	// Options for Server
	IP                 string
	Port               int
	Path               string
	DisableConcurrency bool

	KontrolURL  *url.URL
	KontrolKey  string
	KontrolUser string
}

var defaultConfig = &Config{
	Username:    "unknown",
	Environment: "unknown",
	Region:      "unknown",
	IP:          "0.0.0.0",
	Port:        0,
	// KontrolURL:  &url.URL{Scheme: "ws", Host: "localhost:4000"},
}

// New returns a new Config initialized with defaults.
func New() *Config {
	c := new(Config)
	*c = *defaultConfig
	return c
}

func (c *Config) ReadEnvironmentVariables() {
	if environment := os.Getenv("KITE_ENVIRONMENT"); environment == "" {
		c.Environment = environment
	}

	if region := os.Getenv("KITE_REGION"); region == "" {
		c.Region = region
	}
}

// ReadKiteKey parsed the user's kite key and returns a new Config.
func (c *Config) ReadKiteKey() error {
	key, err := kitekey.Parse()
	if err != nil {
		return err
	}

	c.KiteKey = key.Raw

	if username, ok := key.Claims["sub"].(string); ok {
		c.Username = username
	}

	if kontrolUser, ok := key.Claims["iss"].(string); ok {
		c.KontrolUser = kontrolUser
	}

	if kontrolURL, ok := key.Claims["kontrolURL"].(string); ok {
		c.KontrolURL, err = url.Parse(kontrolURL)
		if err != nil {
			return err
		}
	}

	if kontrolKey, ok := key.Claims["kontrolKey"].(string); ok {
		c.KontrolKey = kontrolKey
	}

	return nil
}

// validate fields of the options struct. It exits if an error is occured.
func (o *Config) validate() {
	// if o.PublicIP == "" {
	// 	o.PublicIP = "0.0.0.0"
	// }

	// if o.Port == "" {
	// 	o.Port = "0" // OS binds to an automatic port
	// }

	// if o.Path == "" {
	// 	o.Path = "/kite/" + o.Kitename
	// }

	// if o.Path[0] != '/' {
	// 	o.Path = "/" + o.Path
	// }
}

// 		// 	k.TrustKontrolKey(k.kiteKey.Claims["iss"].(string), k.kiteKey.Claims["kontrolKey"].(string))
// 		// }
// 	}