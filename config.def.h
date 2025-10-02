/* See LICENSE file for copyright and license details. */

#define SUPKEY  Mod4Mask
#define MODKEY  Mod1Mask
#define LaunchApp(cmd) (const char *[]){cmd, NULL}
#define Shell(cmd)     (const char *[]){"/bin/sh", "-c", cmd, NULL}
#define Termi(cmd)     (const char *[]){"st", "-e", "/bin/sh", "-c", cmd, NULL}

/* appearance */
static const unsigned int borderpx = 1;
static const unsigned int snap     = 0;
static const int scalepreview      = 3;
static const int previewbar        = 1;
static const int showbar           = 1;
static const int topbar            = 1;
static const int barheight         = 15;
static const int vertpad           = 0;
static const int sidepad           = 0;
static const int defaultwinpad     = 1;
static const int swallowfloating   = 1;
// static const char *fonts[]         = { "DejaVuSansMono Nerd Font:style=Book:size=17" };
static const char *fonts[]         = { "DejaVuSansMono Nerd Font:style=Book:size=17" };
static const char dmenufont[]      = "DejaVuSansMono Nerd Font:style=Book:size=24";
static const char col_gray1[]      = "#222222";
static const char col_gray2[]      = "#444444";
static const char col_gray3[]      = "#bbbbbb";
static const char col_gray4[]      = "#eeeeee";
static const char col_cyan[]       = "#023047"; // #005577
static const char col_bg[]         = "#0077b6";
static const char col_fg[]         = "#00b4d8";
static const char *colors[][3]     = {
//               fg         bg        border
// [SchemeNorm] = { col_gray3, col_cyan,  col_gray2 },
// [SchemeSel]  = { col_gray4, col_cyan,  col_cyan  },
   [SchemeNorm] = { col_bg   , col_cyan, col_gray2 },
   [SchemeSel]  = { col_fg   , col_cyan, col_cyan  },
};

static const char *const autostart[] = {
  "dwmblocks"        , NULL ,
  "reset_sys_default", NULL ,
  "daemon"           , NULL ,
  NULL,
};

/* tagging */
static const char *tags[] = { "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix" };

static const Rule rules[] = {
  /* cls                     instance    title      tags mask     isfloating    isterminal     noswallow    monitor */
  {"floatwindow",            NULL,       NULL,      0,            1,            0,             0,           -1 },
  {"st",                     NULL,       NULL,      0,            0,            1,             1,           -1 },
  {"Surf",                   NULL,       NULL,      0,            0,            0,             1,           -1 }, // no swallow for markdown
  {"chrome",                 NULL,       NULL,      0,            0,            0,             0,           -1 },
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
static const int refreshrate        = 120;  /* refresh rate (per second) for client move/resize */

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
{ "S |"                , layout_stack_vert         },
{ "S ―"                , layout_stack_hori         }, // 12
// { "∅"               , NULL                      }, // no layout , abandon
{ NULL                 , NULL                      },
};

static const Layout overviewlayout = { "󰾍",  layout_overview };

/* commands */
static char dmenumon[2]                                     = "0"; /* component of dmenucmd, manipulated in spawn() */
// static const char *dmenucmd[]                            = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray4, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *dmenucmd[]                               = { "rofi", "-show", "drun", "-theme", "fullscreen-preview", "-font", "JetBrainsMono Nerd Font 24", NULL };
static const char *scratchpadcmd[]                          = { "st", "-g", "120x32", "-t", "scratchpad", NULL };

