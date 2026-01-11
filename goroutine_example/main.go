package main

import (
	"fmt"
	"goroutineexample/goroutinefolder"
	"os"
	"sort"
)

func main() {
	funcs := map[string]func(){
		"AtomicOperation":               goroutinefolder.AtomicOperation,
		"Blockselect":                   goroutinefolder.Blockselect,
		"Bufferedchannel":               goroutinefolder.Bufferedchannel,
		"Channelblockpush":              goroutinefolder.Channelblockpush,
		"Channelblockpull":              goroutinefolder.Channelblockpull,
		"Channelblockpullwait":          goroutinefolder.Channelblockpullwait,
		"Closingchannels":               goroutinefolder.Closingchannels,
		"ConcurrentCacheUnsafe":         goroutinefolder.ConcurrentCacheUnsafe,
		"ConcurrentCacheSafe":           goroutinefolder.ConcurrentCacheSafe,
		"ConcurrentCacheRWMutex":        goroutinefolder.ConcurrentCacheRWMutex,
		"ConcurrentCacheComparison":     goroutinefolder.ConcurrentCacheComparison,
		"ContextCancelTimeout":          goroutinefolder.ContextCancelTimeout,
		"Directions":                    goroutinefolder.Directions,
		"DownloadWithTimeout":           goroutinefolder.DownloadWithTimeout,
		"DownloadWithTimeoutV2":         goroutinefolder.DownloadWithTimeoutV2,
		"DownloaderWithAllFeatures":     goroutinefolder.DownloaderWithAllFeatures,
		"DownloaderWithShortTimeout":    goroutinefolder.DownloaderWithShortTimeout,
		"DownloaderWithLongTimeout":     goroutinefolder.DownloaderWithLongTimeout,
		"Foorloop":                      goroutinefolder.Foorloop,
		"MultiGoroutineoneval":          goroutinefolder.MultiGoroutineoneval,
		"PipelineFaninFanout":           goroutinefolder.PipelineFaninFanout,
		"SignalControlAndClosing":       goroutinefolder.SignalControlAndClosing,
		"Syncmutex":                     goroutinefolder.Syncmutex,
		"SyncMutex2":                    goroutinefolder.SyncMutex2,
		"WaitgroupMutex":                goroutinefolder.WaitgroupMutex,
		"WaitgroupMutexAdvanced":        goroutinefolder.WaitgroupMutexAdvanced,
		"WorkerPoolExample":             goroutinefolder.WorkerPoolExample,
		"WorkerPoolWithNames":           goroutinefolder.WorkerPoolWithNames,
		"WorkerPoolConfigurableDefault": goroutinefolder.WorkerPoolConfigurableDefault,
		"WorkerPoolLarge":               goroutinefolder.WorkerPoolLarge,
		"Workerpools":                   goroutinefolder.Workerpools,
		"WorkerpoolV2":                  goroutinefolder.WorkerpoolV2,
		"SumWithChannel":                goroutinefolder.SumWithChannel,
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <function_name>")
		fmt.Println("Available functions:")

		// Sort keys for consistent output
		var keys []string
		for k := range funcs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Println(" - " + k)
		}
		return
	}

	funcName := os.Args[1]
	if f, ok := funcs[funcName]; ok {
		fmt.Printf("Running %s...\n", funcName)
		f()
	} else {
		fmt.Printf("Function '%s' not found.\n", funcName)
	}
}
