#include "drw.h"
#include "util.h"
#include "ds.h"

/* function declarations */
static Atom getatomprop(Client *c, Atom prop);
static Client *nexttiled(Client *c);
static Client *swallowingclient(Window w);
static Client *termforwin(const Client *c);
static Client *wintoclient(Window w);
static Monitor *createmon(void);
static Monitor *dirtomon(int dir);
static Monitor *recttomon(int x, int y, int w, int h);
static Monitor *wintomon(Window w);
static int applysizehints(Client *c, int *x, int *y, int *w, int *h, int interact);
static int getrootptr(int *x, int *y);
static int gettextprop(Window w, Atom atom, char *text, unsigned int size);
static int isdescprocess(pid_t p, pid_t c);
static int sendevent(Client *c, Atom proto);
static int updategeom(void);
static int xerror(Display *dpy, XErrorEvent *ee);
static int xerrordummy(Display *dpy, XErrorEvent *ee);
static int xerrorstart(Display *dpy, XErrorEvent *ee);
static long getstate(Window w);
static pid_t getparentprocess(pid_t p);
static pid_t winpid(Window w);
static void applyrules(Client *c);
static void arrange(Monitor *m);
static void arrangemon(Monitor *m);
static void attach(Client *c);
static void attachstack(Client *c);
static void autostart_exec(void);
static void buttonpress(XEvent *e);
static void checkotherwm(void);
static void cleanup(void);
static void cleanupmon(Monitor *mon);
static void clientmessage(XEvent *e);
static void configure(Client *c);
static void configurenotify(XEvent *e);
static void configurerequest(XEvent *e);
static void cyclelayout(const Arg *arg);
static void destroynotify(XEvent *e);
static void detach(Client *c);
static void detachstack(Client *c);
static void drawbar(Monitor *m);
static void drawbars(void);
static void enternotify(XEvent *e);
static void expose(XEvent *e);
static void focus(Client *c);
static void focusin(XEvent *e);
static void focusmaster(const Arg *arg);
static void focusmon(const Arg *arg);
static void focusstack(const Arg *arg);
static void grabbuttons(Client *c, int focused);
static void grabkeys(void);
static void incnmaster(const Arg *arg);
static void jump_to_sel(const Arg *arg);
static void keypress(XEvent *e);
static void killclient(const Arg *arg);
static void killclient_unsel(const Arg *arg);
static void killclient_unsel(const Arg *arg);
static void freeclasshints(XClassHint *ch);
static void mappingnotify(XEvent *e);
static void maprequest(XEvent *e);
static void motionnotify(XEvent *e);
static void movemouse(const Arg *arg);
static void movestack(const Arg *arg);
static void movewin(const Arg *arg);
static void next_theme(const Arg *arg);
static void pointerfocuswin(Client *c);
static void pop(Client *c);
static void previewtag(const Arg *arg);
static void propertynotify(XEvent *e);
static void quit(const Arg *arg);
static void reset();
static void resize(Client *c, int x, int y, int w, int h, int interact);
static void resizeclient(Client *c, int x, int y, int w, int h);
static void resizemouse(const Arg *arg);
static void resizewin(const Arg *arg);
static void restack(Monitor *m);
static void restoresession();
static void run(void);
static void savesession();
static void scan(void);
static void sendmon(Client *c, Monitor *m);
static void setclientstate(Client *c, long state);
static void sethfact(const Arg *arg);
static void setfocus(Client *c);
static void setfullscreen(Client *c, int fullscreen);
static void setlayout(const Arg *arg);
static void setmfact(const Arg *arg);
static void setup(void);
static void seturgent(Client *c, int urg);
static void shiftview(const Arg *arg);
static void showhide(Client *c);
static void showtagpreview(unsigned int i);
static void sighup(int unused);
static void sigterm(int unused);
static void spawn(const Arg *arg);
static void spawn_or_focus(const Arg *arg);
static void tag(const Arg *arg);
static void tagmon(const Arg *arg);
static void takepreview(void);
static void togglebar(const Arg *arg);
static void togglefloating(const Arg *arg);
static void togglefullscreen(const Arg *arg);
static void toggleoverview(const Arg *arg);
static void toggle_scratchpad(const Arg *arg);
static void scratchpad_to_normal(const Arg *arg);
static void togglesticky(const Arg *arg);
static void toggletag(const Arg *arg);
static void toggleview(const Arg *arg);
static void unfocus(Client *c, int setfocus);
static void unmanage(Client *c, int destroyed);
static void unmapnotify(XEvent *e);
static void updatebarpos(Monitor *m);
static void updatebars(void);
static void updateclientlist(void);
static void updatenumlockmask(void);
static void updatesizehints(Client *c);
static void updatestatus(void);
static void updatetitle(Client *c);
static void updatewindowtype(Client *c);
static void updatewmhints(Client *c);
static void view(const Arg *arg);
static void zoom(const Arg *arg);

/* layout */
static void layout_monocle(Monitor *m);
static void layout_center_free_shape(Monitor *m);
static void layout_center_equal_ratio(Monitor *m);
static void layout_fibonacci(Monitor *m, int s);
static void layout_fib_dwindle(Monitor * m);
static void layout_fib_spiral(Monitor * m);
static void layout_grid(Monitor *m);
static void layout_tile_right(Monitor *m);
static void layout_tile_left(Monitor *m);
static void layout_stack_hori(Monitor *m);
static void layout_stack_vert(Monitor *m);
static void layout_hacker(Monitor *m);
static void layout_grid_gap(Monitor *m);
static void layout_overview(Monitor *m);
static void layout_workflow(Monitor *m);

/* variables */
static const char broken[] = "broken";
static char stext[256];
static int screen;
static int sw, sh;
static int bh;
static int lrpad;
static int vp;
static int sp;
static int (*xerrorxlib)(Display *, XErrorEvent *);
static unsigned int numlockmask = 0;
static void (*handler[LASTEvent]) (XEvent *) = {
  [ButtonPress] = buttonpress,
  [ClientMessage] = clientmessage,
  [ConfigureRequest] = configurerequest,
  [ConfigureNotify] = configurenotify,
  [DestroyNotify] = destroynotify,
  [EnterNotify] = enternotify,
  [Expose] = expose,
  [FocusIn] = focusin,
  [KeyPress] = keypress,
  [MappingNotify] = mappingnotify,
  [MapRequest] = maprequest,
  [MotionNotify] = motionnotify,
  [PropertyNotify] = propertynotify,
  [UnmapNotify] = unmapnotify
};
static Atom wmatom[WMLast], netatom[NetLast];
static int restart = 0;
static int running = 1;
static Cur *cursor[CurLast];
static Clr **scheme;
static Display *dpy;
static Drw *drw;
static Monitor *mons, *selmon;
static Window root, wmcheckwin;
static xcb_connection_t *xcon;
static int winpad = 0;

/* configuration, allows nested code to access above variables */
#include "config.h"

struct Pertag {
  const Layout *ltidxs[LENGTH(tags) + 1][2];
  float hfacts[LENGTH(tags) + 1];
  float mfacts[LENGTH(tags) + 1];
  int nmasters[LENGTH(tags) + 1];
  int showbars[LENGTH(tags) + 1];
  unsigned int curtag, prevtag;
  unsigned int sellts[LENGTH(tags) + 1];
};
struct NumTags { char limitexceeded[LENGTH(tags) > 31 ? -1 : 1]; };
static pid_t *autostart_pids;
static size_t autostart_len;

/* scratchpad 当前等待处理的 scratchpad class */
static const char *scratchpad_class_wait = NULL;

/* execute command from autostart array */
static void
autostart_exec() {
  const char *const *p;
  size_t i = 0;

  for (p = autostart; *p; autostart_len++, p++) {
    while (*++p)
      ;
  }

  autostart_pids = malloc(autostart_len * sizeof(pid_t));
  for (p = autostart; *p; i++, p++) {
    if ((autostart_pids[i] = fork()) == 0) {
      setsid();
      execvp(*p, (char *const *)p);
      fprintf(stderr, "dwm: execvp %s\n", *p);
      perror(" failed");
      _exit(EXIT_FAILURE);
    }
    while (*++p)
      ;
  }
}

/* function implementations */
  void
applyrules(Client *c)
{
  const char *cls, *instance;
  unsigned int i;
  const Rule *r;
  Monitor *m;
  XClassHint ch = { NULL, NULL };

  c->isfloating = 0;
  c->isontop = 0;
  c->tags = 0;
  XGetClassHint(dpy, c->win, &ch);
  cls      = ch.res_class ? ch.res_class : broken;
  instance = ch.res_name  ? ch.res_name  : broken;

  for (i = 0; i < LENGTH(rules); i++) {
    r = &rules[i];
    if ((!r->title || strstr(c->name, r->title)) && (!r->cls || strstr(cls, r->cls)) && (!r->instance || strstr(instance, r->instance))) {
      c->isterminal = r->isterminal;
      c->noswallow  = r->noswallow;
      c->isfloating = r->isfloating;
      c->isontop = r->isontop;
      c->tags |= r->tags;
      for (m = mons; m && m->num != r->monitor; m = m->next)
        ;
      if (m) {
        c->mon = m;
      }
    }
  }
  if (ch.res_class) {
    XFree(ch.res_class);
  }
  if (ch.res_name) {
    XFree(ch.res_name);
  }

  c->tags = c->tags & TAGMASK ? c->tags & TAGMASK : c->mon->tagset[c->mon->seltags];
}

  int
applysizehints(Client *c, int *x, int *y, int *w, int *h, int interact)
{
  int baseismin;
  Monitor *m = c->mon;

  /* set minimum possible */
  *w = MAX(1, *w);
  *h = MAX(1, *h);
  if (interact) {
    if (*x > sw) {
      *x = sw - WIDTH(c);
    }
    if (*y > sh) {
      *y = sh - HEIGHT(c);
    }
    if (*x + *w + 2 * c->bw < 0) {
      *x = 0;
    }
    if (*y + *h + 2 * c->bw < 0) {
      *y = 0;
    }
  } else {
    if (*x >= m->wx + m->ww) {
      *x = m->wx + m->ww - WIDTH(c);
    }
    if (*y >= m->wy + m->wh) {
      *y = m->wy + m->wh - HEIGHT(c);
    }
    if (*x + *w + 2 * c->bw <= m->wx) {
      *x = m->wx;
    }
    if (*y + *h + 2 * c->bw <= m->wy) {
      *y = m->wy;
    }
  }
  if (*h < bh) {
    *h = bh;
  }
  if (*w < bh) {
    *w = bh;
  }
  if (resizehints || c->isfloating || !c->mon->lt[c->mon->sellt]->arrange) {
    if (!c->hintsvalid) {
      updatesizehints(c);
    }
    /* see last two sentences in ICCCM 4.1.2.3 */
    baseismin = c->basew == c->minw && c->baseh == c->minh;
    if (!baseismin) { /* temporarily remove base dimensions */
      *w -= c->basew;
      *h -= c->baseh;
    }
    /* adjust for aspect limits */
    if (c->mina > 0 && c->maxa > 0) {
      if (c->maxa < (float)*w / *h) {
        *w = *h * c->maxa + 0.5;
      } else if (c->mina < (float)*h / *w) {
        *h = *w * c->mina + 0.5;
      }
    }
    if (baseismin) { /* increment calculation requires this */
      *w -= c->basew;
      *h -= c->baseh;
    }
    /* adjust for increment value */
    if (c->incw) {
      *w -= *w % c->incw;
    }
    if (c->inch) {
      *h -= *h % c->inch;
    }
    /* restore base dimensions */
    *w = MAX(*w + c->basew, c->minw);
    *h = MAX(*h + c->baseh, c->minh);
    if (c->maxw) {
      *w = MIN(*w, c->maxw);
    }
    if (c->maxh) {
      *h = MIN(*h, c->maxh);
    }
  }
  return *x != c->x || *y != c->y || *w != c->w || *h != c->h;
}

  void
arrange(Monitor *m)
{
  if (m) {
    showhide(m->stack);
  } else {
    for (m = mons; m; m = m->next) {
      showhide(m->stack);
    }
  }

  if (m) {
    arrangemon(m);
    restack(m);
  } else {
    for (m = mons; m; m = m->next) {
      arrangemon(m);
    }
  }
}

  void
arrangemon(Monitor *m)
{
  if (m->isoverview) {
    strncpy(m->ltsymbol, overviewlayout.symbol, sizeof m->ltsymbol);
    overviewlayout.arrange(m);
    return;
  }
  strncpy(m->ltsymbol, m->lt[m->sellt]->symbol, sizeof m->ltsymbol);
  if (m->lt[m->sellt]->arrange) {
    m->lt[m->sellt]->arrange(m);
  }
}

  void
attach(Client *c)
{
  c->next = c->mon->clients;
  c->mon->clients = c;
}

  void
attachstack(Client *c)
{
  c->snext = c->mon->stack;
  c->mon->stack = c;
}

  void
swallow(Client *p, Client *c)
{
  if (c->noswallow || c->isterminal) {
    return;
  }
  if (c->noswallow && !swallowfloating && c->isfloating) {
    return;
  }

  detach(c);
  detachstack(c);

  setclientstate(c, WithdrawnState);
  XUnmapWindow(dpy, p->win);

  p->swallowing = c;
  c->mon = p->mon;

  Window w = p->win;
  p->win = c->win;
  c->win = w;
  updatetitle(p);
  XMoveResizeWindow(dpy, p->win, p->x, p->y, p->w, p->h);
  arrange(p->mon);
  configure(p);
  updateclientlist();
}

  void
unswallow(Client *c)
{
  c->win = c->swallowing->win;

  free(c->swallowing);
  c->swallowing = NULL;

  setfullscreen(c, 0);
  updatetitle(c);
  arrange(c->mon);
  XMapWindow(dpy, c->win);
  XMoveResizeWindow(dpy, c->win, c->x, c->y, c->w, c->h);
  setclientstate(c, NormalState);
  focus(NULL);
  arrange(c->mon);
}

  void
buttonpress(XEvent *e)
{
  unsigned int i, x, click;
  Arg arg = {0};
  Client *c;
  Monitor *m;
  XButtonPressedEvent *ev = &e->xbutton;

  click = ClkRootWin;
  /* focus monitor if necessary */
  if ((m = wintomon(ev->window)) && m != selmon) {
    unfocus(selmon->sel, 1);
    selmon = m;
    focus(NULL);
  }
  if (ev->window == selmon->barwin) {
    i = x = 0;
    do {
      x += TEXTW(tags[i]);
    } while (ev->x >= x && ++i < LENGTH(tags));
    if (i < LENGTH(tags)) {
      click = ClkTagBar;
      arg.ui = 1 << i;
      if (selmon->previewshow) {
        selmon->previewshow = 0;
        XUnmapWindow(dpy, selmon->tagwin);
      }
    } else if (ev->x < x + TEXTW(selmon->ltsymbol)) {
      click = ClkLtSymbol;
    } else if (ev->x > selmon->ww - (int)TEXTW(stext)) {
      click = ClkStatusText;
    } else {
      click = ClkWinTitle;
    }
  } else if ((c = wintoclient(ev->window))) {
    focus(c);
    restack(selmon);
    XAllowEvents(dpy, ReplayPointer, CurrentTime);
    click = ClkClientWin;
  }
  for (i = 0; i < LENGTH(buttons); i++) {
    if (click == buttons[i].click && buttons[i].func && buttons[i].button == ev->button && CLEANMASK(buttons[i].mask) == CLEANMASK(ev->state)) {
      buttons[i].func(click == ClkTagBar && buttons[i].arg.i == 0 ? &arg : &buttons[i].arg);
    }
  }
}

  void
