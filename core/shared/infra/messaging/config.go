package messaging

import "fmt"

type Config struct {
	Servers      string
	Group        string
	OffsetReset  string
	ConsumeTopic []string
	HandlerName  string
	DLQ          bool
}

func (c *Config) Validate() error {
	if len(c.Servers) == 0 {
		return fmt.Errorf("server cant empty")
	}

	if len(c.ConsumeTopic) == 0 {
		return fmt.Errorf("do not have any topic")
	}

	return nil
}
