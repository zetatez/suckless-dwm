package plugins

import (
	"fmt"
	"sync"

	"cmds/sugar"
)

func SearchBooksOnline() {
	content, err := sugar.GetInput("search books online: ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	urls := []string{
		"https://libgen.is/search.php?req=%s",
		"https:// openlibrary.org/search?q=%s",
		"https://z-lib.id/s?q=%s",
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