checkotherwm(void)
{
  xerrorxlib = XSetErrorHandler(xerrorstart);
  /* this causes an error if some other window manager is running */
  XSelectInput(dpy, DefaultRootWindow(dpy), SubstructureRedirectMask);
  XSync(dpy, False);
  XSetErrorHandler(xerror);
  XSync(dpy, False);
}

  void
cleanup(void)
{
  Arg a = {.ui = ~0};
  Layout foo = { "", NULL };
  Monitor *m;
  size_t i;

  view(&a);
  selmon->lt[selmon->sellt] = &foo;

  for (m = mons; m; m = m->next) {
    while (m->stack) {
      unmanage(m->stack, 0);
    }
  }

  XUngrabKey(dpy, AnyKey, AnyModifier, root);

  while (mons) {
    cleanupmon(mons);
  }

  for (i = 0; i < CurLast; i++) {
    drw_cur_free(drw, cursor[i]);
  }

  for (i = 0; i < LENGTH(colors); i++) {
    drw_scm_free(drw, scheme[i], 3);
  }

  free(scheme);
  XDestroyWindow(dpy, wmcheckwin);
  drw_free(drw);
  XSync(dpy, False);
  XSetInputFocus(dpy, PointerRoot, RevertToPointerRoot, CurrentTime);
  XDeleteProperty(dpy, root, netatom[NetActiveWindow]);
}

  void
cleanupmon(Monitor *mon)
{
  Monitor *m;
  size_t i;

  if (mon == mons) {
    mons = mons->next;
  } else {
    for (m = mons; m && m->next != mon; m = m->next)
      ;
    m->next = mon->next;
  }

  for (i = 0; i < LENGTH(tags); i++) {
    if (mon->tagmap[i]) {
      XFreePixmap(dpy, mon->tagmap[i]);
    }
  }

  free(mon->tagmap);
  XUnmapWindow(dpy, mon->barwin);
  XDestroyWindow(dpy, mon->barwin);
  XUnmapWindow(dpy, mon->tagwin);
  XDestroyWindow(dpy, mon->tagwin);
  free(mon);
}

  void
clientmessage(XEvent *e)
{
  XClientMessageEvent *cme = &e->xclient;
  Client *c = wintoclient(cme->window);

  if (!c) {
    return;
  }
  if (cme->message_type == netatom[NetWMState]) {
    if (cme->data.l[1] == netatom[NetWMFullscreen] || cme->data.l[2] == netatom[NetWMFullscreen]) {
      setfullscreen(c, (cme->data.l[0] == 1 /*_NET_WM_STATE_ADD */ || (cme->data.l[0] == 2 /*_NET_WM_STATE_TOGGLE */ && !c->isfullscreen)));
    }
  } else if (cme->message_type == netatom[NetActiveWindow]) {
    if (c != selmon->sel && !c->isurgent) {
      seturgent(c, 1);
    }
  }
}

  void
configure(Client *c)
{
  XConfigureEvent ce;

  ce.type = ConfigureNotify;
  ce.display = dpy;
  ce.event = c->win;
  ce.window = c->win;
  ce.x = c->x;
  ce.y = c->y;
  ce.width = c->w;
  ce.height = c->h;
  ce.border_width = c->bw;
  ce.above = None;
  ce.override_redirect = False;
  XSendEvent(dpy, c->win, False, StructureNotifyMask, (XEvent *)&ce);
}

  void
configurenotify(XEvent *e)
{
  Monitor *m;
  Client *c;
  XConfigureEvent *ev = &e->xconfigure;
  int dirty;

  /* TODO: updategeom handling sucks, needs to be simplified */
  if (ev->window == root) {
    dirty = (sw != ev->width || sh != ev->height);
    sw = ev->width;
    sh = ev->height;
    if (updategeom() || dirty) {
      drw_resize(drw, sw, bh);
      updatebars();
      for (m = mons; m; m = m->next) {
        for (c = m->clients; c; c = c->next) {
          if (c->isfullscreen) {
            resizeclient(c, m->mx, m->my, m->mw, m->mh);
          }
        }
        XMoveResizeWindow(dpy, m->barwin, m->wx + sp, m->by + vp, m->ww -  2*sp, bh);
      }
      focus(NULL);
      arrange(NULL);
    }
  }
}

  void
configurerequest(XEvent *e)
{
  Client *c;
  Monitor *m;
  XConfigureRequestEvent *ev = &e->xconfigurerequest;
  XWindowChanges wc;

  if ((c = wintoclient(ev->window))) {
    if (ev->value_mask & CWBorderWidth) {
      c->bw = ev->border_width;
    } else if (c->isfloating || !selmon->lt[selmon->sellt]->arrange) {
      m = c->mon;
      if (ev->value_mask & CWX) {
        c->oldx = c->x;
        c->x = m->mx + ev->x;
      }

      if (ev->value_mask & CWY) {
        c->oldy = c->y;
        c->y = m->my + ev->y;
      }

      if (ev->value_mask & CWWidth) {
        c->oldw = c->w;
        c->w = ev->width;
      }

      if (ev->value_mask & CWHeight) {
        c->oldh = c->h;
        c->h = ev->height;
      }

      if ((c->x + c->w) > m->mx + m->mw && c->isfloating) {
        c->x = m->mx + (m->mw / 2 - WIDTH(c) / 2); /* center in x direction */
      }

      if ((c->y + c->h) > m->my + m->mh && c->isfloating) {
        c->y = m->my + (m->mh / 2 - HEIGHT(c) / 2); /* center in y direction */
      }

      if ((ev->value_mask & (CWX|CWY)) && !(ev->value_mask & (CWWidth|CWHeight))) {
        configure(c);
      }

      if (ISVISIBLE(c)) {
        XMoveResizeWindow(dpy, c->win, c->x, c->y, c->w, c->h);
      }
    } else {
      configure(c);
    }
  } else {
    wc.x = ev->x;
    wc.y = ev->y;
    wc.width = ev->width;
    wc.height = ev->height;
    wc.border_width = ev->border_width;
    wc.sibling = ev->above;
    wc.stack_mode = ev->detail;
    XConfigureWindow(dpy, ev->window, ev->value_mask, &wc);
  }
  XSync(dpy, False);
}

  Monitor *
createmon(void)
{
  Monitor *m;
  unsigned int i;

  m = ecalloc(1, sizeof(Monitor));
  m->tagset[0] = m->tagset[1] = 1;
  m->mfact = mfact;
  m->hfact = hfact;
  m->nmaster = nmaster;
  m->showbar = showbar;
  m->topbar = topbar;
  m->lt[0] = &layouts[0];
  m->lt[1] = &layouts[1 % LENGTH(layouts)];
  m->tagmap = ecalloc(LENGTH(tags), sizeof(Pixmap));
  m->isoverview = 0;
  strncpy(m->ltsymbol, layouts[0].symbol, sizeof m->ltsymbol);
  m->pertag = ecalloc(1, sizeof(Pertag));
  m->pertag->curtag = m->pertag->prevtag = 1;

  for (i = 0; i <= LENGTH(tags); i++) {
    m->pertag->nmasters[i] = m->nmaster;
    m->pertag->mfacts[i] = m->mfact;
    m->pertag->hfacts[i] = m->hfact;
    m->pertag->ltidxs[i][0] = m->lt[0];
    m->pertag->ltidxs[i][1] = m->lt[1];
    m->pertag->sellts[i] = m->sellt;
    m->pertag->showbars[i] = m->showbar;
  }

  return m;
}

  void
destroynotify(XEvent *e)
{
  Client *c;
  XDestroyWindowEvent *ev = &e->xdestroywindow;

  if ((c = wintoclient(ev->window))) {
    unmanage(c, 1);
  } else if ((c = swallowingclient(ev->window))) {
    unmanage(c->swallowing, 1);
  }
}

  void
detach(Client *c)
{
  Client **tc;

  for (int i = 1; i < LENGTH(tags); i++) {
    if (c == c->mon->tagmarked[i]) {
      c->mon->tagmarked[i] = NULL;
    }
  }

  for (tc = &c->mon->clients; *tc && *tc != c; tc = &(*tc)->next)
    ;
  *tc = c->next;
}

  void
detachstack(Client *c)
{
  Client **tc, *t;

  for (tc = &c->mon->stack; *tc && *tc != c; tc = &(*tc)->snext)
    ;
  *tc = c->snext;

  if (c == c->mon->sel) {
    for (t = c->mon->stack; t && !ISVISIBLE(t); t = t->snext)
      ;
    c->mon->sel = t;
  }
}

  Monitor *
dirtomon(int dir)
{
  Monitor *m = NULL;

  if (dir > 0) {
    if (!(m = selmon->next)) {
      m = mons;
    }
  } else if (selmon == mons) {
    for (m = mons; m->next; m = m->next)
      ;
  } else {
    for (m = mons; m->next != selmon; m = m->next)
      ;
  }
  return m;
}

  void
drawbar(Monitor *m)
{
  int x, w, tw = 0;
  int boxs = drw->fonts->h / 9;
  int boxw = drw->fonts->h / 6 + 2;
  unsigned int i, occ = 0, urg = 0;
  Client *c;

  if (!m->showbar) {
    return;
  }

  /* draw status first so it can be overdrawn by tags later */
  if (m == selmon) { /* status is only drawn on selected monitor */
    drw_setscheme(drw, scheme[SchemeNorm]);
    tw = TEXTW(stext) - lrpad + 2; /* 2px right padding */
    drw_text(drw, m->ww - tw - 2 * sp, 0, tw, bh, 0, stext, 0);
  }

  for (c = m->clients; c; c = c->next) {
    occ |= c->tags;
    if (c->isurgent) {
      urg |= c->tags;
    }
  }

  x = 0;
  if (m->isoverview) {
    //
  } else {
    for (i = 0; i < LENGTH(tags); i++) {
      w = TEXTW(tags[i]);
      drw_setscheme(drw, scheme[m->tagset[m->seltags] & 1 << i ? SchemeSel : SchemeNorm]);
      drw_text(drw, x, 0, w, bh, lrpad / 2, tags[i], urg & 1 << i);
      if (occ & 1 << i) {
        drw_rect(drw, x + boxs, boxs, boxw, boxw, m == selmon && selmon->sel && selmon->sel->tags & 1 << i, urg & 1 << i);
      }
      x += w;
    }
  }
  w = TEXTW(m->ltsymbol);
  drw_setscheme(drw, scheme[SchemeNorm]);
  x = drw_text(drw, x, 0, w, bh, lrpad / 2, m->ltsymbol, 0);

  if ((w = m->ww - tw - x) > bh) {
    drw_setscheme(drw, scheme[m == selmon ? SchemeSel : SchemeNorm]);
    if (m->sel) {
      drw_text(drw, x, 0, w - 2 * sp, bh, lrpad / 2, m->sel->name, 0);
      if (m->sel->isfloating) {
        drw_rect(drw, x + boxs, boxs, boxw, boxw, m->sel->isfixed, 0);
      }
      if (m->sel->issticky) {
        drw_polygon(drw, x + boxs, m->sel->isfloating ? boxs * 2 + boxw : boxs, stickyiconbb.x, stickyiconbb.y, boxw, boxw * stickyiconbb.y / stickyiconbb.x, stickyicon, LENGTH(stickyicon), Nonconvex, m->sel->tags & m->tagset[m->seltags]);
      }
    } else {
      drw_rect(drw, x, 0, w - 2 * sp, bh, 1, 1);
    }
  }
  drw_map(drw, m->barwin, 0, 0, m->ww, bh);
}

  void
drawbars(void)
{
  Monitor *m;

  for (m = mons; m; m = m->next) {
    drawbar(m);
  }
}

  void
enternotify(XEvent *e)
{
  Client *c;
  Monitor *m;
  XCrossingEvent *ev = &e->xcrossing;

  if ((ev->mode != NotifyNormal || ev->detail == NotifyInferior) && ev->window != root) {
    return;
  }
  c = wintoclient(ev->window);
  m = c ? c->mon : wintomon(ev->window);
  if (m != selmon) {
    unfocus(selmon->sel, 1);
    selmon = m;
  } else if (!c || c == selmon->sel) {
    return;
  }
  focus(c);
}

  void
expose(XEvent *e)
{
  Monitor *m;
  XExposeEvent *ev = &e->xexpose;

  if (ev->count == 0 && (m = wintomon(ev->window))) {
    drawbar(m);
  }
}

  void
focus(Client *c)
{
  if (!c || !ISVISIBLE(c)) {
    for (c = selmon->stack; c && !ISVISIBLE(c); c = c->snext)
      ;
  }
  if (selmon->sel && selmon->sel != c) {
    unfocus(selmon->sel, 0);
  }
  if (c) {
    if (c->mon != selmon) {
      selmon = c->mon;
    }
    if (c->isurgent) {
      seturgent(c, 0);
    }
    detachstack(c);
    attachstack(c);
    grabbuttons(c, 1);
    XSetWindowBorder(dpy, c->win, scheme[SchemeSel][ColBorder].pixel);
    setfocus(c);
  } else {
    XSetInputFocus(dpy, root, RevertToPointerRoot, CurrentTime);
    XDeleteProperty(dpy, root, netatom[NetActiveWindow]);
  }
  selmon->sel = c;
  drawbars();
}

/* there are some broken focus acquiring clients needing extra handling */
  void
focusin(XEvent *e)
{
  XFocusChangeEvent *ev = &e->xfocus;

  if (selmon->sel && ev->window != selmon->sel->win) {
    setfocus(selmon->sel);
  }
}

  void
focusmaster(const Arg *arg)
{
  Client *master;

  if (selmon->nmaster > 1) {
    return;
  }
  if (!selmon->sel || (selmon->sel->isfullscreen && lockfullscreen)) {
    return;
  }

  master = nexttiled(selmon->clients);

  if (!master) {
    return;
  }

  int i;
  for (i = 0; !(selmon->tagset[selmon->seltags] & 1 << i); i++)
    ;
  i++;

  if (selmon->sel == master) {
    if (selmon->tagmarked[i] && ISVISIBLE(selmon->tagmarked[i])) {
      focus(selmon->tagmarked[i]);
    }
  } else {
    selmon->tagmarked[i] = selmon->sel;
    focus(master);
  }
}

  void
focusmon(const Arg *arg)
{
  Monitor *m;

  if (!mons->next) {
    return;
  }
  if ((m = dirtomon(arg->i)) == selmon) {
    return;
  }
  unfocus(selmon->sel, 0);
  selmon = m;
  focus(NULL);
}

  void
