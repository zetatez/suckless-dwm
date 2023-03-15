#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import sys
import signal
import re
import time
import wget
import shutil
import psutil
import socket
import json
import pyperclip
import sqlparse
import datetime
from PyQt5 import QtWidgets

my_home_path = "/home/dionysus"
my_download_path = "/home/dionysus/Downloads"
my_library_path = "/home/dionysus/my-library"
my_play_path = "/home/dionysus/my-play"
my_trojan_path = "/home/dionysus/.trojan"
my_wallpaper_path = "/home/dionysus/Pictures/wallpapers"
my_default_wallpaper = "Van-Gogh-003.jpg"

win_name_float = "00001011"
win_name_scratchpad = "scratchpad"


def empty():
    msg = "empty"
    os.system("notify-send '{}'".format(msg))


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


def get_pids_by_pname(pname):
    pids = []
    for p in psutil.process_iter():
        if p.name() == pname:
            pids.append(p.pid)
    return pids


def get_pids_by_cmd(cmd):
    cmd = cmd.rstrip(" &").replace("'", "").replace('"', "").strip()
    cmd_ps = "ps -ef|grep '{}'".format(cmd) + "|grep -v grep|awk '{print $2}'"
    pids = [int(pid) for pid in popen(cmd_ps).strip().replace("\n", " ").strip().split()]
    return pids


def toggle_by_pname(pname, cmd):
    pids = get_pids_by_pname(pname)
    if pids:
        [os.kill(pid, signal.SIGKILL) for pid in pids]
    else:
        os.system(cmd)
    return


def toggle_by_cmd(cmd):
    pids = get_pids_by_cmd(cmd)
    if pids:
        [os.kill(pid, signal.SIGKILL) for pid in pids]
    else:
        os.system(cmd)
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
    cmd = "echo '{}'|dmenu -p 'keep file?'".format('\n'.join(['yes', 'no']))
    option = popen(cmd)
    if option.strip() == "no":
        os.remove(file_path)
    return


# utils
# ---------------------------------------------------------
# app
# -----------------------
def open_app(app):
    return os.system(app)


def app_passmenu():
    return open_app(app="passmenu")


def app_photoshop():
    return open_app(app="gimp")


def app_wps():
    return open_app(app="wps")


# wf: workflow
# -----------------------
def wf_cal():
    cmd = "dmenu < /dev/null -p 'cal with julia>>'"
    cmd = popen(cmd).strip()
    if not cmd:
        return
    cmd = "julia -E '{}'".format(cmd)
    msg = popen(cmd).strip()
    pyperclip.copy(msg)
    os.system("notify-send '{}'".format(msg))
    return


def wf_clipmenu():
    cmd = "clipmenu"
    os.system(cmd)
    return


def wf_dmenu():
    cmd = "dmenu_run"
    os.system(cmd)
    return


def __get_current_mouse_url():
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


def wf_download_cur_to_download():
    url, ok = __get_current_mouse_url()
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


def wf_download_arxiv_to_lib():
    my_library_paper_path = my_library_path + "/papers"

    url, ok = __get_current_mouse_url()
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


def wf_find_local_area_network_server():
    ip = "127.0.0.1"
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
    except Exception as e:
        msg = "find loacl area network server failed: {}".format(e)
        os.system("notify-send '{}'".format(msg))
        s.close()

    ip = list(ip)
    ip[-1] = '0'
    ip = ''.join(ip)
    cmd = "nmap -sT -p 22 -oG - {}/24".format(ip)
    s = popen(cmd)
    server_list = [x.strip().split()[1] for x in s.split('\n') if 'ssh' in x]
    server_list.sort()
    cmd = "echo '{}'|dmenu -p 'local'".format('\n'.join(server_list))
    option = popen(cmd).strip()
    if not option:
        return

    pyperclip.copy(option)
    msg = "find local area network server success, please check clipboard: {}".format(option)
    os.system("notify-send '{}'".format(msg))

    return


