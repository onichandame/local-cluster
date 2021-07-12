package application

import (
	"os"
	"path/filepath"
	"sync"

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

type CacheStateMap struct {
	lock     sync.Mutex
	stateMap map[uint]*promise.Promise
}

var cacheStateMap CacheStateMap

func PrepareCache(appDef *model.Application) error {
	cachePath := GetCachePath(appDef)
	spec, err := GetSpec(appDef)
	if err != nil {
		return err
	}
	// only one routine can manipulate the map at a moment
	cacheStateMap.lock.Lock()
	defer func() { cacheStateMap.lock.Unlock() }()
	if cacheStateMap.stateMap == nil {
		cacheStateMap.stateMap = make(map[uint]*promise.Promise)
	}
	// check if another process is caching or it is already cached
	if utils.PathExists(cachePath) {
		if spec.Hash != "" {
			if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
				logrus.Infof("cache for app %d is broken! will delete it", appDef.ID)
				if err := os.Remove(cachePath); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if p, ok := cacheStateMap.stateMap[appDef.ID]; ok {
		_, err := p.Await()
		return err
	}
	cacheStateMap.stateMap[appDef.ID] = promise.New(func(resolve func(promise.Any), reject func(error)) {
		logrus.Infof("downloading cache for app %s", appDef.Name)
		tmpFilePath := newTmpFilePath()
		if err := utils.Download(spec.DownloadUrl, tmpFilePath); err != nil {
			reject(err)
		}
		if spec.Hash != "" {
			if err := utils.CheckFileHash(tmpFilePath, spec.Hash); err != nil {
				reject(err)
			}
		}
		logrus.Infof("downloaded cache for app %s", appDef.Name)
		delete(cacheStateMap.stateMap, appDef.ID)
		if err := os.Rename(tmpFilePath, cachePath); err != nil {
			reject(err)
		}
		resolve(nil)
	})
	return nil
}

func WaitCache(appDef *model.Application) error {
	if p, ok := cacheStateMap.stateMap[appDef.ID]; ok {
		_, err := p.Await()
		return err
	}
	return nil
}

func AuditCache() error {
	pattern := filepath.Join(config.ConfigPresets.CacheDir, tmpPrefix+"*")
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
	return filepath.Join(config.ConfigPresets.CacheDir, tmpPrefix+salt.String())
}