focusstack(const Arg *arg)
{
  Client *c = NULL, *i;

  if (!selmon->sel || (selmon->sel->isfullscreen && lockfullscreen))
    return;
  if (arg->i > 0) {
    for (c = selmon->sel->next; c && !ISVISIBLE(c); c = c->next)
      ;

    if (!c) {
      for (c = selmon->clients; c && !ISVISIBLE(c); c = c->next)
        ;
    }
  } else {
    for (i = selmon->clients; i != selmon->sel; i = i->next) {
      if (ISVISIBLE(i))
        c = i;
    }

    if (!c) {
      for (; i; i = i->next) {
        if (ISVISIBLE(i)) {
          c = i;
        }
      }
    }
  }

  if (c) {
    focus(c);
    restack(selmon);
  }
}

  void
pointerfocuswin(Client *c)
{
  if (c) {
    XWarpPointer(dpy, None, root, 0, 0, 0, 0, c->x + c->w / 2, c->y + c->h / 2);
    focus(c);
  } else {
    XWarpPointer(dpy, None, root, 0, 0, 0, 0, selmon->wx + selmon->ww / 3, selmon->wy + selmon->wh / 2);
  }
}

  Atom
getatomprop(Client *c, Atom prop)
{
  int di;
  unsigned long dl;
  unsigned char *p = NULL;
  Atom da, atom = None;

  if (XGetWindowProperty(dpy, c->win, prop, 0L, sizeof atom, False, XA_ATOM, &da, &di, &dl, &dl, &p) == Success && p) {
    if (dl > 0) {
      atom = *(Atom *)p;
    };
    XFree(p);
  }
  return atom;
}

  int
getrootptr(int *x, int *y)
{
  int di;
  unsigned int dui;
  Window dummy;

  return XQueryPointer(dpy, root, &dummy, &dummy, x, y, &di, &di, &dui);
}

  long
getstate(Window w)
{
  int format;
  long result = -1;
  unsigned char *p = NULL;
  unsigned long n, extra;
  Atom real;

  if (XGetWindowProperty(dpy, w, wmatom[WMState], 0L, 2L, False, wmatom[WMState], &real, &format, &n, &extra, (unsigned char **)&p) != Success) {
    return -1;
  }
  if (n != 0) {
    result = *p;
  }
  XFree(p);
  return result;
}

  int
gettextprop(Window w, Atom atom, char *text, unsigned int size)
{
  char **list = NULL;
  int n;
  XTextProperty name;

  if (!text || size == 0) {
    return 0;
  }

  text[0] = '\0';
  if (!XGetTextProperty(dpy, w, &name, atom) || !name.nitems) {
    return 0;
  }

  if (name.encoding == XA_STRING) {
    strncpy(text, (char *)name.value, size - 1);
  } else if (XmbTextPropertyToTextList(dpy, &name, &list, &n) >= Success && n > 0 && *list) {
    strncpy(text, *list, size - 1);
    XFreeStringList(list);
  }

  text[size - 1] = '\0';
  XFree(name.value);
  return 1;
}

  void
grabbuttons(Client *c, int focused)
{
  updatenumlockmask();
  {
    unsigned int i, j;
    unsigned int modifiers[] = { 0, LockMask, numlockmask, numlockmask|LockMask };
    XUngrabButton(dpy, AnyButton, AnyModifier, c->win);
    if (!focused) {
      XGrabButton(dpy, AnyButton, AnyModifier, c->win, False, BUTTONMASK, GrabModeSync, GrabModeSync, None, None);
    }
    for (i = 0; i < LENGTH(buttons); i++) {
      if (buttons[i].click == ClkClientWin) {
        for (j = 0; j < LENGTH(modifiers); j++) {
          XGrabButton(dpy, buttons[i].button, buttons[i].mask | modifiers[j], c->win, False, BUTTONMASK, GrabModeAsync, GrabModeSync, None, None);
        }
      }
    }
  }
}

  void
grabkeys(void)
{
  updatenumlockmask();
  {
    unsigned int i, j, k;
    unsigned int modifiers[] = { 0, LockMask, numlockmask, numlockmask|LockMask };
    int start, end, skip;
    KeySym *syms;

    XUngrabKey(dpy, AnyKey, AnyModifier, root);
    XDisplayKeycodes(dpy, &start, &end);
    syms = XGetKeyboardMapping(dpy, start, end - start + 1, &skip);
    if (!syms) {
      return;
    }
    for (k = start; k <= end; k++) {
      for (i = 0; i < LENGTH(keys); i++) {
        /* skip modifier codes, we do that ourselves */
        if (keys[i].keysym == syms[(k - start) * skip]) {
          for (j = 0; j < LENGTH(modifiers); j++) {
            XGrabKey(dpy, k, keys[i].mod | modifiers[j], root, True, GrabModeAsync, GrabModeAsync);
          }
        }
      }
    }
    XFree(syms);
  }
}

  void
incnmaster(const Arg *arg)
{
  selmon->nmaster = selmon->pertag->nmasters[selmon->pertag->curtag] = MIN(MAX(selmon->nmaster + arg->i, 0), maxnmaster);
  arrange(selmon);
}

#ifdef XINERAMA
  static int
isuniquegeom(XineramaScreenInfo *unique, size_t n, XineramaScreenInfo *info)
{
  while (n--) {
    if (unique[n].x_org == info->x_org && unique[n].y_org == info->y_org && unique[n].width == info->width && unique[n].height == info->height) {
      return 0;
    }
  }
  return 1;
}
#endif

  void
jump_to_sel(const Arg *arg)
{
  Client *c = selmon->sel;
  if (!c) { return; }

  /* 清除 overview 状态 */
  selmon->isoverview = 0;

  /* 跳转到当前窗口所在的 tag */
  view(&(Arg){ .ui = c->tags });
}

  void
keypress(XEvent *e)
{
  unsigned int i;
  KeySym keysym;
  XKeyEvent *ev;

  ev = &e->xkey;
  keysym = XKeycodeToKeysym(dpy, (KeyCode)ev->keycode, 0);
  for (i = 0; i < LENGTH(keys); i++) {
    if (keysym == keys[i].keysym && CLEANMASK(keys[i].mod) == CLEANMASK(ev->state) && keys[i].func) {
      keys[i].func(&(keys[i].arg));
    }
  }
}

  void
killclient(const Arg *arg)
{
  if (!selmon->sel) {
    return;
  }
  if (!sendevent(selmon->sel, wmatom[WMDelete])) {
    XGrabServer(dpy);
    XSetErrorHandler(xerrordummy);
    XSetCloseDownMode(dpy, DestroyAll);
    XKillClient(dpy, selmon->sel->win);
    XSync(dpy, False);
    XSetErrorHandler(xerror);
    XUngrabServer(dpy);
  }
}

  void
killclient_unsel(const Arg *arg)
{
  Client *c;
  if (!selmon->sel) {
    return;
  }
  for (c = selmon->clients; c; c = c->next) {
    if (!ISVISIBLE(c) || (c == selmon->sel)) {
      continue;
    }
    if (!sendevent(c, wmatom[WMDelete])) {
      XGrabServer(dpy);
      XSetErrorHandler(xerrordummy);
      XSetCloseDownMode(dpy, DestroyAll);
      XKillClient(dpy, selmon->sel->win);
      XSync(dpy, False);
      XSetErrorHandler(xerror);
      XUngrabServer(dpy);
    }
  }
}

static void
freeclasshints(XClassHint *ch) {
  if (ch && ch->res_name) XFree(ch->res_name);
  if (ch && ch->res_class) XFree(ch->res_class);
}

  void
manage(Window w, XWindowAttributes *wa)
{
  Client *c, *t = NULL, *term = NULL;
  Window trans = None;
  XWindowChanges wc;

  c = ecalloc(1, sizeof(Client));
  c->win = w;
  c->pid = winpid(w);
  /* geometry */
  c->x = c->oldx = wa->x;
  c->y = c->oldy = wa->y;
  c->w = c->oldw = wa->width;
  c->h = c->oldh = wa->height;
  c->oldbw = wa->border_width;

  updatetitle(c);

  XClassHint ch = { NULL, NULL };
  if (XGetClassHint(dpy, c->win, &ch)) {
    if (ch.res_class) {
      strncpy(c->class, ch.res_class, sizeof(c->class) - 1);
    }
    if (ch.res_name) {
      strncpy(c->instance, ch.res_name, sizeof(c->instance) - 1);
    }
    freeclasshints(&ch);
  }

  if (XGetTransientForHint(dpy, w, &trans) && (t = wintoclient(trans))) {
    c->mon = t->mon;
    c->tags = t->tags;
  } else {
    c->mon = selmon;
    applyrules(c);
    term = termforwin(c);
  }

  if (c->x + WIDTH(c) > c->mon->wx + c->mon->ww) {
    c->x = c->mon->wx + c->mon->ww - WIDTH(c);
  }
  if (c->y + HEIGHT(c) > c->mon->wy + c->mon->wh) {
    c->y = c->mon->wy + c->mon->wh - HEIGHT(c);
  }
  c->x = MAX(c->x, c->mon->wx);
  c->y = MAX(c->y, c->mon->wy);
  c->bw = borderpx;

  /* scratchpad 首次出现的自动居中处理 */
  if (scratchpad_class_wait && strcmp(c->class, scratchpad_class_wait) == 0) {
    // c->tags = 1 << 30;  /* 首次启动不隐藏 */
    c->isfloating = 1;
    int nw = c->mon->ww * scratchpad_width;
    int nh = c->mon->wh * scratchpad_height;
    int nx = c->mon->wx + (c->mon->ww - nw) / 2;
    int ny = c->mon->wy + (c->mon->wh - nh) / 2;
    resize(c, nx, ny, nw, nh, 0);
    /* 处理完清空 */
    scratchpad_class_wait = NULL;
  }

  /* scratchpad */
  selmon->tagset[selmon->seltags] &= ~(1<<30);
  if (!strcmp(c->name, scratchpad_class)) {
    c->mon->tagset[c->mon->seltags] |= c->tags = 1<<30;
    c->isfloating = 1;
    c->x = c->mon->wx + (c->mon->ww / 2 - WIDTH(c) / 2);
    c->y = c->mon->wy + (c->mon->wh / 2 - HEIGHT(c) / 2);
  }

  wc.border_width = c->bw;
  XConfigureWindow(dpy, w, CWBorderWidth, &wc);
  XSetWindowBorder(dpy, w, scheme[SchemeNorm][ColBorder].pixel);
  configure(c); /* propagates border_width, if size doesn't change */
  updatewindowtype(c);
  updatesizehints(c);
  updatewmhints(c);
  XSelectInput(dpy, w, EnterWindowMask|FocusChangeMask|PropertyChangeMask|StructureNotifyMask);
  grabbuttons(c, 0);
  if (!c->isfloating) {
    c->isfloating = c->oldstate = trans != None || c->isfixed;
  }
  if (c->isfloating) {
    XRaiseWindow(dpy, c->win);
  }
  attach(c);
  attachstack(c);
  XChangeProperty(dpy, root, netatom[NetClientList], XA_WINDOW, 32, PropModeAppend, (unsigned char *) &(c->win), 1);
  XMoveResizeWindow(dpy, c->win, c->x + 2 * sw, c->y, c->w, c->h); /* some windows require this */
  setclientstate(c, NormalState);
  if (c->mon == selmon) {
    unfocus(selmon->sel, 0);
  }
  c->mon->sel = c;
  arrange(c->mon);
  XMapWindow(dpy, c->win);
  if (term) {
    swallow(term, c);
  }
  focus(NULL);
}

  void
mappingnotify(XEvent *e)
{
  XMappingEvent *ev = &e->xmapping;

  XRefreshKeyboardMapping(ev);
  if (ev->request == MappingKeyboard) {
    grabkeys();
  }
}

  void
maprequest(XEvent *e)
{
  static XWindowAttributes wa;
  XMapRequestEvent *ev = &e->xmaprequest;

  if (!XGetWindowAttributes(dpy, ev->window, &wa) || wa.override_redirect) {
    return;
  }
  if (!wintoclient(ev->window)) {
    manage(ev->window, &wa);
  }
}

  void
motionnotify(XEvent *e)
{
  static Monitor *mon = NULL;
  Monitor *m;
  XMotionEvent *ev = &e->xmotion;
  unsigned int i, x;

  if (ev->window == selmon->barwin) {
    i = x = 0;
    do {
      x += TEXTW(tags[i]);
    } while (ev->x >= x && ++i < LENGTH(tags));
    /* FIXME when hovering the mouse over the tags and we view the tag,
     * the preview window get's in the preview shot
     * */
    if (i < LENGTH(tags)) {
      if (selmon->previewshow != (i + 1) && !(selmon->tagset[selmon->seltags] & 1 << i)) {
        selmon->previewshow = i + 1;
        showtagpreview(i);
      } else if (selmon->tagset[selmon->seltags] & 1 << i) {
        selmon->previewshow = 0;
        XUnmapWindow(dpy, selmon->tagwin);
      }
    } else if (selmon->previewshow) {
      selmon->previewshow = 0;
      XUnmapWindow(dpy, selmon->tagwin);
    }
  } else if (ev->window == selmon->tagwin) {
    selmon->previewshow = 0;
    XUnmapWindow(dpy, selmon->tagwin);
  } else if (selmon->previewshow) {
    selmon->previewshow = 0;
    XUnmapWindow(dpy, selmon->tagwin);
  }

  if (ev->window != root) {
    return;
  }
  if ((m = recttomon(ev->x_root, ev->y_root, 1, 1)) != mon && mon) {
    unfocus(selmon->sel, 1);
    selmon = m;
    focus(NULL);
  }
  mon = m;
}

  void
movemouse(const Arg *arg)
{
  int x, y, ocx, ocy, nx, ny;
  Client *c;
  Monitor *m;
  XEvent ev;
  Time lasttime = 0;

  if (!(c = selmon->sel)) {
    return;
  }
  if (c->isfullscreen) { /* no support moving fullscreen windows by mouse */
    return;
  }
  restack(selmon);
  ocx = c->x;
  ocy = c->y;
  if (XGrabPointer(dpy, root, False, MOUSEMASK, GrabModeAsync, GrabModeAsync, None, cursor[CurMove]->cursor, CurrentTime) != GrabSuccess) {
    return;
  }
  if (!getrootptr(&x, &y)) {
    return;
  }
  do {
    XMaskEvent(dpy, MOUSEMASK|ExposureMask|SubstructureRedirectMask, &ev);
    switch(ev.type) {
      case ConfigureRequest:
      case Expose:
      case MapRequest:
        handler[ev.type](&ev);
        break;
      case MotionNotify:
        if ((ev.xmotion.time - lasttime) <= (1000 / refreshrate)) {
          continue;
        }
        lasttime = ev.xmotion.time;

        nx = ocx + (ev.xmotion.x - x);
        ny = ocy + (ev.xmotion.y - y);
        if (abs(selmon->wx - nx) < snap) {
          nx = selmon->wx;
        } else if (abs((selmon->wx + selmon->ww) - (nx + WIDTH(c))) < snap) {
          nx = selmon->wx + selmon->ww - WIDTH(c);
        }
        if (abs(selmon->wy - ny) < snap) {
          ny = selmon->wy;
        } else if (abs((selmon->wy + selmon->wh) - (ny + HEIGHT(c))) < snap) {
          ny = selmon->wy + selmon->wh - HEIGHT(c);
        }
        if (!c->isfloating && selmon->lt[selmon->sellt]->arrange && (abs(nx - c->x) > snap || abs(ny - c->y) > snap)) {
          togglefloating(NULL);
        }
        if (!selmon->lt[selmon->sellt]->arrange || c->isfloating) {
          resize(c, nx, ny, c->w, c->h, 1);
        }
        break;
    }
  } while (ev.type != ButtonRelease);
  XUngrabPointer(dpy, CurrentTime);
  if ((m = recttomon(c->x, c->y, c->w, c->h)) != selmon) {
    sendmon(c, m);
    selmon = m;
    focus(NULL);
  }
}

  void
