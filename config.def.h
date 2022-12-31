/* See LICENSE file for copyright and license details. */

/* appearance */
static const unsigned int borderpx  = 1;        /* border pixel of windows */
static const unsigned int snap      = 0;        /* snap pixel */                                                                                            // patch: dwm-tag-preview
static const int scalepreview       = 4;        /* preview scaling (display w and h / scalepreview) */                                                      // patch: dwm-tag-preview
static const int previewbar         = 1;        /* show the bar in the preview window */
static const int swallowfloating    = 1;        /* 1 means swallow floating windows by default */                                                           // patch: dwm-swallow
static const int showbar            = 1;        /* 0 means no bar */
static const int topbar             = 1;        /* 0 means bottom bar */
static const int vertpad            = 8;        /* vertical padding of bar */                                                                               // patch: dwm-barpadding
static const int sidepad            = 1;        /* horizontal padding of bar */                                                                             // patch: dmenu-alpha
static const int barheight          = 12;       /* bh = (barheight > drw->fonts->h ) && (barheight < 3 * drw->fonts->h ) ? barheight : drw->fonts->h + 2 */ // patch: dwm-bar-height
static const char *fonts[]          = {"DejaVuSansMono Nerd Font:style=Book:size=10"};
static const char dmenufont[]       = "DejaVuSansMono Nerd Font:style=Book:size=10";
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

static const char *const autostart[] = {        // patch: dwm-cool-autostart
  "dwmstatus", "2>&1 >>/dev/null &", NULL,      // patch: dwm-cool-autostart
  "/home/dionysus/.dwm/autostart.sh", NULL,     // patch: dwm-cool-autostart
  NULL /* terminate */                          // patch: dwm-cool-autostart
};                                              // patch: dwm-cool-autostart

