/* See LICENSE file for copyright and license details. */

#include <X11/XF86keysym.h>

#define SUPKEY  Mod4Mask
#define MODKEY  Mod1Mask
#define Spawn(cmd)             (const char *[]){cmd, NULL}
#define CmdClass(cmd, cls)     (const char *[]){cmd, cls, NULL}
#define SpawnShellCmd(cmd)     (const char *[]){"/bin/sh", "-c", cmd, NULL}
#define SpawnTermiCmd(cmd)     (const char *[]){"st", "-e", "/bin/sh", "-c", cmd, NULL}

/* appearance */
static const unsigned int borderpx = 2;
static const unsigned int snap     = 24;
static const int scalepreview      = 3;
static const int previewbar        = 1;
static const int showbar           = 1;
static const int topbar            = 1;
static const int barheight         = 18;
static const int vertpad           = 0;
static const int sidepad           = 0;
static const int defaultwinpad     = 0;
static const int swallowfloating   = 1;
static const char *fonts[]         = {
  "VictorMono Nerd Font:style=Medium:size=14",
  "Source Han Serif CN:style=Regular:size=14",
  "Noto Color Emoji:style=Regular:size=12",
};
static const char dmenufont[]      = "VictorMono Nerd Font:style=Medium:size=24";


static const char *colors[][3] = {
  [SchemeNorm] = { "#0077b6", "#023047", "#023047" },
  [SchemeSel]  = { "#00b4d8", "#023047", "#00b4d8" },
};
static int current_theme_idx = 0;
static const char *themes[][SchemeLast][3] = {
  /*                  fg,        bg,        border */
  {
    [SchemeNorm] = { "#0077b6", "#023047", "#023047" },
    [SchemeSel]  = { "#00b4d8", "#023047", "#00b4d8" },
  },
  {
    [SchemeNorm] = { "#118ab2", "#031b34", "#031b34" },
    [SchemeSel]  = { "#06d6a0", "#031b34", "#06d6a0" },
  },
  {
    [SchemeNorm] = { "#f4a261", "#2b1400", "#2b1400" },
    [SchemeSel]  = { "#e76f51", "#2b1400", "#e76f51" },
  },
};

static const char *const autostart[] = {
  "reset_sys_default", NULL,
  "slock", NULL,
  "daemon", NULL,
  NULL,
};

/* tagging */
static const char *tags[] = { "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix" };

static const Rule rules[] = {
  /* cls                     instance    title      tags mask     isfloating    isterminal     noswallow     isontop  monitor */
  {"float",                  NULL,       NULL,      0,            1,            0,             0,            0,       -1 },
  {"st",                     NULL,       NULL,      0,            0,            1,             1,            0,       -1 },
  {"Surf",                   NULL,       NULL,      0,            0,            0,             1,            0,       -1 }, // no swallow for markdown
  {"Onboard",                NULL,       NULL,      0,            1,            0,             0,            1,       -1 },
//{"netease-cloud-music",    NULL,       NULL,      1<<8,         0,            0,             0,            0,       -1 },
};

/* stickyicon */
static const XPoint stickyicon[]    = { { 0, 0 }, { 4, 0 }, { 4, 8 }, { 2, 6 }, { 0, 8 }, { 0, 0 } }; /* stickyicon: represents the icon as an array of vertices */
static const XPoint stickyiconbb    = {4,8};	                                      /* stickyicon: defines the bottom right corner of the polygon's bounding box (speeds up scaling) */

/* layout(s) */
static const float mfact            = 0.50;
static const float hfact            = 0.50;
static const int nmaster            = 1;
static const int maxnmaster         = 16;
static const int resizehints        = 0;
static const int lockfullscreen     = 0;
static const int refreshrate        = 120;

static const Layout layouts[] = {
  { "W 󰴈"                  , layout_workflow           }, // default
  { "F "                  , layout_fib_spiral         },
  { "F "                  , layout_fib_dwindle        },
  { "C ⧈"                  , layout_center_free_shape  },
  { "C ⧅"                  , layout_center_equal_ratio },
  { "T 󱂫"                  , layout_tile_right         }, // 5
  { "T 󱂪"                  , layout_tile_left          },
  { "G 󰝘"                  , layout_grid               },
  { "G 󱇙"                  , layout_grid_gap           },
  { "M 󱣴"                  , layout_monocle            },
  { "H "                  , layout_hacker             }, // 10
  { "S |"                  , layout_stack_vert         },
  { "S ―"                  , layout_stack_hori         },
  { NULL                   , NULL                      },
};