movewin(const Arg *arg)
{
  Client *c;
  int nx, ny;
  c = selmon->sel;
  if (!c) {
    return;
  }
  if (!c->isfloating) {
    togglefloating(NULL);
  }
  nx = c->x;
  ny = c->y;
  switch (arg->ui) {
    case UP:
      ny -= c->mon->wh / 16;
      ny = MAX(ny, c->mon->wy);
      break;
    case DOWN:
      ny += c->mon->wh / 16;
      ny = MIN(ny, c->mon->wy + c->mon->wh - HEIGHT(c));
      break;
    case LEFT:
      nx -= c->mon->ww / 32;
      nx = MAX(nx, c->mon->wx);
      break;
    case RIGHT:
      nx += c->mon->ww / 32;
      nx = MIN(nx, c->mon->wx + c->mon->ww - WIDTH(c));
      break;
  }
  resize(c, nx, ny, c->w, c->h, 1);
  focus(c);
  pointerfocuswin(c);
}

  void
next_theme(const Arg *arg)
{
  current_theme_idx = (current_theme_idx + 1) % (sizeof(themes)/sizeof(themes[0]));

  /* 更新 drw scheme */
  for (int i = 0; i < SchemeLast; i++)
    scheme[i] = drw_scm_create(drw, themes[current_theme_idx][i], 3);

  /* 重绘 bar + arrange */
  for (Monitor *m = mons; m; m = m->next)
    drawbar(m);

  arrange(selmon);
}

  void
resizewin(const Arg *arg)
{
  Client *c;
  int nx, ny, nw, nh, cx, cy;
  c = selmon->sel;
  if (!c) {
    return;
  }
  if (!c->isfloating) {
    togglefloating(NULL);
  }
  nx = c->x;
  ny = c->y;
  nw = c->w;
  nh = c->h;
  cx = c->x + c->w/2;
  cy = c->y + c->h/2;
  switch (arg->ui) {
    case HORINC:
      nx = cx - c->w/2 - c->mon->ww / 32;
      nw = nw + 2 * c->mon->ww / 32;
      break;
    case HORDEC:
      nx = cx - c->w/2 + c->mon->ww / 32;
      nw = nw - 2 * c->mon->ww / 32;
      break;
    case VECINC:
      ny = cy - c->h/2 - c->mon->wh / 32;
      nh = nh + 2 * c->mon->wh / 32;
      break;
    case VECDEC:
      ny = cy - c->h/2 + c->mon->wh / 32;
      nh = nh - 2 * c->mon->wh / 32;
      break;
  }
  nw = MAX(nw, 0);
  nh = MAX(nh, 0);
  nw = MIN(nw, c->mon->ww);
  nh = MIN(nh, c->mon->wh);
  nx = MAX(nx, c->mon->wx);
  ny = MAX(ny, c->mon->wy);
  nx = MIN(nx, c->mon->ww - nw + c->mon->wx);
  ny = MIN(ny, c->mon->wh - nh + c->mon->wy);
  if (nw == 0 || nh == 0) {
    return;
  }

  resize(c, nx, ny, nw, nh, 1);
  focus(c);
  XWarpPointer(dpy, None, root, 0, 0, 0, 0, c->x + c->w - 2 * c->bw, c->y + c->h - 2 * c->bw);
}

  Client *
nexttiled(Client *c)
{
  for (; c && (c->isfloating || !ISVISIBLE(c)); c = c->next)
    ;
  return c;
}

  void
pop(Client *c)
{
  int i;
  for (i = 0; !(selmon->tagset[selmon->seltags] & 1 << i); i++)
    ;
  i++;
  c->mon->tagmarked[i] = nexttiled(c->mon->clients);

  detach(c);
  attach(c);
  focus(c);
  arrange(c->mon);
}

  void
propertynotify(XEvent *e)
{
  Client *c;
  Window trans;
  XPropertyEvent *ev = &e->xproperty;

  if ((ev->window == root) && (ev->atom == XA_WM_NAME)) {
    updatestatus();
  } else if (ev->state == PropertyDelete) {
    return; /* ignore */
  } else if ((c = wintoclient(ev->window))) {
    switch(ev->atom) {
      default: break;
      case XA_WM_TRANSIENT_FOR:
               if (!c->isfloating && (XGetTransientForHint(dpy, c->win, &trans)) &&
                   (c->isfloating = (wintoclient(trans)) != NULL))
                 arrange(c->mon);
               break;
      case XA_WM_NORMAL_HINTS:
               c->hintsvalid = 0;
               break;
      case XA_WM_HINTS:
               updatewmhints(c);
               drawbars();
               break;
    }
    if (ev->atom == XA_WM_NAME || ev->atom == netatom[NetWMName]) {
      updatetitle(c);
      if (c == c->mon->sel) {
        drawbar(c->mon);
      }
    }
    if (ev->atom == netatom[NetWMWindowType]) {
      updatewindowtype(c);
    }
  }
}

  void
quit(const Arg *arg)
{
  size_t i;

  for (i = 0; i < autostart_len; i++) {
    if (0 < autostart_pids[i]) {
      kill(autostart_pids[i], SIGTERM);
      waitpid(autostart_pids[i], NULL, 0);
    }
  }

  if(arg->i) {
    restart = 1;
  }

  running = 0;

  if (restart == 1) {
    savesession();
  }
}

  Monitor *
recttomon(int x, int y, int w, int h)
{
  Monitor *m, *r = selmon;
  int a, area = 0;

  for (m = mons; m; m = m->next) {
    if ((a = INTERSECT(x, y, w, h, m)) > area) {
      area = a;
      r = m;
    }
  }
  return r;
}

  void
resize(Client *c, int x, int y, int w, int h, int interact)
{
  if (applysizehints(c, &x, &y, &w, &h, interact)) {
    resizeclient(c, x, y, w, h);
  }
}

  void
resizeclient(Client *c, int x, int y, int w, int h)
{
  XWindowChanges wc;

  c->oldx = c->x; c->x = wc.x = x;
  c->oldy = c->y; c->y = wc.y = y;
  c->oldw = c->w; c->w = wc.width = w;
  c->oldh = c->h; c->h = wc.height = h;
  wc.border_width = c->bw;
  if (((nexttiled(c->mon->clients) == c && !nexttiled(c->next)) || &layout_monocle == c->mon->lt[c->mon->sellt]->arrange) && !c->isfullscreen && !c->isfloating) {
    c->w = wc.width  += c->bw * 2 + 1;
    c->h = wc.height += c->bw * 2;
    wc.border_width = 0;
  }
  XConfigureWindow(dpy, c->win, CWX|CWY|CWWidth|CWHeight|CWBorderWidth, &wc);
  configure(c);
  XSync(dpy, False);
}

  void
resizemouse(const Arg *arg)
{
  int ocx, ocy, nw, nh;
  Client *c;
  Monitor *m;
  XEvent ev;
  Time lasttime = 0;
  if (!(c = selmon->sel)) {
    return;
  }
  if (c->isfullscreen) { /* no support resizing fullscreen windows by mouse */
    return;
  }
  restack(selmon);
  ocx = c->x;
  ocy = c->y;
  if (XGrabPointer(dpy, root, False, MOUSEMASK, GrabModeAsync, GrabModeAsync, None, cursor[CurResize]->cursor, CurrentTime) != GrabSuccess) {
    return;
  }
  XWarpPointer(dpy, None, c->win, 0, 0, 0, 0, c->w + c->bw - 1, c->h + c->bw - 1);
  do {
    XMaskEvent(dpy, MOUSEMASK|ExposureMask|SubstructureRedirectMask, &ev);
    switch(ev.type) {
      case ConfigureRequest:
      case Expose:
      case MapRequest:
        handler[ev.type](&ev);
        break;
      case MotionNotify:
        if ((ev.xmotion.time - lasttime) <= (1000 / refreshrate)) {
          continue;
        }
        lasttime = ev.xmotion.time;
        nw = MAX(ev.xmotion.x - ocx - 2 * c->bw + 1, 1);
        nh = MAX(ev.xmotion.y - ocy - 2 * c->bw + 1, 1);
        if (c->mon->wx + nw >= selmon->wx && c->mon->wx + nw <= selmon->wx + selmon->ww && c->mon->wy + nh >= selmon->wy && c->mon->wy + nh <= selmon->wy + selmon->wh) {
          if (!c->isfloating && selmon->lt[selmon->sellt]->arrange && (abs(nw - c->w) > snap || abs(nh - c->h) > snap)) {
            togglefloating(NULL);
          }
        }
        if (!selmon->lt[selmon->sellt]->arrange || c->isfloating) {
          resize(c, c->x, c->y, nw, nh, 1);
        }
        break;
    }
  } while (ev.type != ButtonRelease);
  XWarpPointer(dpy, None, c->win, 0, 0, 0, 0, c->w + c->bw - 1, c->h + c->bw - 1);
  XUngrabPointer(dpy, CurrentTime);
  while (XCheckMaskEvent(dpy, EnterWindowMask, &ev))
    ;
  if ((m = recttomon(c->x, c->y, c->w, c->h)) != selmon) {
    sendmon(c, m);
    selmon = m;
    focus(NULL);
  }
}

  void
restack(Monitor *m)
{
  Client *c;
  XEvent ev;
  XWindowChanges wc;

  drawbar(m);
  if (!m->sel) {
    return;
  }
  if (m->sel->isfloating || !m->lt[m->sellt]->arrange) {
    XRaiseWindow(dpy, m->sel->win);
  }
  if (m->lt[m->sellt]->arrange) {
    wc.stack_mode = Below;
    wc.sibling = m->barwin;
    for (c = m->stack; c; c = c->snext) {
      if (!c->isfloating && ISVISIBLE(c)) {
        XConfigureWindow(dpy, c->win, CWSibling|CWStackMode, &wc);
        wc.sibling = c->win;
      }
    }
  }

  /* raise always-on-top windows last */
  for (c = m->stack; c; c = c->snext) {
    if (c->isontop && ISVISIBLE(c)) {
      XRaiseWindow(dpy, c->win);
    }
  }

  XSync(dpy, False);
  while (XCheckMaskEvent(dpy, EnterWindowMask, &ev))
    ;
}

void
reset(void) {
  selmon->mfact = mfact;
  selmon->hfact = hfact;
  selmon->nmaster = nmaster;

  if (selmon->sel) {
    arrange(selmon);
  } else {
    drawbar(selmon);
  }
}

  void
run(void)
{
  XEvent ev;
  /* main event loop */
  XSync(dpy, False);
  while (running && !XNextEvent(dpy, &ev)) {
    if (handler[ev.type]) {
      handler[ev.type](&ev); /* call handler */
    }
  }
}

  void
scan(void)
{
  unsigned int i, num;
  Window d1, d2, *wins = NULL;
  XWindowAttributes wa;

  if (XQueryTree(dpy, root, &d1, &d2, &wins, &num)) {
    for (i = 0; i < num; i++) {
      if (!XGetWindowAttributes(dpy, wins[i], &wa) || wa.override_redirect || XGetTransientForHint(dpy, wins[i], &d1)) {
        continue;
      }
      if (wa.map_state == IsViewable || getstate(wins[i]) == IconicState) {
        manage(wins[i], &wa);
      }
    }
    for (i = 0; i < num; i++) { /* now the transients */
      if (!XGetWindowAttributes(dpy, wins[i], &wa)) {
        continue;
      }
      if (XGetTransientForHint(dpy, wins[i], &d1) && (wa.map_state == IsViewable || getstate(wins[i]) == IconicState)) {
        manage(wins[i], &wa);
      }
    }
    if (wins) {
      XFree(wins);
    }
  }
}

  void
sendmon(Client *c, Monitor *m)
{
  if (c->mon == m) {
    return;
  }
  unfocus(c, 1);
  detach(c);
  detachstack(c);
  c->mon = m;
  c->tags = m->tagset[m->seltags]; /* assign tags of target monitor */
  attach(c);
  attachstack(c);
  focus(NULL);
  arrange(NULL);
}

  void
setclientstate(Client *c, long state)
{
  long data[] = { state, None };

  XChangeProperty(dpy, c->win, wmatom[WMState], wmatom[WMState], 32, PropModeReplace, (unsigned char *)data, 2);
}

  int
sendevent(Client *c, Atom proto)
{
  int n;
  Atom *protocols;
  int exists = 0;
  XEvent ev;

  if (XGetWMProtocols(dpy, c->win, &protocols, &n)) {
    while (!exists && n--) {
      exists = protocols[n] == proto;
    }
    XFree(protocols);
  }
  if (exists) {
    ev.type = ClientMessage;
    ev.xclient.window = c->win;
    ev.xclient.message_type = wmatom[WMProtocols];
    ev.xclient.format = 32;
    ev.xclient.data.l[0] = proto;
    ev.xclient.data.l[1] = CurrentTime;
    XSendEvent(dpy, c->win, False, NoEventMask, &ev);
  }
  return exists;
}

  void
setfocus(Client *c)
{
  if (!c->neverfocus) {
    XSetInputFocus(dpy, c->win, RevertToPointerRoot, CurrentTime);
    XChangeProperty(dpy, root, netatom[NetActiveWindow], XA_WINDOW, 32, PropModeReplace, (unsigned char *) &(c->win), 1);
  }
  sendevent(c, wmatom[WMTakeFocus]);
}

  void
setfullscreen(Client *c, int fullscreen)
{
  if (fullscreen && !c->isfullscreen) {
    XChangeProperty(dpy, c->win, netatom[NetWMState], XA_ATOM, 32, PropModeReplace, (unsigned char*)&netatom[NetWMFullscreen], 1);
    c->isfullscreen = 1;
    c->oldbw = c->bw;
    c->bw = 0;
    resizeclient(c, c->mon->mx, c->mon->my, c->mon->mw, c->mon->mh);
    XRaiseWindow(dpy, c->win);
  } else if (!fullscreen && c->isfullscreen){
    XChangeProperty(dpy, c->win, netatom[NetWMState], XA_ATOM, 32, PropModeReplace, (unsigned char*)0, 0);
    c->isfullscreen = 0;
    c->bw = c->oldbw;
    c->x = c->oldx;
    c->y = c->oldy;
    c->w = c->oldw;
    c->h = c->oldh;
    resizeclient(c, c->x, c->y, c->w, c->h);
    arrange(c->mon);
  }
}

  void
