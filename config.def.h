/* See LICENSE file for copyright and license details. */

#define SUPKEY Mod4Mask
#define MODKEY Mod1Mask
#define SH(cmd)   { "/bin/sh", "-c", cmd, NULL }
#define ST(cmd)   { "st", "-e", "/bin/sh", "-c", cmd, NULL }

/* appearance */
static const unsigned int borderpx  = 1;
static const unsigned int snap      = 0;
static const int scalepreview       = 3;
static const int previewbar         = 1;
static const int showbar            = 1;
static const int topbar             = 1;
static const int barheight          = 26;
static const int vertpad            = 0;
static const int sidepad            = 0;
static const int defaultwinpad      = 0;
static const int swallowfloating    = 1;
static const char *fonts[]          = {
  "Source Han Serif CN,思源宋体 CN,Source Han Serif CN ExtraLight,思源宋体 CN ExtraLight:style=ExtraLight,Regular:size=18",
  // "Hack:style=Regular:size=18:antialias=true:autohint=true",
  // "DejaVuSansMono Nerd Font:style=Book:size=16",
};
static const char dmenufont[]       = "DejaVuSansMono Nerd Font:style=Book:size=14";
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
  "picom", NULL,
  "hhkb", NULL,
  "autostart", NULL, // other stuff
  NULL
};

/* tagging */
static const char *tags[] = { "i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix" };

static const Rule rules[] = {
  /* cls                     instance    title      tags mask     isfloating    isterminal     noswallow    monitor */
  {"floatwindow",            NULL,       NULL,      0,            1,            0,             0,           -1 },
  {"st",                     NULL,       NULL,      0,            0,            1,             1,           -1 },
  {"Vimb",                   NULL,       NULL,      0,            0,            0,             1,           -1 }, // for markdown. Vimb V must be capital
  {"netease-cloud-music",    NULL,       NULL,      1<<8,         0,            0,             0,           -1 },
  {"yesplaymusic",           NULL,       NULL,      1<<8,         0,            0,             0,           -1 },
};

/* stickyicon */
static const XPoint stickyicon[]    = { {0,0}, {4,0}, {4,8}, {2,6}, {0,8}, {0,0} }; /* stickyicon: represents the icon as an array of vertices */
static const XPoint stickyiconbb    = {4,8};	                                      /* stickyicon: defines the bottom right corner of the polygon's bounding box (speeds up scaling) */

/* layout(s) */
static const float mfact            = 0.50;
static const float ffact            = 0.50;
static const int nmaster            = 1;
static const int maxnmaster         = 8;
static const int resizehints        = 0;
static const int lockfullscreen     = 0;

static const Layout layouts[] = {
   { "",  layout_workflow            }, // default
   { "🍁 >",  layout_fibonaccispiral     },
   { "🍁 <",  layout_fibonaccidwindle    },
   { "⧈ *",  layout_centerfreeshape     },
   { "⧈ =",  layout_centerequalratio    },
   { "◨ >",  layout_tileright           },
   { "◧ <",  layout_tileleft            },
   { "󰝘 G",  layout_grid                },
   { " F",  layout_monocle             },
   { "󰓌 S",  layout_hacker              },
   { "󱎞 |",  layout_stairs_vert         },
   { "󱎞 ―",  layout_stairs_hori         },
   { "⬓ |" , layout_bottomstackvert     },
   { "⬓ ―",  layout_bottomstackhori     },
// { "◧",  layout_tileright_vertical  },
// { "󰘸",  layout_deckvert            },
// { "󰘸",  layout_deckhori            },
// { "∅",  NULL                       }, // no layout, abandon
   { NULL, NULL                       },
};

// static const Layout overviewlayout = { "󰾍",  layout_overview };
static const Layout overviewlayout = { "󰾍",  layout_overview };

