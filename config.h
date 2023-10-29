/* See LICENSE file for copyright and license details. */

/* appearance */
static const unsigned int borderpx  = 1;        /* border pixel of windows */
static const unsigned int snap      = 0;        /* snap pixel */                                                                                            // patch: dwm-tag-preview
static const int scalepreview       = 3;        /* preview scaling (display w and h / scalepreview) */                                                      // patch: dwm-tag-preview
static const int previewbar         = 1;        /* show the bar in the preview window */
static const int showbar            = 1;        /* 0 means no bar */
static const int topbar             = 1;        /* 0 means bottom bar */
static const int barheight          = 12;       /* bh = (barheight > drw->fonts->h ) && (barheight < 3 * drw->fonts->h ) ? barheight : drw->fonts->h + 2 */ // patch: dwm-bar-height
static const int vertpad            = 0;        /* default 6  : vertical padding of bar */                                                                               // patch: dwm-barpadding
static const int sidepad            = 0;        /* default 256: horizontal padding of bar */                                                                             // patch: dwm-barpadding
static const int defaultwinpad      = 0;        /* default 8  : window padding of bar */
static const int swallowfloating    = 1;        /* default 1  : 1 means swallow floating windows */                                                           // patch: dwm-swallow
static const char *fonts[]          = {"DejaVuSansMono Nerd Font:style=Book:size=12"};
static const char dmenufont[]       = "DejaVuSansMono Nerd Font:style=Book:size=12";
static const char col_gray1[]       = "#222222";
static const char col_gray2[]       = "#444444";
static const char col_gray3[]       = "#bbbbbb";
static const char col_gray4[]       = "#eeeeee";
static const char col_cyan[]        = "#005577";
static const char *colors[][3]      = {
  /*               fg         bg         border   */
  [SchemeNorm] = { col_gray3, col_cyan,  col_gray2 },
  [SchemeSel]  = { col_gray4, col_cyan,  col_cyan  },
};

static const char *const autostart[] = {        // patch: dwm-cool-autostart
  "dwmblocks", "2>&1 >>/dev/null &", NULL,      // patch: dwm-cool-autostart
  "/home/dionysus/.dwm/autostart.sh", NULL,     // patch: dwm-cool-autostart
  NULL /* terminate */                          // patch: dwm-cool-autostart
};                                              // patch: dwm-cool-autostart

/* tagging */
// static const char *tags[] = { "", "2", "3", "4", "5", "6", "7", "8", "ζ(s)=∑1/n^s" };
static const char *tags[] = { "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ζ(s)=∑1/n^s" };

static const Rule rules[] = {
  /* xprop(1):
   *    WM_CLASS(STRING) = instance, class
   *    WM_NAME(STRING) = title
   */
  /* class                   instance    title    tags mask     isfloating    isterminal     noswallow    monitor */
  {"st",                     NULL,       NULL,    0,            0,            1,             1,           -1 },
  {"music",                  NULL,       NULL,    0,            1,            0,             0,           -1 },
  {"cava",                   NULL,       NULL,    0,            1,            0,             0,           -1 },
  {"00001011",               NULL,       NULL,    0,            1,            1,             0,           -1 },
};

static const char *skipswallow[] = { "vimb", "surf" };   // patch: dwm-swallow: fix dwm-swallow annoying "swallow all problem". by myself. you can specify process name to skip swallow

/* layout(s) */
static const float mfact            = 0.50; /* factor of master area size [0.00..1.00] */                 // limit [0.05..0.95] had been extended to [0.00..1.00].
static const float ffact            = 0.50; /* factor of ffact [0.00..1.00] */                            // ffact, by myself
static const int nmaster            = 1;    /* number of clients in master area */
static const int resizehints        = 0;    /* 1 means respect size hints in tiled resizals */
static const int lockfullscreen     = 0;    /* 1 will force focus on the fullscreen window */
static const unsigned int gapoh     = 24;   /* horiz outer gap between windows and screen edge */ // patch: dwm-overview
static const unsigned int gapow     = 32;   /* vert  outer gap between windows and screen edge */ // patch: dwm-overview
static const unsigned int gapih     = 12;   /* horiz inner gap between windows */                 // patch: dwm-overview
static const unsigned int gapiw     = 16;   /* vert  inner gap between windows */                 // patch: dwm-overview

