#!/bin/bash

# KEYBOARD_CMD="onboard"
#
# is_keyboard_running() {
#     pgrep -x onboard >/dev/null
# }
#
# show_keyboard() {
#     if ! is_keyboard_running; then
#         onboard &
#         sleep 0.4
#     fi
#     xdotool search --class Onboard windowmap %@ 2>/dev/null
# }
#
# hide_keyboard() {
#     xdotool search --class Onboard windowunmap %@ 2>/dev/null
# }
#
# LAST_WID=""
#
# while true; do
#     WID=$(xdotool getwindowfocus 2>/dev/null)
#
#     if [ "$WID" != "$LAST_WID" ]; then
#         LAST_WID="$WID"
#
#         CLASS=$(xprop -id "$WID" WM_CLASS 2>/dev/null | awk -F'"' '{print $2}')
#
#         echo $CLASS
#
#         case "$CLASS" in
#             # ❌ 不弹键盘
#             onboard)
#                 hide_keyboard
#                 ;;
#             # ✅ 会输入的 GUI
#             dwm|st-256color|qutebrowser|google-chrome|obsidian|wechat)
#                 show_keyboard
#                 ;;
#             *)
#                 hide_keyboard
#                 ;;
#         esac
#     fi
#
#     sleep 0.2
# done


while true; do
    if ! pgrep -x onboard >/dev/null; then
        onboard &
    fi
    sleep 2
done
