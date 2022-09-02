/* See LICENSE file for copyright and license details. */

/* appearance */
static const unsigned int borderpx  = 1;        /* border pixel of windows */
static const unsigned int snap      = 0;        /* snap pixel */
static const int swallowfloating    = 1;        /* 1 means swallow floating windows by default */  // dwm-swallow
static const int showbar            = 1;        /* 0 means no bar */
static const int topbar             = 1;        /* 0 means bottom bar */
static const int barheight          = 16;       /* bh = (barheight > drw->fonts->h ) && (barheight < 3 * drw->fonts->h ) ? barheight : drw->fonts->h + 2 */ // dwm-bar-height
static const char *fonts[]          = { "monospace:size=12" };
static const char dmenufont[]       = "monospace:size=12";
static const char col_gray1[]       = "#222222";
static const char col_gray2[]       = "#444444";
static const char col_gray3[]       = "#bbbbbb";
static const char col_gray4[]       = "#eeeeee";
static const char col_cyan[]        = "#005577";
static const char *colors[][3]      = {
    /*               fg         bg         border   */
    [SchemeNorm] = { col_gray3, col_gray1, col_gray2 },
    [SchemeSel]  = { col_gray4, col_cyan,  col_cyan  },
};

static const char *const autostart[] = {                               // dwm-cool-autostart
    "dwmstatus", "2>&1 >>/dev/null &", NULL,                           // dwm-cool-autostart
    "/home/lorenzo/.dwm/autostart.sh", NULL,                           // dwm-cool-autostart
    NULL /* terminate */                                               // dwm-cool-autostart
};                                                                     // dwm-cool-autostart

/* tagging */
static const char *tags[] = { "ζ(s)=∑1/n^s", "-e^iπ=1", "i", "o", "∞", "∫", "∇", "i", "0" };


static const Rule rules[] = {
    /* xprop(1):
     *    WM_CLASS(STRING) = instance, class
     *    WM_NAME(STRING) = title
     */
    /* class                   instance    title    tags mask     isfloating    isterminal     noswallow    monitor */
    { "st",                    NULL,       NULL,    0,            0,            1,             1,           -1 },
    { "netease-cloud-music",   NULL,       NULL,    1 << 8,       0,            0,             0,           -1 },
};

static const char *skipswallow[] = { "vimb", "surf" };   // dwm-swallow: fix dwm-swallow annoying "swallow all problem". by myself. you can specify process name to skip swallow

/* layout(s) */
static const float mfact            = 0.50; /* factor of master area size [0.00..1.00] */                 // limit [0.05..0.95] had been extended to [0.00..1.00].
static const float ffact            = 0.50; /* factor of ffact [0.00..1.00] */                            // ffact, by myself
static const int nmaster            = 1;    /* number of clients in master area */
static const int resizehints        = 0;    /* 1 means respect size hints in tiled resizals */
static const int lockfullscreen     = 1;    /* 1 will force focus on the fullscreen window */
static const unsigned int gappoh    = 24;   /* horiz outer gap between windows and screen edge */ // dwm-overview
static const unsigned int gappow    = 32;   /* vert  outer gap between windows and screen edge */ // dwm-overview
static const unsigned int gappih    = 12;   /* horiz inner gap between windows */                 // dwm-overview
static const unsigned int gappiw    = 16;   /* vert  inner gap between windows */                 // dwm-overview