#include "layouts.c"                          // layouts
static const Layout layouts[] = {
  /* symbol    arrange function */
  { "⧉",        fibonaccispiral },
  { "⧉",       fibonaccidwindle },
  { "⧈",         centeranyshape },
  { "⧈",       centerequalratio },
  { "󰘸",               deckvert },
  { "󰘸",               deckhori },
  { "⬓" ,       bottomstackvert },
  { "⬓",        bottomstackhori },
  { "◨",              tileright },
  { "◧",               tileleft },
  { "󰾍",                   grid },
  { "󰓌",                 hacker },
  { "⬚",                monocle },
  { "∅",                   NULL }, // no layout function means floating behavior
  { NULL,                  NULL }, // patch: dwm-cyclelayouts
};

static const Layout overviewlayout = { "󰕰",  overview }; // patch: dwm-overview: can be any layout

/* key definitions */
#define SUPKEY Mod4Mask
#define MODKEY Mod1Mask
#define TAGKEYS(KEY,TAG) \
  { MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
  { MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask|ShiftMask, KEY,      previewtag,     {.ui = TAG     } }, \

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define PREFIX    "/home/dionysus/.suckless/suckless-dwm/utils"
#define SH(cmd)   { "/bin/sh", "-c", cmd, NULL }
#define ST(cmd)   { "st", "-e", "/bin/sh", "-c", cmd, NULL }
#define STSP(cmd) { "st", "-g", "180x48", "-t", scratchpadname, "-e", "sh", "-c", cmd, NULL }

/* commands */
static char dmenumon[2]                = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]          = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray4, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *termcmd[]           = { "st", NULL };
static const char *scratchpadcmd[]     = { "st", "-g", "180x48", "-t", "scratchpad", NULL }; // patch: dwm-scratchpad
static const char *tabbedtermcmd[]     = { "tabbed", "-r", "2", "st", "-w", "''", NULL };

/* commands */
// static const char *cmd_lazy_exec[]            =  ST("lazy -o exec -f \"$(fd --type f -e sh -e jl -e py -e tex -e c -e cpp -e go -e scala -e java -e rs -e sql --exclude .git . '${HOME}'|fzf --prompt='exec>' --preview 'lazy -o view -f {}' --select-1 --exit-0|xargs lazy -e {}\"");
// static const char *cmd_lazy_copy[]            =  ST("lazy -o copy -f \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='copy>' --preview 'lazy -o view -f {}' --select-1 --exit-0)\"");
// static const char *cmd_lazy_rename[]          =  ST("lazy -o rename -f \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='rename>' --preview 'lazy -o view -f {}' --select-1 --exit-0)\"");
// static const char *cmd_lazy_delete[]          =  ST("lazy -o delete -f \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='delete>' --preview 'lazy -o view -f {}' --select-1 --exit-0)\"");
static const char *cmd_lazy_open[]               =  ST("lazy -o open -f \"$(fd --type f --hidden --exclude .git . '/home/dionysus/'|fzf --prompt='open>' --preview 'lazy -o view -f {}' --select-1 --exit-0)\"");
static const char *cmd_lazy_open_book[]          =  ST("lazy -o open -f \"$(fd --type f -e pdf -e epub -e djvu -e mobi --exclude .git . '/home/dionysus/my-library'|fzf --prompt='books>' --preview 'lazy -o view -f {}' --reverse --select-1 --exit-0)\"");
static const char *cmd_lazy_open_media[]         =  ST("lazy -o open -f \"$(fd --type f -e jpg -e jpeg -e png -e gif -e bmp -e tiff -e mp3 -e flac -e mkv -e avi -e mp4 --exclude .git . '/home/dionysus/'|fzf --prompt='medias>' --preview 'lazy -o view -f {}' --reverse --select-1 --exit-0)\"");
static const char *cmd_lazy_open_wiki[]          =  ST("lazy -o open -f \"$(fd --type f --hidden --exclude .git . '/home/dionysus/my-wiki'|fzf --prompt='wikis>' --preview 'lazy -o view -f {}' --select-1 --exit-0)\"");

