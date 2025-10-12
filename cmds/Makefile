default: clean build

prepare:
	mkdir -p ./bins/
	go mod tidy

build_shell_cmds:
	cp -f ./shell_scripts/* ./bins/

clean:
	rm -rf bins/*

build: prepare build_shell_cmds
	go build -o ./bins/search                                    ./cmds/search/main.go
	go build -o ./bins/format_json                               ./cmds/format_json/main.go
	go build -o ./bins/format_sql                                ./cmds/format_sql/main.go
	go build -o ./bins/format_yaml                               ./cmds/format_yaml/main.go
	go build -o ./bins/handle_copied                             ./cmds/handle_copied/main.go
	go build -o ./bins/jump_to_code_from_log                     ./cmds/jump_to_code_from_log/main.go
	go build -o ./bins/launch_chrome                             ./cmds/launch_chrome/main.go
	go build -o ./bins/launch_qutebrowser                        ./cmds/launch_qutebrowser/main.go
	go build -o ./bins/launch_file_manager                       ./cmds/launch_file_manager/main.go
	go build -o ./bins/note_scripts                              ./cmds/note_scripts/main.go
	go build -o ./bins/note_todo                                 ./cmds/note_todo/main.go
	go build -o ./bins/note_diary                                ./cmds/note_diary/main.go
	go build -o ./bins/note_flash_card                           ./cmds/note_flash_card/main.go
	go build -o ./bins/note_timeline                             ./cmds/note_timeline/main.go
	go build -o ./bins/ssh_to                                    ./cmds/ssh_to/main.go
	go build -o ./bins/sys_bluetooth_connect                     ./cmds/sys_bluetooth_connect/main.go
	go build -o ./bins/sys_bluetooth_disconnect                  ./cmds/sys_bluetooth_disconnect/main.go
	go build -o ./bins/sys_bluetooth_scan_and_connect            ./cmds/sys_bluetooth_scan_and_connect/main.go
	go build -o ./bins/sys_screen                                ./cmds/sys_screen/main.go
	go build -o ./bins/sys_shortcuts                             ./cmds/sys_shortcuts/main.go
	go build -o ./bins/sys_toggle_keyboard_light                 ./cmds/sys_toggle_keyboard_light/main.go
	go build -o ./bins/sys_wifi_connect                          ./cmds/sys_wifi_connect/main.go
	go build -o ./bins/toggle_addressbook                        ./cmds/toggle_addressbook/main.go
	go build -o ./bins/toggle_calendar                           ./cmds/toggle_calendar/main.go
	go build -o ./bins/toggle_calendar_scheduling                ./cmds/toggle_calendar_scheduling/main.go
	go build -o ./bins/toggle_calendar_scheduling_today          ./cmds/toggle_calendar_scheduling_today/main.go
	go build -o ./bins/toggle_clipmenu                           ./cmds/toggle_clipmenu/main.go
	go build -o ./bins/toggle_flameshot                          ./cmds/toggle_flameshot/main.go
	go build -o ./bins/toggle_inkscape                           ./cmds/toggle_inkscape/main.go
	go build -o ./bins/toggle_irssi                              ./cmds/toggle_irssi/main.go
	go build -o ./bins/toggle_julia                              ./cmds/toggle_julia/main.go
	go build -o ./bins/toggle_krita                              ./cmds/toggle_krita/main.go
	go build -o ./bins/toggle_lazydocker                         ./cmds/toggle_lazydocker/main.go
	go build -o ./bins/toggle_lazygit                            ./cmds/toggle_lazygit/main.go
	go build -o ./bins/toggle_lua                                ./cmds/toggle_lua/main.go
	go build -o ./bins/toggle_music                              ./cmds/toggle_music/main.go
	go build -o ./bins/toggle_music_net_cloud                    ./cmds/toggle_music_net_cloud/main.go
	go build -o ./bins/toggle_mutt                               ./cmds/toggle_mutt/main.go
	go build -o ./bins/toggle_obsidian                           ./cmds/toggle_obsidian/main.go
	go build -o ./bins/toggle_passmenu                           ./cmds/toggle_passmenu/main.go
	go build -o ./bins/toggle_python                             ./cmds/toggle_python/main.go
	go build -o ./bins/toggle_rec_audio                          ./cmds/toggle_rec_audio/main.go
	go build -o ./bins/toggle_rec_screen                         ./cmds/toggle_rec_screen/main.go
	go build -o ./bins/toggle_rec_webcam                         ./cmds/toggle_rec_webcam/main.go
	go build -o ./bins/toggle_redshift                           ./cmds/toggle_redshift/main.go
	go build -o ./bins/toggle_scala                              ./cmds/toggle_scala/main.go
	go build -o ./bins/toggle_screenkey                          ./cmds/toggle_screenkey/main.go
	go build -o ./bins/toggle_show                               ./cmds/toggle_show/main.go
	go build -o ./bins/toggle_sublime                            ./cmds/toggle_sublime/main.go
	go build -o ./bins/toggle_top                                ./cmds/toggle_top/main.go
	go build -o ./bins/toggle_tty_clock                          ./cmds/toggle_tty_clock/main.go
	go build -o ./bins/toggle_xournal                            ./cmds/toggle_xournal/main.go
	go build -o ./bins/toggle_yazi                               ./cmds/toggle_yazi/main.go
	go build -o ./bins/open_url_with_chrome                      ./cmds/open_url_with_chrome/main.go
	go build -o ./bins/open_url_with_qutebrowser                 ./cmds/open_url_with_qutebrowser/main.go
	go build -o ./bins/open_url_as_app                           ./cmds/open_url_as_app/main.go
	go build -o ./bins/wallpaper                                 ./cmds/wallpaper/main.go
	go build -o ./bins/daemon                                    ./cmds/daemon/main.go
	chmod +x ./bins/*
	@echo "Done!"

.PHONY: default
