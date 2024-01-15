default: clean build

build:
	mkdir -p ./bins/
	go build -o ./bins/format-json                   ./cmds/format-json/main.go &
	go build -o ./bins/format-sql                    ./cmds/format-sql/main.go &
	go build -o ./bins/handle-copied                 ./cmds/handle-copied/main.go &
	go build -o ./bins/search                        ./cmds/search/main.go &
	go build -o ./bins/toggle-addressbook            ./cmds/toggle-addressbook/main.go &
	go build -o ./bins/toggle-bluetooth              ./cmds/toggle-bluetooth/main.go &
	go build -o ./bins/toggle-calendar-schedule-day  ./cmds/toggle-calendar-schedule-day/main.go &
	go build -o ./bins/toggle-calendar-schedule-week ./cmds/toggle-calendar-schedule-week/main.go &
	go build -o ./bins/toggle-chrome                 ./cmds/toggle-chrome/main.go &
	go build -o ./bins/toggle-clipmenu               ./cmds/toggle-clipmenu/main.go &
	go build -o ./bins/toggle-diary                  ./cmds/toggle-diary/main.go &
	go build -o ./bins/toggle-flameshot              ./cmds/toggle-flameshot/main.go &
	go build -o ./bins/toggle-inkscape               ./cmds/toggle-inkscape/main.go &
	go build -o ./bins/toggle-joshuto                ./cmds/toggle-joshuto/main.go &
	go build -o ./bins/toggle-julia                  ./cmds/toggle-julia/main.go &
	go build -o ./bins/toggle-kb-light               ./cmds/toggle-kb-light/main.go &
	go build -o ./bins/toggle-lazydocker             ./cmds/toggle-lazydocker/main.go &
	go build -o ./bins/toggle-music                  ./cmds/toggle-music/main.go &
	go build -o ./bins/toggle-music-net-cloud        ./cmds/toggle-music-net-cloud/main.go &
	go build -o ./bins/toggle-mutt                   ./cmds/toggle-mutt/main.go &
	go build -o ./bins/toggle-passmenu               ./cmds/toggle-passmenu/main.go &
	go build -o ./bins/toggle-photoshop              ./cmds/toggle-photoshop/main.go &
	go build -o ./bins/toggle-python                 ./cmds/toggle-python/main.go &
	go build -o ./bins/toggle-rec-audio              ./cmds/toggle-rec-audio/main.go &
	go build -o ./bins/toggle-rec-video              ./cmds/toggle-rec-video/main.go &
	go build -o ./bins/toggle-redshift               ./cmds/toggle-redshift/main.go &
	go build -o ./bins/toggle-rss                    ./cmds/toggle-rss/main.go &
	go build -o ./bins/toggle-screen                 ./cmds/toggle-screen/main.go &
	go build -o ./bins/toggle-screenkey              ./cmds/toggle-screenkey/main.go &
	go build -o ./bins/toggle-show                   ./cmds/toggle-show/main.go &
	go build -o ./bins/toggle-sublime                ./cmds/toggle-sublime/main.go &
	go build -o ./bins/toggle-sys-shortcuts          ./cmds/toggle-sys-shortcuts/main.go &
	go build -o ./bins/toggle-top                    ./cmds/toggle-top/main.go &
	go build -o ./bins/toggle-wallpaper              ./cmds/toggle-wallpaper/main.go &
	go build -o ./bins/toggle-wechat                 ./cmds/toggle-wechat/main.go &
	go build -o ./bins/toggle-xournal                ./cmds/toggle-xournal/main.go &
	go build -o ./bins/wifi-connect                  ./cmds/wifi-connect/main.go &
	go build -o ./bins/ssh                           ./cmds/ssh/main.go &
	cp ./cmds/toggle-search-file-content/main        ./bins/toggle-search-file-content &
	cp ./cmds/toggle-git-show-log/main               ./bins/git-show-log &

clean:
	rm -rf bins/*

.PHONY: default

