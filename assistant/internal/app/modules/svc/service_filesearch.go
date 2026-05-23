package svc

import (
	"fmt"
	"strings"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) SearchWeb(query string) error {
	url := "https://www.google.com/search?q=" + strings.ReplaceAll(query, " ", "+")
	return s.OpenURL("chrome", url)
}

func (s *Service) FileSearch() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e bash -c 'fd --type f --hidden --exclude /boot --exclude /bin --exclude /sbin --exclude /dev --exclude /lib --exclude /lib64 --exclude /lost+found --exclude /mnt --exclude /run --exclude /srv --exclude /sys --exclude /usr --exclude /var --exclude .git --exclude .cache --exclude .local . "/" | fzf --prompt="search file>" --preview "lazy -o view -f {}" --select-1 --exit-0 --print0 | xargs -0 -o lazy -o open -f'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) FileSearchContent() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e bash -c 'RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case"; matched=$(FZF_DEFAULT_COMMAND="$RG_PREFIX " fzf --prompt="search file content> " --preview "bat --style=full --color=always --highlight-line {2} {1}" --bind "change:reload:$RG_PREFIX {q} || true" --ansi --disabled --query "" --height=100%% --layout=reverse --delimiter :); filepath=$(echo "$matched" | awk -F: "{print \$1}"); rowno=$(echo "$matched" | awk -F: "{print \$2}"); if [ -n "$filepath" ]; then nvim +"$rowno" "$filepath"; fi'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) FileSearchBook() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e bash -c 'fd --type f --extension pdf --extension epub --extension djvu --extension mobi --exclude .git --exclude .cache --exclude .local --hidden . "$HOME" | fzf --prompt="search book>" --preview "lazy -o view -f {}" --select-1 --exit-0 --print0 | xargs -0 -o lazy -o open -f'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) FileSearchMedia() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e bash -c 'fd --type f --extension jpg --extension jpeg --extension png --extension gif --extension bmp --extension tiff --extension avi --extension flac --extension mkv --extension mp3 --extension mp4 --hidden --exclude .git --exclude .cache --exclude .local . "$HOME" | fzf --prompt="search media>" --preview "lazy -o view -f {}" --select-1 --exit-0 --print0 | xargs -0 -o lazy -o open -f'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) FileSearchWiki() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e bash -c 'fd --type f --extension md --hidden --exclude .git --exclude .cache --exclude .local . "$HOME" | fzf --prompt="search wiki>" --preview "lazy -o view -f {}" --select-1 --exit-0 --print0 | xargs -0 -o lazy -o open -f'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) FileOpenImages(dir string) error {
	if dir == "" {
		dir = "."
	}
	return startScript("bash", fmt.Sprintf("fd --type f --extension jpg --extension jpeg --extension png --extension gif --extension bmp --extension tiff --exclude repos . '%s' | sxiv -ftio", dir))
}
