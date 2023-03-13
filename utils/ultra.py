#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import utils


def empty():
    msg = "action not found"
    os.system("notify-send '{}'".format(msg))


options = {
    "[wf] handle copied": utils.wf_handle_copied,
    "[wf] websites": utils.wf_websites,
    "[wf] map": utils.wf_map,
    "[wf] format json": utils.wf_format_json,
    "[wf] format sql": utils.wf_format_sql,
    "[wf] get host ip": utils.wf_get_host_ip,
    "[wf] get now unix nano sec": utils.wf_get_now_unix_nano_sec,
    "[wf] get now unix sec": utils.wf_get_now_unix_sec,
    "[wf] inkspace": utils.wf_sketchpad,
    "[wf] latex": utils.wf_latex,
    "[wf] note": utils.wf_xournal,
    "[wf] sketchpad": utils.wf_sketchpad,
    "[wf] trans datetime to unix sec": utils.wf_trans_datetime_to_unix_sec,
    "[wf] trans unix sec to datetime": utils.wf_trans_unix_sec_to_datetime,
    "[wf] trans unix sec to datetime": utils.wf_trans_unix_sec_to_datetime,
    "[wf] xournal": utils.wf_xournal,
    "[wf] trans baee 10 to base x": utils.wf_trans_base_10_to_base_x,
    "[wf] trans string to base x": utils.wf_trans_string_to_base_x,
    "[wf] download arxiv to lib": utils.wf_download_arxiv_to_lib,
    "[wf] download cur to download": utils.wf_download_cur_to_download,
    "[tg] flameshot": utils.toggle_flameshot,
    "[tg] screen": utils.toggle_screen,
    "[tg] addressbook": utils.toggle_addressbook,
    "[tg] bluetooth": utils.toggle_bluetooth,
    "[tg] calendar scheduling": utils.toggle_calendar_scheduling,
    "[tg] calendar schedule": utils.toggle_calendar_schedule,
    "[tg] diary": utils.toggle_diary,
    "[tg] top": utils.toggle_top,
    "[tg] trojan": utils.toggle_trojan,
    "[tg] flameshot": utils.toggle_flameshot,
    "[tg] vivaldi": utils.toggle_vivaldi,
    "[tg] chrome with proxy": utils.toggle_chrome_with_proxy,
    "[tg] gitter": utils.toggle_gitter,
    "[tg] irc": utils.toggle_irc,
    "[tg] julia": utils.toggle_julia,
    "[tg] lazydocker": utils.toggle_lazydocker,
    "[tg] mathpix": utils.toggle_mathpix,
    "[tg] music": utils.toggle_music,
    "[tg] music net cloud": utils.toggle_music_net_cloud,
    "[tg] mutt": utils.toggle_mutt,
    "[tg] rss": utils.toggle_rss,
    "[tg] redshift": utils.toggle_redshift,
    "[tg] screenkey": utils.toggle_screenkey,
    "[tg] show": utils.toggle_show,
    "[tg] sublime": utils.toggle_sublime,
    "[tg] vifm": utils.toggle_vifm,
    "[tg] wechat": utils.toggle_wechat,
    "[tg] wifi": utils.toggle_wifi,
    "[tg] wallpaper": utils.toggle_wallpaper,
    "[tg] rec audio": utils.toggle_rec_audio,
    "[tg] rec video": utils.toggle_rec_video,
    "[tg] screen": utils.toggle_screen,
    "[tg] sys shortcuts": utils.toggle_sys_shortcuts,
    "[app] passmenu": utils.app_passmenu,
    "[app] photoshop": utils.app_photoshop,
    "[app] wps": utils.app_wps,
}

cmd = "echo '{}'|dmenu -p ' Ultra'".format("\n".join(options.keys()))
option = utils.popen(cmd).strip()
if option:
    options.get(option, empty)()
