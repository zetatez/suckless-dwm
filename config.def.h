/* See LICENSE file for copyright and license details. */

#define SUPKEY  Mod4Mask
#define MODKEY  Mod1Mask
#define Termi(cmd) (const char *[]){ "st", "-e", "/bin/sh", "-c", cmd, NULL }
#define Shell(cmd)    (const char *[]){ "/bin/sh", "-c", cmd, NULL }

/* appearance */
static const unsigned int borderpx = 1;
static const unsigned int snap     = 0;
static const int scalepreview      = 3;
static const int previewbar        = 1;
static const int showbar           = 1;
static const int topbar            = 1;
static const int barheight         = 18;
static const int vertpad           = 0;
static const int sidepad           = 0;
static const int defaultwinpad     = 0;
static const int swallowfloating   = 1;
static const char *fonts[]         = { "DejaVuSansMono Nerd Font:style=Book:size=18" };
static const char dmenufont[]      = "DejaVuSansMono Nerd Font:style=Book:size=16";
static const char col_gray1[]      = "#222222";
static const char col_gray2[]      = "#444444";
static const char col_gray3[]      = "#bbbbbb";
static const char col_gray4[]      = "#eeeeee";
static const char col_cyan[]       = "#005577";
static const char *colors[][3]     = {
  /*               fg         bg         border   */
  [SchemeNorm] = { col_gray3, col_cyan, col_gray2 },
  [SchemeSel]  = { col_gray4, col_cyan, col_cyan },
};

static const char *const autostart[] = {
  "dwmblocks", NULL,
  "picom"    , NULL,
  "hhkb"     , NULL,
  "autostart", NULL,
  NULL,
};

/* tagging */
static const char *tags[] = { "", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix" };

static const Rule rules[] = {
  /* cls                     instance    title      tags mask     isfloating    isterminal     noswallow    monitor */
  {"floatwindow",            NULL,       NULL,      0,            1,            0,             0,           -1 },
  {"st",                     NULL,       NULL,      0,            0,            1,             1,           -1 },
  {"qutebrowser",            NULL,       NULL,      0,            0,            0,             1,           -1 }, // for markdown.
//{"netease-cloud-music",    NULL,       NULL,      1<<8,         0,            0,             0,           -1 },
};

/* stickyicon */
static const XPoint stickyicon[]    = { {0,0}, {4,0}, {4,8}, {2,6}, {0,8}, {0,0} }; /* stickyicon: represents the icon as an array of vertices */
static const XPoint stickyiconbb    = {4,8};	                                      /* stickyicon: defines the bottom right corner of the polygon's bounding box (speeds up scaling) */

/* layout(s) */
static const float mfact            = 0.50;
static const float hfact            = 0.50;
static const int nmaster            = 1;
static const int maxnmaster         = 16;
static const int resizehints        = 0;
static const int lockfullscreen     = 0;

static const Layout layouts[] = {
{ "🐧"                 , layout_workflow           }, // default
{ ""                  , layout_fib_spiral         },
{ ""                  , layout_fib_dwindle        },
{ "⧈"                  , layout_center_free_shape  },
{ "⧅"                  , layout_center_equal_ratio },
{ "󱂫"                  , layout_tile_right         }, // 5
{ "󱂪"                  , layout_tile_left          },
{ "󰝘"                  , layout_grid               },
{ "󱇙"                  , layout_grid_gap           },
{ "󱣴"                  , layout_monocle            },
{ ""                  , layout_hacker             }, // 10
{ " |"                , layout_stack_vert         },
{ " ―"                , layout_stack_hori         }, // 12
// { "∅"               , NULL                      }, // no layout , abandon
{ NULL                 , NULL                      },
};

static const Layout overviewlayout = { "󰾍",  layout_overview };

/* commands */
static char dmenumon[2]                                     = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]                               = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray4, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *termcmd[]                                = { "st", NULL };
static const char *scratchpadcmd[]                          = { "st", "-g", "120x32", "-t", "scratchpad", NULL };

