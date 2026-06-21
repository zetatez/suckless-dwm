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

### Super Layer – Function Row

| Key            | Action                                                  |
| ---            | ---                                                     |
| `SUP+F1/F2/F3` | Toggle, lower, or raise system volume (`sys_volume_*`)  |
| `SUP+F4`       | Toggle microphone mute (`sys_micro_toggle`)             |
| `SUP+F5/F6`    | Dim or brighten the panel (`sys_display_light_*`)       |
| `SUP+F7`       | Launch display controls (`sys_display`)                 |
| `SUP+F8`       | Open Wi-Fi helper (`sys_wifi_connect`)                  |
| `SUP+F9`       | Open Bluetooth helper (`sys_bluetooth_connect`)         |
| `SUP+F10/F11`  | Adjust microphone levels down/up (`sys_micro_down/up`)  |
| `SUP+F12`      | Toggle keyboard backlight (`sys_toggle_keyboard_light`) |

### Super Layer – Browser Presets

The number row opens curated URLs in Chrome via `open_url_with_chrome`; holding `Shift` launches the same URL as a dedicated app via `open_url_as_app`.

| Key     | Destination                           |
| ---     | ---                                   |
| `SUP+1` | Google Gemini                         |
| `SUP+2` | GitHub profile (`github.com/zetatez`) |
| `SUP+3` | Gmail                                 |
| `SUP+4` | Google Translate                      |
| `SUP+5` | CCTV5 stream                          |
| `SUP+6` | Web WeChat                            |
| `SUP+7` | YouTube                               |
| `SUP+8` | Bilibili                              |
| `SUP+9` | LeetCode CN search                    |
| `SUP+0` | Doubao chat                           |
| `SUP+/` | Toggle opencode                       |

> `SUP+Shift+1..0` opens the same URLs in qutebrowser instead of Chrome.

### Super Layer – Launchers & Toggles

| Key     | Command                    | Purpose                                   |
| ---     | ---                        | ---                                       |
| `SUP+a` | `thunar`                   | Open the file manager                     |
| `SUP+b` | `qutebrowser`              | Launch qutebrowser                        |
| `SUP+c` | `note monthly-work`        | Append to the monthly work log            |
| `SUP+d` | `lazydocker` (toggle)      | Show or hide lazydocker                   |
| `SUP+g` | `chrome`                   | Launch Google Chrome                      |
| `SUP+i` | `flameshot gui` (toggle)   | Toggle Flameshot screenshot UI            |
| `SUP+m` | `file-search-content`      | Search file contents                      |
| `SUP+n` | Python scratchpad (toggle) | Toggle the Python scratchpad terminal     |
| `SUP+o` | `handle-clipboard`         | Process clipboard content                 |
| `SUP+p` | `file-search`              | Search files from a terminal prompt       |
| `SUP+q` | `slock`                    | Lock the screen                           |
| `SUP+r` | `yazi` (toggle)            | Show or hide the `yazi` TUI file manager  |
| `SUP+s` | `search`                   | Invoke the custom search interface        |
| `SUP+t` | `next_theme`               | Cycle through the configured color themes |
| `SUP+u` | `solve-leetcode`           | Solve LeetCode problems                   |
| `SUP+v` | `send-to-feishu`           | Send content to Feishu                    |
| `SUP+w` | `sys-kill`                 | Kill windows/processes                    |
| `SUP+x` | `note scripts`             | Jump to the scripts notebook              |
| `SUP+y` | `toggle-rec-show`          | Toggle recording/show overlay             |
| `SUP+z` | `note todo`                | Open the todo capture note                |

