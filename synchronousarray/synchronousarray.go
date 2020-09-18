package synchronous_array

import (
	"sync"

	"github.com/cirocosta/go-mod-dep-source-finder/result"
)

type SynchronousArray struct {
	Results []result.Result
	sync.Mutex
}

func (array *SynchronousArray) Add(resultToAdd result.Result) {
	array.Lock()
	defer array.Unlock()

	array.Results = append(array.Results, resultToAdd)
}
