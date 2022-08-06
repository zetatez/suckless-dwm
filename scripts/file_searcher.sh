#!/bin/sh

RG() {
    RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
    INITIAL_QUERY="$1"
    matched=$(FZF_DEFAULT_COMMAND="$RG_PREFIX '$INITIAL_QUERY' || true" fzf --bind "change:reload:$RG_PREFIX {q} || true" --ansi --disabled --query "$INITIAL_QUERY" --delimiter : --preview 'bat --style=full --color=always --highlight-line {2} {1}')
    filename=$(echo $matched|awk -F: '{print $1}')
    rownum=$(echo $matched|awk -F: '{print $2}')
    # colnum=$(echo $matched|awk -F: '{print $3}')
    if [ "$filename" != "" ]; then vim +$rownum $filename; fi
}

RG