static const char *cmd_screen_light_dec[]        =  SH("sudo light -U 5");
static const char *cmd_screen_light_inc[]        =  SH("sudo light -A 5");
static const char *cmd_screenslock[]             =  SH("sleep .5 && xset dpms force off && slock");
static const char *cmd_shutdown[]                =  SH("sleep .5 && systemctl poweroff");
static const char *cmd_suspend[]                 =  SH("sleep .5 && systemctl suspend");
static const char *cmd_volume_dec[]              =  SH("amixer set Master 5%-");
static const char *cmd_volume_inc[]              =  SH("amixer set Master 5%+");
static const char *cmd_volume_toggle[]           =  SH("amixer set Master toggle");
static const char *cmd_microphone_dec[]          =  SH("amixer set Capture 5%-");
static const char *cmd_microphone_inc[]          =  SH("amixer set Capture 5%+");
static const char *cmd_microphone_toggle[]       =  SH("amixer set Capture toggle");

/* utils */
static const char *app_passmenu[]                =  SH(PREFIX"/app-passmenu.py");
static const char *toggle_photoshop[]            =  SH(PREFIX"/toggle-photoshop.py");
static const char *toggle_addressbook[]          =  SH(PREFIX"/toggle-addressbook.py");
static const char *toggle_bluetooth[]            =  SH(PREFIX"/toggle-bluetooth.py");
static const char *toggle_calendar_schedule[]    =  SH(PREFIX"/toggle-calendar_schedule.py");
static const char *toggle_calendar_scheduling[]  =  SH(PREFIX"/toggle-calendar_scheduling.py");
static const char *toggle_chrome_with_proxy[]    =  SH(PREFIX"/toggle-chrome-with-proxy.py");
static const char *toggle_diary[]                =  SH(PREFIX"/toggle-diary.py");
static const char *toggle_edge[]                 =  SH(PREFIX"/toggle-edge.py");
static const char *toggle_email[]                =  SH(PREFIX"/toggle-mutt.py");
static const char *toggle_flameshot[]            =  SH(PREFIX"/toggle-flameshot.py");
static const char *toggle_inkscape[]             =  SH(PREFIX"/toggle-inkscape.py");
static const char *toggle_julia[]                =  SH(PREFIX"/toggle-julia.py");
static const char *toggle_kb_light[]             =  SH(PREFIX"/toggle-kb-light");
static const char *toggle_lazydocker[]           =  SH(PREFIX"/toggle-lazydocker.py");
static const char *toggle_music[]                =  SH(PREFIX"/toggle-music.py");
static const char *toggle_music_net_cloud[]      =  SH(PREFIX"/toggle-music-net-cloud.py");
static const char *toggle_rec_audio[]            =  SH(PREFIX"/toggle-rec-audio.py");
static const char *toggle_rec_video[]            =  SH(PREFIX"/toggle-rec-video.py");
static const char *toggle_redshift[]             =  SH(PREFIX"/toggle-redshift.py");
static const char *toggle_screen[]               =  SH(PREFIX"/toggle-screen.py");
static const char *toggle_screenkey[]            =  SH(PREFIX"/toggle-screenkey.py");
static const char *toggle_show[]                 =  SH(PREFIX"/toggle-show.py");
static const char *toggle_sublime[]              =  SH(PREFIX"/toggle-sublime.py");
static const char *toggle_sys_shortcuts[]        =  SH(PREFIX"/toggle-sys-shortcuts.py");
static const char *toggle_top[]                  =  SH(PREFIX"/toggle-top.py");
static const char *toggle_trojan[]               =  SH(PREFIX"/toggle-trojan.py");
static const char *toggle_joshuto[]              =  SH(PREFIX"/toggle-joshuto.py");
static const char *toggle_wallpaper[]            =  SH(PREFIX"/toggle-wallpaper.py");
static const char *toggle_wechat[]               =  SH(PREFIX"/toggle-wechat.py");
static const char *wf_clipmenu[]                 =  SH(PREFIX"/wf-clipmenu.py");
static const char *wf_download_arxiv_to_lib[]    =  SH(PREFIX"/wf-download-arxiv-to-lib.py");
static const char *wf_download_cur_to_download[] =  SH(PREFIX"/wf-download-cur-to-download.py");
static const char *wf_handle_copied[]            =  SH(PREFIX"/wf-handle-copied.py");
static const char *wf_latex[]                    =  SH(PREFIX"/wf-latex.py");
static const char *wf_rg[]                       =  ST(PREFIX"/wf-rg");
static const char *wf_sketchpad[]                =  SH(PREFIX"/wf-sketchpad.py");
static const char *wf_wifi[]                     =  SH(PREFIX"/wf-wifi.py");
static const char *wf_xournal[]                  =  SH(PREFIX"/wf-xournal.py");
static const char *search[]                      =  SH(PREFIX"/search.py");