setlayout(const Arg *arg)
{
  if (!arg || !arg->v || arg->v != selmon->lt[selmon->sellt]) {
    selmon->sellt = selmon->pertag->sellts[selmon->pertag->curtag] ^= 1;
  }
  if (arg && arg->v) {
    selmon->lt[selmon->sellt] = selmon->pertag->ltidxs[selmon->pertag->curtag][selmon->sellt] = (Layout *)arg->v;
  }
  strncpy(selmon->ltsymbol, selmon->lt[selmon->sellt]->symbol, sizeof selmon->ltsymbol);
  if (selmon->sel) {
    arrange(selmon);
  } else {
    drawbar(selmon);
  }
}

  void
setmfact(const Arg *arg)
{
  float f;

  if (!arg || !selmon->lt[selmon->sellt]->arrange) {
    return;
  }
  f = arg->f + selmon->mfact;
  f = f < 0.00 ? 0.001 : f > 1.00 ? 0.999 : f;
  selmon->mfact = selmon->pertag->mfacts[selmon->pertag->curtag] = f;
  arrange(selmon);
}

  void
sethfact(const Arg *arg)
{
  float f;

  if (!arg || !selmon->lt[selmon->sellt]->arrange) {
    return;
  }
  f = arg->f + selmon->hfact;
  f = f < 0.00 ? 0.001 : f > 1.00 ? 0.999 : f;
  selmon->hfact = selmon->pertag->hfacts[selmon->pertag->curtag] = f;
  arrange(selmon);
}

  void
showtagpreview(unsigned int i)
{
  if (!selmon->previewshow || !selmon->tagmap[i]) {
    XUnmapWindow(dpy, selmon->tagwin);
    return;
  }

  XSetWindowBackgroundPixmap(dpy, selmon->tagwin, selmon->tagmap[i]);
  XCopyArea(dpy, selmon->tagmap[i], selmon->tagwin, drw->gc, 0, 0, selmon->mw / scalepreview, selmon->mh / scalepreview, 0, 0);
  XSync(dpy, False);
  XMapRaised(dpy, selmon->tagwin);
}

  void
takepreview(void)
{
  Client *c;
  Imlib_Image image;
  unsigned int occ = 0, i;

  for (c = selmon->clients; c; c = c->next) {
    occ |= c->tags;
  }

  for (i = 0; i < LENGTH(tags); i++) {
    /* searching for tags that are occupied && selected */
    if (!(occ & 1 << i) || !(selmon->tagset[selmon->seltags] & 1 << i)) {
      continue;
    }

    if (selmon->tagmap[i]) { /* tagmap exist, clean it */
      XFreePixmap(dpy, selmon->tagmap[i]);
      selmon->tagmap[i] = 0;
    }

    /* try to unmap the window so it doesn't show the preview on the preview */
    selmon->previewshow = 0;
    XUnmapWindow(dpy, selmon->tagwin);
    XSync(dpy, False);

    if (!(image = imlib_create_image(sw, sh))) {
      fprintf(stderr, "dwm: imlib: failed to create image, skipping.");
      continue;
    }
    imlib_context_set_image(image);
    imlib_context_set_display(dpy);
    /* uncomment if using alpha patch */
    //imlib_image_set_has_alpha(1);
    //imlib_context_set_blend(0);
    //imlib_context_set_visual(visual);
    imlib_context_set_visual(DefaultVisual(dpy, screen));
    imlib_context_set_drawable(root);

    if (previewbar) {
      imlib_copy_drawable_to_image(0, selmon->wx, selmon->wy, selmon->ww, selmon->wh, 0, 0, 1);
    } else {
      imlib_copy_drawable_to_image(0, selmon->mx, selmon->my, selmon->mw ,selmon->mh, 0, 0, 1);
    }
    selmon->tagmap[i] = XCreatePixmap(dpy, selmon->tagwin, selmon->mw / scalepreview, selmon->mh / scalepreview, DefaultDepth(dpy, screen));
    imlib_context_set_drawable(selmon->tagmap[i]);
    imlib_render_image_part_on_drawable_at_size(0, 0, selmon->mw, selmon->mh, 0, 0, selmon->mw / scalepreview, selmon->mh / scalepreview);
    imlib_free_image();
  }
}

  void
previewtag(const Arg *arg)
{
  if (selmon->previewshow != (arg->ui + 1)) {
    selmon->previewshow = arg->ui + 1;
  } else {
    selmon->previewshow = 0;
  }
  showtagpreview(arg->ui);
}

  void
setup(void)
{
  int i;
  XSetWindowAttributes wa;
  Atom utf8string;
  struct sigaction sa;

  /* do not transform children into zombies when they terminate */
  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_NOCLDSTOP | SA_NOCLDWAIT | SA_RESTART;
  sa.sa_handler = SIG_IGN;
  sigaction(SIGCHLD, &sa, NULL);

  /* clean up any zombies (inherited from .xinitrc etc) immediately */
  while (waitpid(-1, NULL, WNOHANG) > 0)
    ;

  signal(SIGHUP, sighup);
  signal(SIGTERM, sigterm);

  /* init screen */
  screen = DefaultScreen(dpy);
  sw = DisplayWidth(dpy, screen);
  sh = DisplayHeight(dpy, screen);
  root = RootWindow(dpy, screen);
  drw = drw_create(dpy, screen, root, sw, sh);
  if (!drw_fontset_create(drw, fonts, LENGTH(fonts))) {
    die("no fonts could be loaded.");
  }
  lrpad = drw->fonts->h;
  bh = (barheight > drw->fonts->h ) && (barheight < 3 * drw->fonts->h ) ? barheight : drw->fonts->h + 2;
  sp = sidepad;
  vp = (topbar ? 1 : -1) * vertpad;
  winpad = defaultwinpad;
  updategeom();
  /* init atoms */
  utf8string = XInternAtom(dpy, "UTF8_STRING", False);
  wmatom[WMProtocols] = XInternAtom(dpy, "WM_PROTOCOLS", False);
  wmatom[WMDelete] = XInternAtom(dpy, "WM_DELETE_WINDOW", False);
  wmatom[WMState] = XInternAtom(dpy, "WM_STATE", False);
  wmatom[WMTakeFocus] = XInternAtom(dpy, "WM_TAKE_FOCUS", False);
  netatom[NetActiveWindow] = XInternAtom(dpy, "_NET_ACTIVE_WINDOW", False);
  netatom[NetSupported] = XInternAtom(dpy, "_NET_SUPPORTED", False);
  netatom[NetWMName] = XInternAtom(dpy, "_NET_WM_NAME", False);
  netatom[NetWMState] = XInternAtom(dpy, "_NET_WM_STATE", False);
  netatom[NetWMCheck] = XInternAtom(dpy, "_NET_SUPPORTING_WM_CHECK", False);
  netatom[NetWMFullscreen] = XInternAtom(dpy, "_NET_WM_STATE_FULLSCREEN", False);
  netatom[NetWMWindowType] = XInternAtom(dpy, "_NET_WM_WINDOW_TYPE", False);
  netatom[NetWMWindowTypeDialog] = XInternAtom(dpy, "_NET_WM_WINDOW_TYPE_DIALOG", False);
  netatom[NetClientList] = XInternAtom(dpy, "_NET_CLIENT_LIST", False);
  /* init cursors */
  cursor[CurNormal] = drw_cur_create(drw, XC_left_ptr);
  cursor[CurResize] = drw_cur_create(drw, XC_sizing);
  cursor[CurMove] = drw_cur_create(drw, XC_fleur);
  /* init appearance */
  scheme = ecalloc(LENGTH(colors), sizeof(Clr *));
  for (i = 0; i < LENGTH(colors); i++) {
    scheme[i] = drw_scm_create(drw, colors[i], 3);
  }
  /* init bars */
  updatebars();
  updatestatus();
  /* supporting window for NetWMCheck */
  wmcheckwin = XCreateSimpleWindow(dpy, root, 0, 0, 1, 1, 0, 0, 0);
  XChangeProperty(dpy, wmcheckwin, netatom[NetWMCheck], XA_WINDOW, 32, PropModeReplace, (unsigned char *) &wmcheckwin, 1);
  XChangeProperty(dpy, wmcheckwin, netatom[NetWMName], utf8string, 8, PropModeReplace, (unsigned char *) "dwm", 3);
  XChangeProperty(dpy, root, netatom[NetWMCheck], XA_WINDOW, 32, PropModeReplace, (unsigned char *) &wmcheckwin, 1);
  /* EWMH support per view */
  XChangeProperty(dpy, root, netatom[NetSupported], XA_ATOM, 32, PropModeReplace, (unsigned char *) netatom, NetLast);
  XDeleteProperty(dpy, root, netatom[NetClientList]);
  /* select events */
  wa.cursor = cursor[CurNormal]->cursor;
  wa.event_mask = SubstructureRedirectMask|SubstructureNotifyMask|ButtonPressMask|PointerMotionMask|EnterWindowMask|LeaveWindowMask|StructureNotifyMask|PropertyChangeMask;
  XChangeWindowAttributes(dpy, root, CWEventMask|CWCursor, &wa);
  XSelectInput(dpy, root, wa.event_mask);
  grabkeys();
  focus(NULL);
}

  void
seturgent(Client *c, int urg)
{
  XWMHints *wmh;

  c->isurgent = urg;
  if (!(wmh = XGetWMHints(dpy, c->win))) {
    return;
  }
  wmh->flags = urg ? (wmh->flags | XUrgencyHint) : (wmh->flags & ~XUrgencyHint);
  XSetWMHints(dpy, c->win, wmh);
  XFree(wmh);
}

  void
showhide(Client *c)
{
  if (!c) {
    return;
  }
  if (ISVISIBLE(c)) {
    /* show clients top down */
    XMoveWindow(dpy, c->win, c->x, c->y);
    if ((!c->mon->lt[c->mon->sellt]->arrange || c->isfloating) && !c->isfullscreen) {
      resize(c, c->x, c->y, c->w, c->h, 0);
    }
    showhide(c->snext);
  } else {
    /* hide clients bottom up */
    showhide(c->snext);
    XMoveWindow(dpy, c->win, WIDTH(c) * -2, c->y);
  }
}

  void
sighup(int unused)
{
  Arg a = {.i = 1};
  quit(&a);
}

  void
sigterm(int unused)
{
  Arg a = {.i = 0};
  quit(&a);
}

  void
spawn(const Arg *arg)
{
  struct sigaction sa;

  if (arg->v == dmenucmd) {
    dmenumon[0] = '0' + selmon->num;
  }
  selmon->tagset[selmon->seltags] &= ~(1<<30); /* scratchpad */
  if (fork() == 0) {
    if (dpy) {
      close(ConnectionNumber(dpy));
    }
    setsid();
    sigemptyset(&sa.sa_mask);
    sa.sa_flags = 0;
    sa.sa_handler = SIG_DFL;
    sigaction(SIGCHLD, &sa, NULL);

    execvp(((char **)arg->v)[0], (char **)arg->v);
    die("dwm: execvp '%s' failed:", ((char **)arg->v)[0]);
  }
}

  static void
spawn_or_focus(const Arg *arg)
{
  const char *const *data = arg->v;  // ← 完美匹配宏类型
  const char *cmd   = data[0];
  const char *class = data[1];

  Client *c;
  XClassHint ch = { NULL, NULL };

  /* 找到窗口 -> focus */
  for (c = selmon->clients; c; c = c->next) {
    if (XGetClassHint(dpy, c->win, &ch)) {
      if (ch.res_class && strcmp(ch.res_class, class) == 0) {
        view(&(Arg){ .ui = c->tags });
        focus(c);
        arrange(selmon);

        if (ch.res_name)  XFree(ch.res_name);
        if (ch.res_class) XFree(ch.res_class);
        return;
      }
    }
    freeclasshints(&ch);
  }

  /* 没找到 -> spawn */
  // spawn(&(Arg){ .v = (const char *[]){ cmd, NULL } });
  spawn(&(Arg){ .v = (const char *[]){ "/bin/sh", "-c", cmd, NULL } });
}

  void
tag(const Arg *arg)
{
  if (selmon->sel && arg->ui & TAGMASK) {
    selmon->sel->tags = arg->ui & TAGMASK;
    focus(NULL);
    arrange(selmon);
  }
}

  void
tagmon(const Arg *arg)
{
  if (!selmon->sel || !mons->next) {
    return;
  }
  sendmon(selmon->sel, dirtomon(arg->i));
}

  void
togglebar(const Arg *arg)
{
  selmon->showbar = selmon->pertag->showbars[selmon->pertag->curtag] = !selmon->showbar;
  if (selmon->showbar) {
    winpad = defaultwinpad;
  } else {
    winpad = 0;
  }

  updatebarpos(selmon);
  XMoveResizeWindow(dpy, selmon->barwin, selmon->wx + sp, selmon->by + vp, selmon->ww - 2*sp, bh);
  arrange(selmon);
}

  void
togglefloating(const Arg *arg)
{
  if (!selmon->sel) {
    return;
  }
  if (selmon->sel->isfullscreen) { /* no support for fullscreen windows */
    return;
  }
  // selmon->sel->isfloating = !selmon->sel->isfloating || selmon->sel->isfixed;
  if (!selmon->sel->isfixed) {
    selmon->sel->isfloating = !selmon->sel->isfloating;
  } else {
    selmon->sel->isfloating = 1;
  }
  if (selmon->sel->isfloating) {
    resize(selmon->sel, selmon->sel->x, selmon->sel->y, selmon->sel->w, selmon->sel->h, 0);
  }
  arrange(selmon);
}

  void
toggle_scratchpad(const Arg *arg)
{
  const char *const *data = arg->v;
  const char *cmd   = data[0];
  const char *class = data[1];

  Client *c;

  for (c = selmon->clients; c; c = c->next) {
    if (strcmp(c->class, class) == 0) {

      /* 置浮动 */
      c->isfloating = 1;

      /* 已可见 -> 隐藏: 移到 scratchpad tag：1<<30 */
      if (ISVISIBLE(c)) {
        c->tags = 1 << 30;
        arrange(selmon);
        return;
      }

      /* 不可见 -> 显示, 放回当前 tag */
      c->tags = selmon->tagset[selmon->seltags];

      /* 居中 + 大小设置 */
      int nw = selmon->ww * scratchpad_width;
      int nh = selmon->wh * scratchpad_height;
      int nx = selmon->wx + (selmon->ww - nw) / 2;
      int ny = selmon->wy + (selmon->wh - nh) / 2;

      resize(c, nx, ny, nw, nh, 0);

      focus(c);
      arrange(selmon);
      return;
    }
  }

  /* 第一次 spawn: 记录 class，等待 manage() 处理 */
  scratchpad_class_wait = class;
  spawn(&(Arg){ .v = (const char *[]){ "/bin/sh", "-c", cmd, NULL } });
}

  void