static const Key keys[] = {
/*  modifier                      key               function           argument                                    */
{ SUPKEY,                       XK_F1,           spawn,             {.v = Shell("sys-volume-toggle")                     } },
{ SUPKEY,                       XK_F2,           spawn,             {.v = Shell("sys-volume-down")                       } },
{ SUPKEY,                       XK_F3,           spawn,             {.v = Shell("sys-volume-up")                         } },
{ SUPKEY,                       XK_F4,           spawn,             {.v = Shell("sys-micro-toggle")                      } },
{ SUPKEY,                       XK_F5,           spawn,             {.v = Shell("sys-micro-down")                        } },
{ SUPKEY,                       XK_F6,           spawn,             {.v = Shell("sys-micro-up")                          } },
{ SUPKEY,                       XK_F7,           spawn,             {.v = Shell("sys-wifi-connect")                      } },
{ SUPKEY,                       XK_F8,           spawn,             {.v = Shell("sys-screen")                            } },
{ SUPKEY,                       XK_F9,           spawn,             {.v = Shell("sys-bluetooth")                         } },
{ SUPKEY,                       XK_F10,          spawn,             {.v = Shell("sys-screen-light-down")                 } },
{ SUPKEY,                       XK_F11,          spawn,             {.v = Shell("sys-screen-light-up")                   } },
{ SUPKEY,                       XK_F12,          spawn,             {.v = Shell("sys-toggle-keyboard-light")             } },

{ SUPKEY,                       XK_1,            spawn,             {.v = Shell("qutebrowser-open-url-chatgpt")          } },
{ SUPKEY,                       XK_2,            spawn,             {.v = Shell("qutebrowser-open-url-youtube")          } },
{ SUPKEY,                       XK_3,            spawn,             {.v = Shell("qutebrowser-open-url-github")           } },
{ SUPKEY,                       XK_4,            spawn,             {.v = Shell("qutebrowser-open-url-google-mail")      } },
{ SUPKEY,                       XK_5,            spawn,             {.v = Shell("qutebrowser-open-url-google-translate") } },
{ SUPKEY,                       XK_6,            spawn,             {.v = Shell("qutebrowser-open-url-doubao")           } },
{ SUPKEY,                       XK_7,            spawn,             {.v = Shell("qutebrowser-open-url-google")           } },
{ SUPKEY,                       XK_8,            spawn,             {.v = Shell("qutebrowser-open-url-instagram")        } },
{ SUPKEY,                       XK_9,            spawn,             {.v = Shell("qutebrowser-open-url-leetcode")         } },
{ SUPKEY,                       XK_0,            spawn,             {.v = Shell("qutebrowser-open-url-wechat")           } },

{ SUPKEY|ShiftMask,             XK_1,            spawn,             {.v = Shell("chrome-open-url-chatgpt")          } },
{ SUPKEY|ShiftMask,             XK_2,            spawn,             {.v = Shell("chrome-open-url-youtube")          } },
{ SUPKEY|ShiftMask,             XK_3,            spawn,             {.v = Shell("chrome-open-url-github")           } },
{ SUPKEY|ShiftMask,             XK_4,            spawn,             {.v = Shell("chrome-open-url-google-mail")      } },
{ SUPKEY|ShiftMask,             XK_5,            spawn,             {.v = Shell("chrome-open-url-google-translate") } },
{ SUPKEY|ShiftMask,             XK_6,            spawn,             {.v = Shell("chrome-open-url-doubao")           } },
{ SUPKEY|ShiftMask,             XK_7,            spawn,             {.v = Shell("chrome-open-url-google")           } },
{ SUPKEY|ShiftMask,             XK_8,            spawn,             {.v = Shell("chrome-open-url-instagram")        } },
{ SUPKEY|ShiftMask,             XK_9,            spawn,             {.v = Shell("chrome-open-url-leetcode")         } },
{ SUPKEY|ShiftMask,             XK_0,            spawn,             {.v = Shell("chrome-open-url-wechat")           } },

{ SUPKEY,                       XK_k,            movewin,           {.ui = UP                                            } },
{ SUPKEY,                       XK_j,            movewin,           {.ui = DOWN                                          } },
{ SUPKEY,                       XK_h,            movewin,           {.ui = LEFT                                          } },
{ SUPKEY,                       XK_l,            movewin,           {.ui = RIGHT                                         } },
{ SUPKEY|ShiftMask,             XK_k,            resizewin,         {.ui = VECINC                                        } },
{ SUPKEY|ShiftMask,             XK_j,            resizewin,         {.ui = VECDEC                                        } },
{ SUPKEY|ShiftMask,             XK_h,            resizewin,         {.ui = HORDEC                                        } },
{ SUPKEY|ShiftMask,             XK_l,            resizewin,         {.ui = HORINC                                        } },

{ SUPKEY,                       XK_a,            spawn,             {.v = Shell("launch-qutebrowser")                    } },
{ SUPKEY,                       XK_b,            spawn,             {.v = Shell("launch-chrome")                         } },
{ SUPKEY,                       XK_c,            spawn,             {.v = Shell("toggle-calendar-scheduling")            } },
{ SUPKEY,                       XK_d,            spawn,             {.v = Shell("toggle-lazydocker")                     } },
{ SUPKEY,                       XK_e,            spawn,             {.v = Shell("toggle-mutt")                           } },
{ SUPKEY,                       XK_f,            togglefullscreen,  {0                                                   } },
{ SUPKEY,                       XK_g,            spawn,             {.v = Shell("toggle-lazygit")                        } },
{ SUPKEY,                       XK_i,            spawn,             {.v = Shell("toggle-flameshot")                      } },
{ SUPKEY,                       XK_m,            spawn,             {.v = Termi("lazy-open-search-file-content")         } },
{ SUPKEY,                       XK_n,            spawn,             {.v = Shell("toggle-python")                         } },
{ SUPKEY,                       XK_o,            spawn,             {.v = Shell("handle-copied")                         } },
{ SUPKEY,                       XK_p,            spawn,             {.v = Termi("lazy-open-search-book")                 } },
{ SUPKEY,                       XK_q,            spawn,             {.v = Shell("systemctl suspend && slock")            } },
{ SUPKEY,                       XK_r,            spawn,             {.v = Shell("toggle-yazi")                           } },
{ SUPKEY,                       XK_s,            spawn,             {.v = Shell("search")                                } },
{ SUPKEY,                       XK_t,            spawn,             {.v = Termi("lazy-open-search-file")                 } },
{ SUPKEY,                       XK_u,            spawn,             {.v = Shell("toggle-calendar-scheduling-today")      } },
{ SUPKEY,                       XK_v,            spawn,             {.v = Termi("lazy-open-search-media")                } },
{ SUPKEY,                       XK_w,            spawn,             {.v = Termi("lazy-open-search-wiki")                 } },
{ SUPKEY,                       XK_x,            spawn,             {.v = Shell("toggle-sublime")                        } },
{ SUPKEY,                       XK_y,            spawn,             {.v = Shell("toggle-show")                           } },
{ SUPKEY,                       XK_z,            spawn,             {.v = Shell("chrome-open-url-google")                } },
{ SUPKEY,                       XK_Escape,       spawn,             {.v = Shell("toggle-top")                            } },
{ SUPKEY,                       XK_Delete,       spawn,             {.v = Shell("sys-shortcuts")                         } },
{ SUPKEY,                       XK_BackSpace,    spawn,             {.v = Shell("toggle-passmenu")                       } },
{ SUPKEY,                       XK_backslash,    spawn,             {.v = Shell("set-keyboard-rate")                     } },
{ SUPKEY,                       XK_semicolon,    spawn,             {.v = Shell("jump-to-code-from-log")                 } },
{ SUPKEY,                       XK_comma,        spawn,             {.v = Shell("note-diary")                            } },
{ SUPKEY,                       XK_period,       spawn,             {.v = Shell("note-scripts")                          } },
{ SUPKEY,                       XK_slash,        spawn,             {.v = Shell("note-timeline")                         } },
{ SUPKEY,                       XK_apostrophe,   spawn,             {.v = Shell("toggle-tty-clock")                      } },
// { SUPKEY,                       XK_Home,         spawn,             {.v =                                             } },
// { SUPKEY,                       XK_bracketleft,  spawn,             {.v =                                             } },
// { SUPKEY,                       XK_bracketleft,  spawn,             {.v =                                             } },

{ SUPKEY|ShiftMask,             XK_a,            spawn,             {.v = Shell("toggle-addressbook")                    } },
{ SUPKEY|ShiftMask,             XK_c,            killclient,        {0                                                   } },
{ SUPKEY|ShiftMask,             XK_i,            spawn,             {.v = Shell("toggle-inkscape")                       } },
{ SUPKEY|ShiftMask,             XK_m,            spawn,             {.v = Shell("toggle-music-net-cloud")                } },
{ SUPKEY|ShiftMask,             XK_n,            spawn,             {.v = Shell("toggle-julia")                          } },
{ SUPKEY|ShiftMask,             XK_o,            spawn,             {.v = Shell("toggle-obsidian")                       } },
{ SUPKEY|ShiftMask,             XK_p,            spawn,             {.v = Shell("toggle-krita")                          } },
{ SUPKEY|ShiftMask,             XK_r,            spawn,             {.v = Shell("toggle-redshift")                       } },
{ SUPKEY|ShiftMask,             XK_x,            spawn,             {.v = Shell("toggle-wallpaper")                      } },
{ SUPKEY|ShiftMask,             XK_z,            spawn,             {.v = Shell("qutebrowser-open-url-google")           } },
{ SUPKEY|ShiftMask,             XK_Delete,       spawn,             {.v = Shell("systemctl poweroff")                    } },
{ SUPKEY|ShiftMask,             XK_apostrophe,   spawn,             {.v = Shell("toggle-screenkey")                      } },
{ SUPKEY|ShiftMask,             XK_comma,        spawn,             {.v = Shell("toggle-rec-audio")                      } },
{ SUPKEY|ShiftMask,             XK_period,       spawn,             {.v = Shell("toggle-rec-screen")                     } },
{ SUPKEY|ShiftMask,             XK_slash,        spawn,             {.v = Shell("toggle-rec-webcam")                     } },
{ SUPKEY|ShiftMask,             XK_Return,       spawn,             {.v = termcmd                                        } },
// { SUPKEY|ShiftMask,             XK_b,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_d,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_e,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_f,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_g,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_q,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_s,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_t,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_u,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_v,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_w,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_y,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_Home,         spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_End,          spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_Escape,       spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_BackSpace,    spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_bracketleft,  spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_bracketleft,  spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_backslash,    spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_semicolon,    spawn,             {.v =                                             } },

// MODKEY, etc
{ MODKEY,                       XK_u,            spawn,             {.v = Termi("lazy-open-file")                        } },
{ MODKEY,                       XK_q,            spawn,             {.v = Shell("slock")                                 } },
{ MODKEY,                       XK_apostrophe,   togglescratch,     {.v = scratchpadcmd                                  } },
{ MODKEY,                       XK_c,            spawn,             {.v = Shell("toggle-clipmenu")                       } },
{ MODKEY,                       XK_p,            spawn,             {.v = dmenucmd                                       } },
{ MODKEY,                       XK_Return,       zoom,              {0                                                   } },
{ MODKEY,                       XK_Tab,          view,              {0                                                   } },
{ MODKEY,                       XK_b,            togglebar,         {0                                                   } },
{ MODKEY,                       XK_f,            togglefullscreen,  {0                                                   } },
{ MODKEY|ShiftMask,             XK_f,            togglefloating,    {0                                                   } },
{ MODKEY,                       XK_o,            toggleoverview,    {0                                                   } },
{ MODKEY,                       XK_s,            reset,             {0                                                   } },
{ MODKEY|ShiftMask,             XK_s,            togglesticky,      {0                                                   } },
{ MODKEY|ShiftMask,             XK_space,        focusmaster,       {0                                                   } },
{ MODKEY,                       XK_minus,        scratchpad_show,   {0                                                   } },
{ MODKEY|ShiftMask,             XK_minus,        scratchpad_hide,   {0                                                   } },
{ MODKEY,                       XK_equal,        scratchpad_remove, {0                                                   } },
{ MODKEY,                       XK_bracketleft,  focusmon,          {.i = +1                                             } }, // monitor related
{ MODKEY,                       XK_bracketleft,  focusmon,          {.i = -1                                             } }, // monitor related, not tested
{ MODKEY|ShiftMask,             XK_bracketleft,  tagmon,            {.i = +1                                             } }, // monitor related
{ MODKEY|ShiftMask,             XK_bracketleft,  tagmon,            {.i = -1                                             } }, // monitor related, not tested
{ MODKEY,                       XK_d,            incnmaster,        {.i = -1                                             } },
{ MODKEY,                       XK_i,            incnmaster,        {.i = +1                                             } },
{ MODKEY,                       XK_h,            movestack,         {.i = -1                                             } },
{ MODKEY,                       XK_l,            movestack,         {.i = +1                                             } },
{ MODKEY,                       XK_comma,        movestack,         {.i = -1                                             } },
{ MODKEY,                       XK_period,       movestack,         {.i = +1                                             } },
{ MODKEY|ShiftMask,             XK_comma,        shiftview,         {.i = -1                                             } },
{ MODKEY|ShiftMask,             XK_period,       shiftview,         {.i = +1                                             } },
{ MODKEY|ControlMask,           XK_comma,        cyclelayout,       {.i = -1                                             } }, // useless
{ MODKEY|ControlMask,           XK_period,       cyclelayout,       {.i = +1                                             } }, // useless
{ MODKEY,                       XK_k,            focusstack,        {.i = -1                                             } },
{ MODKEY,                       XK_j,            focusstack,        {.i = +1                                             } },
{ MODKEY|ShiftMask,             XK_h,            setmfact,          {.f = -0.025                                         } },
{ MODKEY|ShiftMask,             XK_l,            setmfact,          {.f = +0.025                                         } },
{ MODKEY|ShiftMask,             XK_j,            sethfact,          {.f = -0.025                                         } },
{ MODKEY|ShiftMask,             XK_k,            sethfact,          {.f = +0.025                                         } },
{ MODKEY,                       XK_space,        togglefloating,    {0                                                   } },
{ MODKEY,                       XK_u,            setlayout,         {0                                                   } }, // temporary layout switch
{ MODKEY,                       XK_a,            setlayout,         {.v = &layouts[0]                                    } },
{ MODKEY,                       XK_r,            setlayout,         {.v = &layouts[1]                                    } },
{ MODKEY|ShiftMask,             XK_r,            setlayout,         {.v = &layouts[2]                                    } },
{ MODKEY,                       XK_v,            setlayout,         {.v = &layouts[3]                                    } },
{ MODKEY|ShiftMask,             XK_v,            setlayout,         {.v = &layouts[4]                                    } },
{ MODKEY,                       XK_t,            setlayout,         {.v = &layouts[5]                                    } },
{ MODKEY|ShiftMask,             XK_t,            setlayout,         {.v = &layouts[6]                                    } },
{ MODKEY,                       XK_g,            setlayout,         {.v = &layouts[7]                                    } },
{ MODKEY|ShiftMask,             XK_g,            setlayout,         {.v = &layouts[8]                                    } },
{ MODKEY,                       XK_m,            setlayout,         {.v = &layouts[9]                                    } },
{ MODKEY|ShiftMask,             XK_m,            setlayout,         {.v = &layouts[10]                                   } },
{ MODKEY|ShiftMask,             XK_e,            setlayout,         {.v = &layouts[11]                                   } },
{ MODKEY,                       XK_e,            setlayout,         {.v = &layouts[12]                                   } },
{ MODKEY|ShiftMask,             XK_Return,       spawn,             {.v = termcmd                                        } },
{ ControlMask|ShiftMask,        XK_Return,       spawn,             {.v = Shell("kitty")                                 } },
{ MODKEY|ShiftMask,             XK_c,            killclient,        {0                                                   } },
{ MODKEY|ShiftMask,             XK_q,            quit,              {0                                                   } },
{ MODKEY|ShiftMask,             XK_p,            quit,              {1                                                   } },
// { MODKEY,                       XK_n,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_w,            xxxxx,             {.v =                                             } },
// { MODKEY,                       XK_x,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_y,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_z,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_slash,        xxxxx,             {.i =                                             } },
// { MODKEY|ShiftMask,             XK_a,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_b,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_d,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_g,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_i,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_n,            xxxxx,             {.v = x                                           } },
// { MODKEY|ShiftMask,             XK_o,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_u,            xxxxx,             {.i =                                             } },
// { MODKEY|ShiftMask,             XK_w,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_x,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_y,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_y,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_z,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_slash,        xxx,               {.i = +1                                          } },

{ MODKEY,                       XK_1,            view,              {.ui = 1 << 0                                        } }, // view tag 1
{ MODKEY,                       XK_2,            view,              {.ui = 1 << 1                                        } },
{ MODKEY,                       XK_3,            view,              {.ui = 1 << 2                                        } },
{ MODKEY,                       XK_4,            view,              {.ui = 1 << 3                                        } },
{ MODKEY,                       XK_5,            view,              {.ui = 1 << 4                                        } },
{ MODKEY,                       XK_6,            view,              {.ui = 1 << 5                                        } },
{ MODKEY,                       XK_7,            view,              {.ui = 1 << 6                                        } },
{ MODKEY,                       XK_8,            view,              {.ui = 1 << 7                                        } },
{ MODKEY,                       XK_9,            view,              {.ui = 1 << 8                                        } },
{ MODKEY,                       XK_0,            view,              {.ui = ~0                                            } }, // preview all tags
{ MODKEY|ShiftMask,             XK_1,            tag,               {.ui = 1 << 0                                        } }, // move to tag 1
{ MODKEY|ShiftMask,             XK_2,            tag,               {.ui = 1 << 1                                        } },
{ MODKEY|ShiftMask,             XK_3,            tag,               {.ui = 1 << 2                                        } },
{ MODKEY|ShiftMask,             XK_4,            tag,               {.ui = 1 << 3                                        } },
{ MODKEY|ShiftMask,             XK_5,            tag,               {.ui = 1 << 4                                        } },
{ MODKEY|ShiftMask,             XK_6,            tag,               {.ui = 1 << 5                                        } },
{ MODKEY|ShiftMask,             XK_7,            tag,               {.ui = 1 << 6                                        } },
{ MODKEY|ShiftMask,             XK_8,            tag,               {.ui = 1 << 7                                        } },
{ MODKEY|ShiftMask,             XK_9,            tag,               {.ui = 1 << 8                                        } },
{ MODKEY|ShiftMask,             XK_0,            tag,               {.ui = ~0                                            } }, // stick to all tags
{ MODKEY|ControlMask,           XK_1,            toggleview,        {.ui = 1 << 0                                        } }, // toggle view of tag 1
{ MODKEY|ControlMask,           XK_2,            toggleview,        {.ui = 1 << 1                                        } },
{ MODKEY|ControlMask,           XK_3,            toggleview,        {.ui = 1 << 2                                        } },
{ MODKEY|ControlMask,           XK_4,            toggleview,        {.ui = 1 << 3                                        } },
{ MODKEY|ControlMask,           XK_5,            toggleview,        {.ui = 1 << 4                                        } },
{ MODKEY|ControlMask,           XK_6,            toggleview,        {.ui = 1 << 5                                        } },
{ MODKEY|ControlMask,           XK_7,            toggleview,        {.ui = 1 << 6                                        } },
{ MODKEY|ControlMask,           XK_8,            toggleview,        {.ui = 1 << 7                                        } },
{ MODKEY|ControlMask,           XK_9,            toggleview,        {.ui = 1 << 8                                        } },
{ MODKEY|ControlMask|ShiftMask, XK_1,            previewtag,        {.ui = 0                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_2,            previewtag,        {.ui = 1                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_3,            previewtag,        {.ui = 2                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_4,            previewtag,        {.ui = 3                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_5,            previewtag,        {.ui = 4                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_6,            previewtag,        {.ui = 5                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_7,            previewtag,        {.ui = 6                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_8,            previewtag,        {.ui = 7                                             } },
{ MODKEY|ControlMask|ShiftMask, XK_9,            previewtag,        {.ui = 8                                             } },
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
{ ClkTagBar,            0,              Button1,        view,                {0                                          } },
{ ClkTagBar,            0,              Button3,        toggleview,          {0                                          } },
{ ClkTagBar,            MODKEY,         Button1,        tag,                 {0                                          } },
{ ClkTagBar,            MODKEY,         Button3,        toggletag,           {0                                          } },
{ ClkLtSymbol,          0,              Button1,        setlayout,           {0                                          } },
{ ClkLtSymbol,          0,              Button2,        setlayout,           {.v = &layouts[8]                           } },
{ ClkLtSymbol,          0,              Button3,        setlayout,           {.v = &overviewlayout                       } },
{ ClkStatusText,        0,              Button1,        spawn,               {.v = Shell("toggle-tty-clock")             } },
{ ClkStatusText,        0,              Button2,        spawn,               {.v = Shell("sys-shortcuts")                } },
{ ClkStatusText,        0,              Button3,        spawn,               {.v = Shell("toggle-calendar")              } },
{ ClkClientWin,         MODKEY,         Button1,        movemouse,           {0                                          } },
{ ClkClientWin,         MODKEY,         Button2,        togglefloating,      {0                                          } },
{ ClkClientWin,         MODKEY,         Button3,        resizemouse,         {0                                          } },
// { ClkWinTitle,          0,              Button1,        xxxxxxxxx,           {0                                       } },
// { ClkWinTitle,          0,              Button2,        xxxxxxxxx,           {0                                       } },
// { ClkWinTitle,          0,              Button3,        xxxxxxxxx,           {0                                       } },
};