def wf_format_json():
    last_copied_str = pyperclip.paste()
    try:
        s = json.dumps(json.loads(last_copied_str), indent=2)
        pyperclip.copy(s)
        msg = "format json success, please check clipboard:\n{}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "format json failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_format_sql():
    last_copied_str = pyperclip.paste()
    try:
        s = sqlparse.format(last_copied_str, reindent=True, indent=2, keyword_case='upper')
        pyperclip.copy(s)
        msg = "format sql success, please check clipboard:\n{}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "format sql failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_get_now_unix_nano_sec():
    try:
        unix_nano_sec = int(time.time_ns())
        pyperclip.copy(unix_nano_sec)
        msg = "get unix nano sec success, please check clipboard: {}".format(unix_nano_sec)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "get unix nano sec failed: {}".format(e)
        os.system("notify-send '{}'".format(msg))

    return


def wf_get_now_unix_sec():
    try:
        unix_sec = int(time.time())
        pyperclip.copy(unix_sec)
        msg = "get unix sec success, please check clipboard: {}".format(unix_sec)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "get now unix failed: {}".format(e)
        os.system("notify-send '{}'".format(msg))

    return


def wf_handle_copied():
    last_copied_str = pyperclip.paste().strip()

    # if a local file
    if os.path.exists(last_copied_str):
        cmd = "st -e lazy -o {} &".format(last_copied_str)
        os.system(cmd)
        return

    # if a url of: ^(http|https|www|file).+
    if re.match(r'^(http|https|www|file).+', last_copied_str):
        cmd = "vivaldi-stable {} &".format(last_copied_str)
        os.system(cmd)
        return

    # search if match none
    prompt = last_copied_str if len(last_copied_str) < 12 else last_copied_str[:12]
    cmd = "echo '{}'|dmenu -p 'Can not handle: {}. search ?'".format('\n'.join(['yes', 'no']), prompt)
    option = popen(cmd).strip()
    if option == "yes":
        url = "https://cn.bing.com/search?q={}".format(last_copied_str.replace(" ", "+"))
        cmd = "chrome {}".format(url)
        os.system(cmd)

    return


def wf_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
        pyperclip.copy(ip)
        msg = "get host ip success, please check clipboard: {}".format(ip)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "get host ip failed: {}".format(e)
        os.system("notify-send '{}'".format(msg))
        s.close()

    return


def wf_latex():
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

    cmd_tex = "st -g {} -t {} -c {} -e nvim {} &".format(get_geometry_for_st(0.01, 0.05, 94, 56), win_name_float,
                                                         win_name_float, my_play_notes_today_note_tex)

    cmd_xoj = "st -g {} -t {} -c {} -e xournalpp {} &".format(get_geometry_for_st(0.52, 0.16, 88, 38), win_name_float,
                                                              win_name_float, my_play_notes_today_note_xoj)

    pids_tex = get_pids_by_cmd(cmd_tex)
    if pids_tex:
        os.system("notify-send '{}'".format("already in running"))
    else:
        os.system(cmd_tex)

    pids_xoj = get_pids_by_cmd(cmd_xoj)
    if pids_xoj:
        os.system("notify-send '{}'".format("already in running"))
    else:
        os.system(cmd_xoj)

    return


def wf_map():
    cmd = "dmenu < /dev/null -p 'enter location'"
    option = popen(cmd).strip()
    if option:
        url = "https://ditu.amap.com/search?query={}".format(option)
        cmd = "chrome {}".format(url)
        os.system(cmd)

    return


def wf_mount_to_xyz():
    cmd = "ls -1 /dev/sd*"
    s = popen(cmd)
    if "no matches found" in s:
        return

    devices = [x.strip() for x in s.split("\n") if x.strip() != "" and len(x.strip()) >= 9]
    devices.sort()

    print(devices)

    cmd = "echo '{}'|dmenu -p 'mount ?'".format('\n'.join(devices))
    dev = popen(cmd).strip()
    if not dev:
        return

    cmd = "echo '{}'|dmenu -p 'mount {} to ?'".format('\n'.join(["/x", "/y", "/z"]), dev)
    dst = popen(cmd).strip()
    if not dst:
        return

    cmd = "sudo mount {} {}".format(dev, dst)
    os.system(cmd)

    return


def wf_sketchpad():
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

    pids = get_pids_by_cmd(cmd)
    if pids:
        os.system("notify-send '{}'".format("already in running"))
    else:
        os.system(cmd)

    return