#include "movestack.c"
#include "shiftview.c"
static const Key keys[] = {
  /* modifier                     key            function           argument                          */
  // SUPKEY + F1-F12
  { SUPKEY,                       XK_F1,         spawn,             {.v = cmd_volume_toggle           } },
  { SUPKEY,                       XK_F2,         spawn,             {.v = cmd_volume_dec              } },
  { SUPKEY,                       XK_F3,         spawn,             {.v = cmd_volume_inc              } },
  { SUPKEY,                       XK_F4,         spawn,             {.v = cmd_microphone_toggle       } },
  { SUPKEY,                       XK_F5,         spawn,             {.v = cmd_screen_light_dec        } },
  { SUPKEY,                       XK_F6,         spawn,             {.v = cmd_screen_light_inc        } },
  { SUPKEY,                       XK_F7,         spawn,             {.v = wf_wifi                     } },
  { SUPKEY,                       XK_F8,         spawn,             {.v = toggle_screen               } },
  { SUPKEY,                       XK_F9,         spawn,             {.v = toggle_bluetooth            } },
//{ SUPKEY,                       XK_F10,        spawn,             {.v =                             } },
  { SUPKEY,                       XK_F11,        spawn,             {.v = toggle_kb_light             } },
//{ SUPKEY,                       XK_F12,        spawn,             {.v =                             } },

  // SUPKEY-ShiftMask + F1-F12
//{ SUPKEY,                       XK_F1,         spawn,             {.v =                             } },
  { SUPKEY|ShiftMask,             XK_F2,         spawn,             {.v = cmd_microphone_dec          } },
  { SUPKEY|ShiftMask,             XK_F3,         spawn,             {.v = cmd_microphone_inc          } },
//{ SUPKEY,                       XK_F4,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F5,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F6,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F7,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F8,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F9,         spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F10,        spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F11,        spawn,             {.v =                             } },
//{ SUPKEY,                       XK_F12,        spawn,             {.v =                             } },

  // SUPKEY + a-z, etc
  { SUPKEY,                       XK_a,          spawn,             {.v = wf_download_arxiv_to_lib    } },
  { SUPKEY,                       XK_b,          spawn,             {.v = toggle_chrome_with_proxy    } },
  { SUPKEY,                       XK_c,          spawn,             {.v = toggle_calendar_schedule    } },
  { SUPKEY,                       XK_d,          spawn,             {.v = wf_download_cur_to_download } },
  { SUPKEY,                       XK_e,          spawn,             {.v = toggle_email                } },
  { SUPKEY,                       XK_f,          spawn,             {.v = cmd_lazy_open               } },
  { SUPKEY,                       XK_g,          spawn,             {.v = wf_rg                       } },
//{ SUPKEY,                       XK_h,          spawn,             {.v = x                           } },
  { SUPKEY,                       XK_i,          spawn,             {.v = wf_sketchpad                } },
//{ SUPKEY,                       XK_j,          spawn,             {.v = x                           } },
//{ SUPKEY,                       XK_k,          spawn,             {.v = x                           } },
//{ SUPKEY,                       XK_l,          spawn,             {.v = x                           } },
  { SUPKEY,                       XK_m,          spawn,             {.v = toggle_music                } },
  { SUPKEY,                       XK_n,          spawn,             {.v = wf_xournal                  } },
  { SUPKEY,                       XK_o,          spawn,             {.v = wf_handle_copied            } },
  { SUPKEY,                       XK_p,          spawn,             {.v = cmd_lazy_open_book          } },
  { SUPKEY,                       XK_q,          spawn,             {.v = cmd_screenslock             } },
  { SUPKEY,                       XK_r,          spawn,             {.v = toggle_joshuto              } },
  { SUPKEY,                       XK_s,          spawn,             {.v = search                      } },
//{ SUPKEY,                       XK_t,          spawn,             {.v =                             } },
  { SUPKEY,                       XK_u,          spawn,             {.v = toggle_screenkey            } },
  { SUPKEY,                       XK_v,          spawn,             {.v = cmd_lazy_open_media         } },
  { SUPKEY,                       XK_w,          spawn,             {.v = cmd_lazy_open_wiki          } },
  { SUPKEY,                       XK_x,          spawn,             {.v = toggle_wallpaper            } },
  { SUPKEY,                       XK_y,          spawn,             {.v = toggle_show                 } },
//{ SUPKEY,                       XK_z,          spawn,             {.v =                             } },
//{ SUPKEY,                       XK_apostrophe, spawn,             {.v =                             } },
  { SUPKEY,                       XK_BackSpace,  spawn,             {.v = app_passmenu                } },
  { SUPKEY,                       XK_Delete,     spawn,             {.v = toggle_sys_shortcuts        } },
  { SUPKEY,                       XK_Escape,     spawn,             {.v = toggle_top                  } },
  { SUPKEY,                       XK_Print,      spawn,             {.v = toggle_flameshot            } },
  { SUPKEY,                       XK_backslash,  spawn,             {.v = toggle_diary                } },
  { SUPKEY,                       XK_slash,      spawn,             {.v = wf_latex                    } },
//{ SUPKEY,                       XK_comma,      spawn,             {.v =                             } },
//{ SUPKEY,                       XK_period,     spawn,             {.v =                             } },

  // SUPKEY-ShiftMask + a-z, etc
  { SUPKEY|ShiftMask,             XK_a,          spawn,             {.v = toggle_addressbook          } },
  { SUPKEY|ShiftMask,             XK_b,          spawn,             {.v = toggle_edge                 } },
  { SUPKEY|ShiftMask,             XK_c,          spawn,             {.v = toggle_calendar_scheduling  } },
  { SUPKEY|ShiftMask,             XK_d,          spawn,             {.v = toggle_lazydocker           } },
//{ SUPKEY|ShiftMask,             XK_e,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_f,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_g,          spawn,             {.v = x                           } },
//{ SUPKEY|ShiftMask,             XK_h,          spawn,             {.v = x                           } },
  { SUPKEY|ShiftMask,             XK_i,          spawn,             {.v = toggle_inkscape             } },
//{ SUPKEY|ShiftMask,             XK_j,          spawn,             {.v = x                           } },
//{ SUPKEY|ShiftMask,             XK_k,          spawn,             {.v = x                           } },
//{ SUPKEY|ShiftMask,             XK_l,          spawn,             {.v = x                           } },
  { SUPKEY|ShiftMask,             XK_m,          spawn,             {.v = toggle_music_net_cloud      } },
//{ SUPKEY|ShiftMask,             XK_n,          spawn,             {.v =                             } },
  { SUPKEY|ShiftMask,             XK_o,          spawn,             {.v = toggle_julia                } },
  { SUPKEY|ShiftMask,             XK_p,          spawn,             {.v = toggle_photoshop            } },
  { SUPKEY|ShiftMask,             XK_q,          spawn,             {.v = cmd_suspend                 } },
  { SUPKEY|ShiftMask,             XK_r,          spawn,             {.v = toggle_redshift             } },
  { SUPKEY|ShiftMask,             XK_s,          spawn,             {.v = toggle_sublime              } },
  { SUPKEY|ShiftMask,             XK_t,          spawn,             {.v = toggle_trojan               } },
//{ SUPKEY|ShiftMask,             XK_u,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_v,          spawn,             {.v =                             } },
  { SUPKEY|ShiftMask,             XK_w,          spawn,             {.v = toggle_wechat               } },
//{ SUPKEY|ShiftMask,             XK_x,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_y,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_z,          spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_apostrophe, spawn,             {.v =                             } },
  { SUPKEY|ShiftMask,             XK_Delete,     spawn,             {.v = cmd_shutdown                } },
//{ SUPKEY|ShiftMask,             XK_Escape,     spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_Print,      spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_backslash,  spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_BackSpace,  spawn,             {.v =                             } },
//{ SUPKEY|ShiftMask,             XK_slash,      spawn,             {.v =                             } },
  { SUPKEY|ShiftMask,             XK_comma,      spawn,             {.v = toggle_rec_audio            } },
  { SUPKEY|ShiftMask,             XK_period,     spawn,             {.v = toggle_rec_video            } },

  // MODKEY-ShiftMask/ControlMask + a-z, etc
  { MODKEY,                       XK_apostrophe, togglescratch,     {.v = scratchpadcmd               } }, // patch: dwm-scratchpad
  { MODKEY,                       XK_c,          spawn,             {.v = wf_clipmenu                 } },
  { MODKEY,                       XK_p,          spawn,             {.v = dmenucmd                    } },
  { MODKEY,                       XK_Return,     zoom,              {0                                } },
  { MODKEY,                       XK_Tab,        view,              {0                                } }, // switch current tag    with previous tag
  { MODKEY,                       XK_b,          togglebar,         {0                                } },
  { MODKEY,                       XK_f,          togglefullscreen,  {0                                } }, // patch: dwm-actualfullscreen
  { MODKEY,                       XK_o,          toggleoverview,    {0                                } }, // patch: dwm-overview
  { MODKEY,                       XK_s,          reset,             {0                                } }, // reset, by myself
  { MODKEY,                       XK_space,      setlayout,         {0                                } }, // switch current layout with previous layout
  { MODKEY|ControlMask,           XK_space,      focusmaster,       {0                                } }, // patch: dwm-focusmaster
  { MODKEY|ShiftMask,             XK_s,          togglesticky,      {0                                } }, // patch: dwm-sticky
  { MODKEY|ShiftMask,             XK_space,      togglefloating,    {0                                } },
  { MODKEY,                       XK_k,          focusstack,        {.i = -1                          } },
  { MODKEY,                       XK_j,          focusstack,        {.i = +1                          } },
  { MODKEY,                       XK_d,          incnmaster,        {.i = -1                          } },
  { MODKEY,                       XK_i,          incnmaster,        {.i = +1                          } },
  { MODKEY,                       XK_comma,      cyclelayout,       {.i = -1                          } },
  { MODKEY,                       XK_period,     cyclelayout,       {.i = +1                          } },
  { MODKEY|ShiftMask,             XK_comma,      movestack,         {.i = -1                          } },
  { MODKEY|ShiftMask,             XK_period,     movestack,         {.i = +1                          } },
  { MODKEY|ControlMask,           XK_comma,      shiftview,         {.i = -1                          } }, // shiftview
  { MODKEY|ControlMask,           XK_period,     shiftview,         {.i = +1                          } }, // shiftview
  { MODKEY,                       XK_slash,      focusmon,          {.i = +1                          } }, // move focus to another monitor
  { MODKEY|ShiftMask,             XK_slash,      tagmon,            {.i = +1                          } }, // move tag   to another monitor
  { MODKEY|ShiftMask,             XK_h,          setmfact,          {.f = -0.025                      } },
  { MODKEY|ShiftMask,             XK_l,          setmfact,          {.f = +0.025                      } },
  { MODKEY|ShiftMask,             XK_j,          setffact,          {.f = -0.025                      } }, // ffact, by myself
  { MODKEY|ShiftMask,             XK_k,          setffact,          {.f = +0.025                      } }, // ffact, by myself
  { SUPKEY,                       XK_k,          movewin,           {.ui = UP                         } }, // patch: dwm-move-window
  { SUPKEY,                       XK_j,          movewin,           {.ui = DOWN                       } }, // patch: dwm-move-window
  { SUPKEY,                       XK_h,          movewin,           {.ui = LEFT                       } }, // patch: dwm-move-window
  { SUPKEY,                       XK_l,          movewin,           {.ui = RIGHT                      } }, // patch: dwm-move-window
  { SUPKEY|ShiftMask,             XK_k,          resizewin,         {.ui = VECINC                     } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_j,          resizewin,         {.ui = VECDEC                     } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_h,          resizewin,         {.ui = HORDEC                     } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_l,          resizewin,         {.ui = HORINC                     } }, // patch: dwm-resize-window
  { MODKEY,                       XK_r,          setlayout,         {.v = &layouts[0]                 } }, // fibonaccispiral
  { MODKEY|ShiftMask,             XK_r,          setlayout,         {.v = &layouts[1]                 } }, // fibonaccidwindle
  { MODKEY,                       XK_v,          setlayout,         {.v = &layouts[2]                 } }, // centeranyshape
  { MODKEY|ShiftMask,             XK_v,          setlayout,         {.v = &layouts[3]                 } }, // centerequalratio
  { MODKEY,                       XK_y,          setlayout,         {.v = &layouts[4]                 } }, // deckvert
  { MODKEY|ShiftMask,             XK_y,          setlayout,         {.v = &layouts[5]                 } }, // deckhori
  { MODKEY,                       XK_e,          setlayout,         {.v = &layouts[6]                 } }, // bottomstackvert
  { MODKEY|ShiftMask,             XK_e,          setlayout,         {.v = &layouts[7]                 } }, // bottomstackhori
  { MODKEY,                       XK_t,          setlayout,         {.v = &layouts[8]                 } }, // tileright
  { MODKEY|ShiftMask,             XK_t,          setlayout,         {.v = &layouts[9]                 } }, // tileleft
  { MODKEY,                       XK_g,          setlayout,         {.v = &layouts[10]                } }, // grid
  { MODKEY,                       XK_a,          setlayout,         {.v = &layouts[11]                } }, // hacker
  { MODKEY,                       XK_m,          setlayout,         {.v = &layouts[12]                } }, // monocle
  { MODKEY|ShiftMask,             XK_f,          setlayout,         {.v = &layouts[13]                } }, // no layout means floating
    TAGKEYS(XK_1, 0)
    TAGKEYS(XK_2, 1)
    TAGKEYS(XK_3, 2)
    TAGKEYS(XK_4, 3)
    TAGKEYS(XK_5, 4)
    TAGKEYS(XK_6, 5)
    TAGKEYS(XK_7, 6)
    TAGKEYS(XK_8, 7)
    TAGKEYS(XK_9, 8)
  { MODKEY,                       XK_0,          view,              {.ui = ~0                         } },
  { MODKEY|ShiftMask,             XK_0,          tag,               {.ui = ~0                         } },
  { MODKEY|ShiftMask,             XK_Return,     spawn,             {.v = termcmd                     } },
  { SUPKEY|ShiftMask,             XK_Return,     spawn,             {.v = tabbedtermcmd               } },
  { MODKEY|ShiftMask,             XK_c,          killclient,        {0                                } },
  { MODKEY|ShiftMask,             XK_q,          quit,              {0                                } },
  { MODKEY|ShiftMask,             XK_p,          quit,              {1                                } }, // patch: dwm-restartsig
};

