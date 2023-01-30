#!/usr/bin/python
# -*- coding: utf-8 -*-

import os, sys, signal
import re
import time
import wget
import shutil, psutil
from PyQt5 import QtWidgets

my_home_path = "/home/dionysus"
my_download_path = "/home/dionysus/Downloads"
my_library_path = "/home/dionysus/my-library"
my_play_path = "/home/dionysus/my-play"
my_trojan_path = "/home/dionysus/.trojan"
my_wallpaper_path = "/home/dionysus/Pictures/wallpapers"
my_default_wallpaper = "013.jpeg"

win_name_float = "00001011"
win_name_scratchpad = "scratchpad"


def read_file(filename):
    with open(filename, "r", encoding="utf-8") as fh:
        return fh.read()


def write_file(filename, s):
    with open(filename, "w", encoding="utf-8") as fh:
        fh.write(s)


def get_screen_count():
    _ = QtWidgets.QApplication(sys.argv)
    n = QtWidgets.QDesktopWidget().screenCount()
    return n


def get_cur_screen_geometry():
    _ = QtWidgets.QApplication(sys.argv)
    g = QtWidgets.QDesktopWidget().screenGeometry(-1)
    w, h = g.width(), g.height()
    return w, h


def get_xy(xr, yr):
    w, h = get_cur_screen_geometry()
    return int(w * xr), int(h * yr)


def get_geometry_for_st(xr, yr, w, h):
    x, y = get_xy(xr, yr)
    return "{}x{}+{}+{}".format(w, h, x, y)


def get_pids(pname):
    pids = []
    for p in psutil.process_iter():
        if p.name() == pname:
            pids.append(p.pid)
    return pids


def in_running(cmd):
    cmd_ps = "ps -ef|grep '{}'".format(
        cmd.rstrip(" &")) + "|grep -v grep|awk '{print $2}'"
    pids = popen(cmd_ps).strip().replace("\n", " ").strip()
    if pids:
        return pids, True

    return pids, False


def toggle_by_pname(pname, cmd):
    pids = get_pids(pname)
    if pids:
        [os.kill(pid, signal.SIGKILL) for pid in pids]
    else:
        os.system(cmd)


def toggle_by_cmd(cmd):
    pids, ok = in_running(cmd.replace("'", "").replace('"', ""))
    if not ok:
        os.system(cmd)
    else:
        cmd = "kill -9 {}".format(pids)
        os.system(cmd)
    return


def spawn(cmd):
    os.execvp(cmd[0], cmd)
    return


def popen(cmd):
    r = os.popen(cmd)
    text = r.read()
    r.close()
    return text


def open_file_at_foreground(file_path):
    os.system("st -e lazy -o " + file_path)
    return


def open_file_at_background(file_path):
    os.system("st -e lazy -o " + file_path + " &")
    return


def keep_file_or_not(file_path):
    cmd = "echo 'yes\nno'|dmenu -p 'keep file?'"
    option = popen(cmd)
    if option.strip() == "no":
        os.remove(file_path)
    return


# utils
# ---------------------------------------------------------
def open_copied():
    cmd = "xclip -o"
    last_copied_str = popen(cmd)
    last_copied_str = last_copied_str.strip()

    # if a local file
    if os.path.exists(last_copied_str):
        cmd = "st -e lazy -o {} &".format(last_copied_str)
        os.system(cmd)
        return

    # if a url of: ^(http|https|www|file).+
    if re.match(r'^(http|https|www|file).+', last_copied_str):
        cmd = "st -e google {} &".format(last_copied_str)
        os.system(cmd)
        return

    return


def open_note_taking():
    time_str = time.strftime("%Y-%m-%d", time.localtime())
    file_folder = my_library_path + "/xournal"
    file_path = file_folder + "/" + time_str + ".xopp"

    if not os.path.exists(file_folder):
        os.mkdir(file_folder)

    cmd = "xournalpp {} &".format(file_path)

    _, ok = in_running(cmd)
    if not ok:
        os.system(cmd)
    else:
        os.system("notify-send '{}'".format("already in running"))

    return


