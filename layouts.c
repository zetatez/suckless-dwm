// layouts
#include<math.h>

/* dwm-overlaylayer ------------------------------------------------------------ */
void
overlaylayergrid(Monitor *m) {
    unsigned int n, i;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++) ;
    if(n == 0)
        return;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) {
            resize(c, m->wx, m->wy, m->ww - 2 * c->bw, m->wh - 2 * c->bw, False);
        } else if (i < 1 + (n-1)%3) {
            resize(c, m->wx + ((i-1)%3)*m->ww/((n-1)%3), m->wy + m->wh*(1-m->freeh) + ((n-i-1-(n-i-1)%3)/3)*m->wh*m->freeh/((n-1-(n-1-1)%3)/3+1), m->ww/((n-1)%3) - 2*c->bw, m->wh*m->freeh/((n-1-1)/3+1) - 2*c->bw, False);
        } else {
            resize(c, m->wx + ((i-1)%3)*m->ww/3 , m->wy + m->wh*(1-m->freeh) + ((n-i-1-(n-i-1)%3)/3)*m->wh*m->freeh/((n-1-(n-1-1)%3)/3+1), m->ww/3 - 2*c->bw, m->wh*m->freeh/((n-1-1)/3+1) - 2*c->bw, False);
        }
    }
}

void
overlaylayerhorizontal(Monitor *m) {
    unsigned int n, i;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++) ;
    if(n == 0)
        return;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) {
            resize(c, m->wx, m->wy, m->ww - 2 * c->bw, m->wh - 2 * c->bw, False);
        } else {
            resize(c, m->wx, m->wy + m->wh * (1 - m->freeh) + (n - i - 1) * m->wh * m->freeh / (n - 1), m->ww - 2 * c->bw, m->wh * m->freeh / (n - 1) - 2 * c->bw, False);
        }
    }
}

void
overlaylayervertical(Monitor *m) {
    unsigned int n, i;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++) ;
    if(n == 0)
        return;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        if (i == 0) {
            resize(c, m->wx, m->wy, m->ww - 2 * c->bw, m->wh - 2 * c->bw, False);
        } else {
            resize(c, m->wx + (n - i - 1) * m->ww  / (n - 1), m->wy + m->wh * (1 - m->freeh), m->ww / (n - 1) - 2 * c->bw, m->wh * m->freeh - 2 * c->bw, False);
        }
    }
}

/* dwm-center ------------------------------------------------------------ */
void
centerequalratio(Monitor *m) {
	unsigned int n, i;
	Client *c;

	for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++) ;
	if(n == 0)
		return;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        resize(c, m->ww/2 - (m->ww * m->freeh)/2, m->wy + m->wh/2 - (m->wh * m->freeh)/2 , m->ww * m->freeh - 2 * c->bw, m->wh * m->freeh - 2 * c->bw, False);
	}
}

void
centeranyshape(Monitor *m) {
	unsigned int n, i;
	Client *c;

	for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++) ;
	if(n == 0)
		return;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
        resize(c, m->ww/2 - (m->ww * m->mfact)/2, m->wy + m->wh/2 - (m->wh * m->freeh)/2 , m->ww * m->mfact - 2 * c->bw, m->wh * m->freeh - 2 * c->bw, False);
	}
}

/* dwm-fibonacci ------------------------------------------------------------ */
void
fibonacci(Monitor *m, int s) {
	unsigned int i, n, nx, ny, nw, nh;
	Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
    if(n == 0)
        return;

    nx = m->wx;
    ny = 0;
    nw = m->ww;
    nh = m->wh;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
        if((i % 2 && nh / 2 > 2 * c->bw)
                || (!(i % 2) && nw / 2 > 2 * c->bw)) {
            if(i < n - 1) {
                if(i % 2)
                    nh /= 2;
                else
                    nw /= 2;
                if((i % 4) == 2 && !s)
                    nx += nw;
                else if((i % 4) == 3 && !s)
                    ny += nh;
            }
            if((i % 4) == 0) {
                if(s)
                    ny += nh;
                else
                    ny -= nh;
            }
            else if((i % 4) == 1)
                nx += nw;
            else if((i % 4) == 2)
                ny += nh;
            else if((i % 4) == 3) {
                if(s)
                    nx += nw;
                else
                    nx -= nw;
            }
            if(i == 0)
            {
                if(n != 1)
                    nw = m->ww * m->mfact;
                ny = m->wy;
            }
            else if(i == 1)
                nw = m->ww - nw;
            i++;
        }
        resize(c, nx, ny, nw - 2 * c->bw, nh - 2 * c->bw, False);
    }
}