def wf_ssh():
    ssh_list = []

    file_id_rsa_pub = os.path.join(my_home_path, ".ssh/id_rsa.pub")
    if os.path.exists(file_id_rsa_pub):
        id_rsa_pub = list(set([line.strip().split()[-1] for line in read_file(file_id_rsa_pub).split("\n") if line.strip()]))
        ssh_list.extend(id_rsa_pub)

    file_know_hosts = os.path.join(my_home_path, ".ssh/known_hosts")
    if os.path.exists(file_know_hosts):
        know_hosts = [
            "root@" + x
            for x in list(set([line.strip().split()[0] for line in read_file(file_know_hosts).split("\n") if line.strip()]))
        ]
        ssh_list.extend(know_hosts)

    ssh_list = [x for x in ssh_list if not (("127.0.0.1" in x) or ("github" in x) or (x.endswith(".com")))]
    ssh_list.sort(reverse=True)

    cmd = "echo '{}'|dmenu -p 'ssh'".format('\n'.join(ssh_list))
    option = popen(cmd).strip()
    if not option:
        return

    ssh_cmd = "ssh {}".format(option)
    cmd = "st -e /bin/bash -c \"echo '{}' && {}\"".format(ssh_cmd, ssh_cmd)
    os.system(cmd)

    return


def wf_trans_base_10_to_base_x():
    last_copied_str = pyperclip.paste()

    try:
        s = ""
        binary = bin(int(last_copied_str))
        cmd = "echo '{}'|dmenu -p 'base trans string to ?'".format('\n'.join(['2', '8', '10', '16']))
        option = popen(cmd).strip()
        if not option:
            return

        if option == "2":
            s = str(binary)
        elif option == "8":
            s = str(oct(int(binary, 2)))
        elif option == "10":
            s = str(int(binary, 2))
        elif option == "16":
            s = str(hex(int(binary, 2)))
        else:
            return

        pyperclip.copy(s)
        msg = "trans base 10 to base x success, please check clipboard:\n{}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "trans base 10 to base x failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_trans_datetime_to_unix_sec():
    last_copied_str = pyperclip.paste().strip()
    try:
        dt = datetime.datetime.strptime(last_copied_str, "%Y-%m-%d %H:%M:%S")
        s = str(int(datetime.datetime.timestamp(dt)))
        pyperclip.copy(s)
        msg = "trans datetime to unix sec success, please check clipboard: {}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "trans datetime to unix sec failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_trans_string_to_base_x():
    last_copied_str = pyperclip.paste()

    try:
        s = ""
        binary = ''.join(format(ord(i), '08b') for i in last_copied_str)
        cmd = "echo '{}'|dmenu -p 'base trans string to ?'".format('\n'.join(['2', '8', '10', '16']))
        option = popen(cmd).strip()
        if not option:
            return

        if option == "2":
            s = str(binary)
        elif option == "8":
            s = str(oct(int(binary, 2)))
        elif option == "10":
            s = str(int(binary, 2))
        elif option == "16":
            s = str(hex(int(binary, 2)))
        else:
            return

        pyperclip.copy(s)
        msg = "trans string to base x success, please check clipboard:\n{}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "trans string to base x failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_trans_unix_sec_to_datetime():
    last_copied_str = pyperclip.paste().strip()
    try:
        s = str(datetime.datetime.fromtimestamp(int(last_copied_str)))
        pyperclip.copy(s)
        msg = "trans unix sec to datetime success, please check clipboard: {}".format(s)
        os.system("notify-send '{}'".format(msg))
    except Exception as e:
        msg = "trans unix sec to datetime failed: {}\n{}".format(e, last_copied_str)
        os.system("notify-send '{}'".format(msg))

    return


def wf_umount_from_xyz():
    cmd = "echo '{}'|dmenu -p 'umount ?'".format('\n'.join(["/x", "/y", "/z"]))
    dst = popen(cmd).strip()
    if not dst:
        return

    cmd = "sudo umount {}".format(dst)
    os.system(cmd)

    return