/* tagging */
/* static const char *tags[] = { "0", "1", "i", "o", "∞", "∫", "∇", "𝒹𝒮=𝛅𝒬/𝒯", "𝛇(𝓈)" }; */
static const char *tags[] = { "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "ζ(s)=∑1/n^s" };

static const Rule rules[] = {
  /* xprop(1):
   *    WM_CLASS(STRING) = instance, class
   *    WM_NAME(STRING) = title
   */
  /* class                   instance    title    tags mask     isfloating    isterminal     noswallow    monitor */
  { "st",                    NULL,       NULL,    0,            0,            1,             1,           -1 },
  //  { "netease-cloud-music",   NULL,       NULL,    1 << 8,       0,            0,             0,           -1 },
};

static const char *skipswallow[] = { "vimb", "surf" };   // patch: dwm-swallow: fix dwm-swallow annoying "swallow all problem". by myself. you can specify process name to skip swallow

/* layout(s) */
static const float mfact            = 0.50; /* factor of master area size [0.00..1.00] */                 // limit [0.05..0.95] had been extended to [0.00..1.00].
static const float ffact            = 0.50; /* factor of ffact [0.00..1.00] */                            // ffact, by myself
static const int nmaster            = 1;    /* number of clients in master area */
static const int resizehints        = 0;    /* 1 means respect size hints in tiled resizals */
static const int lockfullscreen     = 1;    /* 1 will force focus on the fullscreen window */
static const unsigned int gappoh    = 24;   /* horiz outer gap between windows and screen edge */ // patch: dwm-overview
static const unsigned int gappow    = 32;   /* vert  outer gap between windows and screen edge */ // patch: dwm-overview
static const unsigned int gappih    = 12;   /* horiz inner gap between windows */                 // patch: dwm-overview
static const unsigned int gappiw    = 16;   /* vert  inner gap between windows */                 // patch: dwm-overview

#include "layouts.c"                          // layouts
static const Layout layouts[] = {
  /* symbol               arrange function */
  { "𝒞",                  centerequalratio }, // patch: dwm-center
  { "𝒞",                    centeranyshape }, // patch: dwm-center
  { "𝒞",                           columns }, // patch: dwm-columns
  { "𝒢",                              grid }, // patch: dwm-grid
  { "𝒪",              overlaylayervertical }, // patch: dwm-overlaylayervertical
  { "𝒪",            overlaylayerhorizontal }, // patch: dwm-overlaylayerhorizontal
  { "𝒟",                      deckvertical }, // patch: dwm-deckvertical
  { "𝒟",                    deckhorizontal }, // patch: dwm-deckhorizontal
  { "ℱ",                            spiral }, // patch: dwm-fibonacci
  { "ℱ",                           dwindle }, // patch: dwm-fibonacci
  { "ℬ" ,              bottomstackvertical }, // patch: dwm-bottomstack
  { "ℬ",             bottomstackhorizontal }, // patch: dwm-bottomstack
  { "𝒯",                         tileright }, // tile -> tileright
  { "𝒯",                          tileleft }, // patch: dwm-leftstack
  { "𝒪",                  overlaylayergrid }, // patch: dwm-overlaylayergrid
  { "ℒ",                 logarithmicspiral }, // patch: dwm-logarithmicspiral
  { "ℳ",                           monocle },
  { "⦱",                              NULL }, // no layout function means floating behavior
  { NULL,                             NULL }, // patch: dwm-cyclelayouts
};

static const Layout overviewlayout = { "OVERVIEW",  overview }; // patch: dwm-overview: can be any layout

/* key definitions */
#define SUPKEY Mod4Mask
#define MODKEY Mod1Mask
#define TAGKEYS(KEY,TAG) \
  { MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
  { MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask|ShiftMask, KEY,      previewtag,     {.ui = TAG     } }, \

//{ MODKEY|ControlMask|ShiftMask, KEY,      toggletag,      {.ui = 1 << TAG} }, \  // patch: dwm-tag-preview
//{ MODKEY|ControlMask|ShiftMask, KEY,      previewtag,     {.ui = TAG     } }, \  // patch: dwm-tag-preview

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define SHCMD(cmd) { .v = (const char*[]){ "/bin/sh", "-c", cmd, NULL } }

/* commands */
static char dmenumon[2] = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]          = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray3, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *termcmd[]           = { "st", NULL };
static char scratchpadname[11]         = "scratchpad";                                         // patch: dwm-scratchpad
static const char *scratchpadcmd[]     = { "st", "-g", "180x48", "-t", scratchpadname, NULL }; // patch: dwm-scratchpad

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
static const char *kb_light_toggle[]   = SH("grep 1 /sys/class/leds/tpacpi::kbd_backlight/brightness > /dev/null; sudo sh -c \"echo $? > /sys/class/leds/tpacpi::kbd_backlight/brightness\"");

// sys
// static const char *reboot[]            = SH("systemctl reboot");
static const char *shutdown[]          = SH("systemctl poweroff -i");
static const char *suspend[]           = SH("systemctl suspend");
static const char *screenslock[]       = SH("slock & sleep .5; xset dpms force off");

// lazy
static const char *lazy_open[]         = TM("lazy -o \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='open>' --preview 'lazy -p {}' --select-1 --exit-0)\"");
static const char *lazy_open_media[]   = TM("lazy -o \"$(fd -e jpg -e jpeg -e png -e gif -e bmp -e tiff -e mp3 -e flac -e mkv -e avi -e mp4 --exclude .git . '/home/dionysus/'|fzf --prompt='medias>' --preview 'lazy -p {}' --reverse --select-1 --exit-0)\"");
static const char *lazy_open_book[]    = TM("lazy -o \"$(fd -e pdf -e epub -e djvu -e mobi --exclude .git . '/home/dionysus/obsidian/library/'|fzf --prompt='books>' --preview 'lazy -p {}' --reverse --select-1 --exit-0)\"");
static const char *lazy_open_wiki[]    = TM("lazy -o \"$(fd --type f --hidden  --exclude .git . '/home/dionysus/obsidian/wiki/'|fzf --prompt='wikis>' --preview 'lazy -p {}' --select-1 --exit-0)\"");
// static const char *lazy_copy[]         = TM("lazy -c \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='copy>' --preview 'lazy -p {}' --select-1 --exit-0)\"");
// static const char *lazy_move[]         = TM("lazy -m \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='move>' --preview 'lazy -p {}' --select-1 --exit-0)\"");
// static const char *lazy_exec[]         = TM("lazy -e \"$(fd -e sh -e jl -e py -e tex -e c -e cpp -e go -e scala -e java -e rs -e sql --exclude .git . '${HOME}'|fzf --prompt='exec>' --preview 'lazy -p {}' --select-1 --exit-0|xargs lazy -e {}");
// static const char *lazy_delete[]       = TM("lazy -d \"$(fd --type f --hidden --exclude .git . '/home/dionysus'|fzf --prompt='delete>' --preview 'lazy -p {}' --select-1 --exit-0)\"");

// apps
static const char *addressbook[]       = TM("abook");
static const char *browser[]           = SH("google");
static const char *browser_proxy[]     = SH("google --proxy-server='socks5://127.0.0.1:8000'");
static const char *calendar[]          = TM("nvim +'Calendar -view=week'");
static const char *diary[]             = TMSP("nvim +$ ~/diary/`date +%Y-%m-%d`.md");
static const char *email[]             = TMSP("neomutt");
static const char *gitter[]            = SH("gitter");
static const char *find_file_rg[]      = TMSP("~/.suckless/suckless-dwm/scripts/find_file_rg.sh");
static const char *illustrator[]       = SH("krita");
static const char *irc[]               = TM("irssi");
static const char *julia[]             = TM("julia");
static const char *lazydocker[]        = TM("lazydocker");
static const char *music[]             = SH("netease-cloud-music");
static const char *obsidian[]          = SH("obsidian");
static const char *passmenu[]          = SH("passmenu");
static const char *photoshop[]         = SH("gimp");
static const char *restart_network[]   = SH("sudo systemctl restart NetworkManager.service");
static const char *rss[]               = TM("newsboat");
static const char *screenshot[]        = SH("pkill flameshot; flameshot gui");
static const char *sublime[]           = SH("subl");
static const char *taskwarrior[]       = TM("taskwarrior-tui");
static const char *screenkey_toggle[]  = SH("pgrep -x screenkey > /dev/null; ([ \"$?\" == \"0\" ] && pkill screenkey > /dev/nul) || ([ \"$?\" == \"1\" ] && screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow >>/dev/null 2>&1 &)");
static const char *top[]               = TM("htop");
static const char *trojan[]            = SH("~/.trojan/trojan -c ~/.trojan/config.json >>/dev/null 2>&1 &");
static const char *vifm[]              = TM("vifm");
static const char *wallpaper[]         = SH("feh --bg-fill --recursive --randomize ~/Pictures/wallpapers");
static const char *wechat[]            = SH("wechat-uos");
static const char *wps[      ]         = SH("wps");

// rec: audio/video
static const char *rec_audio[]         = TM("ffmpeg -y -r 60 -f alsa -i default -c:a flac $HOME/Videos/rec-a-$(date '+%F-%H-%M-%S').flac");
static const char *rec_video[]         = TM("ffmpeg -y -s \"$(xdpyinfo|awk '/dimensions/ {print $2;}')\" -r 60 -f x11grab -i \"$DISPLAY\" -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac $HOME/Videos/rec-v-a-$(date '+%F-%H-%M-%S').mkv");

#include "movestack.c"
#include "shiftview.c"
static const Key keys[] = {
  /* modifier                     key            function           argument */
  { MODKEY,                       XK_p,          spawn,             {.v = dmenucmd          } },
  { MODKEY|ShiftMask,             XK_Return,     spawn,             {.v = termcmd           } },

  // SUPKEY + F1-F12
  { SUPKEY,                       XK_F1,         spawn,             {.v = volume_toggle     } },
  { SUPKEY,                       XK_F2,         spawn,             {.v = volume_dec        } },
  { SUPKEY,                       XK_F3,         spawn,             {.v = volume_inc        } },
//{ SUPKEY,                       XK_F4,         spawn,             {.v =                   } },
  { SUPKEY,                       XK_F5,         spawn,             {.v = screen_light_dec  } },
  { SUPKEY,                       XK_F6,         spawn,             {.v = screen_light_inc  } },
//{ SUPKEY,                       XK_F7,         spawn,             {.v =                   } },
  { SUPKEY,                       XK_F8,         spawn,             {.v = wifi              } },
//{ SUPKEY,                       XK_F9,         spawn,             {.v =                   } },
//{ SUPKEY,                       XK_F10,        spawn,             {.v =                   } },
  { SUPKEY,                       XK_F11,        spawn,             {.v = kb_light_toggle   } },
//{ SUPKEY,                       XK_F12,        spawn,             {.v =                   } },

  // SUPKEY + a-z, etc
  { SUPKEY,                       XK_a,          spawn,             {.v = lazy_open_media   } },
  { SUPKEY,                       XK_b,          spawn,             {.v = browser_proxy     } },
  { SUPKEY,                       XK_c,          spawn,             {.v = calendar          } },
  { SUPKEY,                       XK_d,          spawn,             {.v = wallpaper         } },
  { SUPKEY,                       XK_e,          spawn,             {.v = email             } },
  { SUPKEY,                       XK_f,          spawn,             {.v = lazy_open         } },
  { SUPKEY,                       XK_g,          spawn,             {.v = find_file_rg      } },
//{ SUPKEY,                       XK_h,          spawn,             {.v =                   } },
  { SUPKEY,                       XK_i,          spawn,             {.v = irc               } },
//{ SUPKEY,                       XK_j,          spawn,             {.v = x                 } },
//{ SUPKEY,                       XK_k,          spawn,             {.v = x                 } },
//{ SUPKEY,                       XK_l,          spawn,             {.v = x                 } },
//{ SUPKEY,                       XK_m,          spawn,             {.v =                   } },
//{ SUPKEY,                       XK_n,          spawn,             {.v =                   } },
  { SUPKEY,                       XK_o,          spawn,             {.v = julia             } },
  { SUPKEY,                       XK_p,          spawn,             {.v = lazy_open_book    } },
  { SUPKEY,                       XK_q,          spawn,             {.v = screenslock       } },
  { SUPKEY,                       XK_r,          spawn,             {.v = vifm              } },
//{ SUPKEY,                       XK_s,          spawn,             {.v =                   } },
//{ SUPKEY,                       XK_t,          spawn,             {.v =                   } },
  { SUPKEY,                       XK_u,          spawn,             {.v = screenkey_toggle  } },
//{ SUPKEY,                       XK_v,          spawn,             {.v =                   } },
  { SUPKEY,                       XK_w,          spawn,             {.v = lazy_open_wiki    } },
//{ SUPKEY,                       XK_x,          spawn,             {.v =                   } },
//{ SUPKEY,                       XK_y,          spawn,             {.v =                   } },
//{ SUPKEY,                       XK_z,          spawn,             {.v =                   } },
//{ SUPKEY,                       XK_apostrophe, spawn,             {.v =                   } },
  { SUPKEY,                       XK_BackSpace,  spawn,             {.v = passmenu          } },
  { SUPKEY,                       XK_Delete,     spawn,             {.v = shutdown          } },
  { SUPKEY,                       XK_Escape,     spawn,             {.v = top               } },
  { SUPKEY,                       XK_Print,      spawn,             {.v = screenshot        } },
  { SUPKEY,                       XK_backslash,  spawn,             {.v = diary             } },
  { SUPKEY,                       XK_slash,      spawn,             {.v = taskwarrior       } },
//{ SUPKEY,                       XK_comma,      spawn,             {.v =                   } },
//{ SUPKEY,                       XK_period,     spawn,             {.v =                   } },

  // SUPKEY-ShiftMask + a-z, etc
  { SUPKEY|ShiftMask,             XK_a,          spawn,             {.v = addressbook       } },
  { SUPKEY|ShiftMask,             XK_b,          spawn,             {.v = browser           } },
//{ SUPKEY|ShiftMask,             XK_c,          spawn,             {.v =                   } },
  { SUPKEY|ShiftMask,             XK_d,          spawn,             {.v = lazydocker        } },
//{ SUPKEY|ShiftMask,             XK_e,          spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_f,          spawn,             {.v =                   } },
  { SUPKEY|ShiftMask,             XK_g,          spawn,             {.v = gitter            } },
//{ SUPKEY|ShiftMask,             XK_h,          spawn,             {.v = x                 } },
  { SUPKEY|ShiftMask,             XK_i,          spawn,             {.v = illustrator       } },
//{ SUPKEY|ShiftMask,             XK_j,          spawn,             {.v = x                 } },
//{ SUPKEY|ShiftMask,             XK_k,          spawn,             {.v = x                 } },
//{ SUPKEY|ShiftMask,             XK_l,          spawn,             {.v = x                 } },
  { SUPKEY|ShiftMask,             XK_m,          spawn,             {.v = music             } },
  { SUPKEY|ShiftMask,             XK_n,          spawn,             {.v = rss               } },
  { SUPKEY|ShiftMask,             XK_o,          spawn,             {.v = obsidian          } },
  { SUPKEY|ShiftMask,             XK_p,          spawn,             {.v = photoshop         } },
  { SUPKEY|ShiftMask,             XK_q,          spawn,             {.v = suspend           } },
  { SUPKEY|ShiftMask,             XK_r,          spawn,             {.v = wps               } },
  { SUPKEY|ShiftMask,             XK_s,          spawn,             {.v = sublime           } },
  { SUPKEY|ShiftMask,             XK_t,          spawn,             {.v = trojan            } },
  { SUPKEY|ShiftMask,             XK_u,          spawn,             {.v = restart_network   } },
//{ SUPKEY|ShiftMask,             XK_v,          spawn,             {.v =                   } },
  { SUPKEY|ShiftMask,             XK_w,          spawn,             {.v = wechat            } },
//{ SUPKEY|ShiftMask,             XK_x,          spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_y,          spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_z,          spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_apostrophe, spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_Delete,     spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_Escape,     spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_Print,      spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_backslash,  spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_BackSpace,  spawn,             {.v =                   } },
//{ SUPKEY|ShiftMask,             XK_slash,      spawn,             {.v =                   } },
  { SUPKEY|ShiftMask,             XK_comma,      spawn,             {.v = rec_audio         } },
  { SUPKEY|ShiftMask,             XK_period,     spawn,             {.v = rec_video         } },

  { MODKEY,                       XK_b,          togglebar,         {0                      } },
  { MODKEY,                       XK_Return,     zoom,              {0                      } },
  { MODKEY,                       XK_Tab,        view,              {0                      } },
  { MODKEY,                       XK_space,      setlayout,         {0                      } },
  { MODKEY|ShiftMask,             XK_space,      togglefloating,    {0                      } },
  { MODKEY|ShiftMask,             XK_s,          togglesticky,      {0                      } }, // patch: dwm-sticky
  { MODKEY,                       XK_f,          togglefullscreen,  {0                      } }, // patch: dwm-actualfullscreen
  { MODKEY,                       XK_o,          toggleoverview,    {0                      } }, // patch: dwm-overview
  { MODKEY|ControlMask,           XK_space,      focusmaster,       {0                      } }, // patch: dwm-focusmaster

  { MODKEY,                       XK_k,          focusstack,        {.i = -1                } },
  { MODKEY,                       XK_j,          focusstack,        {.i = +1                } },
  { MODKEY,                       XK_d,          incnmaster,        {.i = -1                } },
  { MODKEY,                       XK_i,          incnmaster,        {.i = +1                } },
  { MODKEY,                       XK_comma,      cyclelayout,       {.i = -1                } },
  { MODKEY,                       XK_period,     cyclelayout,       {.i = +1                } },
  { MODKEY|ShiftMask,             XK_comma,      movestack,         {.i = -1                } },
  { MODKEY|ShiftMask,             XK_period,     movestack,         {.i = +1                } },
  { MODKEY|ControlMask,           XK_comma,      shiftview,         {.i = -1                } }, // shiftview
  { MODKEY|ControlMask,           XK_period,     shiftview,         {.i = +1                } }, // shiftview
  { MODKEY,                       XK_slash,      focusmon,          {.i = +1                } }, // move focus to another monitor
  { MODKEY|ShiftMask,             XK_slash,      tagmon,            {.i = +1                } }, // move tag   to another monitor
  { MODKEY|ShiftMask,             XK_h,          setmfact,          {.f = -0.025            } },
  { MODKEY|ShiftMask,             XK_l,          setmfact,          {.f = +0.025            } },
  { MODKEY|ShiftMask,             XK_j,          setffact,          {.f = -0.025            } }, // ffact, by myself
  { MODKEY|ShiftMask,             XK_k,          setffact,          {.f = +0.025            } }, // ffact, by myself
  { MODKEY|ShiftMask,             XK_m,          setlayout,         {.v = &layouts[0]       } }, // centerequalratio
  { MODKEY,                       XK_v,          setlayout,         {.v = &layouts[1]       } }, // centeranyshape
  { MODKEY|ShiftMask,             XK_v,          setlayout,         {.v = &layouts[2]       } }, // columns
  { MODKEY,                       XK_g,          setlayout,         {.v = &layouts[3]       } }, // grid
  { MODKEY,                       XK_e,          setlayout,         {.v = &layouts[4]       } }, // overlaylayervertical
  { MODKEY|ShiftMask,             XK_e,          setlayout,         {.v = &layouts[5]       } }, // overlaylayerhorizontal
  { MODKEY,                       XK_y,          setlayout,         {.v = &layouts[6]       } }, // deckvertical
  { MODKEY|ShiftMask,             XK_y,          setlayout,         {.v = &layouts[7]       } }, // deckhorizontal
  { MODKEY,                       XK_r,          setlayout,         {.v = &layouts[8]       } }, // sprial
  { MODKEY|ShiftMask,             XK_r,          setlayout,         {.v = &layouts[9]       } }, // dwindle
  { MODKEY,                       XK_w,          setlayout,         {.v = &layouts[10]      } }, // bstack
  { MODKEY|ShiftMask,             XK_w,          setlayout,         {.v = &layouts[11]      } }, // bstack
  { MODKEY,                       XK_t,          setlayout,         {.v = &layouts[12]      } }, // tileright
  { MODKEY|ShiftMask,             XK_t,          setlayout,         {.v = &layouts[13]      } }, // lefttile
  { MODKEY|ShiftMask,             XK_g,          setlayout,         {.v = &layouts[14]      } }, // overlaylayergrid
  { MODKEY,                       XK_u,          setlayout,         {.v = &layouts[15]      } }, // logarithmicspiral
  { MODKEY,                       XK_m,          setlayout,         {.v = &layouts[16]      } }, // monocle
  { MODKEY|ShiftMask,             XK_f,          setlayout,         {.v = &layouts[17]      } }, // no layout means floating
  { MODKEY,                       XK_apostrophe, togglescratch,     {.v = scratchpadcmd     } }, // patch: dwm-scratchpad
  { MODKEY,                       XK_0,          view,              {.ui = ~0               } },
  { MODKEY|ShiftMask,             XK_0,          tag,               {.ui = ~0               } },
  { SUPKEY,                       XK_k,          movewin,           {.ui = UP               } }, // patch: dwm-move-window
  { SUPKEY,                       XK_j,          movewin,           {.ui = DOWN             } }, // patch: dwm-move-window
  { SUPKEY,                       XK_h,          movewin,           {.ui = LEFT             } }, // patch: dwm-move-window
  { SUPKEY,                       XK_l,          movewin,           {.ui = RIGHT            } }, // patch: dwm-move-window
  { SUPKEY|ShiftMask,             XK_k,          resizewin,         {.ui = VINCREASE        } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_j,          resizewin,         {.ui = VDECREASE        } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_h,          resizewin,         {.ui = HDECREASE        } }, // patch: dwm-resize-window
  { SUPKEY|ShiftMask,             XK_l,          resizewin,         {.ui = HINCREASE        } }, // patch: dwm-resize-window
  TAGKEYS(                        XK_1,          0)
  TAGKEYS(                        XK_2,          1)
  TAGKEYS(                        XK_3,          2)
  TAGKEYS(                        XK_4,          3)
  TAGKEYS(                        XK_5,          4)
  TAGKEYS(                        XK_6,          5)
  TAGKEYS(                        XK_7,          6)
  TAGKEYS(                        XK_8,          7)
  TAGKEYS(                        XK_9,          8)
  { MODKEY|ShiftMask,             XK_c,          killclient,        {0                      } },
  { MODKEY|ShiftMask,             XK_q,          quit,              {0                      } },
  { MODKEY|ShiftMask,             XK_p,          quit,              {1                      } }, // patch: dwm-restartsig
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