scratchpad_to_normal(const Arg *arg)
{
  Client *c = selmon->sel;
  if (!c) { return; }

  if (c->isfloating == 1) {
    c->tags = selmon->tagset[selmon->seltags];
    c->isfloating = 0;
    arrange(selmon);
    focus(c);
    restack(selmon);
  }
}

  void
togglesticky(const Arg *arg)
{
  if (!selmon || !selmon->sel) {
    return;
  }
  selmon->sel->issticky = !selmon->sel->issticky;
  arrange(selmon);
}

  void
togglefullscreen(const Arg *arg)
{
  if(selmon->sel) {
    setfullscreen(selmon->sel, !selmon->sel->isfullscreen);
  }
}

  void
toggleoverview(const Arg *arg)
{
  static uint prevtag = 0;

  if (!selmon->isoverview) {
    /* 进入 overview，记录当前 tagset */
    prevtag = selmon->tagset[selmon->seltags];
    selmon->isoverview = 1;

    /* 展示所有窗口 */
    view(&(Arg){ .ui = ~0 });
  } else {
    /* 退出 overview，恢复原 tag */
    selmon->isoverview = 0;

    /* 如果 prevtag 合法就恢复，否则回到 1 号 tag */
    view(&(Arg){ .ui = prevtag ? prevtag : 1 });
  }
}

  void
toggletag(const Arg *arg)
{
  unsigned int newtags;

  if (!selmon->sel) {
    return;
  }
  newtags = selmon->sel->tags ^ (arg->ui & TAGMASK);
  if (newtags) {
    selmon->sel->tags = newtags;
    focus(NULL);
    arrange(selmon);
  }
}

  void
toggleview(const Arg *arg)
{
  unsigned int newtagset = selmon->tagset[selmon->seltags] ^ (arg->ui & TAGMASK);
  int i;

  if (newtagset) {
    takepreview();
    selmon->tagset[selmon->seltags] = newtagset;

    if (newtagset == ~0) {
      selmon->pertag->prevtag = selmon->pertag->curtag;
      selmon->pertag->curtag = 0;
    }

    /* test if the user did not select the same tag */
    if (!(newtagset & 1 << (selmon->pertag->curtag - 1))) {
      selmon->pertag->prevtag = selmon->pertag->curtag;
      for (i = 0; !(newtagset & 1 << i); i++)
        ;
      selmon->pertag->curtag = i + 1;
    }

    /* apply settings for this view */
    selmon->nmaster = selmon->pertag->nmasters[selmon->pertag->curtag];
    selmon->mfact = selmon->pertag->mfacts[selmon->pertag->curtag];
    selmon->hfact = selmon->pertag->hfacts[selmon->pertag->curtag];
    selmon->sellt = selmon->pertag->sellts[selmon->pertag->curtag];
    selmon->lt[selmon->sellt] = selmon->pertag->ltidxs[selmon->pertag->curtag][selmon->sellt];
    selmon->lt[selmon->sellt^1] = selmon->pertag->ltidxs[selmon->pertag->curtag][selmon->sellt^1];

    if (selmon->showbar != selmon->pertag->showbars[selmon->pertag->curtag]) {
      togglebar(NULL);
    }

    focus(NULL);
    arrange(selmon);
  }
}

  void
unfocus(Client *c, int setfocus)
{
  if (!c) {
    return;
  }
  grabbuttons(c, 0);
  XSetWindowBorder(dpy, c->win, scheme[SchemeNorm][ColBorder].pixel);
  if (setfocus) {
    XSetInputFocus(dpy, root, RevertToPointerRoot, CurrentTime);
    XDeleteProperty(dpy, root, netatom[NetActiveWindow]);
  }
}

  void
unmanage(Client *c, int destroyed)
{
  Monitor *m = c->mon;
  XWindowChanges wc;

  if (c->swallowing) {
    unswallow(c);
    return;
  }

  Client *s = swallowingclient(c->win);
  if (s) {
    free(s->swallowing);
    s->swallowing = NULL;
    arrange(m);
    focus(NULL);
    return;
  }

  detach(c);
  detachstack(c);
  if (!destroyed) {
    wc.border_width = c->oldbw;
    XGrabServer(dpy); /* avoid race conditions */
    XSetErrorHandler(xerrordummy);
    XSelectInput(dpy, c->win, NoEventMask);
    XConfigureWindow(dpy, c->win, CWBorderWidth, &wc); /* restore border */
    XUngrabButton(dpy, AnyButton, AnyModifier, c->win);
    setclientstate(c, WithdrawnState);
    XSync(dpy, False);
    XSetErrorHandler(xerror);
    XUngrabServer(dpy);
  }

  free(c);
  if (!s) {
    arrange(m);
    focus(NULL);
    updateclientlist();
  }
}

  void
unmapnotify(XEvent *e)
{
  Client *c;
  XUnmapEvent *ev = &e->xunmap;

  if ((c = wintoclient(ev->window))) {
    if (ev->send_event) {
      setclientstate(c, WithdrawnState);
    } else {
      unmanage(c, 0);
    }
  }
}

  void
updatebars(void)
{
  Monitor *m;
  XSetWindowAttributes wa = {
    .override_redirect = True,
    .background_pixmap = ParentRelative,
    .event_mask = ButtonPressMask|ExposureMask|PointerMotionMask
  };

  XClassHint ch = {"dwm", "dwm"};
  for (m = mons; m; m = m->next) {
    if (!m->tagwin) {
      m->tagwin = XCreateWindow(dpy, root, m->wx, m->by + bh, m->mw / scalepreview, m->mh / scalepreview, 0, DefaultDepth(dpy, screen), CopyFromParent, DefaultVisual(dpy, screen), CWOverrideRedirect|CWBackPixmap|CWEventMask, &wa);
      XDefineCursor(dpy, m->tagwin, cursor[CurNormal]->cursor);
      XUnmapWindow(dpy, m->tagwin);
    }
    if (m->barwin) {
      continue;
    }
    m->barwin = XCreateWindow(dpy, root, m->wx + sp, m->by + vp, m->ww - 2 * sp, bh, 0, DefaultDepth(dpy, screen), CopyFromParent, DefaultVisual(dpy, screen), CWOverrideRedirect|CWBackPixmap|CWEventMask, &wa);
    XDefineCursor(dpy, m->barwin, cursor[CurNormal]->cursor);
    XMapRaised(dpy, m->barwin);
    XSetClassHint(dpy, m->barwin, &ch);
  }
}

  void
updatebarpos(Monitor *m)
{
  m->wy = m->my;
  m->wh = m->mh;
  if (m->showbar) {
    m->wh = m->wh - vertpad - bh;
    m->by = m->topbar ? m->wy : m->wy + m->wh + vertpad;
    m->wy = m->topbar ? m->wy + bh + vp : m->wy;
  } else {
    m->by = -bh - vp;
  }
}

  void
updateclientlist(void)
{
  Client *c;
  Monitor *m;

  XDeleteProperty(dpy, root, netatom[NetClientList]);
  for (m = mons; m; m = m->next) {
    for (c = m->clients; c; c = c->next) {
      XChangeProperty(dpy, root, netatom[NetClientList], XA_WINDOW, 32, PropModeAppend, (unsigned char *) &(c->win), 1);
    }
  }
}

  int
updategeom(void)
{
  int dirty = 0;

#ifdef XINERAMA
  if (XineramaIsActive(dpy)) {
    int i, j, n, nn;
    Client *c;
    Monitor *m;
    XineramaScreenInfo *info = XineramaQueryScreens(dpy, &nn);
    XineramaScreenInfo *unique = NULL;

    for (n = 0, m = mons; m; m = m->next, n++)
      ;

    /* only consider unique geometries as separate screens */
    unique = ecalloc(nn, sizeof(XineramaScreenInfo));
    for (i = 0, j = 0; i < nn; i++) {
      if (isuniquegeom(unique, j, &info[i])) {
        memcpy(&unique[j++], &info[i], sizeof(XineramaScreenInfo));
      }
    }

    XFree(info);
    nn = j;

    /* new monitors if nn > n */
    for (i = n; i < nn; i++) {
      for (m = mons; m && m->next; m = m->next)
        ;
      if (m) {
        m->next = createmon();
      } else {
        mons = createmon();
      }
    }
    for (i = 0, m = mons; i < nn && m; m = m->next, i++) {
      if (i >= n || unique[i].x_org != m->mx || unique[i].y_org != m->my || unique[i].width != m->mw || unique[i].height != m->mh) {
        dirty = 1;
        m->num = i;
        m->mx = m->wx = unique[i].x_org;
        m->my = m->wy = unique[i].y_org;
        m->mw = m->ww = unique[i].width;
        m->mh = m->wh = unique[i].height;
        updatebarpos(m);
      }
    }
    /* removed monitors if n > nn */
    for (i = nn; i < n; i++) {
      for (m = mons; m && m->next; m = m->next)
        ;
      while ((c = m->clients)) {
        dirty = 1;
        m->clients = c->next;
        detachstack(c);
        c->mon = mons;
        attach(c);
        attachstack(c);
      }
      if (m == selmon) {
        selmon = mons;
      }
      cleanupmon(m);
    }
    free(unique);
  } else
#endif /* XINERAMA */
  { /* default monitor setup */
    if (!mons)
      mons = createmon();
    if (mons->mw != sw || mons->mh != sh) {
      dirty = 1;
      mons->mw = mons->ww = sw;
      mons->mh = mons->wh = sh;
      updatebarpos(mons);
    }
  }
  if (dirty) {
    selmon = mons;
    selmon = wintomon(root);
  }
  return dirty;
}

  void
updatenumlockmask(void)
{
  unsigned int i, j;
  XModifierKeymap *modmap;

  numlockmask = 0;
  modmap = XGetModifierMapping(dpy);
  for (i = 0; i < 8; i++) {
    for (j = 0; j < modmap->max_keypermod; j++) {
      if (modmap->modifiermap[i * modmap->max_keypermod + j] == XKeysymToKeycode(dpy, XK_Num_Lock)) {
        numlockmask = (1 << i);
      }
    }
  }
  XFreeModifiermap(modmap);
}

  void
updatesizehints(Client *c)
{
  long msize;
  XSizeHints size;

  if (!XGetWMNormalHints(dpy, c->win, &size, &msize)) { /* size is uninitialized, ensure that size.flags aren't used */
    size.flags = PSize;
  }

  if (size.flags & PBaseSize) {
    c->basew = size.base_width;
    c->baseh = size.base_height;
  } else if (size.flags & PMinSize) {
    c->basew = size.min_width;
    c->baseh = size.min_height;
  } else {
    c->basew = c->baseh = 0;
  }

  if (size.flags & PResizeInc) {
    c->incw = size.width_inc;
    c->inch = size.height_inc;
  } else {
    c->incw = c->inch = 0;
  }

  if (size.flags & PMaxSize) {
    c->maxw = size.max_width;
    c->maxh = size.max_height;
  } else {
    c->maxw = c->maxh = 0;
  }

  if (size.flags & PMinSize) {
    c->minw = size.min_width;
    c->minh = size.min_height;
  } else if (size.flags & PBaseSize) {
    c->minw = size.base_width;
    c->minh = size.base_height;
  } else {
    c->minw = c->minh = 0;
  }

  if (size.flags & PAspect) {
    c->mina = (float)size.min_aspect.y / size.min_aspect.x;
    c->maxa = (float)size.max_aspect.x / size.max_aspect.y;
  } else {
    c->maxa = c->mina = 0.0;
  }

  c->isfixed = (c->maxw && c->maxh && c->maxw == c->minw && c->maxh == c->minh);
  c->hintsvalid = 1;
}

  void
updatestatus(void)
{
  if (!gettextprop(root, XA_WM_NAME, stext, sizeof(stext))) {
    strcpy(stext, "dwm-"VERSION);
  }
  drawbar(selmon);
}

  void
updatetitle(Client *c)
{
  if (!gettextprop(c->win, netatom[NetWMName], c->name, sizeof c->name)) {
    gettextprop(c->win, XA_WM_NAME, c->name, sizeof c->name);
  }
  if (c->name[0] == '\0') { /* hack to mark broken clients */
    strcpy(c->name, broken);
  }
}

  void
updatewindowtype(Client *c)
{
  Atom state = getatomprop(c, netatom[NetWMState]);
  Atom wtype = getatomprop(c, netatom[NetWMWindowType]);

  if (state == netatom[NetWMFullscreen]) {
    setfullscreen(c, 1);
  }
  if (wtype == netatom[NetWMWindowTypeDialog]) {
    c->isfloating = 1;
  }
}

  void
updatewmhints(Client *c)
{
  XWMHints *wmh;

  if ((wmh = XGetWMHints(dpy, c->win))) {
    if (c == selmon->sel && wmh->flags & XUrgencyHint) {
      wmh->flags &= ~XUrgencyHint;
      XSetWMHints(dpy, c->win, wmh);
    } else {
      c->isurgent = (wmh->flags & XUrgencyHint) ? 1 : 0;
    }

    if (wmh->flags & InputHint) {
      c->neverfocus = !wmh->input;
    } else {
      c->neverfocus = 0;
    }

    XFree(wmh);
  }
}

  void
view(const Arg *arg)
{
  int i;
  unsigned int tmptag;

  if ((arg->ui & TAGMASK) == selmon->tagset[selmon->seltags]) {
    arrange(selmon);
    return;
  }
  takepreview();

  selmon->seltags ^= 1; /* toggle sel tagset */
  if (arg->ui & TAGMASK) {
    selmon->tagset[selmon->seltags] = arg->ui & TAGMASK;
    selmon->pertag->prevtag = selmon->pertag->curtag;

    if (arg->ui == ~0) {
      selmon->pertag->curtag = 0;
    } else {
      for (i = 0; !(arg->ui & 1 << i); i++)
        ;
      selmon->pertag->curtag = i + 1;
    }
  } else {
    tmptag = selmon->pertag->prevtag;
    selmon->pertag->prevtag = selmon->pertag->curtag;
    selmon->pertag->curtag = tmptag;
  }

  selmon->nmaster = selmon->pertag->nmasters[selmon->pertag->curtag];
  selmon->mfact = selmon->pertag->mfacts[selmon->pertag->curtag];
  selmon->hfact = selmon->pertag->hfacts[selmon->pertag->curtag];
  selmon->sellt = selmon->pertag->sellts[selmon->pertag->curtag];
  selmon->lt[selmon->sellt] = selmon->pertag->ltidxs[selmon->pertag->curtag][selmon->sellt];
  selmon->lt[selmon->sellt^1] = selmon->pertag->ltidxs[selmon->pertag->curtag][selmon->sellt^1];

  if (selmon->showbar != selmon->pertag->showbars[selmon->pertag->curtag]) {
    togglebar(NULL);
  }

  focus(NULL);
  arrange(selmon);
}

  pid_t
