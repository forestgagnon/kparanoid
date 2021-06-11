source <(kubectl completion bash)
[ -f ~/.fzf.bash ] && source ~/.fzf.bash

set_bash_prompt(){
    local ctx=$(kubectl config current-context 2>/dev/null)
    PS1="\n$(tput setaf 1)kparanoid$(tput sgr0)@$(tput setaf 2)\h$(tput sgr0) \w\n$(tput setaf 4)$KPARANOID_CLUSTER_NAME$(tput sgr0) ctx=$(tput setaf 3)$ctx$(tput sgr0)\n\$ "
}

export PROMPT_COMMAND=set_bash_prompt
export HISTCONTROL=ignoreboth:erasedups

umask 0000

alias k='kubectl'
