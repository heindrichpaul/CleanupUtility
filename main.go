package main

import (
	"flag"
	"math/big"
	"os"
	"runtime"
	"sync"

	"github.com/heindrichpaul/cleanupUtil/cleanupTasks"
	"github.com/heindrichpaul/cleanupUtil/config"
	"github.com/heindrichpaul/cleanupUtil/logger"
	"github.com/stretchr/powerwalk"
)

var (
	wG             sync.WaitGroup
	total          big.Float
	configFilename string
	l              *logger.Logger
	c              *config.Config
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configFilename, "configFile", "config.json", "This configures the application with the source path for the walker and until which date to delete.")
	flag.Parse()

	c = config.NewConfig(configFilename)
	l = logger.NewLogger(c.LogLocation)

	wG.Add(len(c.ModulesToClean))

	for index, cleanup := range c.ModulesToClean {
		go cleanupModule(&wG, index, cleanup)
	}
	wG.Wait()
	total.Quo(&total, big.NewFloat(float64(1024*1024*1024)))
	l.Printf("Done deleting files. %s GB deleted.\n", total.Text('f', 8))
}

func cleanupModule(waitGroup *sync.WaitGroup, workerIndex int, cleanup *cleanupTasks.CleanupTask) {
	date, et := cleanup.GetTime(c.GetLocation())
	if et == nil {
		l.Printf("[Worker %d] Cleaning all files in %s before %s\n", workerIndex, cleanup.Directory, date)
		if _, e := os.Stat(cleanup.Directory); !os.IsNotExist(e) {
			err := powerwalk.Walk(cleanup.Directory,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						if info.ModTime().Before(date) {
							l.Printf("[Worker %d] Deleting %s created on %s\n", workerIndex, path, info.ModTime())
							if err1 := os.Remove(path); err1 == nil {
								l.Printf("[Worker %d] Deleted %s, saving %d bytes\n", workerIndex, path, info.Size())
								total.Add(&total, big.NewFloat(float64(info.Size())))
							} else {
								l.Printf("[Worker %d] Failed deleting %s (%s)\n", workerIndex, path, err1)
							}
						}
					} else {
						if err1 := os.Remove(path); err1 == nil {
							l.Printf("[Worker %d] Deleted empty folder %s\n", workerIndex, path)
						}
					}
					return nil
				})
			if err != nil {
				l.Printf("[Worker %d] %s\n", workerIndex, err)
			}

			err = powerwalk.Walk(cleanup.Directory,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						if err1 := os.Remove(path); err1 == nil {
							l.Printf("[Worker %d] Empty folder %s created on %s\n", workerIndex, path, info.ModTime())
							l.Printf("[Worker %d] Deleted %s\n", workerIndex, path)
						}
					}

					return nil
				})
			if err != nil {
				l.Printf("[Worker %d] %s\n", workerIndex, err)
			}
		} else {
			l.Printf("[Worker %d] The folder [%s] does not exist\n", workerIndex, cleanup.Directory)
		}
	}
	waitGroup.Done()
}
