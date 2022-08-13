/* See LICENSE file for copyright and license details. */

/* appearance */
static const unsigned int borderpx  = 1;        /* border pixel of windows */
static const unsigned int snap      = 0;        /* snap pixel */
static const int swallowfloating    = 1;        /* 1 means swallow floating windows by default */  // dwm-swallow
static const int showbar            = 1;        /* 0 means no bar */
static const int topbar             = 1;        /* 0 means bottom bar */
static const int barheight          = 24;       /* bh = (barheight > drw->fonts->h ) && (barheight < 3 * drw->fonts->h ) ? barheight : drw->fonts->h + 2 */ // dwm-bar-height
static const char *fonts[]          = { "monospace:size=10" };
static const char dmenufont[]       = "monospace:size=10";
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
	 *	WM_CLASS(STRING) = instance, class
	 *	WM_NAME(STRING) = title
	 */
	/* class      	         instance    title    tags mask     isfloating    isterminal     noswallow    monitor */
	{ "st",                  NULL,       NULL,    0,            0,     	      1,             1,          -1 }, // dwm-swallow
	{ "netease-cloud-music", NULL,       NULL,    0,            0,     	      0,             0,          -1 }, // dwm-swallow
	{ "Gimp",                NULL,       NULL,    0,            1,            0,             0,          -1 }, // dwm-swallow
	{ "Firefox",             NULL,       NULL,    1 << 8,       0,            0,            -1,          -1 }, // dwm-swallow
    { "kitty",	             NULL,       NULL,	  0,		    0,			  0,             0,          -1 }, // dwm-swallow
    { "vimb",	             NULL,       NULL,	  0,		    0,			  0,             0,          -1 }, // dwm-swallow
};

static const SkipSwallow skipswallow[] = {                            // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    /* fix dwm-swallow annoying "swallow all parrent process problem" // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    * you can specify parrent and child process name to skip swallow  // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    *                                                                 // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    */                                                                // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    {"st", "vimb"},                                                   // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
    {"sh", "vimb"},                                                   // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself
};                                                                    // dwm-swallow: fix dwm-swallow annoying "swallow all parrent process problem". by myself

/* layout(s) */
static const float mfact            = 0.50; /* factor of master area size [0.00..1.00] */                 // limit [0.05..0.95] had been extended to [0.00..1.00].
static const int nmaster            = 1;    /* number of clients in master area */
// static const int resizehints     = 1;    /* 1 means respect size hints in tiled resizals */            // dwm-tatami
static const int resizehints        = 0;    /* 1 means respect size hints in tiled resizals */            // dwm-tatami
static const int lockfullscreen     = 1;    /* 1 will force focus on the fullscreen window */
static const float firstwindowszw   = 0.64; /* factor of center first window size width  [0.20, 0.80] */  // dwm-centerfistwindow
static const float firstwindowszh   = 0.48; /* factor of center first window size height [0.20, 0.80] */  // dwm-centerfistwindow
static const float centerwindowszw  = 0.64; /* factor of center window size width  [0.20, 0.80] */        // dwm-center my layout
static const float centerwindowszh  = 0.48; /* factor of center window size height [0.20, 0.80] */        // dwm-center my layout
static const float freeh            = 0.50; /* factor of free h [0.00..1.00] */                           // free h, by myself