def open_my_play():
    tody_str = time.strftime("%Y-%m-%d", time.localtime())
    my_play_notes_today_path = my_play_path + "/notes/" + tody_str
    template_note_tex = my_play_path + "/templates/note.tex"
    template_note_xoj = my_play_path + "/templates/note.xopp"
    my_play_notes_today_note_tex = my_play_path + "/notes/" + tody_str + "/note.tex"
    my_play_notes_today_note_xoj = my_play_path + "/notes/" + tody_str + "/note.xopp"

    if not os.path.exists(my_play_notes_today_path):
        os.mkdir(my_play_notes_today_path)

    if not os.path.exists(template_note_tex):
        msg = "template: note.tex not found"
        os.system("notify-send '{}'".format(msg))
        exit(-1)

    if not os.path.exists(template_note_xoj):
        msg = "template: note.xoj not found"
        os.system("notify-send '{}'".format(msg))
        exit(-1)

    if not os.path.exists(my_play_notes_today_note_tex):
        shutil.copyfile(template_note_tex, my_play_notes_today_note_tex)

    if not os.path.exists(my_play_notes_today_note_xoj):
        shutil.copyfile(template_note_xoj, my_play_notes_today_note_xoj)

    cmd_tex = "st -g {} -t {} -c {} -e nvim {} &".format(
        get_geometry_for_st(0.01, 0.03, 88, 54), win_name_float,
        win_name_float, my_play_notes_today_note_tex)

    cmd_xoj = "st -g {} -t {} -c {} -e xournalpp {} &".format(
        get_geometry_for_st(0.52, 0.16, 88, 38), win_name_float,
        win_name_float, my_play_notes_today_note_xoj)

    _, ok = in_running(cmd_tex)
    if not ok:
        os.system(cmd_tex)
    else:
        os.system("notify-send '{}'".format("already in running"))

    _, ok = in_running(cmd_xoj)
    if not ok:
        os.system(cmd_xoj)
    else:
        os.system("notify-send '{}'".format("already in running"))

    return


def open_passmenu():
    cmd = "passmenu"
    os.system(cmd)

    return


def open_photoshop():
    pname = "gimp"
    cmd = "gimp"
    toggle_by_pname(pname, cmd)

    return


def open_sketchpad():
    time_str = time.strftime("%Y-%m-%d-%H-%M", time.localtime())
    file_folder = my_library_path + "/inkscape"
    template_folder = file_folder + "/templates"
    template_file = template_folder + "/drawing.svg"

    file_path = file_folder + "/" + time_str + ".svg"

    if not os.path.exists(file_folder):
        os.mkdir(file_folder)

    if not os.path.exists(template_file):
        os.mkdir(template_file)
        msg = "template: drawing.svg not found"
        os.system("notify-send '{}'".format(msg))
        exit(-1)

    if not os.path.exists(file_path):
        shutil.copyfile(template_file, file_path)

    cmd = "inkscape {}".format(file_path)

    _, ok = in_running(cmd)
    if not ok:
        os.system(cmd)
    else:
        os.system("notify-send '{}'".format("already in running"))

    return


def open_sublime():
    cmd = "subl"
    os.system(cmd)

    return


def open_wps():
    cmd = "wps"
    os.system(cmd)

    return


def toggle_addressbook():
    cmd = "st -e abook"
    toggle_by_cmd(cmd)

    return


def toggle_calendar():
    cmd = "st -g {} -t {} -c {} -e nvim +':set laststatus=0' +'Calendar -view=day'".format(
        get_geometry_for_st(0.74, 0.08, 40, 24), win_name_float,
        win_name_float)
    toggle_by_cmd(cmd)

    return


def toggle_diary():
    time_str = time.strftime("%Y-%m-%d", time.localtime())
    cmd = "st -e nvim {}/diary/{}.md".format(my_home_path, time_str)
    toggle_by_cmd(cmd)

    return


def toggle_top():
    cmd = "st -e htop"
    toggle_by_cmd(cmd)

    return


def toggle_trojan():
    cmd = "{}/trojan -c {}/config.json".format(my_trojan_path, my_trojan_path)
    toggle_by_cmd(cmd)

    return


def toggle_flameshot():
    cmd = "flameshot gui"
    toggle_by_cmd(cmd)

    return


def toggle_google():
    cmd = "google"
    toggle_by_cmd(cmd)

    return


