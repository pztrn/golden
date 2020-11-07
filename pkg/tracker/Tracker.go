package tracker

import (
	"github.com/vit1251/golden/pkg/charset"
	"github.com/vit1251/golden/pkg/echomail"
	"github.com/vit1251/golden/pkg/file"
	"github.com/vit1251/golden/pkg/mailer/cache"
	"github.com/vit1251/golden/pkg/netmail"
	"github.com/vit1251/golden/pkg/registry"
	"github.com/vit1251/golden/pkg/setup"
	"github.com/vit1251/golden/pkg/stat"
	"log"
	"os"
	"path"
	"time"
)

type Tracker struct {
	registry *registry.Container
}

func NewTracker(r *registry.Container) *Tracker {
	newTracker := new(Tracker)
	newTracker.registry = r
	return newTracker
}

func (self Tracker) Track() {

	trackerStart := time.Now()
	log.Printf("Start tracker session")

	err1 := self.ProcessInbound()
	if err1 != nil {
		log.Printf("err = %+v", err1)
	}
	err2 := self.ProcessOutbound()
	if err2 != nil {
		log.Printf("err = %+v", err2)
	}

	log.Printf("Stop tracker session")
	elapsed := time.Since(trackerStart)

	log.Printf("Tracker session: %+v", elapsed)
}

func (self *Tracker) ProcessInbound() error {

	/* New mailer inbound */
	mi := cache.NewMailerInbound(self.registry)

	/* Scan inbound */
	items, err2 := mi.Scan()
	if err2 != nil {
		return err2
	}
	log.Printf("items = %+v", items)

	for _, item := range items {
		if item.Type == cache.TypeTICmail {
			log.Printf("Tracker: TIC packet: name = %s", item.Name)
			if err := self.processTICmail(item); err != nil {
				log.Printf("Tracker: process TIC with error: err = %+v", err)
			}
		} else {
			// TODO - message about skip ...
		}
	}

	return nil
}

func (self *Tracker) ProcessOutbound() error {
	return nil
}

func (self *Tracker) processTICmail(item *cache.FileEntry) error {

	fileManager := self.restoreFileManager()
	configManager := self.restoreConfigManager()
	statManager := self.restoreStatManager()

	/* Parse */
	newTicParser := NewTicParser(self.registry)
	tic, err1 := newTicParser.ParseFile(item.AbsolutePath)
	if err1 != nil {
		return err1
	}
	log.Printf("tic = %+v", tic)

	areaName := tic.GetArea()

	/* Search area */
	fa, err1 := fileManager.GetAreaByName(areaName)
	if err1 != nil {
		return err1
	}

	/* Prepare area directory */
	boxBasePath, _ := configManager.Get("main", "FileBox")
	inboxBasePath, _ := configManager.Get("main", "Inbound")

	areaLocation := path.Join(boxBasePath, areaName)
	os.MkdirAll(areaLocation, 0755)

	/* Create area */
	if fa == nil {
		/* Prepare area */
		newFa := file.NewFileArea()
		newFa.SetName(areaName)
		newFa.Path = areaLocation
		/* Create area */
		if err := fileManager.CreateFileArea(newFa); err != nil {
			log.Printf("Fail CreateFileArea on FileManager: area = %s err = %+v", areaName, err)
			return err
		}
	}

	/* Create new path */
	inboxTicLocation := path.Join(inboxBasePath, tic.File)
	areaFileLocation := path.Join(areaLocation, tic.File)
	log.Printf("inboxTicLocation = %s areaFileLocation = %s", inboxTicLocation, areaFileLocation)

	/* Move */
	os.Rename(inboxTicLocation, areaFileLocation)

	/* Register file */
	newFile := file.NewFile()
	newFile.SetArea(tic.GetArea())
	newFile.SetDesc(tic.Desc)
	newFile.SetUnixTime(tic.UnixTime)
	newFile.SetFile(tic.File)
	fileManager.RegisterFile(*newFile)

	/* Register status */
	statManager.RegisterInFile(tic.File)

	/* Move TIC */
	areaTicLocation := path.Join(areaLocation, item.Name)
	log.Printf("areaTicLocation = %s", areaTicLocation)
	os.Rename(item.AbsolutePath, areaTicLocation)

	return nil
}

func (self Tracker) restoreCharsetManager() *charset.CharsetManager {
	managerPtr := self.registry.Get("CharsetManager")
	if manager, ok := managerPtr.(*charset.CharsetManager); ok {
		return manager
	} else {
		panic("no charset manager")
	}
}

func (self Tracker) restoreNetmailManager() *netmail.NetmailManager {
	managerPtr := self.registry.Get("NetmailManager")
	if manager, ok := managerPtr.(*netmail.NetmailManager); ok {
		return manager
	} else {
		panic("no netmail manager")
	}
}

func (self Tracker) restoreAreaManager() *echomail.AreaManager {
	managerPtr := self.registry.Get("AreaManager")
	if manager, ok := managerPtr.(*echomail.AreaManager); ok {
		return manager
	} else {
		panic("no area manager")
	}
}

func (self Tracker) restoreMessageManager() *echomail.MessageManager {
	managerPtr := self.registry.Get("MessageManager")
	if manager, ok := managerPtr.(*echomail.MessageManager); ok {
		return manager
	} else {
		panic("no message manager")
	}
}

func (self Tracker) restoreStatManager() *stat.StatManager {
	managerPtr := self.registry.Get("StatManager")
	if manager, ok := managerPtr.(*stat.StatManager); ok {
		return manager
	} else {
		panic("no stat manager")
	}
}

func (self Tracker) restoreConfigManager() *setup.ConfigManager {
	managerPtr := self.registry.Get("ConfigManager")
	if manager, ok := managerPtr.(*setup.ConfigManager); ok {
		return manager
	} else {
		panic("no config manager")
	}
}

func (self Tracker) restoreFileManager() *file.FileManager {
	managerPtr := self.registry.Get("FileManager")
	if manager, ok := managerPtr.(*file.FileManager); ok {
		return manager
	} else {
		panic("no file manager")
	}
}