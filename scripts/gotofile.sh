#!/bin/sh

RG() {
  INITIAL_QUERY=""
  RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
  matched=$(FZF_DEFAULT_COMMAND="$RG_PREFIX '$INITIAL_QUERY'" fzf --bind "change:reload:$RG_PREFIX {q} || true" --ansi --disabled --query "$INITIAL_QUERY" --height=100% --layout=reverse --delimiter : --preview 'bat --style=full --color=always --highlight-line {2} {1}')
  filename=$(echo $matched|awk -F: '{print $1}')
  rownum=$(echo $matched|awk -F: '{print $2}')
  if [ "$filename" != "" ]; then nvim +$rownum $filename; fi
}

RG