def toggle_google_with_proxy():
    cmd = "google --proxy-server='socks5://127.0.0.1:8000'"
    toggle_by_cmd(cmd)

    return


def toggle_gitter():
    cmd = "gitter"
    toggle_by_cmd(cmd)

    return


def toggle_irc():
    cmd = "st -e irssi"
    toggle_by_cmd(cmd)

    return


def toggle_julia():
    cmd = "st -t {} -c {} -e julia".format(win_name_scratchpad,
                                           win_name_scratchpad)
    toggle_by_cmd(cmd)

    return


def toggle_lazydocker():
    cmd = "st -e lazydocker"
    toggle_by_cmd(cmd)

    return


def toggle_mathpix():
    cmd = "mathpix"
    toggle_by_cmd(cmd)

    return


def toggle_music():
    cmd_cava = "st -g {} -t cava -c cava -e cava &".format(
        get_geometry_for_st(0.74, 0.08, 40, 12))
    cmd_ncmpcpp = "st -g {} -t music -c music -e ncmpcpp &".format(
        get_geometry_for_st(0.52, 0.08, 40, 12))

    toggle_by_cmd(cmd_cava)
    toggle_by_cmd(cmd_ncmpcpp)

    return


def toggle_music_net_cloud():
    cmd = "netease-cloud-music"
    toggle_by_cmd(cmd)

    return


def toggle_mutt():
    cmd = "st -e mutt"
    toggle_by_cmd(cmd)

    return


def toggle_rss():
    cmd = "st -e newsboat"
    toggle_by_cmd(cmd)

    return


def toggle_screenkey():
    cmd = "screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow"
    toggle_by_cmd(cmd)

    return


def toggle_show():
    cmd = "st -g {} -t {} -c {} -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0".format(
        get_geometry_for_st(0.74, 0.08, 40, 12), win_name_float,
        win_name_float)
    toggle_by_cmd(cmd)

    return


def toggle_vifm():
    cmd = "st -e vifm"
    toggle_by_cmd(cmd)

    return


def toggle_wechat():
    cmd = "st -e wechat-uos"
    toggle_by_cmd(cmd)

    return


def toggle_wifi():
    cmd = "st -e nmtui"
    toggle_by_cmd(cmd)

    return


def toggle_wallpaper():
    cmd = "feh --bg-fill --recursive --randomize {}".format(my_wallpaper_path)
    os.system(cmd)

    return


def toggle_rec_audio():
    time_str = time.strftime("%Y-%m-%d-%H-%M-%S", time.localtime())
    cmd = "st  -t {} -c {} -e ffmpeg -y -r 60 -f alsa -i default -c:a flac {}/Videos/rec-a-{}.flac".format(
        win_name_scratchpad, win_name_scratchpad, my_home_path, time_str)
    toggle_by_cmd(cmd)

    return


def toggle_rec_video():
    time_str = time.strftime("%Y-%m-%d-%H-%M-%S", time.localtime())
    w, h = get_cur_screen_geometry()
    dpy = os.environ.get("DISPLAY")
    cmd = "st  -t {} -c {} -e ffmpeg -y -s '{}x{}' -r 60 -f x11grab -i {} -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac {}/Videos/rec-v-a-{}.mkv".format(
        win_name_scratchpad, win_name_scratchpad, w, h, dpy, my_home_path,
        time_str)
    toggle_by_cmd(cmd)

    return


