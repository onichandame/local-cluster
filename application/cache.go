package application

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/chebyrash/promise"
	"github.com/google/uuid"
	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

const (
	tmpPrefix = "tmp-"
)

type CacheManager struct {
	lock   sync.Mutex
	caches map[uint]*promise.Promise
}

var manager *CacheManager
var managerInited int32

func getManager() *CacheManager {
	if manager == nil {
		manager = &CacheManager{}
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&manager)), nil, unsafe.Pointer(&CacheManager{caches: make(map[uint]*promise.Promise)}))
	}
	return manager
}

func cache(id uint, url string, hash string) (err error) {
	defer utils.RecoverFromError(&err)
	cachePath := filepath.Join(config.Config.Path.Cache, strconv.Itoa(int(id)))
	// only one routine can manipulate the manager at a moment
	start := func() (p *promise.Promise) {
		manager := getManager()
		manager.lock.Lock()
		defer manager.lock.Unlock()
		// check if already cached
		if utils.PathExists(cachePath) {
			if hash != "" {
				if err = utils.CheckFileHash(cachePath, hash); err != nil {
					logrus.Infof("cache for app %d is broken! will delete it", id)
					if err = os.RemoveAll(cachePath); err != nil {
						panic(err)
					} else {
						cache(id, url, hash)
					}
				}
			}
		} else {
			if p, ok := manager.caches[id]; ok {
				if _, err := p.Await(); err != nil {
					panic(err)
				}
			}
			p = promise.New(func(resolve func(promise.Any), reject func(error)) {
				logrus.Infof("downloading cache for app %d", id)
				tmpFilePath := newTmpFilePath()
				if err := utils.Download(url, tmpFilePath); err != nil {
					reject(err)
				}
				if hash != "" {
					if err := utils.CheckFileHash(tmpFilePath, hash); err != nil {
						reject(err)
					}
				}
				logrus.Infof("downloaded cache for app %d", id)
				delete(manager.caches, id)
				if err := os.Rename(tmpFilePath, cachePath); err != nil {
					reject(err)
				}
				resolve(nil)
			})
			manager.caches[id] = p
		}
		return p
	}
	p := start()
	_, err = p.Await()
	return err
}

func AuditCache() error {
	pattern := filepath.Join(config.Config.Path.Cache, tmpPrefix+"*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, tmp := range matches {
		if err := os.Remove(tmp); err != nil {
			return err
		}
	}
	return nil
}

func newTmpFilePath() string {
	salt, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return filepath.Join(config.Config.Path.Cache, tmpPrefix+salt.String())
}
