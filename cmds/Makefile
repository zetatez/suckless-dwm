default: clean build

prepare:
	mkdir -p ./bins/
	go mod tidy

build-shell-cmds:
	cp -f ./shell-cmds/* ./bins/

clean:
	rm -rf bins/*

build: prepare build-shell-cmds
	go build -o ./bins/format-json                               ./cmds/format-json/main.go
	go build -o ./bins/format-sql                                ./cmds/format-sql/main.go
	go build -o ./bins/format-yaml                               ./cmds/format-yaml/main.go
	go build -o ./bins/handle-copied                             ./cmds/handle-copied/main.go
	go build -o ./bins/jump-to-code-from-log                     ./cmds/jump-to-code-from-log/main.go
	go build -o ./bins/launch-chrome                             ./cmds/launch-chrome/main.go
	go build -o ./bins/launch-qutebrowser                        ./cmds/launch-qutebrowser/main.go
	go build -o ./bins/note-diary                                ./cmds/note-diary/main.go
	go build -o ./bins/note-flash-card                           ./cmds/note-flash-card/main.go
	go build -o ./bins/note-timeline                             ./cmds/note-timeline/main.go
	go build -o ./bins/search                                    ./cmds/search/main.go
	go build -o ./bins/ssh-to                                    ./cmds/ssh-to/main.go
	go build -o ./bins/sys-bluetooth                             ./cmds/sys-bluetooth/main.go
	go build -o ./bins/sys-screen                                ./cmds/sys-screen/main.go
	go build -o ./bins/sys-shortcuts                             ./cmds/sys-shortcuts/main.go
	go build -o ./bins/sys-toggle-keyboard-light                 ./cmds/sys-toggle-keyboard-light/main.go
	go build -o ./bins/sys-wifi-connect                          ./cmds/sys-wifi-connect/main.go
	go build -o ./bins/toggle-addressbook                        ./cmds/toggle-addressbook/main.go
	go build -o ./bins/toggle-calendar                           ./cmds/toggle-calendar/main.go
	go build -o ./bins/toggle-calendar-scheduling                ./cmds/toggle-calendar-scheduling/main.go
	go build -o ./bins/toggle-calendar-scheduling-today          ./cmds/toggle-calendar-scheduling-today/main.go
	go build -o ./bins/toggle-clipmenu                           ./cmds/toggle-clipmenu/main.go
	go build -o ./bins/toggle-flameshot                          ./cmds/toggle-flameshot/main.go
	go build -o ./bins/toggle-inkscape                           ./cmds/toggle-inkscape/main.go
	go build -o ./bins/toggle-irssi                              ./cmds/toggle-irssi/main.go
	go build -o ./bins/toggle-julia                              ./cmds/toggle-julia/main.go
	go build -o ./bins/toggle-krita                              ./cmds/toggle-krita/main.go
	go build -o ./bins/toggle-lazydocker                         ./cmds/toggle-lazydocker/main.go
	go build -o ./bins/toggle-lazygit                            ./cmds/toggle-lazygit/main.go
	go build -o ./bins/toggle-lua                                ./cmds/toggle-lua/main.go
	go build -o ./bins/toggle-music                              ./cmds/toggle-music/main.go
	go build -o ./bins/toggle-music-net-cloud                    ./cmds/toggle-music-net-cloud/main.go
	go build -o ./bins/toggle-mutt                               ./cmds/toggle-mutt/main.go
	go build -o ./bins/toggle-obsidian                           ./cmds/toggle-obsidian/main.go
	go build -o ./bins/toggle-passmenu                           ./cmds/toggle-passmenu/main.go
	go build -o ./bins/toggle-python                             ./cmds/toggle-python/main.go
	go build -o ./bins/toggle-rec-audio                          ./cmds/toggle-rec-audio/main.go
	go build -o ./bins/toggle-rec-screen                         ./cmds/toggle-rec-screen/main.go
	go build -o ./bins/toggle-rec-webcam                         ./cmds/toggle-rec-webcam/main.go
	go build -o ./bins/toggle-redshift                           ./cmds/toggle-redshift/main.go
	go build -o ./bins/toggle-scala                              ./cmds/toggle-scala/main.go
	go build -o ./bins/toggle-screenkey                          ./cmds/toggle-screenkey/main.go
	go build -o ./bins/toggle-show                               ./cmds/toggle-show/main.go
	go build -o ./bins/toggle-sublime                            ./cmds/toggle-sublime/main.go
	go build -o ./bins/toggle-top                                ./cmds/toggle-top/main.go
	go build -o ./bins/toggle-wallpaper                          ./cmds/toggle-wallpaper/main.go
	go build -o ./bins/toggle-xournal                            ./cmds/toggle-xournal/main.go
	go build -o ./bins/toggle-yazi                               ./cmds/toggle-yazi/main.go
	go build -o ./bins/chrome-open-url-chatgpt                   ./cmds/chrome-open-url-chatgpt/main.go
	go build -o ./bins/chrome-open-url-doubao                    ./cmds/chrome-open-url-doubao/main.go
	go build -o ./bins/chrome-open-url-github                    ./cmds/chrome-open-url-github/main.go
	go build -o ./bins/chrome-open-url-google                    ./cmds/chrome-open-url-google/main.go
	go build -o ./bins/chrome-open-url-google-mail               ./cmds/chrome-open-url-google-mail/main.go
	go build -o ./bins/chrome-open-url-google-translate          ./cmds/chrome-open-url-google-translate/main.go
	go build -o ./bins/chrome-open-url-instagram                 ./cmds/chrome-open-url-instagram/main.go
	go build -o ./bins/chrome-open-url-leetcode                  ./cmds/chrome-open-url-leetcode/main.go
	go build -o ./bins/chrome-open-url-wechat                    ./cmds/chrome-open-url-wechat/main.go
	go build -o ./bins/chrome-open-url-youtube                   ./cmds/chrome-open-url-youtube/main.go
	go build -o ./bins/qutebrowser-open-url-chatgpt              ./cmds/qutebrowser-open-url-chatgpt/main.go
	go build -o ./bins/qutebrowser-open-url-doubao               ./cmds/qutebrowser-open-url-doubao/main.go
	go build -o ./bins/qutebrowser-open-url-github               ./cmds/qutebrowser-open-url-github/main.go
	go build -o ./bins/qutebrowser-open-url-google               ./cmds/qutebrowser-open-url-google/main.go
	go build -o ./bins/qutebrowser-open-url-google-mail          ./cmds/qutebrowser-open-url-google-mail/main.go
	go build -o ./bins/qutebrowser-open-url-google-translate     ./cmds/qutebrowser-open-url-google-translate/main.go
	go build -o ./bins/qutebrowser-open-url-instagram            ./cmds/qutebrowser-open-url-instagram/main.go
	go build -o ./bins/qutebrowser-open-url-leetcode             ./cmds/qutebrowser-open-url-leetcode/main.go
	go build -o ./bins/qutebrowser-open-url-wechat               ./cmds/qutebrowser-open-url-wechat/main.go
	go build -o ./bins/qutebrowser-open-url-youtube              ./cmds/qutebrowser-open-url-youtube/main.go
	@echo "Done!"

.PHONY: default
