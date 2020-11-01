package polldance

import (
	"log"

	"github.com/k0kubun/pp"
)

type PollConfig struct {
	InputFilePaths   []string
	OutputHTTPURL    string
	OutputHTTPMethod string
	FilterCommand    string
	Debug            bool
}

func Poll(c *PollConfig) error {
	evp := NewEventProcessor()

	if c.Debug {
		log.Println("start polldance with debug flag")
		evp.AddHandler(func(e *EventData) error {
			pp.Println(e)
			return nil
		})
	}

	if c.FilterCommand != "" {
		fc := &FilterCommand{
			Command: c.FilterCommand,
			Debug:   c.Debug,
		}
		evp.AddHandler(fc.Handler)
	}

	if c.OutputHTTPURL != "" {
		oh := &OutHTTP{
			URLTmpl: c.OutputHTTPURL,
			Debug:   c.Debug,
			Method:  c.OutputHTTPMethod,
		}
		evp.AddHandler(oh.Handler)
	}

	fp := NewFilePoller(evp)
	for _, path := range c.InputFilePaths {
		log.Printf("File: %s", path)
		err := fp.AddFile(path)
		if err != nil {
			log.Fatalf("failed to init filePoller for %s: %s", path, err)
		}
	}

	evp.Wait()
	return nil
}