def wf_web():
    websites = {
        "translate": "https://cn.bing.com/translator?ref=TThis&text=&from=zh-Hans&to=en",
        "scholar": "https://scholar.google.com",
        "arxiv": "https://arxiv.org",
        "wolframalpha": "https://www.wolframalpha.com",
        "bing": "https://cn.bing.com",
        "github": "https://github.com/zetatez?tab=repositories",
        "arch wiki": "https://wiki.archlinux.org",
        "suckless": "https://dwm.suckless.org",
        "map": "https://ditu.amap.com",
        "ocr": "http://ocr.space",
        "regex": "https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference",
        "bilibili": "https://www.bilibili.com",
        "cctv5": "https://tv.cctv.com/live/cctv5",
        "mall": "https://www.jd.com",
        "news of finance": "https://news.futunn.com/en/main/live?lang=zh-CN",
        "runoob": "https://www.runoob.com",
        "mirror": "https://developer.aliyun.com/mirror",
    }

    cmd = "echo '{}'|dmenu -p 'websites'".format('\n'.join(list(websites.keys())))
    option = popen(cmd).strip()
    if option:
        cmd = "chrome {}".format(websites.get(option, "https://cn.bing.com/search?q={}".format(option)))
        os.system(cmd)

    return


def wf_xournal():
    time_str = time.strftime("%Y-%m-%d", time.localtime())
    file_folder = my_library_path + "/xournal"
    file_path = file_folder + "/" + time_str + ".xopp"

    if not os.path.exists(file_folder):
        os.mkdir(file_folder)

    cmd = "xournalpp {} &".format(file_path)

    pids = get_pids_by_cmd(cmd)
    if pids:
        os.system("notify-send '{}'".format("already in running"))
    else:
        os.system(cmd)

    return


# toggle
# -----------------------
def toggle_addressbook():
    cmd = "st -e abook"
    toggle_by_cmd(cmd)

    return


def toggle_bluetooth():
    cmd = "bluetoothctl devices"
    devices = popen(cmd).strip()

    if not devices:
        msg = "bluetoothctl devices returned empty"
        os.system("notify-send '{}'".format(msg))
        return

    # sort by name
    devices = dict([(" ".join(y[1:]), y[0]) for y in [x.strip().lstrip("Device ").split(" ") for x in devices.split("\n")]])
    keys = list(devices.keys())
    keys.sort()
    devices = [[devices.get(k), k] for k in keys]
    devices = "\n".join([" ".join(x) for x in devices])

    cmd = "echo '{}'|dmenu -p 'bluetoosh device>'".format(devices)
    option = popen(cmd).strip()
    if not option:
        return

    id = option.strip().split(" ")[0].strip()

    cmd = "bluetoothctl disconnect"
    os.system(cmd)

    cmd = "bluetoothctl connect {}".format(id)

    res = popen(cmd).strip()
    if "successful" not in res:
        msg = "{} failed: \n{}".format(cmd, res)
        os.system("notify-send '{}'".format(msg))
        return

    return


def toggle_calendar_scheduling():
    cmd = "st -t {} -c {} -e nvim +':set laststatus=0' +'Calendar -view=week'".format("shceduling", "shceduling")
    toggle_by_cmd(cmd)

    return


def toggle_calendar_schedule():
    cmd = "st -g {} -t {} -c {} -e nvim +':set laststatus=0' +'Calendar -view=day'".format(
        get_geometry_for_st(0.80, 0.05, 40, 32), win_name_float, win_name_float)
    toggle_by_cmd(cmd)

    return


def toggle_diary():
    time_str = time.strftime("%Y-%m-%d", time.localtime())
    diary = "{}/diary/{}.md".format(my_home_path, time_str)

    if not os.path.exists(diary):
        s = "### {}\n".format(time.strftime("%a %b %d %H:%M:%S %p CST %Y", time.localtime()))
        write_file(diary, s)

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


def toggle_vivaldi():
    toggle_by_pname(pname="vivaldi-bin", cmd="vivaldi-stable")

    return


def toggle_chrome_with_proxy():
    cmd = "chrome --proxy-server='socks5://127.0.0.1:8000'"
    toggle_by_cmd(cmd)

    return


def toggle_gitter():
    toggle_by_pname(pname="gitter", cmd="gitter")

    return


def toggle_irc():
    cmd = "st -e irssi"
    toggle_by_cmd(cmd)

    return


def toggle_julia():
    cmd = "st -t {} -c {} -e julia".format(win_name_scratchpad, win_name_scratchpad)
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
    cmd_cava = "st -g {} -t cava -c cava -e cava &".format(get_geometry_for_st(0.74, 0.08, 40, 12))
    cmd_ncmpcpp = "st -g {} -t music -c music -e ncmpcpp &".format(get_geometry_for_st(0.52, 0.08, 40, 12))

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


