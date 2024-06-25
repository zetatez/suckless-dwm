#!/bin/sh

go build -o ./bins/search                         ./cmds/search/main.go &
go build -o ./bins/handle-copied                  ./cmds/handle-copied/main.go &
go build -o ./bins/format-yaml                    ./cmds/format-yaml/main.go &
go build -o ./bins/format-json                    ./cmds/format-json/main.go &
go build -o ./bins/format-sql                     ./cmds/format-sql/main.go &
go build -o ./bins/ssh-to                         ./cmds/ssh-to/main.go &
go build -o ./bins/wifi-connect                   ./cmds/wifi-connect/main.go &
go build -o ./bins/jump-to-code-from-log          ./cmds/jump-to-code-from-log/main.go &
go build -o ./bins/launch-chrome                  ./cmds/launch-chrome/main.go &
go build -o ./bins/launch-edge                    ./cmds/launch-edge/main.go &
wait
go build -o ./bins/note-timeline                  ./cmds/note-timeline/main.go &
go build -o ./bins/note-diary                     ./cmds/note-diary/main.go &
go build -o ./bins/note-flash-card                ./cmds/note-flash-card/main.go &
go build -o ./bins/toggle-addressbook             ./cmds/toggle-addressbook/main.go &
go build -o ./bins/toggle-bluetooth               ./cmds/toggle-bluetooth/main.go &
go build -o ./bins/toggle-calendar-scheduling     ./cmds/toggle-calendar-scheduling/main.go &
go build -o ./bins/toggle-calendar-today-schedule ./cmds/toggle-calendar-today-schedule/main.go &
wait
go build -o ./bins/toggle-flameshot               ./cmds/toggle-flameshot/main.go &
go build -o ./bins/toggle-inkscape                ./cmds/toggle-inkscape/main.go &
go build -o ./bins/toggle-irssi                   ./cmds/toggle-irssi/main.go &
go build -o ./bins/toggle-joshuto                 ./cmds/toggle-joshuto/main.go &
go build -o ./bins/toggle-julia                   ./cmds/toggle-julia/main.go &
go build -o ./bins/toggle-keyboard-light          ./cmds/toggle-keyboard-light/main.go &
go build -o ./bins/toggle-lazydocker              ./cmds/toggle-lazydocker/main.go &
go build -o ./bins/toggle-music                   ./cmds/toggle-music/main.go &
wait
go build -o ./bins/toggle-music-net-cloud         ./cmds/toggle-music-net-cloud/main.go &
go build -o ./bins/toggle-music-yes-play-music    ./cmds/toggle-music-yes-play-music/main.go &
go build -o ./bins/toggle-mutt                    ./cmds/toggle-mutt/main.go &
go build -o ./bins/toggle-krita                   ./cmds/toggle-krita/main.go &
go build -o ./bins/toggle-python                  ./cmds/toggle-python/main.go &
go build -o ./bins/toggle-scala                   ./cmds/toggle-scala/main.go &
go build -o ./bins/toggle-lua                     ./cmds/toggle-lua/main.go &
go build -o ./bins/toggle-screen                  ./cmds/toggle-screen/main.go &
wait
go build -o ./bins/toggle-screenkey               ./cmds/toggle-screenkey/main.go &
go build -o ./bins/toggle-sublime                 ./cmds/toggle-sublime/main.go &
go build -o ./bins/toggle-show                    ./cmds/toggle-show/main.go &
go build -o ./bins/toggle-sys-shortcuts           ./cmds/toggle-sys-shortcuts/main.go &
go build -o ./bins/toggle-top                     ./cmds/toggle-top/main.go &
go build -o ./bins/toggle-wallpaper               ./cmds/toggle-wallpaper/main.go &
go build -o ./bins/toggle-clipmenu                ./cmds/toggle-clipmenu/main.go &
wait
go build -o ./bins/toggle-passmenu                ./cmds/toggle-passmenu/main.go &
go build -o ./bins/toggle-redshift                ./cmds/toggle-redshift/main.go &
go build -o ./bins/toggle-xournal                 ./cmds/toggle-xournal/main.go &
go build -o ./bins/toggle-rec-audio               ./cmds/toggle-rec-audio/main.go &
go build -o ./bins/toggle-rec-screen              ./cmds/toggle-rec-screen/main.go &
go build -o ./bins/toggle-rec-webcam              ./cmds/toggle-rec-webcam/main.go &
go build -o ./bins/toggle-obsidian                ./cmds/toggle-obsidian/main.go &
go build -o ./bins/toggle-termius                 ./cmds/toggle-termius/main.go &
go build -o ./bins/openweb-chatgpt                ./cmds/openweb-chatgpt/main.go &
go build -o ./bins/openweb-github                 ./cmds/openweb-github/main.go &
go build -o ./bins/openweb-google-mail            ./cmds/openweb-google-mail/main.go &
go build -o ./bins/openweb-google-translate       ./cmds/openweb-google-translate/main.go &
go build -o ./bins/openweb-instagram              ./cmds/openweb-instagram/main.go &
go build -o ./bins/openweb-leetcode               ./cmds/openweb-leetcode/main.go &
go build -o ./bins/openweb-youtube                ./cmds/openweb-youtube/main.go &
go build -o ./bins/openweb-wechat                 ./cmds/openweb-wechat/main.go &
go build -o ./bins/openweb-gemini                 ./cmds/openweb-gemini/main.go &
wait