static const Key keys[] = {
/*  modifier                    key              function           argument                                              */

{ SUPKEY,                       XK_F1,           spawn,             {.v = Shell("sys_volume_toggle")                     } },
{ SUPKEY,                       XK_F2,           spawn,             {.v = Shell("sys_volume_down")                       } },
{ SUPKEY,                       XK_F3,           spawn,             {.v = Shell("sys_volume_up")                         } },
{ SUPKEY,                       XK_F4,           spawn,             {.v = Shell("sys_micro_toggle")                      } },
{ SUPKEY,                       XK_F5,           spawn,             {.v = Shell("sys_screen_light_down")                 } },
{ SUPKEY,                       XK_F6,           spawn,             {.v = Shell("sys_screen_light_up")                   } },
{ SUPKEY,                       XK_F7,           spawn,             {.v = Shell("sys_screen")                            } },
{ SUPKEY,                       XK_F8,           spawn,             {.v = Shell("sys_wifi_connect")                      } },
{ SUPKEY,                       XK_F9,           spawn,             {.v = Shell("sys_bluetooth_connect")                 } },
{ SUPKEY,                       XK_F10,          spawn,             {.v = Shell("sys_micro_down")                        } },
{ SUPKEY,                       XK_F11,          spawn,             {.v = Shell("sys_micro_up")                          } },
{ SUPKEY,                       XK_F12,          spawn,             {.v = Shell("sys_toggle_keyboard_light")             } },

{ SUPKEY,                       XK_1,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://chatgpt.com/'")                           } },
{ SUPKEY,                       XK_2,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://www.youtube.com'")                        } },
{ SUPKEY,                       XK_3,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://github.com/zetatez'")                     } },
{ SUPKEY,                       XK_4,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://mail.google.com/mail'")                   } },
{ SUPKEY,                       XK_5,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://translate.google.com/?sl=auto&tl=zh-CN'") } },
{ SUPKEY,                       XK_6,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://www.doubao.com/chat/'")                   } },
{ SUPKEY,                       XK_7,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://www.bilibili.com/'")                      } },
{ SUPKEY,                       XK_8,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://web.wechat.com/'")                        } },
{ SUPKEY,                       XK_9,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='https://leetcode.cn/search/?q=%E6%9C%80'")        } },
{ SUPKEY,                       XK_0,            spawn,             {.v = Shell("open_url_with_qutebrowser --url='http://www.google.com/'")                         } },

{ SUPKEY|ShiftMask,             XK_1,            spawn,             {.v = Shell("open_url_with_chrome --url='https://chatgpt.com/'")                                } },
{ SUPKEY|ShiftMask,             XK_2,            spawn,             {.v = Shell("open_url_with_chrome --url='https://www.youtube.com'")                             } },
{ SUPKEY|ShiftMask,             XK_3,            spawn,             {.v = Shell("open_url_with_chrome --url='https://github.com/zetatez'")                          } },
{ SUPKEY|ShiftMask,             XK_4,            spawn,             {.v = Shell("open_url_with_chrome --url='https://mail.google.com/mail'")                        } },
{ SUPKEY|ShiftMask,             XK_5,            spawn,             {.v = Shell("open_url_with_chrome --url='https://translate.google.com/?sl=auto&tl=zh-CN'")      } },
{ SUPKEY|ShiftMask,             XK_6,            spawn,             {.v = Shell("open_url_with_chrome --url='https://www.doubao.com/chat/'")                        } },
{ SUPKEY|ShiftMask,             XK_7,            spawn,             {.v = Shell("open_url_with_chrome --url='https://www.bilibili.com/'")                           } },
{ SUPKEY|ShiftMask,             XK_8,            spawn,             {.v = Shell("open_url_with_chrome --url='https://web.wechat.com/'")                             } },
{ SUPKEY|ShiftMask,             XK_9,            spawn,             {.v = Shell("open_url_with_chrome --url='https://leetcode.cn/search/?q=%E6%9C%80'")             } },
{ SUPKEY|ShiftMask,             XK_0,            spawn,             {.v = Shell("open_url_with_chrome --url='http://www.google.com/'")                              } },

{ SUPKEY,                       XK_k,            movewin,           {.ui = UP                                            } },
{ SUPKEY,                       XK_j,            movewin,           {.ui = DOWN                                          } },
{ SUPKEY,                       XK_h,            movewin,           {.ui = LEFT                                          } },
{ SUPKEY,                       XK_l,            movewin,           {.ui = RIGHT                                         } },
{ SUPKEY|ShiftMask,             XK_k,            resizewin,         {.ui = VECINC                                        } },
{ SUPKEY|ShiftMask,             XK_j,            resizewin,         {.ui = VECDEC                                        } },
{ SUPKEY|ShiftMask,             XK_h,            resizewin,         {.ui = HORDEC                                        } },
{ SUPKEY|ShiftMask,             XK_l,            resizewin,         {.ui = HORINC                                        } },

{ SUPKEY,                       XK_a,            spawn,             {.v = Shell("launch_file_manager")                   } },
{ SUPKEY,                       XK_b,            spawn,             {.v = Shell("launch_qutebrowser")                    } }, // launch_chrome
{ SUPKEY,                       XK_c,            spawn,             {.v = Shell("note_timeline")                         } },
{ SUPKEY,                       XK_d,            spawn,             {.v = Shell("toggle_lazydocker")                     } },
{ SUPKEY,                       XK_e,            spawn,             {.v = Shell("toggle_mutt")                           } },
{ SUPKEY,                       XK_f,            togglefullscreen,  {0                                                   } },
{ SUPKEY,                       XK_g,            spawn,             {.v = Shell("toggle_lazygit")                        } },
{ SUPKEY,                       XK_i,            spawn,             {.v = Shell("toggle_flameshot")                      } },
{ SUPKEY,                       XK_m,            spawn,             {.v = Termi("lazy_open_search_file_content")         } },
{ SUPKEY,                       XK_n,            spawn,             {.v = Shell("toggle_python")                         } },
{ SUPKEY,                       XK_o,            spawn,             {.v = Shell("handle_copied")                         } },
{ SUPKEY,                       XK_p,            spawn,             {.v = Termi("lazy_open_search_book")                 } },
{ SUPKEY,                       XK_q,            spawn,             {.v = Shell("slock")                                 } },
{ SUPKEY,                       XK_r,            spawn,             {.v = Shell("toggle_yazi")                           } },
{ SUPKEY,                       XK_s,            spawn,             {.v = Shell("search")                                } },
{ SUPKEY,                       XK_t,            spawn,             {.v = Termi("lazy_open_search_file")                 } },
{ SUPKEY,                       XK_u,            spawn,             {.v = Termi("lazy_open_search_media")                } },
{ SUPKEY,                       XK_v,            spawn,             {.v = Shell("note_diary")                            } },
{ SUPKEY,                       XK_w,            spawn,             {.v = Termi("lazy_open_search_wiki")                 } },
{ SUPKEY,                       XK_x,            spawn,             {.v = Shell("note_scripts")                          } },
{ SUPKEY,                       XK_y,            spawn,             {.v = Shell("toggle_show")                           } },
{ SUPKEY,                       XK_z,            spawn,             {.v = Shell("note_todo")                             } },
{ SUPKEY,                       XK_Escape,       spawn,             {.v = Shell("toggle_top")                            } },
{ SUPKEY,                       XK_Delete,       spawn,             {.v = Shell("sys_shortcuts")                         } },
{ SUPKEY,                       XK_BackSpace,    spawn,             {.v = Shell("toggle_passmenu")                       } },
{ SUPKEY,                       XK_backslash,    spawn,             {.v = Shell("reset_sys_default")                     } },
{ SUPKEY,                       XK_semicolon,    spawn,             {.v = Shell("jump_to_code_from_log")                 } },
{ SUPKEY,                       XK_apostrophe,   spawn,             {.v = Shell("toggle_tty_clock")                      } },
{ SUPKEY,                       XK_bracketleft,  spawn,             {.v = Shell("toggle_calendar_scheduling")            } },
{ SUPKEY,                       XK_bracketright, spawn,             {.v = Shell("toggle_calendar_scheduling_today")      } },

// { SUPKEY,                       XK_Home,         spawn,             {.v =                                             } },
// { SUPKEY,                       XK_comma,        spawn,             {.v =                                             } },
// { SUPKEY,                       XK_period,       spawn,             {.v =                                             } },
// { SUPKEY,                       XK_slash,        spawn,             {.v =                                             } },
{ SUPKEY|ShiftMask,             XK_b,            spawn,             {.v = Shell("launch_chrome")                         } },
{ SUPKEY|ShiftMask,             XK_c,            killclient,        {0                                                   } },
{ SUPKEY|ShiftMask,             XK_s,            spawn,             {.v = Shell("toggle_sublime")                        } },
{ SUPKEY|ShiftMask,             XK_i,            spawn,             {.v = Shell("toggle_inkscape")                       } },
{ SUPKEY|ShiftMask,             XK_m,            spawn,             {.v = Shell("toggle_music_net_cloud")                } },
{ SUPKEY|ShiftMask,             XK_n,            spawn,             {.v = Shell("toggle_julia")                          } },
{ SUPKEY|ShiftMask,             XK_o,            spawn,             {.v = Shell("toggle_obsidian")                       } },
{ SUPKEY|ShiftMask,             XK_p,            spawn,             {.v = Shell("toggle_krita")                          } },
{ SUPKEY|ShiftMask,             XK_x,            spawn,             {.v = Shell("toggle_wallpaper")                      } },
{ SUPKEY|ShiftMask,             XK_Delete,       spawn,             {.v = Shell("systemctl poweroff")                    } },
{ SUPKEY|ShiftMask,             XK_comma,        spawn,             {.v = Shell("toggle_rec_audio")                      } },
{ SUPKEY|ShiftMask,             XK_apostrophe,   spawn,             {.v = Shell("toggle_screenkey")                      } },
{ SUPKEY|ShiftMask,             XK_period,       spawn,             {.v = Shell("toggle_rec_screen")                     } },
{ SUPKEY|ShiftMask,             XK_slash,        spawn,             {.v = Shell("toggle_rec_webcam")                     } },
{ SUPKEY|ShiftMask,             XK_Return,       spawn,             {.v = LaunchApp("kitty")                             } },
// { SUPKEY|ShiftMask,             XK_a,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_d,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_e,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_f,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_g,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_q,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_r,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_s,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_t,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_u,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_v,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_w,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_y,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_z,            spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_Home,         spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_End,          spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_Escape,       spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_BackSpace,    spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_bracketleft,  spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_bracketright, spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_semicolon,    spawn,             {.v =                                             } },
// { SUPKEY|ShiftMask,             XK_backslash,    spawn,             {.v =                                             } },

// MODKEY, etc
{ MODKEY,                       XK_semicolon,    spawn,             {.v = Shell("rofi -show run -theme fullscreen-preview -font 'JetBrainsMono Nerd Font 24'") } },
{ MODKEY,                       XK_p,            spawn,             {.v = dmenucmd                                       } },
{ MODKEY,                       XK_apostrophe,   togglescratch,     {.v = scratchpadcmd                                  } },
{ MODKEY,                       XK_q,            spawn,             {.v = Shell("slock")                                 } },
{ MODKEY,                       XK_c,            spawn,             {.v = Shell("toggle_clipmenu")                       } },
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
{ MODKEY,                       XK_bracketleft,  focusmon,          {.i = -1                                             } }, // multi monitors: focus on which one -1
{ MODKEY,                       XK_bracketright, focusmon,          {.i = +1                                             } }, // multi monitors: focus on which one +1
{ MODKEY|ShiftMask,             XK_bracketleft,  tagmon,            {.i = -1                                             } }, // multi monitors: move win to monitor prev
{ MODKEY|ShiftMask,             XK_bracketright, tagmon,            {.i = +1                                             } }, // multi monitors: move win to monitor next
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
{ MODKEY,                       XK_u,            setlayout,         {0                                                   } }, // temporary layout switch
{ MODKEY,                       XK_space,        togglefloating,    {0                                                   } },
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
{ MODKEY,                       XK_w,            setlayout,         {.v = &layouts[10]                                   } },
{ MODKEY|ShiftMask,             XK_e,            setlayout,         {.v = &layouts[11]                                   } },
{ MODKEY,                       XK_e,            setlayout,         {.v = &layouts[12]                                   } },
{ MODKEY|ShiftMask,             XK_Return,       spawn,             {.v = LaunchApp("st")                                } },
{ MODKEY|ShiftMask,             XK_c,            killclient,        {0                                                   } },
{ MODKEY|ShiftMask|ControlMask, XK_c,            killclient_unsel,  {0                                                   } },
{ MODKEY|ShiftMask,             XK_q,            quit,              {0                                                   } },
{ MODKEY|ShiftMask,             XK_p,            quit,              {1                                                   } },
{ MODKEY,                       XK_slash,        spawn,             {.v = Termi("lazy_open_search_file")                 } },
// { MODKEY,                       XK_n,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_w,            xxxxx,             {.v =                                             } },
// { MODKEY,                       XK_x,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_y,            xxxxx,             {.v = x                                           } },
// { MODKEY,                       XK_z,            xxxxx,             {.v = x                                           } },
// { MODKEY|ShiftMask,             XK_a,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_b,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_d,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_g,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_i,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_n,            xxxxx,             {.v = x                                           } },
// { MODKEY|ShiftMask,             XK_o,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_u,            xxxxx,             {.i =                                             } },
// { MODKEY|ShiftMask,             XK_w,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_m,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_x,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_y,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_y,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_z,            xxxxx,             {.v =                                             } },
// { MODKEY|ShiftMask,             XK_slash,        xxx,               {.i =                                             } },

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
{ ClkTagBar,            0,              Button1,        view,                {0                                          } },
{ ClkTagBar,            0,              Button3,        toggleview,          {0                                          } },
{ ClkTagBar,            MODKEY,         Button1,        tag,                 {0                                          } },
{ ClkTagBar,            MODKEY,         Button3,        toggletag,           {0                                          } },
{ ClkLtSymbol,          0,              Button1,        setlayout,           {0                                          } },
{ ClkLtSymbol,          0,              Button2,        setlayout,           {.v = &layouts[8]                           } },
{ ClkLtSymbol,          0,              Button3,        setlayout,           {.v = &overviewlayout                       } },
{ ClkStatusText,        0,              Button1,        spawn,               {.v = Shell("toggle_tty_clock")             } },
{ ClkStatusText,        0,              Button2,        spawn,               {.v = Shell("sys_shortcuts")                } },
{ ClkStatusText,        0,              Button3,        spawn,               {.v = Shell("toggle_calendar")              } },
{ ClkClientWin,         MODKEY,         Button1,        movemouse,           {0                                          } },
{ ClkClientWin,         MODKEY,         Button2,        togglefloating,      {0                                          } },
{ ClkClientWin,         MODKEY,         Button3,        resizemouse,         {0                                          } },
// { ClkWinTitle,          0,              Button1,        xxxxxxxxx,           {0                                       } },
// { ClkWinTitle,          0,              Button2,        xxxxxxxxx,           {0                                       } },
// { ClkWinTitle,          0,              Button3,        xxxxxxxxx,           {0                                       } },
};