# systemctl --user enable redshift.service now
def toggle_redshift():
    cmd = "systemctl --user status redshift.service|grep 'active (running)'"
    text = popen(cmd)
    if text:
        cmd = "systemctl --user stop redshift.service"
    else:
        cmd = "systemctl --user start redshift.service"
    os.system(cmd)

    return


def toggle_screenkey():
    cmd = "screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow"
    toggle_by_cmd(cmd)

    return


def toggle_show():
    cmd = "st -g {} -t {} -c {} -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0".format(
        get_geometry_for_st(0.74, 0.08, 40, 12), win_name_float, win_name_float)
    toggle_by_cmd(cmd)

    return


def toggle_sublime():
    toggle_by_pname(pname="subl", cmd="subl")

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
    toggle_by_pname(pname="ffmpeg", cmd=cmd)

    return


def toggle_rec_video():
    time_str = time.strftime("%Y-%m-%d-%H-%M-%S", time.localtime())
    w, h = get_cur_screen_geometry()
    dpy = os.environ.get("DISPLAY")
    cmd = "st  -t {} -c {} -e ffmpeg -y -s '{}x{}' -r 60 -f x11grab -i {} -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac {}/Videos/rec-v-a-{}.mkv".format(
        win_name_scratchpad, win_name_scratchpad, w, h, dpy, my_home_path, time_str)
    toggle_by_pname(pname="ffmpeg", cmd=cmd)

    return


def toggle_screen():
    primary_screen = "eDP-1"
    cmd = "xrandr|grep ' connected'|grep -v 'eDP-1'|awk '{print $1}'"
    second_screen = popen(cmd).strip()

    if not second_screen:
        msg = "have no second screen"
        os.system("notify-send '{}'".format(msg))
        return

    cmds = {
        "only": "xrandr --output {} --auto --output {} --off".format(second_screen, primary_screen),
        "primary only": "xrandr --output {} --auto --output {} --off".format(primary_screen, second_screen),
        "left of": "xrandr --output {} --auto --left-of {} --auto".format(second_screen, primary_screen),
        "right of": "xrandr --output {} --auto --right-of {} --auto".format(second_screen, primary_screen),
        "above": "xrandr --output {} --auto --above {} --auto".format(second_screen, primary_screen),
        "below": "xrandr --output {} --auto --below {} --auto".format(second_screen, primary_screen),
        "roate left": "xrandr --output {} --auto --rotate left --output {} --off".format(second_screen, primary_screen),
        "roate right": "xrandr --output {} --auto --rotate right --output {} --off".format(second_screen, primary_screen),
    }

    cmd = "echo '{}'|dmenu -p '󱣴'".format('\n'.join(list(cmds.keys())))
    option = popen(cmd).strip()
    if not option:
        return

    cmd = cmds.get(option, cmds.get("primary only", "echo"))
    os.system(cmd)

    cmd = "feh --bg-fill {}/{}".format(my_wallpaper_path, my_default_wallpaper)
    os.system(cmd)

    return


def toggle_sys_shortcuts():
    cmds = {
        "󰒲 suspend": "systemctl suspend",
        " poweroff": "systemctl poweroff",
        "ﰇ reboot": "systemctl reboot",
        "󰷛 slock": "slock & sleep 0.5 & xset dpms force off",
        "󰶐 off-display": "sleep .5; xset dpms force off",
    }
    cmd = "echo '{}'|dmenu -p ''".format('\n'.join(list(cmds.keys())))
    option = popen(cmd).strip()
    if not option:
        return

    cmd = cmds.get(option, "")
    if cmd:
        os.system(cmd)

    return


# ultra
# -----------------------
def exec_shell_cmd():
    cmd = "dmenu < /dev/null -p 'shell>>'"
    cmd = popen(cmd).strip()
    if not cmd:
        return
    cmd = "/bin/bash -c '{}'".format(cmd)
    msg = popen(cmd).strip()
    pyperclip.copy(msg)
    os.system("notify-send '{}'".format(msg))
    return


def exec_julia_cmd():
    cmd = "dmenu < /dev/null -p 'julia>>'"
    cmd = popen(cmd).strip()
    if not cmd:
        return
    cmd = "julia -E '{}'".format(cmd)
    msg = popen(cmd).strip()
    pyperclip.copy(msg)
    os.system("notify-send '{}'".format(msg))
    return