/* commands */
static char dmenumon[2]                = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[]          = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray4, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *termcmd[]           = { "st", NULL };
static const char *termcmd_kitty[]     = { "kitty", NULL };
static const char *scratchpadcmd[]     = { "st", "-g", "120x32", "-t", "scratchpad", NULL }; // patch: dwm-scratchpad
static const char *screen_light_dec[]                 = SH("sudo light -N 8; sudo light -U 5");
static const char *screen_light_inc[]                 = SH("sudo light -A 5");
static const char *screen_lock[]                      = SH("slock");
static const char *keyboard_rate[]                    = SH("xset r rate 158 128 && notify-send 'keyboard rate changed to 158 128'");
static const char *sys_suspend[]                      = SH("systemctl suspend && slock");
static const char *sys_shutdown[]                     = SH("systemctl poweroff");
static const char *volume_dec[]                       = SH("amixer set Speaker unmute; amixer set Speaker 98; amixer set Master 5%-");
static const char *volume_inc[]                       = SH("amixer set Speaker unmute; amixer set Speaker 98; amixer set Master 5%+");
static const char *volume_toggle[]                    = SH("amixer set Speaker unmute; amixer set Speaker 98; amixer set Master toggle");
static const char *microphone_dec[]                   = SH("amixer set Capture 5%-");
static const char *microphone_inc[]                   = SH("amixer set Capture 5%+");
static const char *microphone_toggle[]                = SH("amixer set Capture toggle");
static const char *wifi_connect[]                     = SH("wifi-connect");
static const char *toggle_addressbook[]               = SH("toggle-addressbook");
static const char *toggle_bluetooth[]                 = SH("toggle-bluetooth");
static const char *toggle_calendar[]                  = SH("toggle-calendar");
static const char *toggle_calendar_scheduling[]       = SH("toggle-calendar-scheduling");
static const char *toggle_calendar_scheduling_today[] = SH("toggle-calendar-scheduling-today");
static const char *toggle_clipmenu[]                  = SH("toggle-clipmenu");
static const char *toggle_flameshot[]                 = SH("toggle-flameshot");
static const char *toggle_inkscape[]                  = SH("toggle-inkscape");
static const char *toggle_joshuto[]                   = SH("toggle-joshuto");
static const char *toggle_julia[]                     = SH("toggle-julia");
static const char *toggle_python[]                    = SH("toggle-python");
static const char *toggle_keyboard_light[]            = SH("toggle-keyboard-light");
static const char *toggle_krita[]                     = SH("toggle-krita");
static const char *toggle_lazydocker[]                = SH("toggle-lazydocker");
static const char *toggle_music_net_cloud[]           = SH("toggle-music-net-cloud");
static const char *toggle_mutt[]                      = SH("toggle-mutt");
static const char *toggle_obsidian[]                  = SH("toggle-obsidian");
static const char *toggle_passmenu[]                  = SH("toggle-passmenu");
static const char *toggle_rec_audio[]                 = SH("toggle-rec-audio");
static const char *toggle_rec_screen[]                = SH("toggle-rec-screen");
static const char *toggle_rec_webcam[]                = SH("toggle-rec-webcam");
static const char *toggle_redshift[]                  = SH("toggle-redshift");
static const char *toggle_screen[]                    = SH("toggle-screen");
static const char *toggle_screenkey[]                 = SH("toggle-screenkey");
static const char *toggle_show[]                      = SH("toggle-show");
static const char *toggle_sublime[]                   = SH("toggle-sublime");
static const char *toggle_sys_shortcuts[]             = SH("toggle-sys-shortcuts");
static const char *toggle_termius[]                   = SH("toggle-termius");
static const char *toggle_top[]                       = SH("toggle-top");
static const char *toggle_wallpaper[]                 = SH("toggle-wallpaper");
static const char *launch_chrome[]                    = SH("launch-chrome");
static const char *launch_edge[]                      = SH("launch-edge");
static const char *chrome_open_url_google[]           = SH("chrome-open-url-google");
static const char *chrome_open_url_bing[]             = SH("chrome-open-url-bing");
static const char *chrome_open_url_chatgpt[]          = SH("chrome-open-url-chatgpt");
static const char *chrome_open_url_codeium[]          = SH("chrome-open-url-codeium");
static const char *chrome_open_url_doubao[]           = SH("chrome-open-url-doubao");
static const char *chrome_open_url_github[]           = SH("chrome-open-url-github");
static const char *chrome_open_url_google_mail[]      = SH("chrome-open-url-google-mail");
static const char *chrome_open_url_google_translate[] = SH("chrome-open-url-google-translate");
static const char *chrome_open_url_instagram[]        = SH("chrome-open-url-instagram");
static const char *chrome_open_url_leetcode[]         = SH("chrome-open-url-leetcode");
static const char *chrome_open_url_wechat[]           = SH("chrome-open-url-wechat");
static const char *chrome_open_url_youtube[]          = SH("chrome-open-url-youtube");
// static const char *edge_open_url_google[]             = SH("edge-open-url-google");
// static const char *edge_open_url_bing[]               = SH("edge-open-url-bing");
static const char *edge_open_url_chatgpt[]            = SH("edge-open-url-chatgpt");
static const char *edge_open_url_codeium[]            = SH("edge-open-url-codeium");
static const char *edge_open_url_doubao[]             = SH("edge-open-url-doubao");
static const char *edge_open_url_github[]             = SH("edge-open-url-github");
static const char *edge_open_url_google_mail[]        = SH("edge-open-url-google-mail");
static const char *edge_open_url_google_translate[]   = SH("edge-open-url-google-translate");
static const char *edge_open_url_instagram[]          = SH("edge-open-url-instagram");
static const char *edge_open_url_leetcode[]           = SH("edge-open-url-leetcode");
static const char *edge_open_url_wechat[]             = SH("edge-open-url-wechat");
static const char *edge_open_url_youtube[]            = SH("edge-open-url-youtube");
static const char *lazy_open_file[]                   = ST("lazy-open-search-file");
static const char *lazy_open_search_book[]            = ST("lazy-open-search-book");
static const char *lazy_open_search_file_content[]    = ST("lazy-open-search-file-content");
static const char *lazy_open_search_media[]           = ST("lazy-open-search-media");
static const char *lazy_open_search_wiki[]            = ST("lazy-open-search-wiki");
static const char *note_diary[]                       = SH("note-diary");
static const char *note_flash_card[]                  = SH("note-flash-card");
static const char *note_timeline[]                    = SH("note-timeline");
static const char *jump_to_code_from_log[]            = SH("jump-to-code-from-log");
static const char *handle_copied[]                    = SH("handle-copied");
static const char *search[]                           = SH("search");

