package file

import (
	"sync"
)

type Service struct {
	mu    *sync.RWMutex
	cache map[string]*AssetFile
}

func NewService(filePaths ...string) *Service {
	var filesMap = make(map[string]*AssetFile)

	for _, path := range filePaths {
		assetFile := NewAssetFile(path)
		filesMap[path] = assetFile
	}

	return &Service{
		mu:    &sync.RWMutex{},
		cache: filesMap,
	}
}

func (f *Service) GetAssetFile(path string) (*AssetFile, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.getFile(path)
}

func (f *Service) getFile(path string) (*AssetFile, error) {
	if file, exists := f.cache[path]; exists {
		err := file.Open()
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	assetF := NewAssetFile(path)
	f.cache[path] = assetF

	err := assetF.Open()
	if err != nil {
		return nil, err
	}

	return assetF, nil
}