winpid(Window w)
{
  pid_t result = 0;

#ifdef __linux__
  xcb_res_client_id_spec_t spec = {0};
  spec.client = w;
  spec.mask = XCB_RES_CLIENT_ID_MASK_LOCAL_CLIENT_PID;

  xcb_generic_error_t *e = NULL;
  xcb_res_query_client_ids_cookie_t c = xcb_res_query_client_ids(xcon, 1, &spec);
  xcb_res_query_client_ids_reply_t *r = xcb_res_query_client_ids_reply(xcon, c, &e);

  if (!r) {
    return (pid_t)0;
  }

  xcb_res_client_id_value_iterator_t i = xcb_res_query_client_ids_ids_iterator(r);
  for (; i.rem; xcb_res_client_id_value_next(&i)) {
    spec = i.data->spec;
    if (spec.mask & XCB_RES_CLIENT_ID_MASK_LOCAL_CLIENT_PID) {
      uint32_t *t = xcb_res_client_id_value_value(i.data);
      result = *t;
      break;
    }
  }

  free(r);

  if (result == (pid_t)-1) {
    result = 0;
  }
#endif /* __linux__ */

#ifdef __OpenBSD__
  Atom type;
  pid_t ret;
  int format;
  unsigned long len, bytes;
  unsigned char *prop;

  if (XGetWindowProperty(dpy, w, XInternAtom(dpy, "_NET_WM_PID", 0), 0, 1, False, AnyPropertyType, &type, &format, &len, &bytes, &prop) != Success || !prop) {
    return 0;
  }

  ret = *(pid_t*)prop;
  XFree(prop);
  result = ret;
#endif /* __OpenBSD__ */

  return result;
}

  pid_t
getparentprocess(pid_t p)
{
  unsigned int v = 0;

#ifdef __linux__
  FILE *f;
  char buf[256];
  snprintf(buf, sizeof(buf) - 1, "/proc/%u/stat", (unsigned)p);

  if (!(f = fopen(buf, "r"))) {
    return 0;
  }

  fscanf(f, "%*u %*s %*c %u", &v);
  fclose(f);
#endif /* __linux__*/

#ifdef __OpenBSD__
  int n;
  kvm_t *kd;
  struct kinfo_proc *kp;

  kd = kvm_openfiles(NULL, NULL, NULL, KVM_NO_FILES, NULL);
  if (!kd) {
    return 0;
  }

  kp = kvm_getprocs(kd, KERN_PROC_PID, p, sizeof(*kp), &n);
  v = kp->p_ppid;
#endif /* __OpenBSD__ */

  return (pid_t)v;
}

  int
isdescprocess(pid_t p, pid_t c)
{
  while (p != c && c != 0) {
    c = getparentprocess(c);
  }

  return (int)c;
}

  Client *
termforwin(const Client *w)
{
  Client *c;
  Monitor *m;

  if (!w->pid || w->isterminal) {
    return NULL;
  }

  for (m = mons; m; m = m->next) {
    for (c = m->clients; c; c = c->next) {
      if (c->isterminal && !c->swallowing && c->pid && isdescprocess(c->pid, w->pid)) {
        return c;
      }
    }
  }

  return NULL;
}

  Client *
swallowingclient(Window w)
{
  Client *c;
  Monitor *m;

  for (m = mons; m; m = m->next) {
    for (c = m->clients; c; c = c->next) {
      if (c->swallowing && c->swallowing->win == w) {
        return c;
      }
    }
  }

  return NULL;
}

  Client *
wintoclient(Window w)
{
  Client *c;
  Monitor *m;

  for (m = mons; m; m = m->next) {
    for (c = m->clients; c; c = c->next) {
      if (c->win == w) {
        return c;
      }
    }
  }
  return NULL;
}

  Monitor *
wintomon(Window w)
{
  int x, y;
  Client *c;
  Monitor *m;

  if (w == root && getrootptr(&x, &y)) {
    return recttomon(x, y, 1, 1);
  }

  for (m = mons; m; m = m->next) {
    if (w == m->barwin) {
      return m;
    }
  }

  if ((c = wintoclient(w))) {
    return c->mon;
  }

  return selmon;
}

/* There's no way to check accesses to destroyed windows, thus those cases are
 * ignored (especially on UnmapNotify's). Other types of errors call Xlibs
 * default error handler, which may call exit. */
  int
xerror(Display *dpy, XErrorEvent *ee)
{
  if (ee->error_code == BadWindow
      || (ee->request_code == X_SetInputFocus && ee->error_code == BadMatch)
      || (ee->request_code == X_PolyText8 && ee->error_code == BadDrawable)
      || (ee->request_code == X_PolyFillRectangle && ee->error_code == BadDrawable)
      || (ee->request_code == X_PolySegment && ee->error_code == BadDrawable)
      || (ee->request_code == X_ConfigureWindow && ee->error_code == BadMatch)
      || (ee->request_code == X_GrabButton && ee->error_code == BadAccess)
      || (ee->request_code == X_GrabKey && ee->error_code == BadAccess)
      || (ee->request_code == X_CopyArea && ee->error_code == BadDrawable)) {
    return 0;
  }

  fprintf(stderr, "dwm: fatal error: request code=%d, error code=%d\n", ee->request_code, ee->error_code);

  return xerrorxlib(dpy, ee); /* may call exit */
}

  int
xerrordummy(Display *dpy, XErrorEvent *ee)
{
  return 0;
}

/* Startup Error handler to check if another window manager
 * is already running. */
  int
xerrorstart(Display *dpy, XErrorEvent *ee)
{
  die("dwm: another window manager is already running");
  return -1;
}

  void
zoom(const Arg *arg)
{
  Client *c = selmon->sel;

  if (!selmon->lt[selmon->sellt]->arrange || !c || c->isfloating) {
    return;
  }
  if (c == nexttiled(selmon->clients) && !(c = nexttiled(c->next))) {
    return;
  }
  pop(c);
}

  void
cyclelayout(const Arg *arg)
{
  Layout *l;
  for(l = (Layout *)layouts; l != selmon->lt[selmon->sellt]; l++)
    ;
  if(arg->i > 0) {
    if(l->symbol && (l + 1)->symbol) {
      setlayout(&((Arg) { .v = (l + 1) }));
    } else {
      setlayout(&((Arg) { .v = layouts }));
    }
  } else {
    if(l != layouts && (l - 1)->symbol) {
      setlayout(&((Arg) { .v = (l - 1) }));
    } else {
      setlayout(&((Arg) { .v = &layouts[LENGTH(layouts) - 2] }));
    }
  }
}

  void
savesession(void)
{
  FILE *fw = fopen(SESSION_FILE, "w");
  for (Client *c = selmon->clients; c != NULL; c = c->next) {
    fprintf(fw, "%lu %u\n", c->win, c->tags);
  }
  fclose(fw);
}

  void
restoresession(void)
{
  FILE *fr = fopen(SESSION_FILE, "r");
  if (!fr) {
    return;
  }

  char *str = malloc(23 * sizeof(char));
  while (fscanf(fr, "%[^\n] ", str) != EOF) {
    long unsigned int winId;
    unsigned int tagsForWin;
    int check = sscanf(str, "%lu %u", &winId, &tagsForWin);
    if (check != 2) {
      break;
    }

    for (Client *c = selmon->clients; c ; c = c->next) {
      if (c->win == winId) {
        c->tags = tagsForWin;
        break;
      }
    }
  }

  for (Client *c = selmon->clients; c ; c = c->next) {
    focus(c);
    restack(c->mon);
  }

  for (Monitor *m = selmon; m; m = m->next) {
    arrange(m);
  }

  free(str);
  fclose(fr);
  remove(SESSION_FILE);
}

  void
movestack(const Arg *arg)
{
  Client *c = NULL, *p = NULL, *pc = NULL, *i;

  /* early exit if no selected client. panic if not check. fix movestack patch bug */
  if (!selmon->sel) return;

  if(arg->i > 0) {
    /* find the client after selmon->sel */
    for(c = selmon->sel->next; c && (!ISVISIBLE(c) || c->isfloating); c = c->next);
    if(!c) {
      for(c = selmon->clients; c && (!ISVISIBLE(c) || c->isfloating); c = c->next);
    }
  } else {
    /* find the client before selmon->sel */
    for(i = selmon->clients; i != selmon->sel; i = i->next) {
      if(ISVISIBLE(i) && !i->isfloating) {
        c = i;
      }
    }
    if(!c) {
      for(; i; i = i->next) {
        if(ISVISIBLE(i) && !i->isfloating) {
          c = i;
        }
      }
    }
  }

  /* no client to swap with or selmon->sel is the only client */
  if (!c || c == selmon->sel) return;

  /* find the client p that before selmon->sel and c */
  for(i = selmon->clients; i && (!p || !pc); i = i->next) {
    if(i->next == selmon->sel) {
      p = i;
    }
    if(i->next == c) {
      pc = i;
    }
  }

  /* swap c and selmon->sel selmon->clients in the selmon->clients list */
  if(c && c != selmon->sel) {
    Client *temp = selmon->sel->next==c ? selmon->sel : selmon->sel->next;
    selmon->sel->next = c->next==selmon->sel ? c : c->next;
    c->next = temp;

    if(p && p != c) {
      p->next = c;
    }
    if(pc && pc != selmon->sel) {
      pc->next = selmon->sel;
    }

    if(selmon->sel == selmon->clients) {
      selmon->clients = c;
    } else if(c == selmon->clients) {
      selmon->clients = selmon->sel;
    }

    arrange(selmon);
  }
}

  static void
shiftview(const Arg *arg)
{
  Arg shifted;

  if(arg->i > 0) {
    shifted.ui = (selmon->tagset[selmon->seltags] << arg->i)   | selmon->tagset[selmon->seltags] >> (LENGTH(tags) - arg->i);
  } else {
    shifted.ui = selmon->tagset[selmon->seltags] >> (- arg->i) | selmon->tagset[selmon->seltags] << (LENGTH(tags) + arg->i);
  }

  view(&shifted);
}

/* layout begin */
  void
layout_monocle(Monitor *m)
{
  unsigned int n = 0;
  Client *c;

  for (c = m->clients; c; c = c->next) {
    if (ISVISIBLE(c)) {
      n++;
    }
  }

  if (n > 0) { /* override layout symbol */
    snprintf(m->ltsymbol, sizeof m->ltsymbol, "%s %d", selmon->lt[selmon->sellt]->symbol, n);
  }

  for (c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
    resize(
        c,
        m->wx,
        m->wy,
        m->ww - 2*c->bw,
        m->wh - 2*c->bw,
        False
        );
  }
}

  void
layout_center_free_shape(Monitor *m)
{
  unsigned int n, i;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > 0) { snprintf(m->ltsymbol, sizeof m->ltsymbol, "%s %d", selmon->lt[selmon->sellt]->symbol, n); }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    resize(
        c,
        m->ww / 2 - (m->ww * m->mfact) / 2,
        m->wy + m->wh / 2 - (m->wh * m->hfact) / 2,
        m->ww * m->mfact - 2*c->bw,
        m->wh * m->hfact - 2*c->bw,
        False
        );
  }
}

  void
layout_center_equal_ratio(Monitor *m)
{
  unsigned int n, i;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > 0) { snprintf(m->ltsymbol, sizeof m->ltsymbol, "%s %d", selmon->lt[selmon->sellt]->symbol, n); }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    resize(
        c,
        m->ww/2 - (m->ww*m->hfact)/2,
        m->wy + m->wh/2 - (m->wh*m->hfact)/2 + (topbar ? 1 : 0)*winpad,
        m->ww*m->hfact - 2*c->bw,
        (m->wh - (topbar ? 1 : 0)*winpad)*m->hfact - 2*c->bw,
        False
        );
  }
}

  void
layout_fibonacci(Monitor *m, int s)
{
  unsigned int i, n;
  Client *c;
  unsigned int nx, ny, nw, nh;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  // initialize window positions and sizes
  nx = m->wx;
  ny = 0;
  nw = m->ww;
  nh = m->wh - (topbar ? 1 : 0) * winpad;

  // main layout loop
  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    // adjust width or height if necessary
    if ((i % 2 && nh / 2 > 2 * c->bw) || (!(i % 2) && nw / 2 > 2 * c->bw)) {
      if (i < n - 1) {
        if (i % 2) {
          nh /= 2;
        } else {
          nw /= 2;
        }

        if ((i % 4) == 2 && !s) {
          nx += nw;
        } else if ((i % 4) == 3 && !s) {
          ny += nh;
        }
      }

      // adjust position based on fibonacci sequence logic
      switch (i % 4) {
        case 0:
          ny = s ? ny + nh : ny - nh;
          break;
        case 1:
          nx += nw;
          break;
        case 2:
          ny += nh;
          break;
        case 3:
          nx = s ? nx + nw : nx - nw;
          break;
      }

      // update sizes for the first two clients
      if (i == 0) {
        if (n != 1) {
          nw = m->ww * m->mfact;
        }
        ny = m->wy;
      } else if (i == 1) {
        nw = m->ww - nw;
      }
    }

    // resize the client window
    resize(c, nx, ny + (topbar ? 1 : 0) * winpad, nw - 2 * c->bw, nh - 2 * c->bw, False);
  }
}

  void
layout_fib_dwindle(Monitor *m)
{
  layout_fibonacci(m, 1);
}

  void
layout_fib_spiral(Monitor *m)
{
  layout_fibonacci(m, 0);
}

  void
layout_grid(Monitor *m)
{
  unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  cols = 1;
  while (cols * cols < n) {
    cols++;
  }
  rows = (n + cols - 1) / cols;

  ch = (m->wh - (topbar ? 1 : 0) * winpad) / rows;
  cw = m->ww / cols;

  ah = (m->wh - rows * ch) / 2;
  aw = (m->ww - cols * cw) / 2;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    cx = m->wx + aw + (i % cols) * cw;
    cy = m->wy + ah + (i / cols) * ch;

    if (i >= cols * (rows - 1) && n != cols * rows) {
      cx += ((cw + aw) * (cols - n % cols)) / 2;
    }

    resize(c, cx, cy + (topbar ? 1 : 0) * winpad, cw - 2 * c->bw, ch - 2 * c->bw, False);
  }
}

  void
layout_tile_right(Monitor *m)
{
  unsigned int i, n, h, mw, my = 0, ty = 0;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  mw = (n > m->nmaster) ? (m->nmaster ? m->ww * m->mfact : 0) : m->ww;

  unsigned int topbar_offset = topbar ? 1 : 0;
  unsigned int winpad_offset = topbar_offset * winpad;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      h = (m->wh - my - winpad_offset) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx, m->wy + my + winpad_offset, mw - 2 * c->bw, h - 2 * c->bw, False);
      if (my + HEIGHT(c) < m->wh) {
        my += HEIGHT(c);
      }
    } else {
      h = (m->wh - ty - winpad_offset) / (n - i);
      resize(c, m->wx + mw, m->wy + ty + winpad_offset, m->ww - mw - 2 * c->bw, h - 2 * c->bw, False);
      if (ty + HEIGHT(c) < m->wh) {
        ty += HEIGHT(c);
      }
    }
  }
}

  void
