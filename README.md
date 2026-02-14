[dwm](https://dwm.suckless.org/)

![](https://github.com/zetatez/suckless-dwm/blob/master/dwm.png)

# dwm - dynamic window manager
dwm is an extremely fast, small, and dynamic window manager for X.


## Requirements
In order to build dwm you need the Xlib header files.


## Installation
Edit config.mk to match your local setup (dwm is installed into
the /usr/local namespace by default).

Afterwards enter the following command to build and install dwm (if
necessary as root):

    sh build.sh

## Running dwm
Add the following line to your .xinitrc to start dwm using startx:

    exec dwm

In order to connect dwm to a specific display, make sure that
the DISPLAY environment variable is set correctly, e.g.:

    DISPLAY=foo.bar:1 exec dwm

(This will start dwm on display :1 of the host foo.bar.)

In order to display status info in the bar, you can use dwmblocks

    git clone https://github.com/zetatez/suckless-dwmblocks.git
    cd suckless-dwmblocks && sh build.sh

## Configuration
The configuration of dwm is done by creating a custom config.h
and (re)compiling the source code.

## Key Bindings

`SUPKEY` maps to `Mod4Mask` (Super/Win) and `MODKEY` maps to `Mod1Mask` (Alt). All shortcuts below are defined in `config.def.h` and rely on those two modifiers unless noted.

### Hardware Media Keys (no modifier)

 | Key                                             | Action                                    |
 | ---                                             | ---                                       |
 | `XF86AudioMute`                                 | Toggle speaker mute via `pactl`           |
 | `XF86AudioLowerVolume` / `XF86AudioRaiseVolume` | Adjust sink volume by ±5%                 |
 | `XF86AudioMicMute`                              | Toggle microphone mute                    |
 | `XF86MonBrightnessDown` / `XF86MonBrightnessUp` | Adjust display brightness by ±5%          |
 | `XF86AudioPlay` / `XF86AudioPause`              | Toggle media playback through `playerctl` |
 | `XF86AudioStop`                                 | Stop playback via `playerctl stop`        |
 | `XF86AudioPrev` / `XF86AudioNext`               | Previous/next track through `playerctl`   |

### Super Layer – Function Row

 | Key              | Action                                                             |
 | ---              | ---                                                                |
 | `SUP+F1/F2/F3`   | Toggle, lower, or raise system volume (`sys_volume_*`)             |
 | `SUP+F4/F10/F11` | Toggle microphone mute or adjust microphone levels (`sys_micro_*`) |
 | `SUP+F5/F6`      | Dim or brighten the panel (`sys_display_light_*`)                  |
 | `SUP+F7`         | Launch display controls (`sys_display`)                            |
 | `SUP+F8`         | Open Wi-Fi helper (`sys_wifi_connect`)                             |
 | `SUP+F9`         | Open Bluetooth helper (`sys_bluetooth_connect`)                    |
 | `SUP+F12`        | Toggle keyboard backlight (`sys_toggle_keyboard_light`)            |

### Super Layer – Browser Presets

The number row opens curated URLs in qutebrowser; holding `Shift` launches the same URL in Google Chrome.

 | Key     | Destination                           |
 | ---     | ---                                   |
 | `SUP+1` | ChatGPT                               |
 | `SUP+2` | YouTube                               |
 | `SUP+3` | GitHub profile (`github.com/zetatez`) |
 | `SUP+4` | Gmail                                 |
 | `SUP+5` | Google Translate                      |
 | `SUP+6` | Web WeChat                            |
 | `SUP+7` | LeetCode CN search                    |
 | `SUP+8` | CCTV5 stream                          |
 | `SUP+9` | Bilibili                              |
 | `SUP+0` | Doubao chat                           |
 | `SUP+/` | Google Gemini                         |

### Super Layer – Launchers & Toggles

 | Key               | Command                         | Purpose                                                 |
 | ---               | ---                             | ---                                                     |
 | `SUP+a`           | `launch_file_manager`           | Open the default file manager                           |
 | `SUP+b`           | `launch_qutebrowser`            | Focus or start qutebrowser                              |
 | `SUP+c`           | `note_monthly_work`             | Append to the monthly work log                          |
 | `SUP+d`           | `toggle_lazydocker`             | Show or hide lazydocker                                 |
 | `SUP+e`           | `toggle_mutt`                   | Toggle the terminal mail client                         |
 | `SUP+f`           | `lazy_open_search_file`         | Search files from a terminal prompt                     |
 | `SUP+g`           | `st_dir_fzf_lazygit`            | Jump into lazygit via fzf                               |
 | `SUP+i`           | `toggle_flameshot`              | Toggle Flameshot screenshot UI                          |
 | `SUP+m`           | `lazy_open_search_file_content` | Search file contents                                    |
 | `SUP+n`           | Scratchpad (Python)             | Toggle the `sp-python` scratchpad terminal              |
 | `SUP+o`           | `handle_copied`                 | Process clipboard content                               |
 | `SUP+p`           | `lazy_open_search_book`         | Search in the book library                              |
 | `SUP+q`           | `slock`                         | Lock the screen                                         |
 | `SUP+r`           | `toggle_yazi`                   | Show or hide the `yazi` TUI file manager                |
 | `SUP+s`           | `search`                        | Invoke the custom search interface                      |
 | `SUP+t`           | `next_theme`                    | Cycle through the configured color themes               |
 | `SUP+u`           | `lazy_open_search_media`        | Search media assets                                     |
 | `SUP+v`           | `st_dir_fzf_opencode`           | Open the OpenCode workspace selector                    |
 | `SUP+w`           | `lazy_open_search_wiki`         | Search the wiki knowledge base                          |
 | `SUP+x`           | `note_scripts`                  | Jump to the scripts notebook                            |
 | `SUP+y`           | `toggle_show`                   | Toggle on-screen widgets for streaming/presentation     |
 | `SUP+z`           | `note_todo`                     | Open the todo capture note                              |

 | Key               | Command                         | Purpose                                                 |
 | ---               | ---                             | ---                                                     |
 | `SUP+Backspace`   | `toggle_passmenu`               | Display the password picker                             |
 | `SUP+Delete`      | `sys_shortcuts`                 | Show the global shortcut helper                         |
 | `SUP+Escape`      | `toggle_top`                    | Toggle the system monitor overlay                       |
 | `SUP+'`           | `toggle_tty_clock`              | Show or hide the fullscreen terminal clock              |
 | `SUP+\`           | `reset_sys_default`             | Reset desktop defaults                                  |
 | `SUP+[` / `SUP+]` | `toggle_calendar_scheduling(*)` | Show the weekly or today view of the calendar scheduler |
 | `SUP+.`           | `jump_to_code_from_log`         | Jump from log output to source code                     |

### Super Layer – Floating Window Controls

| Key | Action |
| --- | --- |
| `SUP+h/j/k/l` | Move the focused floating client left/down/up/right |
| `SUP+Shift+h/l` | Shrink or expand the client width |
| `SUP+Shift+j/k` | Shrink or expand the client height |

### Super Layer – Shifted Shortcuts

 | Key                | Action                                          |
 | ---                | ---                                             |
 | `SUP+Shift+a`      | Launch `gamescope` in Steam Big Picture mode    |
 | `SUP+Shift+b`      | Launch or focus Google Chrome                   |
 | `SUP+Shift+c`      | Kill the focused client                         |
 | `SUP+Shift+d`      | Launch/focus DingTalk                           |
 | `SUP+Shift+f`      | Launch/focus Feishu                             |
 | `SUP+Shift+i`      | Launch/focus Inkscape                           |
 | `SUP+Shift+m`      | Launch/focus NetEase Cloud Music                |
 | `SUP+Shift+n`      | Toggle the Julia scratchpad (`sp-julia`)        |
 | `SUP+Shift+o`      | Launch/focus Obsidian                           |
 | `SUP+Shift+p`      | Launch/focus Krita                              |
 | `SUP+Shift+s`      | Launch/focus Sublime Text                       |
 | `SUP+Shift+x`      | Launch/focus Xournal++                          |
 | `SUP+Shift+z`      | Launch/focus Zoom                               |
 | `SUP+Shift+Delete` | Immediately power off via `systemctl poweroff`  |
 | `SUP+Shift+Return` | Open an fzf-powered `st` session (`st_dir_fzf`) |
 | `SUP+Shift+'`      | Toggle Screenkey overlay                        |
 | `SUP+Shift+,`      | Toggle audio recording                          |
 | `SUP+Shift+.`      | Toggle screen recording                         |

### Alt Layer – Core Controls

 | Key                     | Action                                           |
 | ---                     | ---                                              |
 | `MOD+Return`            | Promote the focused client to master (`zoom`)    |
 | `MOD+Tab`               | Cycle to the previously viewed tag (`view`)      |
 | `MOD+'`                 | Toggle the primary scratchpad (`sp-st`)          |
 | `MOD+;`                 | Run `rofi -show run` in fullscreen preview theme |
 | `MOD+/` / `MOD+Shift+/` | Launch `snip_fzf` / `snip_create`                |
 | `MOD+Space`             | Toggle floating for the focused client           |
 | `MOD+b`                 | Toggle the status bar                            |
 | `MOD+c`                 | Toggle the clipboard manager                     |
 | `MOD+f`                 | Toggle fullscreen                                |
 | `MOD+o`                 | Enter or leave the overview layout               |
 | `MOD+p`                 | Launch the dmenu/rofi application launcher       |
 | `MOD+q`                 | Lock the session via `slock`                     |
 | `MOD+s`                 | Reset layouts and factors (`reset`)              |
 | `MOD+u`                 | Jump to the selected client in the stack         |
 | `MOD+Shift+Return`      | Spawn a vanilla `st` terminal                    |
 | `MOD+Shift+'`           | Move any scratchpad back to normal tiling        |
 | `MOD+Shift+c`           | Kill the focused client                          |
 | `MOD+Shift+f`           | Toggle floating (duplicate for convenience)      |
 | `MOD+Shift+p`           | Restart dwm (`quit 1`)                           |
 | `MOD+Shift+q`           | Quit dwm (`quit 0`)                              |
 | `MOD+Shift+s`           | Toggle sticky state                              |
 | `MOD+Shift+Space`       | Focus the master area                            |
 | `MOD+Shift+u`           | Restore the previous layout (temporary switch)   |
 | `MOD+Shift+Ctrl+c`      | Kill every unfocused client                      |

### Alt Layer – Layout Selection

 | Key                     | Layout                                 |
 | ---                     | ---                                    |
 | `MOD+a`                 | Workflow `[W]`                         |
 | `MOD+r` / `MOD+Shift+r` | Fibonacci spiral / dwindle             |
 | `MOD+v` / `MOD+Shift+v` | Center free shape / center equal ratio |
 | `MOD+t` / `MOD+Shift+t` | Tile right / tile left                 |
 | `MOD+g` / `MOD+Shift+g` | Grid / grid with gaps                  |
 | `MOD+m`                 | Monocle                                |
 | `MOD+w`                 | Hacker layout                          |
 | `MOD+e` / `MOD+Shift+e` | Stack horizontal / stack vertical      |

### Alt Layer – Window & Stack Management

 | Key                           | Action                                         |
 | ---                           | ---                                            |
 | `MOD+d` / `MOD+i`             | Decrease/increase the number of master clients |
 | `MOD+h` / `MOD+l`             | Move the focused client within the stack       |
 | `MOD+,` / `MOD+.`             | Cycle layouts backward/forward                 |
 | `MOD+Shift+,` / `MOD+Shift+.` | Shift the visible tag set backward/forward     |
 | `MOD+k` / `MOD+j`             | Focus the previous/next client in the stack    |
 | `MOD+Shift+h` / `MOD+Shift+l` | Shrink/grow master width factor (`mfact`)      |
 | `MOD+Shift+j` / `MOD+Shift+k` | Shrink/grow stack height factor (`hfact`)      |

### Alt Layer – Tag Management

 | Key                   | Action                                         |
 | ---                   | ---                                            |
 | `MOD+1..9`            | View the corresponding tag                     |
 | `MOD+0`               | View every tag (all workspace preview)         |
 | `MOD+Shift+1..9`      | Move the focused client to a tag               |
 | `MOD+Shift+0`         | Tag the client with every tag                  |
 | `MOD+Ctrl+1..9`       | Toggle the visibility of individual tags       |
 | `MOD+Shift+Ctrl+1..9` | Preview a tag without switching (`previewtag`) |

### Monitor Focus & Tag Transfer

 | Key               | Action                                               |
 | ---               | ---                                                  |
 | `MOD+[ / ]`       | Focus previous/next monitor                          |
 | `MOD+Shift+[ / ]` | Send the focused client to the previous/next monitor |

### Mouse Bindings

 | Target                     | Button                | Action                                                       |
 | ---                        | ---                   | ---                                                          |
 | Tag bar                    | Left click            | View the clicked tag                                         |
 | Tag bar                    | Right click           | Toggle the tag visibility                                    |
 | Tag bar (`MOD` held)       | Left / right          | Assign the focused client to the tag / toggle tag assignment |
 | Layout symbol              | Left / middle         | Cycle layouts / jump directly to grid-gap layout             |
 | Layout symbol              | Right                 | Activate the overview layout                                 |
 | Client window (`MOD` held) | Left / middle / right | Move / toggle floating / resize the client                   |

This list mirrors the current `keys[]` and `buttons[]` definitions; update it whenever `config.def.h` changes so the README stays authoritative.