void
dwindle(Monitor *m) {
	fibonacci(m, 1);
}

void
spiral(Monitor *m) {
	fibonacci(m, 0);
}

/* dwm-gridmode ------------------------------------------------------------ */
void
grid(Monitor *m) {
    unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next))
        n++;

    /* grid dimensions */
    for(rows = 0; rows <= n/2; rows++)
        if(rows*rows >= n)
            break;
    cols = (rows && (rows - 1) * rows >= n) ? rows - 1 : rows;

    /* window geoms (cell height/width) */
    ch = m->wh / (rows ? rows : 1);
    cw = m->ww / (cols ? cols : 1);
    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
        cx = m->wx + (i / rows) * cw;
        cy = m->wy + (i % rows) * ch;
        /* adjust height/width of last row/column's windows */
        ah = ((i + 1) % rows == 0) ? m->wh - ch * rows : 0;
        aw = (i >= rows * (cols - 1)) ? m->ww - cw * cols : 0;
        resize(c, cx, cy, cw - 2 * c->bw + aw, ch - 2 * c->bw + ah, False);
        i++;
    }
}

/* dwm-lefttile ------------------------------------------------------------ */
void
tileleft(Monitor *m)
{
	unsigned int i, n, h, mw, my, ty;
	Client *c;

	for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
	if (n == 0)
		return;

	if (n > m->nmaster)
		mw = m->nmaster ? m->ww * (1 - m->mfact) : 0;
	else
		mw = m->ww;
	for (i = my = ty = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++)
		if (i < m->nmaster) {
			h = (m->wh - my) / (MIN(n, m->nmaster) - i);
			resize(c, m->wx + m->ww - mw, m->wy + my, mw - (2*c->bw), h - (2*c->bw), 0);
			if (my + HEIGHT(c) < m->wh)
				my += HEIGHT(c);
		} else {
			h = (m->wh - ty) / (n - i);
			resize(c, m->wx, m->wy + ty, m->ww - mw - (2*c->bw), h - (2*c->bw), 0);
			if (ty + HEIGHT(c) < m->wh)
				ty += HEIGHT(c);
		}
}

/* dwm-deck ------------------------------------------------------------ */
void
deckvertical(Monitor *m) {
    unsigned int i, n, mw;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
    if(n == 0)
        return;

    if(n > m->nmaster)
        mw = m->nmaster ? m->ww * m->mfact : 0;
    else
        mw = m->ww;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++)
        if(i < m->nmaster)
            resize(c, m->wx, m->wy, mw - 2*c->bw, m->wh - 2*c->bw, c->bw);
        else
            resize(c, m->wx + mw + (i - m->nmaster) * (m->ww - mw) / (n - m->nmaster), m->wy, m->ww - (mw + (i - m->nmaster) * (m->ww - mw) / (n - m->nmaster)) - 2 * c->bw, m->wh - 2 * c->bw, c->bw);
}

void
deckhorizontal(Monitor *m) {
    unsigned int i, n, mh;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
    if(n == 0)
        return;

    if(n > m->nmaster)
        mh = m->nmaster ? m->wh * (1-m->freeh) : 0;
    else
        mh = m->wh;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++)
        if(i < m->nmaster)
            resize(c, m->wx, m->wy, m->ww - 2*c->bw, mh - 2*c->bw, c->bw);
        else
            resize(c, m->wx, m->wy + mh + (i - m->nmaster) * (m->wh - mh) / (n - m->nmaster), m->ww - 2 * c->bw, m->wh - (mh +  (i - m->nmaster) * (m->wh - mh) / (n - m->nmaster)) - 2 * c->bw, c->bw);
}

