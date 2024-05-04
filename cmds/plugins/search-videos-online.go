package plugins

import (
	"fmt"
	"sync"

	"cmds/sugar"
)

func SearchVideosOnline() {
	content, err := sugar.GetInput("search videos online: ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	urls := []string{
		"https://search.bilibili.com/all?keyword=%s",
		"https://www.youtube.com/results?search_query=%v",
	}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sugar.NewExecService().RunScriptShell(
				fmt.Sprintf(
					"chrome --proxy-server=%s %s",
					ProxyServer,
					fmt.Sprintf(url, content),
				),
			)
		}(url)
	}
	wg.Wait()
}