static const Key keys[] = {
/*  modifier                      key            function           argument                                */
// SUPKEY + F1-F12
   { SUPKEY,                       XK_F1,         spawn,             {.v = volume_toggle                     } },
   { SUPKEY,                       XK_F2,         spawn,             {.v = volume_dec                        } },
   { SUPKEY,                       XK_F3,         spawn,             {.v = volume_inc                        } },
   { SUPKEY,                       XK_F4,         spawn,             {.v = microphone_toggle                 } },
   { SUPKEY,                       XK_F5,         spawn,             {.v = microphone_dec                    } },
   { SUPKEY,                       XK_F6,         spawn,             {.v = microphone_inc                    } },
   { SUPKEY,                       XK_F7,         spawn,             {.v = wifi_connect                      } },
   { SUPKEY,                       XK_F8,         spawn,             {.v = toggle_screen                     } },
   { SUPKEY,                       XK_F9,         spawn,             {.v = toggle_bluetooth                  } },
   { SUPKEY,                       XK_F10,        spawn,             {.v = screen_light_dec                  } },
   { SUPKEY,                       XK_F11,        spawn,             {.v = screen_light_inc                  } },
   { SUPKEY,                       XK_F12,        spawn,             {.v = toggle_keyboard_light             } },

// SUPKEY|ShiftMask + F1-F12
// { SUPKEY|ShiftMask,             XK_F1,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F2,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F3,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F4,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F5,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F6,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F7,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F8,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F9,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F10,        spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F11,        spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_F12,        spawn,             {.v =                                   } },

// SUPKEY + 1-9-0
   { SUPKEY,                       XK_1,          spawn,             {.v = chrome_open_url_chatgpt           } },
   { SUPKEY,                       XK_2,          spawn,             {.v = chrome_open_url_codeium           } },
   { SUPKEY,                       XK_3,          spawn,             {.v = chrome_open_url_google_mail       } },
   { SUPKEY,                       XK_4,          spawn,             {.v = chrome_open_url_google_translate  } },
   { SUPKEY,                       XK_5,          spawn,             {.v = chrome_open_url_youtube           } },
   { SUPKEY,                       XK_6,          spawn,             {.v = chrome_open_url_wechat            } },
   { SUPKEY,                       XK_7,          spawn,             {.v = chrome_open_url_instagram         } },
   { SUPKEY,                       XK_8,          spawn,             {.v = chrome_open_url_github            } },
   { SUPKEY,                       XK_9,          spawn,             {.v = chrome_open_url_leetcode          } },
   { SUPKEY,                       XK_0,          spawn,             {.v = chrome_open_url_doubao            } },

// SUPKEY|ShiftMask + 1-9-0
   { SUPKEY|ShiftMask,             XK_1,            spawn,             {.v = edge_open_url_chatgpt             } },
   { SUPKEY|ShiftMask,             XK_2,            spawn,             {.v = edge_open_url_codeium             } },
   { SUPKEY|ShiftMask,             XK_3,            spawn,             {.v = edge_open_url_google_mail         } },
   { SUPKEY|ShiftMask,             XK_4,            spawn,             {.v = edge_open_url_google_translate    } },
   { SUPKEY|ShiftMask,             XK_5,            spawn,             {.v = edge_open_url_youtube             } },
   { SUPKEY|ShiftMask,             XK_6,            spawn,             {.v = edge_open_url_wechat              } },
   { SUPKEY|ShiftMask,             XK_7,            spawn,             {.v = edge_open_url_instagram           } },
   { SUPKEY|ShiftMask,             XK_8,            spawn,             {.v = edge_open_url_github              } },
   { SUPKEY|ShiftMask,             XK_9,            spawn,             {.v = edge_open_url_leetcode            } },
   { SUPKEY|ShiftMask,             XK_0,            spawn,             {.v = edge_open_url_doubao              } },

// SUPKEY|ShiftMask,+ a-z, etc
   { SUPKEY,                       XK_a,            spawn,             {.v = toggle_python                     } },
   { SUPKEY,                       XK_b,            spawn,             {.v = launch_chrome                     } },
   { SUPKEY,                       XK_c,            spawn,             {.v = toggle_mutt                       } },
   { SUPKEY,                       XK_d,            spawn,             {.v = note_diary                        } },
   { SUPKEY,                       XK_e,            spawn,             {.v = launch_edge                       } },
   { SUPKEY,                       XK_f,            spawn,             {.v = lazy_open_file                    } },
   { SUPKEY,                       XK_g,            spawn,             {.v = lazy_open_search_file_content     } },
// { SUPKEY,                       XK_h,            spawn,             {.v = x                                 } },
// { SUPKEY,                       XK_i,            spawn,             {.v =                                   } },
// { SUPKEY,                       XK_j,            spawn,             {.v = x                                 } },
// { SUPKEY,                       XK_k,            spawn,             {.v = x                                 } },
// { SUPKEY,                       XK_l,            spawn,             {.v = x                                 } },
   { SUPKEY,                       XK_m,            spawn,             {.v = toggle_julia                      } },
   { SUPKEY,                       XK_n,            spawn,             {.v = toggle_obsidian                   } },
   { SUPKEY,                       XK_o,            spawn,             {.v = handle_copied                     } },
   { SUPKEY,                       XK_p,            spawn,             {.v = lazy_open_search_book             } },
   { SUPKEY,                       XK_q,            spawn,             {.v = screen_lock                       } },
   { SUPKEY,                       XK_r,            spawn,             {.v = toggle_joshuto                    } },
   { SUPKEY,                       XK_s,            spawn,             {.v = search                            } },
// { SUPKEY,                       XK_t,            spawn,             {.v =                                   } },
   { SUPKEY,                       XK_u,            spawn,             {.v = toggle_screenkey                  } },
   { SUPKEY,                       XK_v,            spawn,             {.v = lazy_open_search_media            } },
   { SUPKEY,                       XK_w,            spawn,             {.v = lazy_open_search_wiki             } },
   { SUPKEY,                       XK_x,            spawn,             {.v = toggle_wallpaper                  } },
   { SUPKEY,                       XK_y,            spawn,             {.v = toggle_show                       } },
   { SUPKEY,                       XK_z,            spawn,             {.v = chrome_open_url_google            } },
   { SUPKEY,                       XK_apostrophe,   spawn,             {.v = toggle_termius                    } },
   { SUPKEY,                       XK_BackSpace,    spawn,             {.v = toggle_passmenu                   } },
   { SUPKEY,                       XK_Delete,       spawn,             {.v = toggle_sys_shortcuts              } },
   { SUPKEY,                       XK_Escape,       spawn,             {.v = toggle_top                        } },
// { SUPKEY,                       XK_Print,        spawn,             {.v =                                   } },
   { SUPKEY,                       XK_Home,         spawn,             {.v = toggle_flameshot                  } },
   { SUPKEY,                       XK_backslash,    spawn,             {.v = keyboard_rate                     } },
   { SUPKEY,                       XK_slash,        spawn,             {.v = note_flash_card                   } },
   { SUPKEY,                       XK_comma,        spawn,             {.v = jump_to_code_from_log             } },
   { SUPKEY,                       XK_period,       spawn,             {.v = note_timeline                     } },

// SUPKEY-ShiftMask + a-z, etc
   { SUPKEY|ShiftMask,             XK_a,            spawn,             {.v = toggle_addressbook                } },
// { SUPKEY|ShiftMask,             XK_b,            spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_c,            spawn,             {.v = toggle_calendar_scheduling        } },
   { SUPKEY|ShiftMask,             XK_d,            spawn,             {.v = toggle_lazydocker                 } },
// { SUPKEY|ShiftMask,             XK_e,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_f,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_g,            spawn,             {.v = x                                 } },
// { SUPKEY|ShiftMask,             XK_h,            spawn,             {.v = x                                 } },
   { SUPKEY|ShiftMask,             XK_i,            spawn,             {.v = toggle_inkscape                   } },
// { SUPKEY|ShiftMask,             XK_j,            spawn,             {.v = x                                 } },
// { SUPKEY|ShiftMask,             XK_k,            spawn,             {.v = x                                 } },
// { SUPKEY|ShiftMask,             XK_l,            spawn,             {.v = x                                 } },
   { SUPKEY|ShiftMask,             XK_m,            spawn,             {.v = toggle_music_net_cloud            } },
// { SUPKEY|ShiftMask,             XK_n,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_o,            spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_p,            spawn,             {.v = toggle_krita                      } },
   { SUPKEY|ShiftMask,             XK_q,            spawn,             {.v = sys_suspend                       } },
   { SUPKEY|ShiftMask,             XK_r,            spawn,             {.v = toggle_redshift                   } },
   { SUPKEY|ShiftMask,             XK_s,            spawn,             {.v = toggle_sublime                    } },
// { SUPKEY|ShiftMask,             XK_t,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_u,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_v,            spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_w,            spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_x,            spawn,             {.v = toggle_calendar_scheduling_today  } },
// { SUPKEY|ShiftMask,             XK_y,            spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_z,            spawn,             {.v = chrome_open_url_bing              } },
// { SUPKEY|ShiftMask,             XK_apostrophe,   spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_Delete,       spawn,             {.v = sys_shutdown                      } },
// { SUPKEY|ShiftMask,             XK_Escape,       spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_Print,        spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_Home,         spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_End,          spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_backslash,    spawn,             {.v =                                   } },
// { SUPKEY|ShiftMask,             XK_BackSpace,    spawn,             {.v =                                   } },
   { SUPKEY|ShiftMask,             XK_slash,        spawn,             {.v = toggle_rec_webcam                 } },
   { SUPKEY|ShiftMask,             XK_comma,        spawn,             {.v = toggle_rec_audio                  } },
   { SUPKEY|ShiftMask,             XK_period,       spawn,             {.v = toggle_rec_screen                 } },
   { SUPKEY,                       XK_k,            movewin,           {.ui = UP                               } },
   { SUPKEY,                       XK_j,            movewin,           {.ui = DOWN                             } },
   { SUPKEY,                       XK_h,            movewin,           {.ui = LEFT                             } },
   { SUPKEY,                       XK_l,            movewin,           {.ui = RIGHT                            } },
   { SUPKEY|ShiftMask,             XK_k,            resizewin,         {.ui = VECINC                           } },
   { SUPKEY|ShiftMask,             XK_j,            resizewin,         {.ui = VECDEC                           } },
   { SUPKEY|ShiftMask,             XK_h,            resizewin,         {.ui = HORDEC                           } },
   { SUPKEY|ShiftMask,             XK_l,            resizewin,         {.ui = HORINC                           } },

// MODKEY, etc
// { MODKEY,                       XK_n,            xxxxx,             {.v = x                                 } },
// { MODKEY,                       XK_x,            xxxxx,             {.v = x                                 } },
// { MODKEY,                       XK_y,            xxxxx,             {.v = x                                 } },
// { MODKEY,                       XK_z,            xxxxx,             {.v = x                                 } },
// { MODKEY|ShiftMask,             XK_b,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_d,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_g,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_i,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_n,            xxxxx,             {.v = x                                 } },
// { MODKEY|ShiftMask,             XK_o,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_x,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_y,            xxxxx,             {.v =                                   } },
// { MODKEY|ShiftMask,             XK_z,            xxxxx,             {.v =                                   } },
// { MODKEY,                       XK_u,            xxxxx,             {.i =                                   } },
// { MODKEY|ShiftMask,             XK_u,            xxxxx,             {.i =                                   } },
// { MODKEY,                       XK_bracketleft,  xxxxx,             {.i =                                   } },
// { MODKEY,                       XK_bracketright, xxxxx,             {.i =                                   } },
   { MODKEY,                       XK_apostrophe,   togglescratch,     {.v = scratchpadcmd                     } },
   { MODKEY,                       XK_c,            spawn,             {.v = toggle_clipmenu                   } },
   { MODKEY,                       XK_p,            spawn,             {.v = dmenucmd                          } },
   { MODKEY,                       XK_Return,       zoom,              {0                                      } },
   { MODKEY,                       XK_Tab,          view,              {0                                      } },
   { MODKEY,                       XK_b,            togglebar,         {0                                      } },
   { MODKEY,                       XK_f,            togglefullscreen,  {0                                      } },
   { MODKEY|ShiftMask,             XK_f,            togglefloating,    {0                                      } },
   { MODKEY,                       XK_o,            toggleoverview,    {0                                      } },
   { MODKEY,                       XK_s,            reset,             {0                                      } },
   { MODKEY,                       XK_a,            setlayout,         {0                                      } }, // teemporary layout switch
// { MODKEY|ShiftMask,             XK_a,            xxxxxxxxx,         {0                                      } },
   { MODKEY|ShiftMask,             XK_s,            togglesticky,      {0                                      } },
   { MODKEY|ShiftMask,             XK_space,        focusmaster,       {0                                      } },
   { MODKEY,                       XK_minus,        scratchpad_show,   {0                                      } },
   { MODKEY|ShiftMask,             XK_minus,        scratchpad_hide,   {0                                      } },
   { MODKEY,                       XK_equal,        scratchpad_remove, {0                                      } },
   { MODKEY,                       XK_slash,        focusmon,          {.i = +1                                } }, // monitor related
   { MODKEY|ShiftMask,             XK_slash,        tagmon,            {.i = +1                                } }, // monitor related
   { MODKEY,                       XK_d,            incnmaster,        {.i = -1                                } },
   { MODKEY,                       XK_i,            incnmaster,        {.i = +1                                } },
   { MODKEY,                       XK_h,            movestack,         {.i = -1                                } },
   { MODKEY,                       XK_l,            movestack,         {.i = +1                                } },
   { MODKEY,                       XK_comma,        movestack,         {.i = -1                                } },
   { MODKEY,                       XK_period,       movestack,         {.i = +1                                } },
   { MODKEY|ShiftMask,             XK_comma,        shiftview,         {.i = -1                                } },
   { MODKEY|ShiftMask,             XK_period,       shiftview,         {.i = +1                                } },
   { MODKEY|ControlMask,           XK_comma,        cyclelayout,       {.i = -1                                } }, // useless
   { MODKEY|ControlMask,           XK_period,       cyclelayout,       {.i = +1                                } }, // useless
   { MODKEY,                       XK_k,            focusstack,        {.i = -1                                } },
   { MODKEY,                       XK_j,            focusstack,        {.i = +1                                } },
   { MODKEY|ShiftMask,             XK_h,            setmfact,          {.f = -0.025                            } },
   { MODKEY|ShiftMask,             XK_l,            setmfact,          {.f = +0.025                            } },
   { MODKEY|ShiftMask,             XK_j,            setffact,          {.f = -0.025                            } },
   { MODKEY|ShiftMask,             XK_k,            setffact,          {.f = +0.025                            } },
   { MODKEY,                       XK_space,        setlayout,         {.v = &layouts[0]                       } },
   { MODKEY,                       XK_r,            setlayout,         {.v = &layouts[1]                       } },
   { MODKEY|ShiftMask,             XK_r,            setlayout,         {.v = &layouts[2]                       } },
   { MODKEY,                       XK_v,            setlayout,         {.v = &layouts[3]                       } },
   { MODKEY|ShiftMask,             XK_v,            setlayout,         {.v = &layouts[4]                       } },
   { MODKEY,                       XK_t,            setlayout,         {.v = &layouts[5]                       } },
   { MODKEY|ShiftMask,             XK_t,            setlayout,         {.v = &layouts[6]                       } },
   { MODKEY,                       XK_g,            setlayout,         {.v = &layouts[7]                       } },
   { MODKEY,                       XK_m,            setlayout,         {.v = &layouts[8]                       } },
   { MODKEY|ShiftMask,             XK_m,            setlayout,         {.v = &layouts[9]                       } },
   { MODKEY,                       XK_w,            setlayout,         {.v = &layouts[10]                      } },
   { MODKEY|ShiftMask,             XK_w,            setlayout,         {.v = &layouts[11]                      } },
   { MODKEY|ShiftMask,             XK_e,            setlayout,         {.v = &layouts[13]                      } },
   { MODKEY,                       XK_e,            setlayout,         {.v = &layouts[12]                      } },
   { MODKEY|ShiftMask,             XK_Return,       spawn,             {.v = termcmd                           } },
   { MODKEY|ControlMask|ShiftMask, XK_Return,       spawn,             {.v = termcmd_kitty                     } },
   { MODKEY|ShiftMask,             XK_c,            killclient,        {0                                      } },
   { MODKEY|ShiftMask,             XK_q,            quit,              {0                                      } },
   { MODKEY|ShiftMask,             XK_p,            quit,              {1                                      } },

   { MODKEY,                       XK_1,            view,              {.ui = 1 << 0                           } }, // view tag 1
   { MODKEY,                       XK_2,            view,              {.ui = 1 << 1                           } },
   { MODKEY,                       XK_3,            view,              {.ui = 1 << 2                           } },
   { MODKEY,                       XK_4,            view,              {.ui = 1 << 3                           } },
   { MODKEY,                       XK_5,            view,              {.ui = 1 << 4                           } },
   { MODKEY,                       XK_6,            view,              {.ui = 1 << 5                           } },
   { MODKEY,                       XK_7,            view,              {.ui = 1 << 6                           } },
   { MODKEY,                       XK_8,            view,              {.ui = 1 << 7                           } },
   { MODKEY,                       XK_9,            view,              {.ui = 1 << 8                           } },
   { MODKEY,                       XK_0,            view,              {.ui = ~0                               } }, // preview all tags
   { MODKEY|ShiftMask,             XK_1,            tag,               {.ui = 1 << 0                           } }, // move to tag 1
   { MODKEY|ShiftMask,             XK_2,            tag,               {.ui = 1 << 1                           } },
   { MODKEY|ShiftMask,             XK_3,            tag,               {.ui = 1 << 2                           } },
   { MODKEY|ShiftMask,             XK_4,            tag,               {.ui = 1 << 3                           } },
   { MODKEY|ShiftMask,             XK_5,            tag,               {.ui = 1 << 4                           } },
   { MODKEY|ShiftMask,             XK_6,            tag,               {.ui = 1 << 5                           } },
   { MODKEY|ShiftMask,             XK_7,            tag,               {.ui = 1 << 6                           } },
   { MODKEY|ShiftMask,             XK_8,            tag,               {.ui = 1 << 7                           } },
   { MODKEY|ShiftMask,             XK_9,            tag,               {.ui = 1 << 8                           } },
   { MODKEY|ShiftMask,             XK_0,            tag,               {.ui = ~0                               } }, // stick to all tags
   { MODKEY|ControlMask,           XK_1,            toggleview,        {.ui = 1 << 0                           } }, // toggle view of tag 1
   { MODKEY|ControlMask,           XK_2,            toggleview,        {.ui = 1 << 1                           } },
   { MODKEY|ControlMask,           XK_3,            toggleview,        {.ui = 1 << 2                           } },
   { MODKEY|ControlMask,           XK_4,            toggleview,        {.ui = 1 << 3                           } },
   { MODKEY|ControlMask,           XK_5,            toggleview,        {.ui = 1 << 4                           } },
   { MODKEY|ControlMask,           XK_6,            toggleview,        {.ui = 1 << 5                           } },
   { MODKEY|ControlMask,           XK_7,            toggleview,        {.ui = 1 << 6                           } },
   { MODKEY|ControlMask,           XK_8,            toggleview,        {.ui = 1 << 7                           } },
   { MODKEY|ControlMask,           XK_9,            toggleview,        {.ui = 1 << 8                           } },
   { MODKEY|ControlMask|ShiftMask, XK_1,            previewtag,        {.ui = 0                                } },
   { MODKEY|ControlMask|ShiftMask, XK_2,            previewtag,        {.ui = 1                                } },
   { MODKEY|ControlMask|ShiftMask, XK_3,            previewtag,        {.ui = 2                                } },
   { MODKEY|ControlMask|ShiftMask, XK_4,            previewtag,        {.ui = 3                                } },
   { MODKEY|ControlMask|ShiftMask, XK_5,            previewtag,        {.ui = 4                                } },
   { MODKEY|ControlMask|ShiftMask, XK_6,            previewtag,        {.ui = 5                                } },
   { MODKEY|ControlMask|ShiftMask, XK_7,            previewtag,        {.ui = 6                                } },
   { MODKEY|ControlMask|ShiftMask, XK_8,            previewtag,        {.ui = 7                                } },
   { MODKEY|ControlMask|ShiftMask, XK_9,            previewtag,        {.ui = 8                                } },
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
   { ClkTagBar,            0,              Button1,        view,                {0                               } },
   { ClkTagBar,            0,              Button3,        toggleview,          {0                               } },
   { ClkTagBar,            MODKEY,         Button1,        tag,                 {0                               } },
   { ClkTagBar,            MODKEY,         Button3,        toggletag,           {0                               } },
   { ClkLtSymbol,          0,              Button1,        setlayout,           {0                               } },
   { ClkLtSymbol,          0,              Button2,        setlayout,           {.v = &layouts[8]                } },
   { ClkLtSymbol,          0,              Button3,        setlayout,           {.v = &overviewlayout            } },
   { ClkStatusText,        0,              Button1,        spawn,               {.v = toggle_calendar_scheduling } },
   { ClkStatusText,        0,              Button2,        spawn,               {.v = toggle_sys_shortcuts       } },
   { ClkStatusText,        0,              Button3,        spawn,               {.v = toggle_calendar            } },
// { ClkWinTitle,          0,              Button1,        xxxxxxxxx,           {0                               } },
// { ClkWinTitle,          0,              Button2,        xxxxxxxxx,           {0                               } },
// { ClkWinTitle,          0,              Button3,        xxxxxxxxx,           {0                               } },
   { ClkClientWin,         MODKEY,         Button1,        movemouse,           {0                               } },
   { ClkClientWin,         MODKEY,         Button2,        togglefloating,      {0                               } },
   { ClkClientWin,         MODKEY,         Button3,        resizemouse,         {0                               } },
};
