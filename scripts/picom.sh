#!/bin/sh

GRAYSCALE=$(cat <<EOF
uniform float opacity;
uniform bool invert_color;
uniform sampler2D tex;
void main() {
    vec4 c = texture2D(tex, gl_TexCoord[0].xy);
    float g = 0.2126 * c.r + 0.7152 * c.g + 0.0722 * c.b; // CIELAB luma, based on human tristimulus.
    c = vec4(vec3(g), c.a);
    if (invert_color)
        c = vec4(vec3(c.a, c.a, c.a) - vec3(c), c.a);
    c *= opacity;
    gl_FragColor = c;
}
EOF
)

killall -q picom
# wait for picom to be killed, or can not be started
sleep 0.05

if [ "xgrayscale" == "x$1" ]; then
    nohup picom --glx-fshader-win "$GRAYSCALE" --backend glx >>/dev/null 2>&1 &
elif [ "xnormal" == "x$1" ]; then
    nohup picom >>/dev/null 2>&1 &
fi
