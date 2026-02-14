default: clean build

prepare:
	mkdir -p ./bins/
	go mod tidy

build_shell_cmds:
	cp -f ./shell_scripts/* ./bins/

clean:
	rm -rf bins/*

build: prepare build_shell_cmds
	go build -o ./bins/daemon                                    ./cmds/daemon/main.go
	go build -o ./bins/format_go                                 ./cmds/format_go/main.go
	go build -o ./bins/format_json                               ./cmds/format_json/main.go
	go build -o ./bins/format_sql                                ./cmds/format_sql/main.go
	go build -o ./bins/format_yaml                               ./cmds/format_yaml/main.go
	go build -o ./bins/handle_copied                             ./cmds/handle_copied/main.go
	go build -o ./bins/launch_chrome                             ./cmds/launch_chrome/main.go
	go build -o ./bins/launch_file_manager                       ./cmds/launch_file_manager/main.go
	go build -o ./bins/launch_inkscape                           ./cmds/launch_inkscape/main.go
	go build -o ./bins/launch_krita                              ./cmds/launch_krita/main.go
	go build -o ./bins/launch_obsidian                           ./cmds/launch_obsidian/main.go
	go build -o ./bins/launch_qutebrowser                        ./cmds/launch_qutebrowser/main.go
	go build -o ./bins/launch_sublime                            ./cmds/launch_sublime/main.go
	go build -o ./bins/launch_xournal                            ./cmds/launch_xournal/main.go
	go build -o ./bins/note_monthly_work                         ./cmds/note_monthly_work/main.go
	go build -o ./bins/note_scripts                              ./cmds/note_scripts/main.go
	go build -o ./bins/note_todo                                 ./cmds/note_todo/main.go
	go build -o ./bins/open_url_as_app                           ./cmds/open_url_as_app/main.go
	go build -o ./bins/open_url_with_chrome                      ./cmds/open_url_with_chrome/main.go
	go build -o ./bins/open_url_with_qutebrowser                 ./cmds/open_url_with_qutebrowser/main.go
	go build -o ./bins/search                                    ./cmds/search/main.go
	go build -o ./bins/snip_fzf                                  ./cmds/snip_fzf/main.go
	go build -o ./bins/snip_create                               ./cmds/snip_create/main.go
	go build -o ./bins/ssh_to                                    ./cmds/ssh_to/main.go
	go build -o ./bins/sys_bluetooth_connect                     ./cmds/sys_bluetooth_connect/main.go
	go build -o ./bins/sys_bluetooth_disconnect                  ./cmds/sys_bluetooth_disconnect/main.go
	go build -o ./bins/sys_bluetooth_scan_and_connect            ./cmds/sys_bluetooth_scan_and_connect/main.go
	go build -o ./bins/sys_display                               ./cmds/sys_display/main.go
	go build -o ./bins/sys_shortcuts                             ./cmds/sys_shortcuts/main.go
	go build -o ./bins/sys_toggle_keyboard_light                 ./cmds/sys_toggle_keyboard_light/main.go
	go build -o ./bins/sys_wallpaper                             ./cmds/sys_wallpaper/main.go
	go build -o ./bins/sys_wifi_connect                          ./cmds/sys_wifi_connect/main.go
	go build -o ./bins/toggle_calendar                           ./cmds/toggle_calendar/main.go
	go build -o ./bins/toggle_calendar_scheduling                ./cmds/toggle_calendar_scheduling/main.go
	go build -o ./bins/toggle_calendar_scheduling_today          ./cmds/toggle_calendar_scheduling_today/main.go
	go build -o ./bins/toggle_clipmenu                           ./cmds/toggle_clipmenu/main.go
	go build -o ./bins/toggle_flameshot                          ./cmds/toggle_flameshot/main.go
	go build -o ./bins/toggle_irssi                              ./cmds/toggle_irssi/main.go
	go build -o ./bins/toggle_julia                              ./cmds/toggle_julia/main.go
	go build -o ./bins/toggle_keyboard                           ./cmds/toggle_keyboard/main.go
	go build -o ./bins/toggle_lazydocker                         ./cmds/toggle_lazydocker/main.go
	go build -o ./bins/toggle_lazygit                            ./cmds/toggle_lazygit/main.go
	go build -o ./bins/toggle_music                              ./cmds/toggle_music/main.go
	go build -o ./bins/toggle_mutt                               ./cmds/toggle_mutt/main.go
	go build -o ./bins/toggle_netease_cloud_music                ./cmds/toggle_netease_cloud_music/main.go
	go build -o ./bins/toggle_passmenu                           ./cmds/toggle_passmenu/main.go
	go build -o ./bins/toggle_python                             ./cmds/toggle_python/main.go
	go build -o ./bins/toggle_rec_audio                          ./cmds/toggle_rec_audio/main.go
	go build -o ./bins/toggle_rec_screen                         ./cmds/toggle_rec_screen/main.go
	go build -o ./bins/toggle_rec_show                           ./cmds/toggle_rec_show/main.go
	go build -o ./bins/toggle_rec_webcam                         ./cmds/toggle_rec_webcam/main.go
	go build -o ./bins/toggle_scala                              ./cmds/toggle_scala/main.go
	go build -o ./bins/toggle_screenkey                          ./cmds/toggle_screenkey/main.go
	go build -o ./bins/toggle_top                                ./cmds/toggle_top/main.go
	go build -o ./bins/toggle_tty_clock                          ./cmds/toggle_tty_clock/main.go
	go build -o ./bins/toggle_yazi                               ./cmds/toggle_yazi/main.go
	chmod +x ./bins/*
	@echo "Done!"

.PHONY: default