layout_tile_left(Monitor *m)
{
  unsigned int i, n, h, mw, my = 0, ty = 0;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  mw = (n > m->nmaster) ? (m->nmaster ? m->ww * (1 - m->mfact) : 0) : m->ww;

  unsigned int topbar_offset = topbar ? 1 : 0;
  unsigned int winpad_offset = topbar_offset * winpad;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      h = (m->wh - my - winpad_offset) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + m->ww - mw, m->wy + my + winpad_offset, mw - 2 * c->bw, h - 2 * c->bw, False);
      if (my + HEIGHT(c) < m->wh) {
        my += HEIGHT(c);
      }
    } else {
      h = (m->wh - ty - winpad_offset) / (n - i);
      resize(c, m->wx, m->wy + ty + winpad_offset, m->ww - mw - 2 * c->bw, h - 2 * c->bw, False);
      if (ty + HEIGHT(c) < m->wh) {
        ty += HEIGHT(c);
      }
    }
  }
}

  void
layout_stack_hori(Monitor *m)
{
  int w, mh, mx = 0, tx, ty, th;
  unsigned int i, n;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  unsigned int topbar_offset = topbar ? 1 : 0;
  unsigned int winpad_offset = topbar_offset * winpad;

  if (n > m->nmaster) {
    mh = m->nmaster ? (1 - m->hfact) * m->wh : 0;
    th = (m->wh - mh - winpad_offset) / (n - m->nmaster);
    ty = m->wy + mh;
  } else {
    th = mh = m->wh - winpad_offset;
    ty = m->wy;
  }

  for (i = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + mx, m->wy + winpad_offset, w - 2 * c->bw, mh - 2 * c->bw, False);
      mx += WIDTH(c);
    } else {
      resize(c, tx, ty + winpad_offset, m->ww - 2 * c->bw, th - 2 * c->bw, False);
      if (th != m->wh - winpad_offset) {
        ty += HEIGHT(c);
      }
    }
  }
}

  void
layout_stack_vert(Monitor *m)
{
  int w, h, mh, mx = 0, tx, ty, tw;
  unsigned int i, n;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  unsigned int topbar_offset = topbar ? 1 : 0;
  unsigned int winpad_offset = topbar_offset * winpad;

  if (n > m->nmaster) {
    mh = m->nmaster ? (1 - m->hfact) * (m->wh - winpad_offset) : 0;
    tw = m->ww / (n - m->nmaster);
    ty = m->wy + mh;
  } else {
    mh = m->wh - winpad_offset;
    tw = m->ww;
    ty = m->wy;
  }

  for (i = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + mx, m->wy + winpad_offset, w - 2 * c->bw, mh - 2 * c->bw, False);
      mx += WIDTH(c);
    } else {
      h = m->wh - winpad_offset - mh;
      resize(c, tx, ty + winpad_offset, tw - 2 * c->bw, h - 2 * c->bw, False);
      if (tw != m->ww) {
        tx += WIDTH(c);
      }
    }
  }
}

  void
layout_hacker(Monitor *m)
{
  unsigned int i, n;
  int cx, cy, cw, ch;
  int offset_x, offset_y, initial_offset_x, initial_offset_y;
  int center_x, center_y;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) return;

  cw = m->ww * 3 / 5;
  ch = m->wh * 3 / 5;

  center_x = m->wx + (m->ww - cw) / 2;
  center_y = m->wy + (m->wh - ch) / 2;

  initial_offset_x = m->wx + m->ww * 0.01;
  initial_offset_y = m->wy + m->wh * 0.01;

  offset_x = m->ww / 32;
  offset_y = m->wh / 32;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    cx = initial_offset_x + (n - i - 1) * offset_x;
    cy = initial_offset_y + (n - i - 1) * offset_y;

    if (cy + ch > m->wh) {
      cx = center_x;
      cy = center_y;
    }

    resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
  }
}

  void
layout_grid_gap(Monitor *m)
{

  unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
  unsigned int gapoh     = 24;
  unsigned int gapow     = 32;
  unsigned int gapih     = 12;
  unsigned int gapiw     = 16;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  for (cols = 0; cols <= n / 2; cols++) {
    if (cols * cols >= n) {
      break;
    }
  }

  rows = (cols && (cols - 1) * cols >= n) ? cols - 1 : cols;
  ch = (m->wh - 2 * gapoh) / (rows ? rows : 1);
  cw = (m->ww - 2 * gapow) / (cols ? cols : 1);
  ah = rows ? (m->wh - 2 * gapoh - rows * ch) / 2 : 0;
  aw = cols ? (m->ww - 2 * gapow - cols * cw) / 2 : 0;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    cx = m->wx + gapow + aw + (i % cols) * cw;
    cy = m->wy + gapoh + ah + (i / cols) * ch;
    if (i > cols * (rows - 1) - 1 && n != cols * rows) {
      cx = m->wx + gapow + aw + (i % cols) * cw + ((cw + aw) * (cols - n % cols))/2;
      cy = m->wy + gapoh + ah + (i / cols) * ch;
    }
    resize(
        c,
        cx,
        cy,
        cw - gapiw / 2 - 2 * c->bw,
        ch - gapih / 2 - 2 * c->bw,
        False
        );
  }
}

  void
layout_overview(Monitor *m)
{
  unsigned int gapoh     = 24;
  unsigned int gapow     = 32;
  unsigned int gapih     = 12;
  unsigned int gapiw     = 16;

  unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  switch (n) {
    case 1: {
              for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                cw = m->ww * 13 / 20;
                ch = m->wh * 99 / 100;
                cx = m->wx + (m->ww - cw) / 2;
                cy = m->wy + (m->wh - ch) / 2;
                resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
              }
              break;
            }
    case 2: {
              for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                switch (i) {
                  case 0:
                    cw = m->ww * 6 / 10 - gapow - gapih;
                    ch = m->wh * 93 / 100;
                    cx = m->wx + m->ww * 1 / 60 + gapow;
                    cy = m->wy + (m->wh - ch) / 2;
                    break;
                  default:
                    cw = m->ww * 3 / 10 - gapow;
                    ch = m->wh * 6 / 20 - gapoh;
                    cx = m->wx + m->ww - gapow - cw - m->ww * 1 / 60;
                    cy = m->wy + gapoh + m->wh * 2 / 20;
                    break;
                };
                resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
              }
              break;
            }
    case 3: {
              for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                switch (i) {
                  case 2: // left
                    cw = m->ww * 16 / 100 - gapow - gapiw;
                    ch = m->wh * 3 / 10;
                    cx = m->wx + gapow + m->ww * 1 / 200;
                    cy = m->wy + (m->wh - ch) / 2;
                    break;
                  case 0: // center -- master
                    cw = m->ww * 60 / 100;
                    ch = m->wh * 70 / 100;
                    cx = m->wx + m->ww * 18 / 100;
                    cy = m->wy + m->wh * 15 / 100;
                    break;
                  default: // right
                    cw = m->ww * 18 / 100 - gapow - gapiw;
                    ch = m->wh * 16 / 100;
                    cx = m->wx + m->ww - cw - gapow - m->ww * 3 / 100;
                    cy = m->wy + m->wh - ch - m->wh * 3 / 100;
                    break;
                };
                resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
              }
              break;
            }
    case 4: {
              for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                switch (i) {
                  case 3: // left
                    cw = m->ww / 8 - gapow - gapiw;
                    ch = m->wh / 5;
                    cx = m->wx + gapow;
                    cy = m->wy + (m->wh - ch) / 2;
                    break;
                  case 0: // center-1 -- master
                    cw = m->ww * 4 / 8 - 2 * gapiw;
                    ch = m->wh * 4 / 5;
                    cx = m->wx + gapow + (m->ww / 8 - gapow - gapiw) + gapiw + gapiw;
                    cy = m->wy + (m->wh - ch) / 2;
                    break;
                  case 1: // center-2
                    cw = m->ww * 7 / 20 - 2 * gapiw;
                    ch = m->wh * 2 / 5;
                    cx = m->wx + gapow + (m->ww / 8 - gapow - gapiw) + gapiw + gapiw + (m->ww * 4 / 8 - 2 * gapiw) + gapiw + gapiw;
                    cy = m->wy + m->wh * 4 / 100;
                    break;
                  default: // right
                    cw = m->ww / 8 - gapow - gapiw;
                    ch = m->wh / 5;
                    cx = m->wx + m->ww - gapow - cw - m->ww / 8;
                    cy = m->wy + m->wh - gapoh - ch - m->wh * 4 / 100;
                    break;
                };
                resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
              }
              break;
            }
    case 5: {
              for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                switch (i) {
                  case 4:
                    cw = m->ww * 2 / 20 - gapow - gapiw;
                    ch = m->wh * 2 / 20;
                    cx = m->wx + gapow;
                    cy = m->wy + (m->wh - ch) / 2 + m->wh / 4;
                    break;
                  case 0: // center master
                    cw = m->ww * 8 / 20;
                    ch = m->wh * 14 / 20;
                    cx = m->wx + gapow + (m->ww * 2 / 20 - gapow - gapiw) + gapiw + gapiw;
                    cy = m->wy + (m->wh - ch) / 2;
                    break;
                  case 1: // top
                    cw = m->ww * 4 / 20;
                    ch = m->wh * 4 / 20;
                    cx = m->wx + gapow + (m->ww * 2 / 20 - gapow - gapiw) + gapiw + gapiw + (m->ww * 8 / 20) + gapiw + gapiw;
                    cy = m->wy + m->wh * 1 / 20;
                    break;
                  case 2: // bottom
                    cw = m->ww * 6 / 20;
                    ch = m->wh * 6 / 20 - gapoh - gapih;
                    cx = m->wx + gapow + (m->ww * 2 / 20 - gapow - gapiw) + gapiw + gapiw + (m->ww * 8 / 20) + gapiw + gapiw + (m->ww * 1 / 20) + gapiw + gapiw;
                    cy = m->wy + m->wh - m->wh * 1 / 20 - ch;
                    break;
                  default: // right
                    cw = m->ww * 2 / 20 - gapow - gapiw;
                    ch = m->wh * 2 / 20 - gapoh - gapih;
                    cx = m->wx + m->ww - gapow - cw - gapiw;
                    cy = m->wy + m->wh - gapoh - ch - gapih;
                    break;
                };
                resize(c, cx, cy, cw - 2 * c->bw, ch - 2 * c->bw, False);
              }
              break;
            }
    default: { // grid layout
               gapow = m->ww * 4 / 100;
               gapoh = m->wh * 16 / 100;

               for (rows = 0, cols = 0; rows * cols <= n; rows++) {
                 cols = 4 * rows;
               }

               ch = (m->wh - 2 * gapoh) / (rows ? rows : 1);
               cw = (m->ww - 2 * gapow) / (cols ? cols : 1);
               ah = rows ? (m->wh - 2 * gapoh - rows * ch) / 2 : 0;
               aw = cols ? (m->ww - 2 * gapow - cols * cw) / 2 : 0;

               for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
                 cx = m->wx + gapow + aw + (i % cols) * cw;
                 cy = m->wy + gapoh + ah + (i / cols) * ch;
                 if (i > cols * (rows - 1) - 1 && n != cols * rows) {
                   cx = m->wx + gapow + aw + (i % cols) * cw + ((cw + aw) * (cols - n % cols)) / 2;
                   cy = m->wy + gapoh + ah + (i / cols) * ch;
                 }
                 resize(
                     c,
                     cx,
                     cy,
                     cw - gapiw / 2 - 2 * c->bw,
                     ch - gapih / 2 - 2 * c->bw,
                     False
                     );
               }
               break;
             }
  }
}

  void
layout_workflow(Monitor *m)
{
  unsigned int i, n, cx, cy, cw, ch;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  switch (n) {
    case 1:
      for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        cw = m->ww - 2*c->bw;
        ch = m->wh - 2*c->bw;
        cx = m->wx;
        cy = m->wy;
        resize(c, cx, cy, cw, ch, False);
      }
      break;
    case 2:
      for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) {
          cw = m->ww * m->mfact - 2*c->bw;
          ch = m->wh - 2*c->bw;
          cx = m->wx;
          cy = m->wy + (m->wh - ch) / 2;
        } else {
          cw = m->ww * (1 - m->mfact) - 2*c->bw;
          ch = m->wh - 2*c->bw;
          cx = m->wx + m->ww - cw;
          cy = m->wy + (m->wh - ch) / 2;
        }
        resize(c, cx, cy, cw, ch, False);
      }
      break;
    case 3: // left-1 + right-top + right-bottom
      for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) { // Left
          cw = m->ww * m->mfact - 2*c->bw;
          ch = m->wh - 2*c->bw;
          cx = m->wx;
          cy = m->wy + (m->wh - ch) / 2;
        } else if (i == 1) { // RightTop
          cw = m->ww * (1 - m->mfact) - 2*c->bw;
          ch = m->wh * (1-m->hfact) - 2*c->bw;
          cx = m->wx + m->ww * m->mfact;
          cy = m->wy;
        } else { // RightBottom
          cw = m->ww * (1 - m->mfact) - 2*c->bw;
          ch = m->wh * m->hfact - 2*c->bw;
          cx = m->wx + m->ww * m->mfact;
          cy = m->wy + m->wh * (1-m->hfact);
        }
        resize(c, cx, cy, cw, ch, False);
      }
      break;
    case 4:
      for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) { // TopLeft
          cw = m->ww * m->mfact - 2*c->bw;
          ch = m->wh * (1-m->hfact) - 2*c->bw;
          cx = m->wx;
          cy = m->wy;
        } else if (i == 1) { // TopRight
          cw = m->ww * (1 - m->mfact) - 2*c->bw;
          ch = m->wh * (1-m->hfact) - 2*c->bw;
          cx = m->wx + m->ww * m->mfact;
          cy = m->wy;
        } else if (i == 3) { // BottomLeft
          cw = m->ww * m->mfact - 2*c->bw;
          ch = m->wh * m->hfact - 2*c->bw;
          cx = m->wx;
          cy = m->wy + m->wh * (1-m->hfact);
        } else {  // BottomRight
          cw = m->ww * (1 - m->mfact) - 2*c->bw;
          ch = m->wh * m->hfact - 2*c->bw;
          cx = m->wx + m->ww * m->mfact;
          cy = m->wy + m->wh * (1-m->hfact);
        }
        resize(c, cx, cy, cw, ch, False);
      }
      break;
    default:
      layout_fib_spiral(m);
      break;
  };
}
/* layout end */

  int
main(int argc, char *argv[])
{
  if (argc == 2 && !strcmp("-v", argv[1])) {
    die("dwm-"VERSION);
  }

  if (argc != 1) {
    die("usage: dwm [-v]");
  }

  if (!setlocale(LC_CTYPE, "") || !XSupportsLocale()) {
    fputs("warning: no locale support\n", stderr);
  }

  if (!(dpy = XOpenDisplay(NULL))) {
    die("dwm: cannot open display");
  }

  if (!(xcon = XGetXCBConnection(dpy))) {
    die("dwm: cannot get xcb connection\n");
  }

  checkotherwm();
  autostart_exec();
  setup();

#ifdef __OpenBSD__
  if (pledge("stdio rpath proc exec ps", NULL) == -1) {
    die("pledge");
  }
#endif

  scan();
  restoresession();
  run();

  if(restart) { execvp(argv[0], argv); }

  cleanup();
  XCloseDisplay(dpy);
  return EXIT_SUCCESS;
}