| Key             | Command                            | Purpose                                    |
| ---             | ---                                | ---                                        |
| `SUP+Backspace` | `passmenu`                         | Display the password picker                |
| `SUP+Delete`    | `sys-shortcut`                     | Show the global shortcut helper            |
| `SUP+Escape`    | `htop` (toggle)                    | Toggle the system monitor overlay          |
| `SUP+'`         | `toggle-tty-clock`                 | Show or hide the fullscreen terminal clock |
| `SUP+\`         | `sys-reset`                        | Reset desktop defaults                     |
| `SUP+[`         | `nvim +Calendar -view=month`       | Show the weekly calendar scheduling view   |
| `SUP+]`         | `nvim +Calendar -view=day`         | Show today's calendar scheduling view      |

### Super Layer – Floating Window Controls

| Key             | Action                                              |
| ---             | ---                                                 |
| `SUP+h/j/k/l`   | Move the focused floating client left/down/up/right |
| `SUP+Shift+h/l` | Shrink or expand the client width                   |
| `SUP+Shift+j/k` | Shrink or expand the client height                  |

### Super Layer – Shifted Shortcuts

| Key           | Action                                 |
| ---           | ---                                    |
| `SUP+Shift+1` | Open Google Gemini in qutebrowser      |
| `SUP+Shift+2` | Open GitHub in qutebrowser             |
| `SUP+Shift+3` | Open Gmail in qutebrowser              |
| `SUP+Shift+4` | Open Google Translate in qutebrowser   |
| `SUP+Shift+5` | Open CCTV5 in qutebrowser              |
| `SUP+Shift+6` | Open Web WeChat in qutebrowser         |
| `SUP+Shift+7` | Open YouTube in qutebrowser            |
| `SUP+Shift+8` | Open Bilibili in qutebrowser           |
| `SUP+Shift+9` | Open LeetCode in qutebrowser           |
| `SUP+Shift+0` | Open Doubao in qutebrowser             |
| `SUP+Shift+/` | Toggle opencode                        |
| `SUP+Shift+d` | Launch/focus DingTalk                  |
| `SUP+Shift+e` | Toggle the terminal mail client (mutt) |
| `SUP+Shift+f` | Launch/focus Feishu                    |
| `SUP+Shift+i` | Launch/focus Inkscape                  |
| `SUP+Shift+m` | Launch/focus NetEase Cloud Music       |
| `SUP+Shift+n` | Toggle the Julia scratchpad            |
| `SUP+Shift+o` | Launch/focus Obsidian                  |
| `SUP+Shift+s` | Launch/focus Sublime Text              |
| `SUP+Shift+u` | Solve LeetCode screenshot              |
| `SUP+Shift+w` | Send to Feishu                         |
| `SUP+Shift+x` | Launch/focus Xournal++                 |
| `SUP+Shift+z` | Launch/focus Zoom                      |
| `SUP+Shift+'` | Toggle Screenkey overlay               |
| `SUP+Shift+,` | Toggle audio recording                 |
| `SUP+Shift+.` | Toggle screen recording                |

### Alt Layer – Core Controls

| Key                | Action                                           |
| ---                | ---                                              |
| `MOD+Return`       | Promote the focused client to master (`zoom`)    |
| `MOD+Tab`          | Cycle to the previously viewed tag (`view`)      |
| `MOD+'`            | Toggle the primary scratchpad (`sp-st`)          |
| `MOD+;`            | Run `rofi -show run` in fullscreen preview theme |
| `MOD+/`            | Launch `snip_fzf`                                |
| `MOD+Shift+/`      | Launch `snip_create`                             |
| `MOD+b`            | Toggle the status bar                            |
| `MOD+c`            | Toggle the clipboard manager                     |
| `MOD+f`            | Toggle fullscreen                                |
| `MOD+o`            | Enter or leave the overview layout               |
| `MOD+p`            | Launch the dmenu/rofi application launcher       |
| `MOD+q`            | Lock the session via `slock`                     |
| `MOD+s`            | Reset layouts and factors (`reset`)              |
| `MOD+u`            | Jump to the selected client in the stack         |
| `MOD+Shift+Return` | Spawn a vanilla `st` terminal                    |
| `MOD+Shift+'`      | Move any scratchpad back to normal tiling        |
| `MOD+Shift+c`      | Kill the focused client                          |
| `MOD+Shift+f`      | Toggle floating (duplicate for convenience)      |
| `MOD+Shift+p`      | Restart dwm (`quit 1`)                           |
| `MOD+Shift+q`      | Quit dwm (`quit 0`)                              |
| `MOD+Shift+s`      | Toggle sticky state                              |
| `MOD+Shift+Space`  | Focus the master area                            |
| `MOD+Shift+u`      | Restore the previous layout (temporary switch)   |
| `MOD+Shift+Ctrl+c` | Kill every unfocused client                      |

### Alt Layer – Layout Selection

| Key           | Layout                   |
| ---           | ---                      |
| `MOD+a`       | Workflow `[W]`           |
| `MOD+r`       | Fibonacci spiral `[F]`   |
| `MOD+Shift+r` | Fibonacci dwindle `[F]`  |
| `MOD+v`       | Center free shape `[C]`  |
| `MOD+Shift+v` | Center equal ratio `[C]` |
| `MOD+t`       | Tile right `[T]`         |
| `MOD+Shift+t` | Tile left `[T]`          |
| `MOD+g`       | Grid `[G]`               |
| `MOD+Shift+g` | Grid with gaps `[G]`     |
| `MOD+m`       | Monocle `[M]`            |
| `MOD+w`       | Hacker `[H]`             |
| `MOD+e`       | Stack horizontal `[S]`   |
| `MOD+Shift+e` | Stack vertical `[S]`     |

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
| Client window (`MOD` held) | Left / middle / right | Move / toggle floating / resize the client                   |

This list mirrors the current `keys[]` and `buttons[]` definitions; update it whenever `config.def.h` changes so the README stays authoritative.
