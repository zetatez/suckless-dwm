#include "layouts.h"

void
centerequalratio(Monitor *m)
{
  unsigned int n, i;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > 0) { snprintf(m->ltsymbol, sizeof m->ltsymbol, "%s %d", selmon->lt[selmon->sellt]->symbol, n); }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    resize(c, m->ww/2 - (m->ww*m->ffact)/2, m->wy + m->wh/2 - (m->wh*m->ffact)/2 + (topbar ? 1 : 0)*winpad, m->ww*m->ffact - 2*c->bw, (m->wh - (topbar ? 1 : 0)*winpad) * m->ffact - 2*c->bw, False);
  }
}

void
centeranyshape(Monitor *m)
{
  unsigned int n, i;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > 0) { snprintf(m->ltsymbol, sizeof m->ltsymbol, "%s %d", selmon->lt[selmon->sellt]->symbol, n); }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    resize(c, m->ww / 2 - (m->ww * m->mfact) / 2, m->wy + m->wh / 2 - (m->wh * m->ffact) / 2, m->ww * m->mfact - 2 * c->bw, m->wh * m->ffact - 2 * c->bw, False);
  }
}

void
fibonacci(Monitor *m, int s)
{
  unsigned int i, n, nx, ny, nw, nh;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  nx = m->wx;
  ny = 0;
  nw = m->ww;
  nh = m->wh - (topbar ? 1 : 0)*winpad;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
    if ((i % 2 && nh / 2 > 2 * c->bw) || (!(i % 2) && nw / 2 > 2 * c->bw)) {
      if (i < n - 1) {
        if (i % 2)
          nh /= 2;
        else
          nw /= 2;
        if ((i % 4) == 2 && !s)
          nx += nw;
        else if ((i % 4) == 3 && !s)
          ny += nh;
      }

      if ((i % 4) == 0) {
        if (s)
          ny += nh;
        else
          ny -= nh;
      } else if ((i % 4) == 1) {
        nx += nw;
      } else if ((i % 4) == 2) {
        ny += nh;
      } else if ((i % 4) == 3) {
        if (s)
          nx += nw;
        else
          nx -= nw;
      }

      if (i == 0) {
        if (n != 1)
          nw = m->ww * m->mfact;
        ny = m->wy;
      } else if (i == 1) {
        nw = m->ww - nw;
      }

      i++;
    }

    resize(c, nx, ny + (topbar ? 1 : 0)*winpad, nw - 2 * c->bw, nh - 2 * c->bw, False);
  }
}

void
fibonaccidwindle(Monitor *m)
{
  fibonacci(m, 1);
}

void
fibonaccispiral(Monitor *m)
{
  fibonacci(m, 0);
}

void
grid(Monitor *m)
{
  unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  for (cols = 0; cols <= n / 2; cols++) {
    if (cols * cols >= n) {
      break;
    }
  }

  rows = (cols && (cols - 1) * cols >= n) ? cols - 1 : cols;
  ch = (m->wh - (topbar ? 1 : 0)*winpad) / (rows ? rows : 1);
  cw = m->ww / (cols ? cols : 1);
  ah = rows ? (m->wh - (topbar ? 1 : 0)*winpad - rows * ch) / 2 : 0;
  aw = cols ? (m->ww - cols * cw) / 2 : 0;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
    cx = m->wx + aw + (i % cols) * cw;
    cy = m->wy + ah + (i / cols) * ch;
    if (i > cols * (rows - 1) - 1 && n != cols * rows) {
      cx = m->wx + aw + (i % cols) * cw + ((cw + aw) * (cols - n % cols))/2;
      cy = m->wy + ah + (i / cols) * ch;
    }
    resize(c, cx, cy + (topbar ? 1 : 0)*winpad, cw - 2 * c->bw, ch - 2 * c->bw, False);
    i++;
  }
}

void
tileright(Monitor *m)
{
  unsigned int i, n, h, mw, my, ty;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mw = m->nmaster ? m->ww * m->mfact : 0;
  } else {
    mw = m->ww;
  }

  for (i = my = ty = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      h = (m->wh - my - (topbar ? 1 : 0)*winpad) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx, m->wy + my + (topbar ? 1 : 0)*winpad, mw - (2*c->bw), h - (2*c->bw), 0);
      if (my + HEIGHT(c) < m->wh) {
        my += HEIGHT(c);
      }
    } else {
      h = (m->wh - ty - (topbar ? 1 : 0)*winpad) / (n - i);
      resize(c, m->wx + mw, m->wy + ty + (topbar ? 1 : 0)*winpad, m->ww - mw - (2*c->bw), h - (2*c->bw), 0);
      if (ty + HEIGHT(c) < m->wh) {
        ty += HEIGHT(c);
      }
    }
  }
}