#include "layouts.c"                                   // layouts
static const Layout layouts[] = {
	/* symbol     arrange function */
	{ "Center ER",                 centerequalratio }, // dwm-center
	{ "Center AS",                   centeranyshape }, // dwm-center
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

// Funcktion Keys: F1-F12
static const char *volume_toggle[]     = { "sh", "-c", "amixer set Master toggle", NULL };
static const char *volume_dec[]        = { "sh", "-c", "amixer -qM set Master 5%- umute", NULL };
static const char *volume_inc[]        = { "sh", "-c", "amixer -qM set Master 5%+ umute", NULL };
static const char *shutdown[]          = { "sh", "-c", "sudo shutdown now", NULL };
static const char *reboot[]            = { "sh", "-c", "sudo reboot", NULL };
static const char *screen_light_dec[]  = { "sh", "-c", "sudo light -U 5", NULL };
static const char *screen_light_inc[]  = { "sh", "-c", "sudo light -A 5", NULL };
static const char *bluetoothctl[]      = { "st", "-e", "sh", "-c", "bluetoothctl", NULL };
static const char *toggle_kb_light[]   = { "sh", "-c", "flag=$(cat /sys/class/leds/tpacpi::kbd_backlight/brightness); ([ \"$flag\" == \"0\" ] && sudo sh -c 'echo 1 > /sys/class/leds/tpacpi::kbd_backlight/brightness') || ([ \"$flag\" == \"1\" ] && sudo sh -c 'echo 0 > /sys/class/leds/tpacpi::kbd_backlight/brightness')", NULL };
static const char *weather[]           = { "st", "-e", "sh", "-c", "curl wttr.in/ShangHai; sleep 2", NULL };

// Chopin: open, exec, copy, move, remove, open wiki, open book, open media
static const char *chopin_open[]       = { "st", "-e", "sh", "-c", "fd --type f --hidden --exclude .git . '/home/lorenzo'|fzf --prompt='open>' --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -o {}", NULL };
static const char *chopin_copy[]       = { "st", "-e", "sh", "-c", "chopin -c \"$(fd --type f --hidden --exclude .git . './'|fzf --prompt='copy>'  --preview 'bat --color=always {}' --select-1 --exit-0)\"", NULL };
static const char *chopin_move[]       = { "st", "-e", "sh", "-c", "chopin -m \"$(fd --type f --hidden --exclude .git . './'|fzf --prompt='move>' --preview 'bat --color=always {}' --select-1 --exit-0)\"", NULL };
static const char *chopin_exec[]       = { "st", "-e", "sh", "-c", "fd -e sh -e jl -e py -e tex -e c -e cpp -e go -e scala -e java -e rs -e sql --exclude .git . './'|fzf --prompt='exec>'  --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -e {}", NULL };
static const char *chopin_remove[]     = { "st", "-e", "sh", "-c", "chopin -r \"$(fd --type f --hidden --exclude .git . './'|fzf --prompt='remove>' --preview 'bat --color=always {}' --select-1 --exit-0)\"", NULL };
static const char *chopin_open_media[] = { "st", "-e", "sh", "-c", "fd -e jpg -e jpeg -e png -e gif -e bmp -e tiff -e mp3 -e flac -e mkv -e avi -e mp4 --exclude .git . '/home/lorenzo/'|fzf --prompt='medias>' --reverse --select-1 --exit-0|xargs chopin -o {}", NULL };
static const char *chopin_open_book[]  = { "st", "-e", "sh", "-c", "fd -e pdf -e epub -e djvu -e mobi --exclude .git . '/home/lorenzo/obsidian/docs/'|fzf --prompt='books>' --reverse --select-1 --exit-0|xargs chopin -o {}", NULL };
static const char *chopin_open_wiki[]  = { "st", "-e", "sh", "-c", "fd --type f --hidden --exclude .git . '/home/lorenzo/obsidian/wiki/'|fzf --prompt='wikis>' --preview 'bat --color=always {}' --select-1 --exit-0|xargs chopin -o {}", NULL };

// System
static const char *slock[]             = { "slock", NULL };
static const char *suspend[]           = { "sh", "-c", "systemctl suspend", NULL };

// Picom
static const char *picom_grayscale[]   = { "sh", "-c", "~/.suckless/arch-dwm/scripts/picom.sh grayscale", NULL };
static const char *picom_normal[]      = { "sh", "-c", "~/.suckless/arch-dwm/scripts/picom.sh normal", NULL };

// Rec audio, video
static const char *rec_audio[]         = { "st", "-e", "sh", "-c", "ffmpeg -y -r 60 -f alsa -i default -c:a flac $HOME/Videos/rec-a-$(date '+%F-%H-%M-%S').flac", NULL };
static const char *rec_video[]         = { "st", "-e", "sh", "-c", "ffmpeg -y -s \"$(xdpyinfo | awk '/dimensions/ {print $2;}')\" -r 60 -f x11grab -i \"$DISPLAY\" -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac $HOME/Videos/rec-v-a-$(date '+%F-%H-%M-%S').mkv", NULL };

// Tools
static const char *dynamic_wallpaper[] = { "sh", "-c", "feh --bg-fill --recursive --randomize ~/Pictures/wallpapers", NULL };
static const char *calendar[]          = { "st", "-e", "sh", "-c", "vim -c 'Calendar -view=clock'", NULL };
static const char *ranger[]            = { "st", "-e", "sh", "-c", "ranger", NULL };
static const char *gotofile[]          = { "st", "-e", "sh", "-c", "~/.suckless/arch-dwm/scripts/gotofile.sh", NULL };
static const char *task[]              = { "st", "-e", "sh", "-c", "task calendar; task list; sleep 1", NULL };
static const char *screenshot[]        = { "sh", "-c", "pkill flameshot; flameshot gui", NULL };
static const char *togglescreenkey[]   = { "sh", "-c", "ps -ef|grep screenkey|grep -v grep >>/dev/null; ([ \"$?\" == \"0\" ] && pkill screenkey) || ([ \"$?\" != \"0\" ] && nohup screenkey --opacity 0 -s small --font-color yellow >>/dev/null 2>&1 &)", NULL };

// Applications
static const char *browser[]           = { "sh", "-c", "google --proxy-server='socks5://127.0.0.1:1080'", NULL };
static const char *lazydocker[]        = { "st", "-e", "sh", "-c", "lazydocker", NULL };
static const char *sublime[]           = { "subl", NULL };
static const char *email[]             = { "thunderbird", NULL };
static const char *illustrator[]       = { "krita", NULL };
static const char *music[]             = { "netease-cloud-music", NULL };
static const char *newsboat[]          = { "st", "-e", "sh", "-c", "newsboat", NULL };
static const char *obsidian[]          = { "obsidian", NULL };
static const char *photoshop[]         = { "gimp", NULL };
static const char *trojan[]            = { "sh", "-c", "nohup ~/.trojan/trojan -c ~/.trojan/config.json >>/dev/null 2>&1 &", NULL };
static const char *wechat[]            = { "wechat-uos", NULL };
static const char *zeal[]              = { "zeal", NULL };

#include "movestack.c"
#include "shiftview.c"
static Key keys[] = {
	/* modifier                     key            function           argument */
    { MODKEY,                       XK_p,          spawn,             {.v = dmenucmd          } },
	{ MODKEY,                       XK_apostrophe, togglescratch,     {.v = scratchpadcmd     } }, // dwm-scratchpad
    { MODKEY|ShiftMask,             XK_Return,     spawn,             {.v = termcmd           } },

    // Funcktion Keys: F1-F12
    { SUPKEY,                       XK_F1,         spawn,             {.v = volume_toggle     } },
    { SUPKEY,                       XK_F2,         spawn,             {.v = volume_dec        } },
    { SUPKEY,                       XK_F3,         spawn,             {.v = volume_inc        } },
    { SUPKEY,                       XK_F4,         spawn,             {.v = shutdown          } },
    { SUPKEY|ShiftMask,             XK_F4,         spawn,             {.v = reboot            } },
    { SUPKEY,                       XK_F5,         spawn,             {.v = screen_light_dec  } },
    { SUPKEY,                       XK_F6,         spawn,             {.v = screen_light_inc  } },
    { SUPKEY,                       XK_F10,        spawn,             {.v = bluetoothctl      } },
    { SUPKEY,                       XK_F11,        spawn,             {.v = toggle_kb_light   } },
    { SUPKEY,                       XK_F12,        spawn,             {.v = weather           } },

    // Chopin: open, exec, copy, move, remove, open wiki, open book, open media
    { SUPKEY,                       XK_f,          spawn,             {.v = chopin_open       } },
    { SUPKEY,                       XK_n,          spawn,             {.v = chopin_copy       } },
    { SUPKEY,                       XK_v,          spawn,             {.v = chopin_move       } },
    { SUPKEY,                       XK_x,          spawn,             {.v = chopin_exec       } },
    { SUPKEY,                       XK_z,          spawn,             {.v = chopin_remove     } },
    { SUPKEY,                       XK_a,          spawn,             {.v = chopin_open_media } },
    { SUPKEY,                       XK_p,          spawn,             {.v = chopin_open_book  } },
    { SUPKEY,                       XK_w,          spawn,             {.v = chopin_open_wiki  } },

    // System
    { SUPKEY,                       XK_q,          spawn,             {.v = slock             } },
    { SUPKEY|ShiftMask,             XK_q,          spawn,             {.v = suspend           } },

    // Picom
    { SUPKEY,                       XK_period,     spawn,             {.v = picom_normal      } },
    { SUPKEY,                       XK_comma,      spawn,             {.v = picom_grayscale   } },

    // Rec audio, video
    { SUPKEY,                       XK_backslash,  spawn,             {.v = rec_audio         } },
    { SUPKEY|ShiftMask,             XK_backslash,  spawn,             {.v = rec_video         } },

    // Tools
    { SUPKEY,                       XK_d,          spawn,             {.v = dynamic_wallpaper  } },
    { SUPKEY,                       XK_c,          spawn,             {.v = calendar          } },
    { SUPKEY,                       XK_r,          spawn,             {.v = ranger            } },
    { SUPKEY,                       XK_g,          spawn,             {.v = gotofile          } },
    { SUPKEY,                       XK_t,          spawn,             {.v = task              } },
    { SUPKEY,                       XK_Print,      spawn,             {.v = screenshot        } },
    { SUPKEY,                       XK_slash,      spawn,             {.v = togglescreenkey   } },

    // Applications
    { SUPKEY,                       XK_b,          spawn,             {.v = browser           } },
    { SUPKEY|ShiftMask,             XK_d,          spawn,             {.v = lazydocker        } },
    { SUPKEY|ShiftMask,             XK_s,          spawn,             {.v = sublime           } },
    { SUPKEY|ShiftMask,             XK_e,          spawn,             {.v = email             } },
    { SUPKEY|ShiftMask,             XK_i,          spawn,             {.v = illustrator       } },
    { SUPKEY|ShiftMask,             XK_m,          spawn,             {.v = music             } },
    { SUPKEY|ShiftMask,             XK_n,          spawn,             {.v = newsboat          } },
    { SUPKEY|ShiftMask,             XK_o,          spawn,             {.v = obsidian          } },
    { SUPKEY|ShiftMask,             XK_p,          spawn,             {.v = photoshop         } },
    { SUPKEY|ShiftMask,             XK_t,          spawn,             {.v = trojan            } },
    { SUPKEY|ShiftMask,             XK_w,          spawn,             {.v = wechat            } },
    { SUPKEY|ShiftMask,             XK_z,          spawn,             {.v = zeal              } },

	{ MODKEY,                       XK_b,          togglebar,         {0} },
	{ MODKEY,                       XK_Return,     zoom,              {0} },
	{ MODKEY,                       XK_Tab,        view,              {0} },
	{ MODKEY,                       XK_space,      setlayout,         {0} },
	{ MODKEY|ShiftMask,             XK_space,      togglefloating,    {0} },
	{ MODKEY|ShiftMask,             XK_s,          togglesticky,      {0} },                 // dwm-sticky
	{ MODKEY,                       XK_f,          togglefullscreen,  {0} },                 // dwm-actualfullscreen
	{ MODKEY,                       XK_j,          focusstack,        {.i = +1 } },
	{ MODKEY,                       XK_k,          focusstack,        {.i = -1 } },
	{ MODKEY,                       XK_i,          incnmaster,        {.i = +1 } },
	{ MODKEY,                       XK_d,          incnmaster,        {.i = -1 } },
    { MODKEY,                       XK_comma,      cyclelayout,       {.i = -1 } },
	{ MODKEY,                       XK_period,     cyclelayout,       {.i = +1 } },
	{ MODKEY|ShiftMask,             XK_comma,      movestack,         {.i = +1 } },
	{ MODKEY|ShiftMask,             XK_period,     movestack,         {.i = -1 } },
    { MODKEY|ControlMask,           XK_comma,      shiftview,         {.i = -1 } },          // shiftview
    { MODKEY|ControlMask,           XK_period,     shiftview,         {.i = +1 } },          // shiftview
	{ MODKEY,                       XK_slash,      focusmon,          {.i = +1 } },          // move cursor to another monitor
	{ MODKEY|ShiftMask,             XK_slash,      tagmon,            {.i = +1 } },          // move tag    to another monitor
	{ MODKEY|ShiftMask,             XK_h,          setmfact,          {.f = -0.025} },
	{ MODKEY|ShiftMask,             XK_l,          setmfact,          {.f = +0.025} },
	{ MODKEY|ShiftMask,             XK_j,          setfreeh,          {.f = -0.025} },       // free h, by myself
	{ MODKEY|ShiftMask,             XK_k,          setfreeh,          {.f = +0.025} },       // free h, by myself
	/* { MODKEY|ShiftMask,             XK_o,          setfrees,          {.f = -0.025} },       // free s, by myself */
	/* { MODKEY|ShiftMask,             XK_i,          setfrees,          {.f = +0.025} },       // free s, by myself */
	{ MODKEY|ShiftMask,             XK_m,          setlayout,         {.v = &layouts[0]} },  // centerequalratio         dwm-layouts
    { MODKEY,                       XK_v,          setlayout,         {.v = &layouts[1]} },  // centeranyshape           dwm-layouts
    { MODKEY,                       XK_g,          setlayout,         {.v = &layouts[2]} },  // grid                     dwm-layouts
    { MODKEY,                       XK_w,          setlayout,         {.v = &layouts[3]} },  // overlaylayervertical     dwm-layouts
    { MODKEY|ShiftMask,             XK_w,          setlayout,         {.v = &layouts[4]} },  // overlaylayerhorizontal   dwm-layouts
	{ MODKEY,                       XK_y,          setlayout,         {.v = &layouts[5]} },  // deckhorizontal           dwm-layouts
	{ MODKEY|ShiftMask,             XK_y,          setlayout,         {.v = &layouts[6]} },  // deckvertical             dwm-layouts
	{ MODKEY,                       XK_r,          setlayout,         {.v = &layouts[7]} },  // sprial                   dwm-fibonacci
	{ MODKEY|ShiftMask,             XK_r,          setlayout,         {.v = &layouts[8]} },  // dwindle                  dwm-fibonacci
	{ MODKEY,                       XK_e,          setlayout,         {.v = &layouts[9]} },  // bstack                   dwm-bottomstackhorizontal
	{ MODKEY|ShiftMask,             XK_e,          setlayout,         {.v = &layouts[10]} }, // bstack                   dwm-bottomstackvertical
	{ MODKEY,                       XK_t,          setlayout,         {.v = &layouts[11]} }, // tileright                default tile
    { MODKEY|ShiftMask,             XK_t,          setlayout,         {.v = &layouts[12]} }, // lefttile                 dwm-lefttile
    { MODKEY|ShiftMask,             XK_g,          setlayout,         {.v = &layouts[13]} }, // overlaylayergrid         dwm-overlaylayergrid
    { MODKEY|ShiftMask,             XK_u,          setlayout,         {.v = &layouts[14]} }, // logarithmicspiral        dwm-layouts
	{ MODKEY,                       XK_m,          setlayout,         {.v = &layouts[15]} }, // monocle
	{ MODKEY|ShiftMask,             XK_f,          setlayout,         {.v = &layouts[17]} }, // no layout means floating
	{ MODKEY,                       XK_0,          view,              {.ui = ~0 } },
	{ MODKEY|ShiftMask,             XK_0,          tag,               {.ui = ~0 } },

    { SUPKEY,                       XK_k,          movewin,           {.ui = UP} },          // dwm-move-window
    { SUPKEY,                       XK_j,          movewin,           {.ui = DOWN} },        // dwm-move-window
    { SUPKEY,                       XK_h,          movewin,           {.ui = LEFT} },        // dwm-move-window
    { SUPKEY,                       XK_l,          movewin,           {.ui = RIGHT} },       // dwm-move-window
                                                                                             //
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
    { ClkTagBar,            0,              Button4,        shiftview,      { .i = -1 } },           //                                                         // dwm-shiftview
    { ClkTagBar,            0,              Button5,        shiftview,      { .i = +1 } },           //                                                         // dwm-shiftview
};