/* dwm-bottomstack ------------------------------------------------------------ */
static void
bottomstackhorizontal(Monitor *m) {
	int w, mh, mx, tx, ty, th;
	unsigned int i, n;
	Client *c;

	for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
	if (n == 0)
		return;

	if (n > m->nmaster) {
		mh = m->nmaster ? (1 - m->freeh) * m->wh : 0;
		th = (m->wh - mh) / (n - m->nmaster);
		ty = m->wy + mh;
	} else {
		th = mh = m->wh;
		ty = m->wy;
	}
	for (i = mx = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
		if (i < m->nmaster) {
			w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
			resize(c, m->wx + mx, m->wy, w - (2 * c->bw), mh - (2 * c->bw), 0);
			mx += WIDTH(c);
		} else {
			resize(c, tx, ty, m->ww - (2 * c->bw), th - (2 * c->bw), 0);
			if (th != m->wh)
				ty += HEIGHT(c);
		}
	}
}

static void
bottomstackvertical(Monitor *m) {
	int w, h, mh, mx, tx, ty, tw;
	unsigned int i, n;
	Client *c;

	for (n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
	if (n == 0)
		return;

	if (n > m->nmaster) {
		mh = m->nmaster ? (1 - m->freeh) * m->wh : 0;
		tw = m->ww / (n - m->nmaster);
		ty = m->wy + mh;
	} else {
		mh = m->wh;
		tw = m->ww;
		ty = m->wy;
	}
	for (i = mx = 0, tx = m->wx, c = nexttiled(m->clients); c; c = nexttiled(c->next), i++) {
		if (i < m->nmaster) {
			w = (m->ww - mx) / (MIN(n, m->nmaster) - i);
			resize(c, m->wx + mx, m->wy, w - (2 * c->bw), mh - (2 * c->bw), 0);
			mx += WIDTH(c);
		} else {
			h = m->wh - mh;
			resize(c, tx, ty, tw - (2 * c->bw), h - (2 * c->bw), 0);
			if (tw != m->ww)
				tx += WIDTH(c);
		}
	}
}

/* dwm-logarithmic-spiral ------------------------------------------------------------ */
// control the shape of logarithmic spiral
static const float logarithmicspiralstart = -50;
static const float logarithmicspiralstop  = 50;
static const float logarithmicspiralstep  = 0.1;    // control the interval of each window
static const float logarithmicspiralalpha = 1;
static const float logarithmicspiralkapa  = 0.2;    // control the interval of each window cycle: 0.2, 0.025, 0.05, 0.3063489(golden LS)
static const int   logarithmicspirallen   = (const int) ((logarithmicspiralstop - logarithmicspiralstart )/logarithmicspiralstep);

void
logarithmicspiral(Monitor *m) {
	unsigned int n, idx;
    float i, v, minx, maxx, miny, maxy;
    float phi[logarithmicspirallen];

    float x[logarithmicspirallen];
    float y[logarithmicspirallen];

    float ww[logarithmicspirallen];
    float wh[logarithmicspirallen];
    float wx[logarithmicspirallen];
    float wy[logarithmicspirallen];

	Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next), n++);
    if(n == 0)
        return;

    for (idx = 0, i = logarithmicspiralstart; i < logarithmicspiralstop && idx < sizeof(phi) / sizeof(phi[0]); i += logarithmicspiralstep, phi[idx] = i, idx++);
    for (idx = 0; idx < sizeof(phi) / sizeof(phi[0]); idx++) {
        v = logarithmicspiralalpha * exp(logarithmicspiralkapa * phi[idx]);
        x[idx] = v * cos(phi[idx]);
        y[idx] = v * sin(phi[idx]);
    }

    // min max
    minx = maxx = x[0];
    miny = maxy = y[0];

    for (idx = 1; idx < sizeof(phi) / sizeof(phi[0]); idx++) {
        if (x[idx] < minx) { minx = x[idx]; }
        if (x[idx] > maxx) { maxx = x[idx]; }
        if (y[idx] < miny) { miny = y[idx]; }
        if (y[idx] > maxy) { maxy = y[idx]; }
    }

    for (idx = 0; idx < sizeof(phi) / sizeof(phi[0]); idx++) {
        // min max normal
        x[idx] = (x[idx] - minx)/(maxx-minx);
        y[idx] = (y[idx] - miny)/(maxy-miny);

        // allocate window size
        ww[idx] = 96;
        wh[idx] = 32;
        wx[idx] = (m->ww - 2*ww[idx]/2) * x[idx] - ww[idx]/2;
        wy[idx] = (m->wh - 2*wh[idx]/2) * y[idx] - wh[idx]/2;
    }

    // last -5 window center
    idx = logarithmicspirallen-5;
    ww[idx] = 320;
    wh[idx] = 120;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -10 window center
    idx = logarithmicspirallen-10;
    ww[idx] = 360;
    wh[idx] = 120;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -16 window center
    idx = logarithmicspirallen-16;
    ww[idx] = 640;
    wh[idx] = 240;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -23 window center
    idx = logarithmicspirallen-23;
    ww[idx] = 240;
    wh[idx] = 80;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -31 window center
    idx = logarithmicspirallen-31;
    ww[idx] = 320;
    wh[idx] = 120;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -36 window center
    idx = logarithmicspirallen-36;
    ww[idx] = 240;
    wh[idx] = 80;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

    // last -48 window center
    idx = logarithmicspirallen-48;
    ww[idx] = 320;
    wh[idx] = 120;
    wx[idx] = wx[idx] - ww[idx]/2;
    wy[idx] = wy[idx];

	for(i = 0, c = nexttiled(m->clients); c && i < logarithmicspirallen; c = nexttiled(c->next), i++) {
        if (i < 1) {
            resize(c, m->ww/2 - (m->ww * m->mfact)/2, m->wy + m->wh/2 - (m->wh * m->freeh)/2 , m->ww * m->mfact - 2 * c->bw, m->wh * m->freeh - 2 * c->bw, False);
        } else {
            idx = logarithmicspirallen - 1 - i;
            wx[idx] = wx[idx] < 0 ? 0: wx[idx];
            wy[idx] = wy[idx] < 0 ? 0: wy[idx];
            wx[idx] = wx[idx] + ww[idx] > m->ww ? m->ww - ww[idx]: wx[idx];
            wy[idx] = wy[idx] + wh[idx] > m->wh ? m->wh - wh[idx]: wy[idx];
            resize(c, m->wx + wx[idx], m->wy + wy[idx], ww[idx] - 2 * c->bw, wh[idx] - 2 * c->bw, False);
        }
	}
}

/* dwm-overview ------------------------------------------------------------ */
void
overview(Monitor *m) {
    unsigned int i, n, cx, cy, cw, ch, aw, ah, cols, rows;
    Client *c;

    for(n = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next))
        n++;

    /* grid dimensions */
    for(rows = 0; rows <= n/2; rows++)
        if(rows*rows >= n)
            break;
    cols = (rows && (rows - 1) * rows >= n) ? rows - 1 : rows;

    /* window geoms (cell height/width) */
    ch = (m->wh - 2 * gappoh) / (rows ? rows : 1);
    cw = (m->ww - 2 * gappow) / (cols ? cols : 1);

    /* round err adjust cx/cy */
    ah = rows ? (m->wh - 2 * gappoh - rows * ch)/2: 0;
    aw = cols ? (m->ww - 2 * gappow - cols * cw)/2: 0;

    for(i = 0, c = nexttiled(m->clients); c; c = nexttiled(c->next)) {
        cx = m->wx + gappow + aw + (i / rows) * cw;
        cy = m->wy + gappoh + ah + (i % rows) * ch;
        resize(c, cx, cy, cw - gappiw/2 - 2 * c->bw, ch - gappih/2 - 2 * c->bw, False);
        i++;
    }
}
