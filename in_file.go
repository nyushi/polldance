package polldance

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FilePoller struct {
	evp     *EventProcessor
	watcher *fsnotify.Watcher
}

func NewFilePoller(evp *EventProcessor) *FilePoller {
	fp := &FilePoller{
		evp: evp,
	}
	fp.Start()
	return fp
}

func (f *FilePoller) Start() error {
	f.evp.AddSource()
	w, err := fsnotify.NewWatcher()
	f.watcher = w
	if err != nil {
		return fmt.Errorf("erorr at NewWatcher: %w", err)
	}

	go f.watchLoop()
	return nil
}

func (f *FilePoller) AddFile(path string) error {
	if err := f.watcher.Add(path); err != nil {
		return fmt.Errorf("error at watcher.Add: %w", err)
	}
	ev, err := f.newEventFromFile(path)
	if err != nil {
		return fmt.Errorf("initial event not created: %s", err)
	}
	f.evp.Push(ev)
	return nil
}

func (f *FilePoller) newEventFromFile(path string) (*EventData, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error at ReadFile: %w", err)
	}
	return &EventData{
		Source: path,
		Data:   string(b),
	}, nil

}
func (f *FilePoller) watchLoop() {
	defer f.evp.RemoveSource()
	defer f.watcher.Close()
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				log.Printf("no new event")
				return
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				go func() {
					for {
						if err := f.AddFile(event.Name); err == nil {
							return
						}
						time.Sleep(time.Second)
					}
				}()
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				ev, err := f.newEventFromFile(event.Name)
				if err != nil {
					log.Printf("error at write event: %s", err)
					continue
				}
				f.evp.Push(ev)
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				log.Printf("no new watcher error")
				return
			}
			log.Printf("error from watcher: %s", err)
		}
	}
}