# ultra
# -----------------------
def ultra():

    options = {
        "[wf] search": search,
        "[wf] handle copied": wf_handle_copied,
        "[wf] web": wf_web,
        "[wf] ssh": wf_ssh,
        "[wf] cal": wf_cal,
        "[wf] dmenu": wf_dmenu,
        "[wf] find local area network server": wf_find_local_area_network_server,
        "[wf] format json": wf_format_json,
        "[wf] format sql": wf_format_sql,
        "[wf] map": wf_map,
        "[wf] download arxiv to lib": wf_download_arxiv_to_lib,
        "[wf] download cur to download": wf_download_cur_to_download,
        "[wf] ip": wf_ip,
        "[wf] mount": wf_mount_to_xyz,
        "[wf] umount": wf_umount_from_xyz,
        "[wf] get now unix nano sec": wf_get_now_unix_nano_sec,
        "[wf] get now unix sec": wf_get_now_unix_sec,
        "[wf] trans baee 10 to base x": wf_trans_base_10_to_base_x,
        "[wf] trans datetime to unix sec": wf_trans_datetime_to_unix_sec,
        "[wf] trans string to base x": wf_trans_string_to_base_x,
        "[wf] clipmenu": wf_clipmenu,
        "[wf] trans unix sec to datetime": wf_trans_unix_sec_to_datetime,
        "[wf] trans unix sec to datetime": wf_trans_unix_sec_to_datetime,
        "[wf] inkspace": wf_sketchpad,
        "[wf] sketchpad": wf_sketchpad,
        "[wf] latex": wf_latex,
        "[wf] note": wf_xournal,
        "[wf] xournal": wf_xournal,
        "[exec] shell cmd": exec_shell_cmd,
        "[exec] julia cmd": exec_julia_cmd,
        "[tg] addressbook": toggle_addressbook,
        "[tg] bluetooth": toggle_bluetooth,
        "[tg] calendar schedule": toggle_calendar_schedule,
        "[tg] calendar scheduling": toggle_calendar_scheduling,
        "[tg] chrome with proxy": toggle_chrome_with_proxy,
        "[tg] diary": toggle_diary,
        "[tg] flameshot": toggle_flameshot,
        "[tg] flameshot": toggle_flameshot,
        "[tg] gitter": toggle_gitter,
        "[tg] irc": toggle_irc,
        "[tg] julia": toggle_julia,
        "[tg] lazydocker": toggle_lazydocker,
        "[tg] mathpix": toggle_mathpix,
        "[tg] music net cloud": toggle_music_net_cloud,
        "[tg] music": toggle_music,
        "[tg] mutt": toggle_mutt,
        "[tg] rec audio": toggle_rec_audio,
        "[tg] rec video": toggle_rec_video,
        "[tg] redshift": toggle_redshift,
        "[tg] rss": toggle_rss,
        "[tg] screen": toggle_screen,
        "[tg] screen": toggle_screen,
        "[tg] screenkey": toggle_screenkey,
        "[tg] show": toggle_show,
        "[tg] sublime": toggle_sublime,
        "[tg] sys shortcuts": toggle_sys_shortcuts,
        "[tg] top": toggle_top,
        "[tg] trojan": toggle_trojan,
        "[tg] vifm": toggle_vifm,
        "[tg] vivaldi": toggle_vivaldi,
        "[tg] wallpaper": toggle_wallpaper,
        "[tg] wechat": toggle_wechat,
        "[tg] wifi": toggle_wifi,
        "[app] passmenu": app_passmenu,
        "[app] photoshop": app_photoshop,
        "[app] wps": app_wps,
    }

    cmd = "echo '{}'|dmenu -p ' Ultra'".format("\n".join(options.keys()))
    option = popen(cmd).strip()
    if option:
        options.get(option, empty)()