def toggle_screen():
    primary_screen = "eDP-1"
    cmd = "xrandr|grep ' connected'|grep -v 'eDP-1'|awk '{print $1}'"
    second_screen = popen(cmd).strip()

    if not second_screen:
        msg = "have no second screen"
        os.system("notify-send '{}'".format(msg))
        return

    cmd = "echo 'only\nleft-of\nright-of\nabove\nbelow\nrotate left\nrotate right'|dmenu -p 'please select screen position'"
    option = popen(cmd).strip()
    if not option:
        return

    msg = "setting monitor: {}".format(option)
    os.system("notify-send '{}'".format(msg))

    cmd = "xrandr --output {} --auto --output {} --off".format(
        second_screen, primary_screen)
    if option == "only":
        cmd = "xrandr --output {} --auto --output {} --off".format(
            second_screen, primary_screen)
    elif option == "left-of":
        cmd = "xrandr --output {} --auto --{} {} --auto".format(
            second_screen, option, primary_screen)
    elif option == "right-of":
        cmd = "xrandr --output {} --auto --{} {} --auto".format(
            second_screen, option, primary_screen)
    elif option == "above":
        cmd = "xrandr --output {} --auto --{} {} --auto".format(
            second_screen, option, primary_screen)
    elif option == "below":
        cmd = "xrandr --output {} --auto --{} {} --auto".format(
            second_screen, option, primary_screen)
    elif option == "rotate left":
        cmd = "xrandr --output {} --auto --{} --output {} --off".format(
            second_screen, option, primary_screen)
    elif option == "rotate right":
        cmd = "xrandr --output {} --auto --{} --output {} --off".format(
            second_screen, option, primary_screen)

    os.system(cmd)

    cmd = "feh --bg-fill {}/{}".format(my_wallpaper_path, my_default_wallpaper)
    os.system(cmd)

    return


def get_current_mouse_url():
    dx, dy = 4, 130
    text = popen("xdotool getmouselocation")
    mouse_pos = dict([x.split(":") for x in text.split()])
    cx, cy = int(mouse_pos.get("x", 0)), int(mouse_pos.get("y", 0))

    cmd = "xdotool mousemove {} {} click 3".format(cx, cy)
    os.system(cmd)

    cmd = "xdotool mousemove {} {} click 1".format(cx + dx, cy + dy)
    os.system(cmd)

    cmd = "xdotool mousemove {} {}".format(cx, cy)
    os.system(cmd)
    url = popen("xclip -o")
    url = url.strip()

    if re.match(r'^(http|https|www|file).+', url):
        return url, True

    return '', False


def download_arxiv_to_lib():
    my_library_paper_path = my_library_path + "/papers"

    url, ok = get_current_mouse_url()
    if not ok:
        msg = "Error: get current mouse url failed"
        os.system("notify-send '{}'".format(msg))
        return

    if "arxiv" not in url:
        msg = "Error: not a lawful arxiv url"
        os.system("notify-send '{}'".format(msg))
        return

    file_name = url.split("/")[-1] + ".pdf"
    file_path = my_library_paper_path + "/" + file_name

    if not os.path.exists(file_path):
        msg = "downloading {} \nto          {}".format(url, file_path)
        os.system("notify-send '{}'".format(msg))

        wget.download(url, file_path)

        msg = "job done!"
        os.system("notify-send '{}'".format(msg))

        open_file_at_foreground(file_path)
        keep_file_or_not(file_path)
    else:
        open_file_at_background(file_path)

    return


def download_cur_to_download():
    url, ok = get_current_mouse_url()
    if not ok:
        msg = "Error: get current mouse url failed"
        os.system("notify-send '{}'".format(msg))
        return

    # github repo url: automatically redirect to raw file url: download always fail, so just open it with the chrome
    # need:
    # - turn off vpn
    # - /etc/hosts
    #   github raw
    #   199.232.28.133 raw.githubusercontent.com
    if re.match(r'^https://github.com.+blob.*', url):
        url = url.replace("github.com", "raw.githubusercontent.com")
        url = url.replace("blob/", "")
        cmd = "google {}".format(url)
        os.system(cmd)
        return

    # the file name may not consistent with the real name
    file_name = url.split("/")[-1]
    file_path = my_download_path + "/" + file_name

    # - wget is too slow, but it can return the file name
    # - aria2c is not fast either
    if not os.path.exists(file_path):
        msg = "downloading {} \nto          {}".format(url, my_download_path)
        os.system("notify-send '{}'".format(msg))

        if re.match(r'.*(tar|rar|zip|gzip|xz|7z|bz|bz2|tgz|pkg|exe).*', url):
            cmd = "google {}".format(url)
            os.system(cmd)
        else:
            file_path = wget.download(url, my_download_path)

            msg = "job done!"
            os.system("notify-send '{}'".format(msg))

            open_file_at_foreground(file_path)
            keep_file_or_not(file_path)
    else:
        open_file_at_background(file_path)

    return


if __name__ == '__main__':
    pass
