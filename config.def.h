
/* key definitions */
#define SUPKEY Mod4Mask
#define MODKEY Mod1Mask
#define TAGKEYS(KEY,TAG) \
  { MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
  { MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
  { MODKEY|ControlMask|ShiftMask, KEY,      previewtag,     {.ui = TAG     } }, \

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define PREFIX    "/home/dionysus/.suckless/suckless-dwm/bin"
#define SH(cmd)   { "/bin/sh", "-c", cmd, NULL }
#define ST(cmd)   { "st", "-e", "/bin/sh", "-c", cmd, NULL }
#define STSP(cmd) { "st", "-g", "180x48", "-t", scratchpadname, "-e", "sh", "-c", cmd, NULL }

/* appearance */
static const unsigned int borderpx  = 1;
static const unsigned int snap      = 0;
static const int scalepreview       = 3;
static const int previewbar         = 1;
static const int showbar            = 1;
static const int topbar             = 1;
static const int barheight          = 12;
static const int vertpad            = 0;
static const int sidepad            = 0;
static const int defaultwinpad      = 0;
static const int swallowfloating    = 1;
static const char *fonts[]          = { "DejaVuSansMono Nerd Font:style=Book:size=12" };
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

static const char *const autostart[] = {
  "dwmblocks", NULL,
  "/home/dionysus/.dwm/autostart.sh", NULL,
  NULL
};

/* tagging */
// static const char *tags[] = { "", "2", "3", "4", "5", "6", "7", "8", "ζ(s)=∑1/n^s" };
static const char *tags[] = { "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ζ(s)=∑1/n^s" };

static const Rule rules[] = {
  /* cls                     instance    title    tags mask     isfloating    isterminal     noswallow    monitor */
  {"st",                     NULL,       NULL,    0,            0,            1,             1,           -1 },
  {"music",                  NULL,       NULL,    0,            1,            0,             0,           -1 },
  {"cava",                   NULL,       NULL,    0,            1,            0,             0,           -1 },
  {"00001011",               NULL,       NULL,    0,            1,            1,             0,           -1 },
};

static const char *skipswallow[] = { "vimb", "surf" };

/* layout(s) */
static const float mfact            = 0.50;
static const float ffact            = 0.50;
static const int nmaster            = 1;
static const int resizehints        = 0;
static const int lockfullscreen     = 0;
static const unsigned int gapoh     = 24;
static const unsigned int gapow     = 32;
static const unsigned int gapih     = 12;
static const unsigned int gapiw     = 16;

#include "layouts.c"
static const Layout layouts[] = {
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
  { "∅",                   NULL },
  { NULL,                  NULL },
};

static const Layout overviewlayout = { "󰕰",  overview };

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
static const char *app_passmenu[]                =  SH(PREFIX"/app-passmenu");
static const char *toggle_photoshop[]            =  SH(PREFIX"/toggle-photoshop");
static const char *toggle_addressbook[]          =  SH(PREFIX"/toggle-addressbook");
static const char *toggle_bluetooth[]            =  SH(PREFIX"/toggle-bluetooth");
static const char *toggle_calendar_schedule[]    =  SH(PREFIX"/toggle-calendar-schedule");
static const char *toggle_calendar_scheduling[]  =  SH(PREFIX"/toggle-calendar-scheduling");
static const char *toggle_chrome_with_proxy[]    =  SH(PREFIX"/toggle-chrome-with-proxy");
static const char *toggle_diary[]                =  SH(PREFIX"/toggle-diary");
static const char *toggle_edge[]                 =  SH(PREFIX"/toggle-edge");
static const char *toggle_email[]                =  SH(PREFIX"/toggle-mutt");
static const char *toggle_flameshot[]            =  SH(PREFIX"/toggle-flameshot");
static const char *toggle_inkscape[]             =  SH(PREFIX"/toggle-inkscape");
static const char *toggle_julia[]                =  SH(PREFIX"/toggle-julia");
static const char *toggle_kb_light[]             =  SH(PREFIX"/toggle-kb-light");
static const char *toggle_lazydocker[]           =  SH(PREFIX"/toggle-lazydocker");
static const char *toggle_music[]                =  SH(PREFIX"/toggle-music");
static const char *toggle_music_net_cloud[]      =  SH(PREFIX"/toggle-music-net-cloud");
static const char *toggle_rec_audio[]            =  SH(PREFIX"/toggle-rec-audio");
static const char *toggle_rec_video[]            =  SH(PREFIX"/toggle-rec-video");
static const char *toggle_redshift[]             =  SH(PREFIX"/toggle-redshift");
static const char *toggle_screen[]               =  SH(PREFIX"/toggle-screen");
static const char *toggle_screenkey[]            =  SH(PREFIX"/toggle-screenkey");
static const char *toggle_show[]                 =  SH(PREFIX"/toggle-show");
static const char *toggle_sublime[]              =  SH(PREFIX"/toggle-sublime");
static const char *toggle_sys_shortcuts[]        =  SH(PREFIX"/toggle-sys-shortcuts");
static const char *toggle_top[]                  =  SH(PREFIX"/toggle-top");
static const char *toggle_trojan[]               =  SH(PREFIX"/toggle-trojan");
static const char *toggle_joshuto[]              =  SH(PREFIX"/toggle-joshuto");
static const char *toggle_wallpaper[]            =  SH(PREFIX"/toggle-wallpaper");
static const char *toggle_wechat[]               =  SH(PREFIX"/toggle-wechat");
static const char *wf_clipmenu[]                 =  SH(PREFIX"/wf-clipmenu");
static const char *wf_download_arxiv_to_lib[]    =  SH(PREFIX"/wf-download-arxiv-to-lib");
static const char *wf_download_cur_to_download[] =  SH(PREFIX"/wf-download-cur-to-download");
static const char *wf_handle_copied[]            =  SH(PREFIX"/wf-handle-copied");
static const char *wf_latex[]                    =  SH(PREFIX"/wf-latex");
static const char *wf_rg[]                       =  ST(PREFIX"/wf-rg");
static const char *wf_sketchpad[]                =  SH(PREFIX"/wf-sketchpad");
static const char *wf_wifi[]                     =  SH(PREFIX"/wf-wifi");
static const char *wf_xournal[]                  =  SH(PREFIX"/wf-xournal");
static const char *search[]                      =  SH(PREFIX"/search");

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
//{ SUPKEY|ShiftMask,             XK_t,          spawn,             {.v =                             } },
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
  { MODKEY,                       XK_apostrophe, togglescratch,     {.v = scratchpadcmd               } },
  { MODKEY,                       XK_c,          spawn,             {.v = wf_clipmenu                 } },
  { MODKEY,                       XK_p,          spawn,             {.v = dmenucmd                    } },
  { MODKEY,                       XK_Return,     zoom,              {0                                } },
  { MODKEY,                       XK_Tab,        view,              {0                                } },
  { MODKEY,                       XK_b,          togglebar,         {0                                } },
  { MODKEY,                       XK_f,          togglefullscreen,  {0                                } },
  { MODKEY,                       XK_o,          toggleoverview,    {0                                } },
  { MODKEY,                       XK_s,          reset,             {0                                } },
  { MODKEY,                       XK_space,      setlayout,         {0                                } },
  { MODKEY|ControlMask,           XK_space,      focusmaster,       {0                                } },
  { MODKEY|ShiftMask,             XK_s,          togglesticky,      {0                                } },
  { MODKEY|ShiftMask,             XK_space,      togglefloating,    {0                                } },
  { MODKEY,                       XK_k,          focusstack,        {.i = -1                          } },
  { MODKEY,                       XK_j,          focusstack,        {.i = +1                          } },
  { MODKEY,                       XK_d,          incnmaster,        {.i = -1                          } },
  { MODKEY,                       XK_i,          incnmaster,        {.i = +1                          } },
  { MODKEY,                       XK_comma,      cyclelayout,       {.i = -1                          } },
  { MODKEY,                       XK_period,     cyclelayout,       {.i = +1                          } },
  { MODKEY|ShiftMask,             XK_comma,      movestack,         {.i = -1                          } },
  { MODKEY|ShiftMask,             XK_period,     movestack,         {.i = +1                          } },
  { MODKEY|ControlMask,           XK_comma,      shiftview,         {.i = -1                          } },
  { MODKEY|ControlMask,           XK_period,     shiftview,         {.i = +1                          } },
  { MODKEY,                       XK_slash,      focusmon,          {.i = +1                          } },
  { MODKEY|ShiftMask,             XK_slash,      tagmon,            {.i = +1                          } },
  { MODKEY|ShiftMask,             XK_h,          setmfact,          {.f = -0.025                      } },
  { MODKEY|ShiftMask,             XK_l,          setmfact,          {.f = +0.025                      } },
  { MODKEY|ShiftMask,             XK_j,          setffact,          {.f = -0.025                      } },
  { MODKEY|ShiftMask,             XK_k,          setffact,          {.f = +0.025                      } },
  { SUPKEY,                       XK_k,          movewin,           {.ui = UP                         } },
  { SUPKEY,                       XK_j,          movewin,           {.ui = DOWN                       } },
  { SUPKEY,                       XK_h,          movewin,           {.ui = LEFT                       } },
  { SUPKEY,                       XK_l,          movewin,           {.ui = RIGHT                      } },
  { SUPKEY|ShiftMask,             XK_k,          resizewin,         {.ui = VECINC                     } },
  { SUPKEY|ShiftMask,             XK_j,          resizewin,         {.ui = VECDEC                     } },
  { SUPKEY|ShiftMask,             XK_h,          resizewin,         {.ui = HORDEC                     } },
  { SUPKEY|ShiftMask,             XK_l,          resizewin,         {.ui = HORINC                     } },
  { MODKEY,                       XK_r,          setlayout,         {.v = &layouts[0]                 } },
  { MODKEY|ShiftMask,             XK_r,          setlayout,         {.v = &layouts[1]                 } },
  { MODKEY,                       XK_v,          setlayout,         {.v = &layouts[2]                 } },
  { MODKEY|ShiftMask,             XK_v,          setlayout,         {.v = &layouts[3]                 } },
  { MODKEY,                       XK_y,          setlayout,         {.v = &layouts[4]                 } },
  { MODKEY|ShiftMask,             XK_y,          setlayout,         {.v = &layouts[5]                 } },
  { MODKEY,                       XK_e,          setlayout,         {.v = &layouts[6]                 } },
  { MODKEY|ShiftMask,             XK_e,          setlayout,         {.v = &layouts[7]                 } },
  { MODKEY,                       XK_t,          setlayout,         {.v = &layouts[8]                 } },
  { MODKEY|ShiftMask,             XK_t,          setlayout,         {.v = &layouts[9]                 } },
  { MODKEY,                       XK_g,          setlayout,         {.v = &layouts[10]                } },
  { MODKEY,                       XK_a,          setlayout,         {.v = &layouts[11]                } },
  { MODKEY,                       XK_m,          setlayout,         {.v = &layouts[12]                } },
  { MODKEY|ShiftMask,             XK_f,          setlayout,         {.v = &layouts[13]                } },
  { MODKEY,                       XK_0,          view,              {.ui = ~0                         } },
  { MODKEY|ShiftMask,             XK_0,          tag,               {.ui = ~0                         } },
  { MODKEY|ShiftMask,             XK_Return,     spawn,             {.v = termcmd                     } },
  { SUPKEY|ShiftMask,             XK_Return,     spawn,             {.v = tabbedtermcmd               } },
  { MODKEY|ShiftMask,             XK_c,          killclient,        {0                                } },
  { MODKEY|ShiftMask,             XK_q,          quit,              {0                                } },
  { MODKEY|ShiftMask,             XK_p,          quit,              {1                                } },
    TAGKEYS(XK_1, 0)
    TAGKEYS(XK_2, 1)
    TAGKEYS(XK_3, 2)
    TAGKEYS(XK_4, 3)
    TAGKEYS(XK_5, 4)
    TAGKEYS(XK_6, 5)
    TAGKEYS(XK_7, 6)
    TAGKEYS(XK_8, 7)
    TAGKEYS(XK_9, 8)
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
  { ClkLtSymbol,          0,              Button1,        setlayout,      {0                         } },
  { ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[2]          } },
  { ClkWinTitle,          0,              Button2,        zoom,           {0                         } },
  { ClkStatusText,        0,              Button1,        spawn,          {.v = termcmd              } },
  { ClkStatusText,        0,              Button2,        spawn,          {.v = toggle_screen        } },
  { ClkStatusText,        0,              Button3,        spawn,          {.v = toggle_sys_shortcuts } },
  { ClkClientWin,         MODKEY,         Button1,        movemouse,      {0                         } },
  { ClkClientWin,         MODKEY,         Button2,        togglefloating, {0                         } },
  { ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0                         } },
  { ClkTagBar,            0,              Button1,        view,           {0                         } },
  { ClkTagBar,            0,              Button3,        toggleview,     {0                         } },
  { ClkTagBar,            MODKEY,         Button1,        tag,            {0                         } },
  { ClkTagBar,            MODKEY,         Button3,        toggletag,      {0                         } },
};