void
tileleft(Monitor *m) {
  unsigned int i, n, h, mw, my, ty;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mw = m->nmaster ? m->ww*(1-m->mfact) : 0;
  } else {
    mw = m->ww;
  }

  for (i = my = ty = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      h = (m->wh - my - (topbar ? 1 : 0)*winpad) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + m->ww - mw, m->wy + my + (topbar ? 1 : 0)*winpad, mw - (2 * c->bw), h - (2 * c->bw), 0);
      if (my + HEIGHT(c) < m->wh) {
        my += HEIGHT(c);
      }
    } else {
      h = (m->wh - ty - (topbar ? 1 : 0) * winpad) / (n - i);
      resize(c, m->wx, m->wy + ty + (topbar ? 1 : 0) * winpad, m->ww - mw - (2 * c->bw), h - (2 * c->bw), 0);
      if (ty + HEIGHT(c) < m->wh) {
        ty += HEIGHT(c);
      }
    }
  }
}

void
deckvert(Monitor *m) {
  unsigned int i, n, mw;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mw = m->nmaster ? m->ww*m->mfact : 0;
  } else {
    mw = m->ww;
  }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster)
      resize(c, m->wx, m->wy + (topbar ? 1 : 0) * winpad, mw - 2 * c->bw, m->wh - 2*c->bw - (topbar ? 1 : 0) * winpad, c->bw);
    else
      resize(c, m->wx + mw + (i - m->nmaster)*(m->ww - mw)/(n - m->nmaster), m->wy + (topbar ? 1 : 0)*winpad, m->ww - (mw + (i - m->nmaster) * (m->ww - mw)/(n - m->nmaster)) - 2*c->bw, m->wh - 2*c->bw - (topbar ? 1 : 0)*winpad, c->bw);
  }
}

void
deckhori(Monitor *m)
{
  unsigned int i, n, mh;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mh = m->nmaster ? m->wh * (1 - m->ffact) : 0;
  } else {
    mh = m->wh;
  }

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      resize(c, m->wx, m->wy + (topbar ? 1 : 0)*winpad, m->ww - 2 * c->bw, mh - 2*c->bw - (topbar ? 1 : 0)*winpad, c->bw);
    } else {
      resize(c, m->wx, m->wy + mh + (i - m->nmaster)*(m->wh - mh)/(n - m->nmaster), m->ww - 2*c->bw, m->wh - (mh + (i - m->nmaster)*(m->wh - (topbar ? 1 : 0)*winpad - mh)/(n - m->nmaster)) - 2*c->bw - (topbar ? 1 : 0)*winpad, c->bw);
    }
  }
}

void
bottomstackhori(Monitor *m) {
  int w, mh, mx, tx, ty, th;
  unsigned int i, n;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mh = m->nmaster ? (1 - m->ffact) * m->wh : 0;
    th = (m->wh - mh - (topbar ? 1 : 0)*winpad) / (n - m->nmaster);
    ty = m->wy + mh;
  } else {
    th = mh = m->wh - (topbar ? 1 : 0)*winpad;
    ty = m->wy;
  }

  for (i = mx = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + mx, m->wy + (topbar ? 1 : 0)*winpad, w - (2*c->bw), mh - (2*c->bw), 0);
      mx += WIDTH(c);
    } else {
      resize(c, tx, ty + (topbar ? 1 : 0)*winpad, m->ww - (2*c->bw), th - (2*c->bw), 0);
      if (th != m->wh - (topbar ? 1 : 0)*winpad) {
        ty += HEIGHT(c);
      }
    }
  }
}

void
bottomstackvert(Monitor *m)
{
  int w, h, mh, mx, tx, ty, tw;
  unsigned int i, n;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  if (n > m->nmaster) {
    mh = m->nmaster ? (1 - m->ffact) * (m->wh - (topbar ? 1 : 0)*winpad) : 0;
    tw = m->ww / (n - m->nmaster);
    ty = m->wy + mh;
  } else {
    mh = m->wh - (topbar ? 1 : 0)*winpad;
    tw = m->ww;
    ty = m->wy;
  }

  for (i = mx = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    if (i < m->nmaster) {
      w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
      resize(c, m->wx + mx, m->wy + (topbar ? 1 : 0)*winpad, w - (2*c->bw), mh - (2*c->bw), 0);
      mx += WIDTH(c);
    } else {
      h = m->wh - (topbar ? 1 : 0)*winpad - mh;
      resize(c, tx, ty + (topbar ? 1 : 0)*winpad, tw - (2*c->bw), h - (2*c->bw), 0);
      if (tw != m->ww) {
        tx += WIDTH(c);
      }
    }
  }
}

void
overview(Monitor *m)
{
  unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

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
    resize(c, cx, cy, cw - gapiw / 2 - 2 * c->bw, ch - gapih / 2 - 2 * c->bw, False);
  }
}

void
hacker(Monitor *m)
{
  unsigned int i, n, cx, cy, cw, ch;
  Client *c;

  for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++)
    ;

  if (n == 0) { return; }

  cw = (m->ww - 2*gapow)*3/5;
  ch = (m->wh - 2*gapoh)*3/5;

  for (i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
    cx = m->wx + gapow + (n-i-1)*(m->ww/34);
    cy = m->wy + gapoh + (n-i-1)*(m->wh/34);
    if (cy + ch - 2*c->bw > m->wh) {
      cx = (m->ww - 2*gapow)/2 - cw/2;
      cy = (m->wh - 2*gapoh)/2 - ch/2;
    }
    resize(c, cx, cy, cw - 2*c->bw, ch - 2*c->bw, False);
  }
}