#include "layouts.c"                                   // layouts
static const Layout layouts[] = {
    /* symbol     arrange function */
    { "Center ER",                 centerequalratio }, // dwm-center
    { "Center AS",                   centeranyshape }, // dwm-center
    { "Columns",                            columns }, // dwm-columns
    { "Grid",                                  grid }, // dwm-grid
    { "Overlaylayer",          overlaylayervertical }, // dwm-overlaylayervertical
    { "Overlaylayer",        overlaylayerhorizontal }, // dwm-overlaylayerhorizontal
    { "Deck",                          deckvertical }, // dwm-deckvertical
    { "Deck",                        deckhorizontal }, // dwm-deckhorizontal
    { "Fibonacci",                           spiral }, // dwm-fibonacci
    { "Fibonacci",                          dwindle }, // dwm-fibonacci
    { "Bottom Stack",           bottomstackvertical }, // dwm-bottomstack
    { "Bottom Stack",         bottomstackhorizontal }, // dwm-bottomstack
    { "Tile Right",                       tileright }, // tile -> tileright
    { "Tile Left",                         tileleft }, // dwm-leftstack
    { "Overlaylayer",              overlaylayergrid }, // dwm-overlaylayergrid
    { "Logarithmic Spiral",       logarithmicspiral }, // dwm-logarithmicspiral
    { "Monocle",                            monocle },
    { "∅",                                     NULL }, /* no layout function means floating behavior */
    { NULL,                                    NULL }, // dwm-cyclelayouts
};

static const Layout overviewlayout = { "OVERVIEW",  overview }; // dwm-overview // can be any layout