def search():
    cmd = "dmenu < /dev/null -p 'search'"
    search = popen(cmd).strip()
    if not search:
        return

    # if a local file
    if os.path.exists(search):
        cmd = "st -e lazy -o {} &".format(search)
        os.system(cmd)
        return

    # if a url of: ^(http|https|www|file).+
    if re.match(r'^(http|https|www|file).+', search):
        cmd = "vivaldi-stable {} &".format(search)
        os.system(cmd)
        return

    # if websites keyword
    websites = {
        "arch wiki": "https://wiki.archlinux.org",
        "arxiv": "https://arxiv.org",
        "bilibili": "https://www.bilibili.com",
        "bing": "https://cn.bing.com",
        "cctv5": "https://tv.cctv.com/live/cctv5",
        "github": "https://github.com/zetatez?tab=repositories",
        "mall": "https://www.jd.com",
        # "map": "https://ditu.amap.com",
        "mirror": "https://developer.aliyun.com/mirror",
        "news": "https://news.futunn.com/en/main/live?lang=zh-CN",
        "ocr": "http://ocr.space",
        "regex": "https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference",
        "runoob": "https://www.runoob.com",
        "scholar": "https://scholar.google.com",
        "suckless": "https://dwm.suckless.org",
        "trans": "https://cn.bing.com/translator?ref=TThis&text=&from=zh-Hans&to=en",
        "translate": "https://cn.bing.com/translator?ref=TThis&text=&from=zh-Hans&to=en",
        "wolframalpha": "https://www.wolframalpha.com",
    }

    # TODO: <00:46:36 2023-03-15- Author: Dion>: fuzzy match ?
    url = websites.get(search, False)
    if url:
        cmd = "chrome {}".format(url)
        os.system(cmd)
        return

    # if workflow keyword
    workflows = {
        "addressbook": toggle_addressbook,
        "bluetooth": toggle_bluetooth,
        "calendar schedule": toggle_calendar_schedule,
        "calendar scheduling": toggle_calendar_scheduling,
        "chrome with proxy": toggle_chrome_with_proxy,
        "diary": toggle_diary,
        "cal": wf_cal,
        "clipmenu": wf_clipmenu,
        "shell": exec_shell_cmd,
        "julia": exec_julia_cmd,
        "dmenu": wf_dmenu,
        "download arxiv to lib": wf_download_arxiv_to_lib,
        "download cur to download": wf_download_cur_to_download,
        "find local area network server": wf_find_local_area_network_server,
        "flameshot": toggle_flameshot,
        "flameshot": toggle_flameshot,
        "format json": wf_format_json,
        "format sql": wf_format_sql,
        "get now unix nano sec": wf_get_now_unix_nano_sec,
        "get now unix sec": wf_get_now_unix_sec,
        "gitter": toggle_gitter,
        "handle copied": wf_handle_copied,
        "inkspace": wf_sketchpad,
        "ip": wf_ip,
        "irc": toggle_irc,
        "julia": toggle_julia,
        "latex": wf_latex,
        "lazydocker": toggle_lazydocker,
        "map": wf_map,
        "mathpix": toggle_mathpix,
        "mount": wf_mount_to_xyz,
        "music net cloud": toggle_music_net_cloud,
        "music": toggle_music,
        "mutt": toggle_mutt,
        "note": wf_xournal,
        "passmenu": app_passmenu,
        "photoshop": app_photoshop,
        "rec audio": toggle_rec_audio,
        "rec video": toggle_rec_video,
        "redshift": toggle_redshift,
        "rss": toggle_rss,
        "screen": toggle_screen,
        "screen": toggle_screen,
        "screenkey": toggle_screenkey,
        "search": search,
        "shot": toggle_flameshot,
        "show": toggle_show,
        "sketchpad": wf_sketchpad,
        "ssh": wf_ssh,
        "sublime": toggle_sublime,
        "sys shortcuts": toggle_sys_shortcuts,
        "top": toggle_top,
        "trans baee 10 to base x": wf_trans_base_10_to_base_x,
        "trans datetime to unix sec": wf_trans_datetime_to_unix_sec,
        "trans string to base x": wf_trans_string_to_base_x,
        "trans unix sec to datetime": wf_trans_unix_sec_to_datetime,
        "trans unix sec to datetime": wf_trans_unix_sec_to_datetime,
        "trojan": toggle_trojan,
        "ultra": ultra,
        "umount": wf_umount_from_xyz,
        "vifm": toggle_vifm,
        "vivaldi": toggle_vivaldi,
        "wallpaper": toggle_wallpaper,
        "web": wf_web,
        "wechat": toggle_wechat,
        "wifi": toggle_wifi,
        "wps": app_wps,
        "xournal": wf_xournal,
    }
    workflow = workflows.get(search, empty)
    if workflow != empty:
        workflow()
        return

    url = "https://cn.bing.com/search?q={}".format(search.replace(" ", "+"))
    cmd = "chrome {}".format(url)
    os.system(cmd)
    return


if __name__ == '__main__':
    wf_cal()
    pass