static const Layout overviewlayout = { "󰾍",  layout_overview };

/* commands */
static char dmenumon[2]                                     = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]                               = { "rofi", "-show", "drun", "-theme", "fullscreen-preview", "-font", "JetBrainsMono Nerd Font 24", NULL };

/* scratchpad */
static const char scratchpad_class[] = "scratchpad";
static const float scratchpad_width  = 0.75;
static const float scratchpad_height = 0.60;

static const Key keys[] = {
   /* modifier                     key                      function   argument                                 */
   { 0,                    XF86XK_AudioMute,        spawn,     { .v = SpawnShellCmd("pactl set-sink-mute @DEFAULT_SINK@ toggle") } },
   { 0,                    XF86XK_AudioLowerVolume, spawn,     { .v = SpawnShellCmd("pactl set-sink-volume @DEFAULT_SINK@ -5%") } },
   { 0,                    XF86XK_AudioRaiseVolume, spawn,     { .v = SpawnShellCmd("pactl set-sink-volume @DEFAULT_SINK@ +5%") } },
   { 0,                    XF86XK_AudioMicMute,     spawn,     { .v = SpawnShellCmd("pactl set-source-mute @DEFAULT_SOURCE@ toggle") } },
   { 0,                    XF86XK_MonBrightnessDown,spawn,     { .v = SpawnShellCmd("brightnessctl set 5%-") } },
   { 0,                    XF86XK_MonBrightnessUp,  spawn,     { .v = SpawnShellCmd("brightnessctl set +5%") } },
   { 0,                    XF86XK_AudioPlay,        spawn,     { .v = SpawnShellCmd("playerctl play-pause") } },
   { 0,                    XF86XK_AudioPause,       spawn,     { .v = SpawnShellCmd("playerctl play-pause") } },
   { 0,                    XF86XK_AudioStop,        spawn,     { .v = SpawnShellCmd("playerctl stop") } },
   { 0,                    XF86XK_AudioPrev,        spawn,     { .v = SpawnShellCmd("playerctl previous") } },
   { 0,                    XF86XK_AudioNext,        spawn,     { .v = SpawnShellCmd("playerctl next") } },

   { SUPKEY,                       XK_F1,           spawn,             { .v = Spawn("sys_volume_toggle")         } },
   { SUPKEY,                       XK_F2,           spawn,             { .v = Spawn("sys_volume_down")           } },
   { SUPKEY,                       XK_F3,           spawn,             { .v = Spawn("sys_volume_up")             } },
   { SUPKEY,                       XK_F4,           spawn,             { .v = Spawn("sys_micro_toggle")          } },
   { SUPKEY,                       XK_F5,           spawn,             { .v = Spawn("sys_display_light_down")    } },
   { SUPKEY,                       XK_F6,           spawn,             { .v = Spawn("sys_display_light_up")      } },
   { SUPKEY,                       XK_F7,           spawn,             { .v = Spawn("sys_display")               } },
   { SUPKEY,                       XK_F8,           spawn,             { .v = Spawn("sys_wifi_connect")          } },
   { SUPKEY,                       XK_F9,           spawn,             { .v = Spawn("sys_bluetooth_connect")     } },
   { SUPKEY,                       XK_F10,          spawn,             { .v = Spawn("sys_micro_down")            } },
   { SUPKEY,                       XK_F11,          spawn,             { .v = Spawn("sys_micro_up")              } },
   { SUPKEY,                       XK_F12,          spawn,             { .v = Spawn("sys_toggle_keyboard_light") } },

   { SUPKEY,                       XK_1,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://chatgpt.com/'")                                              } },
   { SUPKEY,                       XK_2,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://www.youtube.com'")                                           } },
   { SUPKEY,                       XK_3,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://github.com/zetatez'")                                        } },
   { SUPKEY,                       XK_4,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://mail.google.com/mail'")                                      } },
   { SUPKEY,                       XK_5,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://translate.google.com/?sl=auto&tl=zh-CN'")                    } },
   { SUPKEY,                       XK_6,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://web.wechat.com/'")                                           } },
   { SUPKEY,                       XK_7,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://leetcode.cn/search/?q=%E6%9C%80'")                           } },
   { SUPKEY,                       XK_8,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://tv.cctv.com/live/cctv5/?spm=C28340.P2qo7O8Q1Led.S87602.57'") } },
   { SUPKEY,                       XK_9,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://www.bilibili.com/'")                                         } },
   { SUPKEY,                       XK_0,            spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://www.doubao.com/chat/'")                                      } },
   { SUPKEY,                       XK_slash,        spawn,             { .v = SpawnShellCmd("open_url_with_qutebrowser --url='https://gemini.google.com/app'")                                     } },

   { SUPKEY|ShiftMask,             XK_1,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://chatgpt.com/'")                                                   } },
   { SUPKEY|ShiftMask,             XK_2,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://www.youtube.com'")                                                } },
   { SUPKEY|ShiftMask,             XK_3,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://github.com/zetatez'")                                             } },
   { SUPKEY|ShiftMask,             XK_4,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://mail.google.com/mail'")                                           } },
   { SUPKEY|ShiftMask,             XK_5,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://translate.google.com/?sl=auto&tl=zh-CN'")                         } },
   { SUPKEY|ShiftMask,             XK_6,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://web.wechat.com/'")                                                } },
   { SUPKEY|ShiftMask,             XK_7,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://leetcode.cn/search/?q=%E6%9C%80'")                                } },
   { SUPKEY|ShiftMask,             XK_8,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://tv.cctv.com/live/cctv5/?spm=C28340.P2qo7O8Q1Led.S87602.57'")      } },
   { SUPKEY|ShiftMask,             XK_9,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://www.bilibili.com/'")                                              } },
   { SUPKEY|ShiftMask,             XK_0,            spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://www.doubao.com/chat/'")                                           } },
   { SUPKEY|ShiftMask,             XK_slash,        spawn,             { .v = SpawnShellCmd("open_url_with_chrome --url='https://gemini.google.com/app'")                                          } },

   { SUPKEY,                       XK_k,            movewin,           { .ui = UP                                                } },
   { SUPKEY,                       XK_j,            movewin,           { .ui = DOWN                                              } },
   { SUPKEY,                       XK_h,            movewin,           { .ui = LEFT                                              } },
   { SUPKEY,                       XK_l,            movewin,           { .ui = RIGHT                                             } },

   { SUPKEY,                       XK_a,            spawn,             { .v = Spawn("launch_file_manager")                       } },
   { SUPKEY,                       XK_b,            spawn_or_focus,    { .v = CmdClass("launch_qutebrowser", "qutebrowser")      } },
   { SUPKEY,                       XK_c,            spawn,             { .v = Spawn("note_monthly_work")                         } },
   { SUPKEY,                       XK_d,            spawn,             { .v = Spawn("toggle_lazydocker")                         } },
   { SUPKEY,                       XK_e,            spawn,             { .v = Spawn("toggle_mutt")                               } },
   { SUPKEY,                       XK_f,            spawn,             { .v = SpawnTermiCmd("lazy_open_search_file")             } },
   { SUPKEY,                       XK_g,            spawn_or_focus,    { .v = CmdClass("launch_chrome", "Google-chrome")         } },
   { SUPKEY,                       XK_i,            spawn,             { .v = Spawn("toggle_flameshot")                          } },
   { SUPKEY,                       XK_m,            spawn,             { .v = SpawnTermiCmd("lazy_open_search_file_content")     } },

   { SUPKEY,                       XK_n,            toggle_scratchpad, { .v = CmdClass("st -c sp-python -e python -i -c 'import os, sys, datetime as dt, re, json, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt; print(dir())'", "sp-python") } },
   { SUPKEY,                       XK_o,            spawn,             { .v = Spawn("handle_copied")                             } },
   { SUPKEY,                       XK_p,            spawn,             { .v = SpawnTermiCmd("lazy_open_search_book")             } },
   { SUPKEY,                       XK_q,            spawn,             { .v = Spawn("slock")                                     } },
   { SUPKEY,                       XK_r,            spawn,             { .v = Spawn("toggle_yazi")                               } },
   { SUPKEY,                       XK_s,            spawn,             { .v = Spawn("search")                                    } },
   { SUPKEY,                       XK_t,            next_theme,        { 0                                                       } },
   { SUPKEY,                       XK_u,            spawn,             { .v = SpawnTermiCmd("lazy_open_search_media")            } },
   { SUPKEY,                       XK_v,            spawn,             { .v = SpawnTermiCmd("opencode")                          } },
   { SUPKEY,                       XK_w,            spawn,             { .v = SpawnTermiCmd("lazy_open_search_wiki")             } },
   { SUPKEY,                       XK_x,            spawn,             { .v = Spawn("note_scripts")                              } },
   { SUPKEY,                       XK_y,            spawn,             { .v = Spawn("toggle_show")                               } },
   { SUPKEY,                       XK_z,            spawn,             { .v = Spawn("note_todo")                                 } },

   { SUPKEY,                       XK_BackSpace,    spawn,             { .v = Spawn("toggle_passmenu")                           } },
   { SUPKEY,                       XK_Delete,       spawn,             { .v = Spawn("sys_shortcuts")                             } },
   { SUPKEY,                       XK_Escape,       spawn,             { .v = Spawn("toggle_top")                                } },
   { SUPKEY,                       XK_apostrophe,   spawn,             { .v = Spawn("toggle_tty_clock")                          } },
   { SUPKEY,                       XK_backslash,    spawn,             { .v = Spawn("reset_sys_default")                         } },
   { SUPKEY,                       XK_bracketleft,  spawn,             { .v = Spawn("toggle_calendar_scheduling")                } },
   { SUPKEY,                       XK_bracketright, spawn,             { .v = Spawn("toggle_calendar_scheduling_today")          } },
   { SUPKEY,                       XK_period,       spawn,             { .v = Spawn("jump_to_code_from_log")                     } },

   { SUPKEY|ShiftMask,             XK_k,            resizewin,         { .ui = VECINC                                                 } },
   { SUPKEY|ShiftMask,             XK_j,            resizewin,         { .ui = VECDEC                                                 } },
   { SUPKEY|ShiftMask,             XK_h,            resizewin,         { .ui = HORDEC                                                 } },
   { SUPKEY|ShiftMask,             XK_l,            resizewin,         { .ui = HORINC                                                 } },

   { SUPKEY|ShiftMask,             XK_a,            spawn,             { .v = SpawnShellCmd("gamescope -e -f -- steam -bigpicture")   } },
// { SUPKEY|ShiftMask,             XK_b,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_c,            killclient,        { 0                                                            } },
   { SUPKEY|ShiftMask,             XK_d,            spawn_or_focus,    { .v = CmdClass("dingtalk", "com.alibabainc.dingtalk")         } },
// { SUPKEY|ShiftMask,             XK_e,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_f,            spawn_or_focus,    { .v = CmdClass("feishu", "Feishu")                            } },
// { SUPKEY|ShiftMask,             XK_g,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_i,            spawn_or_focus,    { .v = CmdClass("inkscape", "Inkscape")                        } },
   { SUPKEY|ShiftMask,             XK_m,            spawn_or_focus,    { .v = CmdClass("netease-cloud-music", "netease-cloud-music")  } },
   { SUPKEY|ShiftMask,             XK_n,            toggle_scratchpad, { .v = CmdClass("st -c sp-julia -e julia", "sp-julia")         } },
   { SUPKEY|ShiftMask,             XK_o,            spawn_or_focus,    { .v = CmdClass("obsidian", "obsidian")                        } },
   { SUPKEY|ShiftMask,             XK_p,            spawn_or_focus,    { .v = CmdClass("krita", "krita")                              } },
// { SUPKEY|ShiftMask,             XK_q,            spawn,             { .v =                                                         } },
// { SUPKEY|ShiftMask,             XK_r,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_s,            spawn_or_focus,    { .v = CmdClass("subl", "Sublime_text")                        } },
// { SUPKEY|ShiftMask,             XK_t,            spawn,             { .v =                                                         } },
// { SUPKEY|ShiftMask,             XK_u,            spawn,             { .v =                                                         } },
// { SUPKEY|ShiftMask,             XK_v,            spawn,             { .v =                                                         } },
// { SUPKEY|ShiftMask,             XK_w,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_x,            spawn_or_focus,    { .v = CmdClass("xournalpp", "com.github.xournalpp.xournalpp") } },
// { SUPKEY|ShiftMask,             XK_y,            spawn,             { .v =                                                         } },
   { SUPKEY|ShiftMask,             XK_z,            spawn_or_focus,    { .v = CmdClass("zoom", "zoom")                                } },
   { SUPKEY|ShiftMask,             XK_Delete,       spawn,             { .v = SpawnShellCmd("systemctl poweroff")                     } },
   { SUPKEY|ShiftMask,             XK_Return,       spawn,             { .v = Spawn("st_dir_fzf")                                     } },
   { SUPKEY|ShiftMask,             XK_apostrophe,   spawn,             { .v = Spawn("toggle_screenkey")                               } },
   { SUPKEY|ShiftMask,             XK_comma,        spawn,             { .v = Spawn("toggle_rec_audio")                               } },
   { SUPKEY|ShiftMask,             XK_period,       spawn,             { .v = Spawn("toggle_rec_screen")                              } },

  // MODKEY, etc
   { MODKEY,                       XK_Return,       zoom,                 { 0                                            } },
   { MODKEY,                       XK_Tab,          view,                 { 0                                            } },
   { MODKEY,                       XK_apostrophe,   toggle_scratchpad,    { .v = CmdClass("st -c sp-st", "sp-st")        } },
   { MODKEY,                       XK_semicolon,    spawn,                { .v = SpawnShellCmd("rofi -show run -theme fullscreen-preview -font 'JetBrainsMono Nerd Font 24'") } },
   { MODKEY,                       XK_slash,        spawn,                { .v = Spawn("snip_fzf")                       } },
   { MODKEY,                       XK_space,        togglefloating,       { 0                                            } },
   { MODKEY,                       XK_b,            togglebar,            { 0                                            } },
   { MODKEY,                       XK_c,            spawn,                { .v = Spawn("toggle_clipmenu")                } },
   { MODKEY,                       XK_f,            togglefullscreen,     { 0                                            } },
   { MODKEY,                       XK_o,            toggleoverview,       { 0                                            } },
   { MODKEY,                       XK_p,            spawn,                { .v = dmenucmd                                } },
   { MODKEY,                       XK_q,            spawn,                { .v = Spawn("slock")                          } },
   { MODKEY,                       XK_s,            reset,                { 0                                            } },
   { MODKEY,                       XK_u,            jump_to_sel,          { 0                                            } },
   { MODKEY|ShiftMask,             XK_Return,       spawn,                { .v = Spawn("st")                             } },
   { MODKEY|ShiftMask,             XK_apostrophe,   scratchpad_to_normal, { 0                                            } },
   { MODKEY|ShiftMask,             XK_c,            killclient,           { 0                                            } },
   { MODKEY|ShiftMask,             XK_f,            togglefloating,       { 0                                            } },
   { MODKEY|ShiftMask,             XK_p,            quit,                 { 1                                            } },
   { MODKEY|ShiftMask,             XK_q,            quit,                 { 0                                            } },
   { MODKEY|ShiftMask,             XK_s,            togglesticky,         { 0                                            } },
   { MODKEY|ShiftMask,             XK_space,        focusmaster,          { 0                                            } },
   { MODKEY|ShiftMask,             XK_u,            setlayout,            { 0                                            } }, // temporary layout switch
   { MODKEY|ShiftMask|ControlMask, XK_c,            killclient_unsel,     { 0                                            } },

   { MODKEY,                       XK_bracketleft,  focusmon,             { .i = -1                                      } }, // multi monitors: focus on which one -1
   { MODKEY,                       XK_bracketright, focusmon,             { .i = +1                                      } }, // multi monitors: focus on which one +1
   { MODKEY|ShiftMask,             XK_bracketleft,  tagmon,               { .i = -1                                      } }, // multi monitors: move win to monitor prev
   { MODKEY|ShiftMask,             XK_bracketright, tagmon,               { .i = +1                                      } }, // multi monitors: move win to monitor next
   { MODKEY,                       XK_d,            incnmaster,           { .i = -1                                      } },
   { MODKEY,                       XK_i,            incnmaster,           { .i = +1                                      } },
   { MODKEY,                       XK_h,            movestack,            { .i = -1                                      } },
   { MODKEY,                       XK_l,            movestack,            { .i = +1                                      } },
   { MODKEY,                       XK_comma,        cyclelayout,          { .i = -1                                      } }, // useless
   { MODKEY,                       XK_period,       cyclelayout,          { .i = +1                                      } }, // useless
   { MODKEY|ShiftMask,             XK_comma,        shiftview,            { .i = -1                                      } },
   { MODKEY|ShiftMask,             XK_period,       shiftview,            { .i = +1                                      } },
   { MODKEY,                       XK_k,            focusstack,           { .i = -1                                      } },
   { MODKEY,                       XK_j,            focusstack,           { .i = +1                                      } },
   { MODKEY|ShiftMask,             XK_h,            setmfact,             { .f = -0.025                                  } },
   { MODKEY|ShiftMask,             XK_l,            setmfact,             { .f = +0.025                                  } },
   { MODKEY|ShiftMask,             XK_j,            sethfact,             { .f = -0.025                                  } },
   { MODKEY|ShiftMask,             XK_k,            sethfact,             { .f = +0.025                                  } },
   { MODKEY,                       XK_a,            setlayout,            { .v = &layouts[0]                             } },
   { MODKEY,                       XK_r,            setlayout,            { .v = &layouts[1]                             } },
   { MODKEY|ShiftMask,             XK_r,            setlayout,            { .v = &layouts[2]                             } },
   { MODKEY,                       XK_v,            setlayout,            { .v = &layouts[3]                             } },
   { MODKEY|ShiftMask,             XK_v,            setlayout,            { .v = &layouts[4]                             } },
   { MODKEY,                       XK_t,            setlayout,            { .v = &layouts[5]                             } },
   { MODKEY|ShiftMask,             XK_t,            setlayout,            { .v = &layouts[6]                             } },
   { MODKEY,                       XK_g,            setlayout,            { .v = &layouts[7]                             } },
   { MODKEY|ShiftMask,             XK_g,            setlayout,            { .v = &layouts[8]                             } },
   { MODKEY,                       XK_m,            setlayout,            { .v = &layouts[9]                             } },
   { MODKEY,                       XK_w,            setlayout,            { .v = &layouts[10]                            } },
   { MODKEY|ShiftMask,             XK_e,            setlayout,            { .v = &layouts[11]                            } },
   { MODKEY,                       XK_e,            setlayout,            { .v = &layouts[12]                            } },

   { MODKEY,                       XK_1,            view,              {.ui = 1 << 0                                           } }, // view tag 1
   { MODKEY,                       XK_2,            view,              {.ui = 1 << 1                                           } },
   { MODKEY,                       XK_3,            view,              {.ui = 1 << 2                                           } },
   { MODKEY,                       XK_4,            view,              {.ui = 1 << 3                                           } },
   { MODKEY,                       XK_5,            view,              {.ui = 1 << 4                                           } },
   { MODKEY,                       XK_6,            view,              {.ui = 1 << 5                                           } },
   { MODKEY,                       XK_7,            view,              {.ui = 1 << 6                                           } },
   { MODKEY,                       XK_8,            view,              {.ui = 1 << 7                                           } },
   { MODKEY,                       XK_9,            view,              {.ui = 1 << 8                                           } },
   { MODKEY,                       XK_0,            view,              {.ui = ~0                                               } }, // preview all tags
   { MODKEY|ShiftMask,             XK_1,            tag,               {.ui = 1 << 0                                           } }, // move to tag 1
   { MODKEY|ShiftMask,             XK_2,            tag,               {.ui = 1 << 1                                           } },
   { MODKEY|ShiftMask,             XK_3,            tag,               {.ui = 1 << 2                                           } },
   { MODKEY|ShiftMask,             XK_4,            tag,               {.ui = 1 << 3                                           } },
   { MODKEY|ShiftMask,             XK_5,            tag,               {.ui = 1 << 4                                           } },
   { MODKEY|ShiftMask,             XK_6,            tag,               {.ui = 1 << 5                                           } },
   { MODKEY|ShiftMask,             XK_7,            tag,               {.ui = 1 << 6                                           } },
   { MODKEY|ShiftMask,             XK_8,            tag,               {.ui = 1 << 7                                           } },
   { MODKEY|ShiftMask,             XK_9,            tag,               {.ui = 1 << 8                                           } },
   { MODKEY|ShiftMask,             XK_0,            tag,               {.ui = ~0                                               } }, // stick to all tags
   { MODKEY|ShiftMask|ControlMask, XK_1,            previewtag,        {.ui = 0                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_2,            previewtag,        {.ui = 1                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_3,            previewtag,        {.ui = 2                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_4,            previewtag,        {.ui = 3                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_5,            previewtag,        {.ui = 4                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_6,            previewtag,        {.ui = 5                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_7,            previewtag,        {.ui = 6                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_8,            previewtag,        {.ui = 7                                                } },
   { MODKEY|ShiftMask|ControlMask, XK_9,            previewtag,        {.ui = 8                                                } },
   { MODKEY|ControlMask,           XK_1,            toggleview,        {.ui = 1 << 0                                           } }, // toggle view of tag 1
   { MODKEY|ControlMask,           XK_2,            toggleview,        {.ui = 1 << 1                                           } },
   { MODKEY|ControlMask,           XK_3,            toggleview,        {.ui = 1 << 2                                           } },
   { MODKEY|ControlMask,           XK_4,            toggleview,        {.ui = 1 << 3                                           } },
   { MODKEY|ControlMask,           XK_5,            toggleview,        {.ui = 1 << 4                                           } },
   { MODKEY|ControlMask,           XK_6,            toggleview,        {.ui = 1 << 5                                           } },
   { MODKEY|ControlMask,           XK_7,            toggleview,        {.ui = 1 << 6                                           } },
   { MODKEY|ControlMask,           XK_8,            toggleview,        {.ui = 1 << 7                                           } },
   { MODKEY|ControlMask,           XK_9,            toggleview,        {.ui = 1 << 8                                           } },
};

/* button definitions */
/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
// Button1: left   click
// Button2: middle click
// Button3: right  click
// Button4:
// Button5:
static const Button buttons[] = {
  /* click                event mask      button          function             argument */
  { ClkTagBar,            0,              Button1,        view,                {0                              } },
  { ClkTagBar,            0,              Button3,        toggleview,          {0                              } },
  { ClkTagBar,            MODKEY,         Button1,        tag,                 {0                              } },
  { ClkTagBar,            MODKEY,         Button3,        toggletag,           {0                              } },
  { ClkLtSymbol,          0,              Button1,        setlayout,           {0                              } },
  { ClkLtSymbol,          0,              Button2,        setlayout,           {.v = &layouts[8]               } },
  { ClkLtSymbol,          0,              Button3,        setlayout,           {.v = &overviewlayout           } },
  { ClkStatusText,        0,              Button1,        spawn,               {.v = Spawn("toggle_tty_clock") } },
  { ClkStatusText,        0,              Button2,        spawn,               {.v = Spawn("toggle_keyboard")  } },
  { ClkStatusText,        0,              Button3,        spawn,               {.v = Spawn("toggle_calendar")  } },
  { ClkClientWin,         MODKEY,         Button1,        movemouse,           {0                              } },
  { ClkClientWin,         MODKEY,         Button2,        togglefloating,      {0                              } },
  { ClkClientWin,         MODKEY,         Button3,        resizemouse,         {0                              } },
  // { ClkWinTitle,          0,              Button1,        xxxxxxxxx,           {0                           } },
  // { ClkWinTitle,          0,              Button2,        xxxxxxxxx,           {0                           } },
  // { ClkWinTitle,          0,              Button3,        xxxxxxxxx,           {0                           } },
};