/* button definitions */
/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
// Button1: left   click
// Button2: middle click
// Button3: right  click
// Button4:
// Button5:
static const Button buttons[] = {
  /* click                event mask      button          function        argument */
  { ClkLtSymbol,          0,              Button1,        setlayout,      {0                         } }, //          left   click layout symbol: change layout to previous
  { ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[2]          } }, //          right  click layout symbol: change layout to x
  { ClkWinTitle,          0,              Button2,        zoom,           {0                         } }, //          middle click win title    : zoom
  { ClkStatusText,        0,              Button1,        spawn,          {.v = termcmd              } }, //          left   click status text  : open open st
  { ClkStatusText,        0,              Button2,        spawn,          {.v = toggle_screen        } }, //          middle click status text  : open open st
  { ClkStatusText,        0,              Button3,        spawn,          {.v = toggle_sys_shortcuts } }, //          right  click status text  : open
  { ClkClientWin,         MODKEY,         Button1,        movemouse,      {0                         } }, // modkey + left   click client win   : move window with mouse
  { ClkClientWin,         MODKEY,         Button2,        togglefloating, {0                         } }, // modkey + middle click client win   : togglefloating
  { ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0                         } }, // modkey + right  click client win   : resize window with mouse
  { ClkTagBar,            0,              Button1,        view,           {0                         } }, //          left   click tag bar      : view tag
  { ClkTagBar,            0,              Button3,        toggleview,     {0                         } }, //          right  click tag bar      : toggle view, view multiple tags
  { ClkTagBar,            MODKEY,         Button1,        tag,            {0                         } }, // modkey + left   click tag bar      : move window to tag clicked
  { ClkTagBar,            MODKEY,         Button3,        toggletag,      {0                         } }, // modkey + right  click tag bar      : toggle tag
};