/* key definitions */
#define SUPKEY Mod4Mask
#define MODKEY Mod1Mask
#define TAGKEYS(KEY,TAG) \
    { MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
    { MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
    { MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
    { MODKEY|ControlMask|ShiftMask, KEY,      toggletag,      {.ui = 1 << TAG} },

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define SHCMD(cmd) { .v = (const char*[]){ "/bin/sh", "-c", cmd, NULL } }

/* commands */
static char dmenumon[2] = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]          = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray3, "-sb", col_cyan, "-sf", col_gray4, NULL };
static char scratchpadname[11]         = "scratchpad";                                         // dwm-scratchpad
static const char *scratchpadcmd[]     = { "st", "-g", "180x48", "-t", scratchpadname, NULL }; // dwm-scratchpad
static const char *termcmd[]           = { "st", NULL };

#define SH(cmd)    { "/bin/sh", "-c", cmd, NULL }
#define TM(cmd)    { "st", "-e", "/bin/sh", "-c", cmd, NULL }
#define TMSP(cmd)  { "st", "-g", "180x48", "-t", scratchpadname, "-e", "sh", "-c", cmd, NULL }

// SUPKEY + F1-F12
static const char *volume_toggle[]     = SH("amixer set Master toggle");
static const char *volume_dec[]        = SH("amixer -qM set Master 5%- umute");
static const char *volume_inc[]        = SH("amixer -qM set Master 5%+ umute");
static const char *screen_light_dec[]  = SH("sudo light -U 5");
static const char *screen_light_inc[]  = SH("sudo light -A 5");
static const char *wifi[]              = TM("nmtui");
static const char *bluetoothctl[]      = TM("bluetoothctl");
static const char *toggle_kb_light[]   = SH("flag=$(cat /sys/class/leds/tpacpi::kbd_backlight/brightness); ([ \"$flag\" == \"0\" ] && sudo sh -c 'echo 1 > /sys/class/leds/tpacpi::kbd_backlight/brightness') || ([ \"$flag\" == \"1\" ] && sudo sh -c 'echo 0 > /sys/class/leds/tpacpi::kbd_backlight/brightness')");
static const char *weather[]           = TMSP("curl wttr.in/ShangHai; sleep 2");

// Chopin: open, exec, copy, move, remove, open wiki, open book, open media
static const char *chopin_open[]       = TM("fd --type f --hidden --exclude .git . '/home/lorenzo'|fzf --prompt='open>' --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -o {}");
static const char *chopin_copy[]       = TM("chopin -c \"$(fd --type f --hidden --exclude .git . '/home/lorenzo'|fzf --prompt='copy>'  --preview 'bat --color=always {}' --select-1 --exit-0)\"");
static const char *chopin_move[]       = TM("chopin -m \"$(fd --type f --hidden --exclude .git . '/home/lorenzo'|fzf --prompt='move>' --preview 'bat --color=always {}' --select-1 --exit-0)\"");
static const char *chopin_exec[]       = TM("fd -e sh -e jl -e py -e tex -e c -e cpp -e go -e scala -e java -e rs -e sql --exclude .git . '/home/lorenzo'|fzf --prompt='exec>'  --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -e {}");
static const char *chopin_remove[]     = TM("chopin -r \"$(fd --type f --hidden --exclude .git . '/home/lorenzo'|fzf --prompt='remove>' --preview 'bat --color=always {}' --select-1 --exit-0)\"");
static const char *chopin_open_media[] = TM("fd -e jpg -e jpeg -e png -e gif -e bmp -e tiff -e mp3 -e flac -e mkv -e avi -e mp4 --exclude .git . '/home/lorenzo/'|fzf --prompt='medias>' --reverse --select-1 --exit-0|xargs chopin -o {}");
static const char *chopin_open_book[]  = TM("fd -e pdf -e epub -e djvu -e mobi --exclude .git . '/home/lorenzo/obsidian/docs/'|fzf --prompt='books>' --reverse --select-1 --exit-0|xargs chopin -o {}");
static const char *chopin_open_wiki[]  = TM("fd --type f --hidden --exclude .git . '/home/lorenzo/obsidian/wiki/'|fzf --prompt='wikis>' --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -o {}");

// SUPKEY + a-z
static const char *browser[]           = SH("google --proxy-server='socks5://127.0.0.1:1080'");
static const char *calendar[]          = TM("vim -c 'Calendar -view=clock'");
static const char *dynamic_wallpaper[] = SH("feh --bg-fill --recursive --randomize ~/Pictures/wallpapers");
static const char *email[]             = TM("mutt");
static const char *gotofile[]          = TM("~/.suckless/arch-dwm/scripts/gotofile.sh");
static const char *irc[]               = TM("weechat");
static const char *calculator[]        = TM("julia");
static const char *slock[]             = SH("slock");
static const char *vifm[]              = TM("vifm");
static const char *task[]              = TM("task calendar; task list; sleep 1");
static const char *togglescreenkey[]   = SH("ps -ef|grep screenkey|grep -v grep >>/dev/null; ([ \"$?\" == \"0\" ] && pkill screenkey) || ([ \"$?\" != \"0\" ] && nohup screenkey --opacity 0 -s small --font-color yellow >>/dev/null 2>&1 &)");
static const char *trans_en2zh[]       = TM("echo 'Translate EN to ZH > '; trans en:zh ");

// SUPKEY + etc
static const char *passmenu[]          = SH("passmenu");
static const char *shutdown[]          = SH("sudo shutdown now");
static const char *htop[]              = TM("htop");
static const char *screenshot[]        = SH("pkill flameshot; flameshot gui");
static const char *diary[]             = TMSP("vim +$ ~/diary/`date +%Y-%m-%d`.md");
static const char *todo[]              = TMSP("taskell ~/privacy/.taskell.md");
static const char *picom_grayscale[]   = SH("~/.suckless/arch-dwm/scripts/picom.sh grayscale");
static const char *picom_normal[]      = SH("~/.suckless/arch-dwm/scripts/picom.sh normal");

// SUPKEY-ShiftMask + a-z
static const char *addressbook[]       = TM("abook");
static const char *lazydocker[]        = TM("lazydocker");
static const char *illustrator[]       = SH("krita");
static const char *music[]             = SH("netease-cloud-music");
static const char *rss[]               = TM("newsboat");
static const char *obsidian[]          = SH("obsidian");
static const char *photoshop[]         = SH("gimp");
static const char *suspend[]           = SH("systemctl suspend");
static const char *wps[      ]         = SH("wps");
static const char *sublime[]           = SH("subl");
static const char *trojan[]            = SH("nohup ~/.trojan/trojan -c ~/.trojan/config.json >>/dev/null 2>&1 &");
static const char *nudoku[]            = TM("nudoku -d hard");
static const char *wechat[]            = SH("wechat-uos");
static const char *zeal[]              = SH("zeal");
static const char *trans_zh2en[]       = TM("echo 'Translate ZH to EN > '; trans zh:en ");

// SUPKEY-ShiftMask + etc
static const char *reboot[]            = SH("sudo reboot");
static const char *vit[]               = TMSP("vit");
static const char *rec_audio[]         = TM("ffmpeg -y -r 60 -f alsa -i default -c:a flac $HOME/Videos/rec-a-$(date '+%F-%H-%M-%S').flac");
static const char *rec_video[]         = TM("ffmpeg -y -s \"$(xdpyinfo | awk '/dimensions/ {print $2;}')\" -r 60 -f x11grab -i \"$DISPLAY\" -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac $HOME/Videos/rec-v-a-$(date '+%F-%H-%M-%S').mkv");

#include "movestack.c"
#include "shiftview.c"
static Key keys[] = {
    /* modifier                     key            function           argument */
    { MODKEY,                       XK_p,          spawn,             {.v = dmenucmd          } },
    { MODKEY|ShiftMask,             XK_Return,     spawn,             {.v = termcmd           } },

    // SUPKEY + F1-F12
    { SUPKEY,                       XK_F1,         spawn,             {.v = volume_toggle     } },
    { SUPKEY,                       XK_F2,         spawn,             {.v = volume_dec        } },
    { SUPKEY,                       XK_F3,         spawn,             {.v = volume_inc        } },
//  { SUPKEY,                       XK_F4,         spawn,             {.v =                   } },
    { SUPKEY,                       XK_F5,         spawn,             {.v = screen_light_dec  } },
    { SUPKEY,                       XK_F6,         spawn,             {.v = screen_light_inc  } },
//  { SUPKEY,                       XK_F7,         spawn,             {.v =                   } },
    { SUPKEY,                       XK_F8,         spawn,             {.v = wifi              } },
//  { SUPKEY,                       XK_F9,         spawn,             {.v =                   } },
    { SUPKEY,                       XK_F10,        spawn,             {.v = bluetoothctl      } },
    { SUPKEY,                       XK_F11,        spawn,             {.v = toggle_kb_light   } },
    { SUPKEY,                       XK_F12,        spawn,             {.v = weather           } },

    // SUPKEY + a-z, etc
    { SUPKEY,                       XK_a,          spawn,             {.v = chopin_open_media } },
    { SUPKEY,                       XK_b,          spawn,             {.v = browser           } },
    { SUPKEY,                       XK_c,          spawn,             {.v = calendar          } },
    { SUPKEY,                       XK_d,          spawn,             {.v = dynamic_wallpaper } },
    { SUPKEY,                       XK_e,          spawn,             {.v = email             } },
    { SUPKEY,                       XK_f,          spawn,             {.v = chopin_open       } },
    { SUPKEY,                       XK_g,          spawn,             {.v = gotofile          } },
//  { SUPKEY,                       XK_h,          spawn,             {.v = x                 } },
    { SUPKEY,                       XK_i,          spawn,             {.v = irc               } },
//  { SUPKEY,                       XK_j,          spawn,             {.v = x                 } },
//  { SUPKEY,                       XK_k,          spawn,             {.v = x                 } },
//  { SUPKEY,                       XK_l,          spawn,             {.v = x                 } },
//  { SUPKEY,                       XK_m,          spawn,             {.v =                   } },
    { SUPKEY,                       XK_n,          spawn,             {.v = chopin_copy       } },
    { SUPKEY,                       XK_o,          spawn,             {.v = calculator        } },
    { SUPKEY,                       XK_p,          spawn,             {.v = chopin_open_book  } },
    { SUPKEY,                       XK_q,          spawn,             {.v = slock             } },
    { SUPKEY,                       XK_r,          spawn,             {.v = vifm              } },
//  { SUPKEY,                       XK_s,          spawn,             {.v =                   } },
    { SUPKEY,                       XK_t,          spawn,             {.v = task              } },
    { SUPKEY,                       XK_u,          spawn,             {.v = togglescreenkey   } },
    { SUPKEY,                       XK_v,          spawn,             {.v = chopin_move       } },
    { SUPKEY,                       XK_w,          spawn,             {.v = chopin_open_wiki  } },
    { SUPKEY,                       XK_x,          spawn,             {.v = chopin_exec       } },
    { SUPKEY,                       XK_y,          spawn,             {.v = trans_en2zh       } },
    { SUPKEY,                       XK_z,          spawn,             {.v = chopin_remove     } },
    { SUPKEY,                       XK_apostrophe, togglescratch,     {.v = scratchpadcmd     } }, // dwm-scratchpad
    { SUPKEY,                       XK_BackSpace,  spawn,             {.v = passmenu          } },
    { SUPKEY,                       XK_Delete,     spawn,             {.v = shutdown          } },
    { SUPKEY,                       XK_Escape,     spawn,             {.v = htop              } },
    { SUPKEY,                       XK_Print,      spawn,             {.v = screenshot        } },
    { SUPKEY,                       XK_backslash,  spawn,             {.v = diary             } },
    { SUPKEY,                       XK_slash,      spawn,             {.v = todo              } },
    { SUPKEY,                       XK_comma,      spawn,             {.v = picom_grayscale   } },
    { SUPKEY,                       XK_period,     spawn,             {.v = picom_normal      } },

    // SUPKEY-ShiftMask + a-z, etc
    { SUPKEY|ShiftMask,             XK_a,          spawn,             {.v = addressbook       } },
//  { SUPKEY|ShiftMask,             XK_b,          spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_c,          spawn,             {.v =                   } },
    { SUPKEY|ShiftMask,             XK_d,          spawn,             {.v = lazydocker        } },
//  { SUPKEY|ShiftMask,             XK_e,          spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_f,          spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_g,          spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_h,          spawn,             {.v = x                 } },
    { SUPKEY|ShiftMask,             XK_i,          spawn,             {.v = illustrator       } },
//  { SUPKEY|ShiftMask,             XK_j,          spawn,             {.v = x                 } },
//  { SUPKEY|ShiftMask,             XK_k,          spawn,             {.v = x                 } },
//  { SUPKEY|ShiftMask,             XK_l,          spawn,             {.v = x                 } },
    { SUPKEY|ShiftMask,             XK_m,          spawn,             {.v = music             } },
    { SUPKEY|ShiftMask,             XK_n,          spawn,             {.v = rss               } },
    { SUPKEY|ShiftMask,             XK_o,          spawn,             {.v = obsidian          } },
    { SUPKEY|ShiftMask,             XK_p,          spawn,             {.v = photoshop         } },
    { SUPKEY|ShiftMask,             XK_q,          spawn,             {.v = suspend           } },
    { SUPKEY|ShiftMask,             XK_r,          spawn,             {.v = wps               } },
    { SUPKEY|ShiftMask,             XK_s,          spawn,             {.v = sublime           } },
    { SUPKEY|ShiftMask,             XK_t,          spawn,             {.v = trojan            } },
//  { SUPKEY|ShiftMask,             XK_u,          spawn,             {.v =                   } },
    { SUPKEY|ShiftMask,             XK_v,          spawn,             {.v = nudoku            } },
    { SUPKEY|ShiftMask,             XK_w,          spawn,             {.v = wechat            } },
//  { SUPKEY|ShiftMask,             XK_x,          spawn,             {.v =                   } },
    { SUPKEY|ShiftMask,             XK_y,          spawn,             {.v = trans_zh2en       } },
    { SUPKEY|ShiftMask,             XK_z,          spawn,             {.v = zeal              } },
//  { SUPKEY|ShiftMask,             XK_apostrophe, spawn,             {.v =                   } },
    { SUPKEY|ShiftMask,             XK_Delete,     spawn,             {.v = reboot            } },
//  { SUPKEY|ShiftMask,             XK_Escape,     spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_Print,      spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_backslash,  spawn,             {.v =                   } },
//  { SUPKEY|ShiftMask,             XK_BackSpace,  spawn,             {.v =                   } },
    { SUPKEY|ShiftMask,             XK_slash,      spawn,             {.v = vit               } },
    { SUPKEY|ShiftMask,             XK_comma,      spawn,             {.v = rec_audio         } },
    { SUPKEY|ShiftMask,             XK_period,     spawn,             {.v = rec_video         } },

    { MODKEY,                       XK_b,          togglebar,         {0} },
    { MODKEY,                       XK_Return,     zoom,              {0} },
    { MODKEY,                       XK_Tab,        view,              {0} },
    { MODKEY,                       XK_space,      setlayout,         {0} },
    { MODKEY|ShiftMask,             XK_space,      togglefloating,    {0} },
    { MODKEY|ShiftMask,             XK_s,          togglesticky,      {0} },                 // dwm-sticky
    { MODKEY,                       XK_f,          togglefullscreen,  {0} },                 // dwm-actualfullscreen
    { MODKEY,                       XK_o,          toggleoverview,    {0} },                 // dwm-overview
    { MODKEY|ControlMask,           XK_space,      focusmaster,       {0} },                 // dwm-focusmaster
    { MODKEY,                       XK_k,          focusstack,        {.i = -1 } },
    { MODKEY,                       XK_j,          focusstack,        {.i = +1 } },
    { MODKEY,                       XK_d,          incnmaster,        {.i = -1 } },
    { MODKEY,                       XK_i,          incnmaster,        {.i = +1 } },
    { MODKEY,                       XK_comma,      cyclelayout,       {.i = -1 } },
    { MODKEY,                       XK_period,     cyclelayout,       {.i = +1 } },
    { MODKEY|ShiftMask,             XK_comma,      movestack,         {.i = -1 } },
    { MODKEY|ShiftMask,             XK_period,     movestack,         {.i = +1 } },
    { MODKEY|ControlMask,           XK_comma,      shiftview,         {.i = -1 } },          // shiftview
    { MODKEY|ControlMask,           XK_period,     shiftview,         {.i = +1 } },          // shiftview
    { MODKEY,                       XK_slash,      focusmon,          {.i = +1 } },          // move cursor to another monitor
    { MODKEY|ShiftMask,             XK_slash,      tagmon,            {.i = +1 } },          // move tag    to another monitor
    { MODKEY|ShiftMask,             XK_h,          setmfact,          {.f = -0.025} },
    { MODKEY|ShiftMask,             XK_l,          setmfact,          {.f = +0.025} },
    { MODKEY|ShiftMask,             XK_j,          setffact,          {.f = -0.025} },       // ffact, by myself
    { MODKEY|ShiftMask,             XK_k,          setffact,          {.f = +0.025} },       // ffact, by myself
    { MODKEY|ShiftMask,             XK_m,          setlayout,         {.v = &layouts[0]} },  // centerequalratio
    { MODKEY,                       XK_v,          setlayout,         {.v = &layouts[1]} },  // centeranyshape
    { MODKEY|ShiftMask,             XK_v,          setlayout,         {.v = &layouts[2]} },  // columns
    { MODKEY,                       XK_g,          setlayout,         {.v = &layouts[3]} },  // grid
    { MODKEY,                       XK_w,          setlayout,         {.v = &layouts[4]} },  // overlaylayervertical
    { MODKEY|ShiftMask,             XK_w,          setlayout,         {.v = &layouts[5]} },  // overlaylayerhorizontal
    { MODKEY,                       XK_y,          setlayout,         {.v = &layouts[6]} },  // deckvertical
    { MODKEY|ShiftMask,             XK_y,          setlayout,         {.v = &layouts[7]} },  // deckhorizontal
    { MODKEY,                       XK_r,          setlayout,         {.v = &layouts[8]} },  // sprial
    { MODKEY|ShiftMask,             XK_r,          setlayout,         {.v = &layouts[9]} },  // dwindle
    { MODKEY,                       XK_e,          setlayout,         {.v = &layouts[10]} },  // bstack
    { MODKEY|ShiftMask,             XK_e,          setlayout,         {.v = &layouts[11]} }, // bstack
    { MODKEY,                       XK_t,          setlayout,         {.v = &layouts[12]} }, // tileright
    { MODKEY|ShiftMask,             XK_t,          setlayout,         {.v = &layouts[13]} }, // lefttile
    { MODKEY|ShiftMask,             XK_g,          setlayout,         {.v = &layouts[14]} }, // overlaylayergrid
    { MODKEY,                       XK_u,          setlayout,         {.v = &layouts[15]} }, // logarithmicspiral
    { MODKEY,                       XK_m,          setlayout,         {.v = &layouts[16]} }, // monocle
    { MODKEY|ShiftMask,             XK_f,          setlayout,         {.v = &layouts[17]} }, // no layout means floating
    { MODKEY,                       XK_0,          view,              {.ui = ~0 } },
    { MODKEY|ShiftMask,             XK_0,          tag,               {.ui = ~0 } },

    { SUPKEY,                       XK_k,          movewin,           {.ui = UP} },          // dwm-move-window
    { SUPKEY,                       XK_j,          movewin,           {.ui = DOWN} },        // dwm-move-window
    { SUPKEY,                       XK_h,          movewin,           {.ui = LEFT} },        // dwm-move-window
    { SUPKEY,                       XK_l,          movewin,           {.ui = RIGHT} },       // dwm-move-window
    { SUPKEY|ShiftMask,             XK_k,          resizewin,         {.ui = VINCREASE} },   // dwm-resize-window
    { SUPKEY|ShiftMask,             XK_j,          resizewin,         {.ui = VDECREASE} },   // dwm-resize-window
    { SUPKEY|ShiftMask,             XK_h,          resizewin,         {.ui = HDECREASE} },   // dwm-resize-window
    { SUPKEY|ShiftMask,             XK_l,          resizewin,         {.ui = HINCREASE} },   // dwm-resize-window
    TAGKEYS(                        XK_1,          0)
    TAGKEYS(                        XK_2,          1)
    TAGKEYS(                        XK_3,          2)
    TAGKEYS(                        XK_4,          3)
    TAGKEYS(                        XK_5,          4)
    TAGKEYS(                        XK_6,          5)
    TAGKEYS(                        XK_7,          6)
    TAGKEYS(                        XK_8,          7)
    TAGKEYS(                        XK_9,          8)
    { MODKEY|ShiftMask,             XK_c,          killclient,        {0} },
    { MODKEY|ShiftMask,             XK_q,          quit,              {0} },
};

/* button definitions */
/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
// Button1: left   click
// Button2: middle click
// Button3: right  click
// Button4:
// Button5:
static Button buttons[] = {
    /* click                event mask      button          function        argument */
    { ClkLtSymbol,          0,              Button1,        setlayout,      {0} },                   //          left   click : change layout to
    { ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[16]} },   //          right  click : change layout to x
    { ClkWinTitle,          0,              Button2,        zoom,           {0} },                   //          middle click : zoom
    { ClkStatusText,        0,              Button2,        spawn,          {.v = termcmd } },       //          middle click : open open st
    { ClkClientWin,         MODKEY,         Button1,        movemouse,      {0} },                   // modkey + left   click : move window with mouse
    { ClkClientWin,         MODKEY,         Button2,        togglefloating, {0} },                   // modkey + middle click : togglefloating
    { ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0} },                   // modkey + right  click : resize window with mouse
    { ClkTagBar,            0,              Button1,        view,           {0} },                   //          left   click : change tag
    { ClkTagBar,            0,              Button3,        toggleview,     {0} },                   // modkey + right  click : toggleview
    { ClkTagBar,            MODKEY,         Button1,        tag,            {0} },                   // modkey + left   click : move window to click tag
    { ClkTagBar,            MODKEY,         Button3,        toggletag,      {0} },                   // modkey + right  click : toggle tag
};
