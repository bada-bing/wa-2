#!/usr/bin/env bash

WPR_projects=( "cart-ui" "cart-ui-next" "form-configuration-ui" "shopping-profile-ui" "wpr-form-generator" )
GO_projects=( "go" )
WA_projects=( "wa-1" "wa-2")

# TODO change structure of projects directories - raycast extension and go applications are not direct children of ~/src
# Consider to reorganize the projects in a more structured way (not necessarily to update the script)

all_projects=("${WPR_projects[@]}" "${GO_projects[@]}" ) # tmuxinator projects

# 0. Select session name
# - Get first level directories from the ~/src directory and select one
session=$(find ~/src -maxdepth 1 -mindepth 1 -type d | fzf)
session_name=$(basename $session | tr . -)

# 1. Setup a new tmux session with the selected name
# - Only create the session if it doesn't exist already
# - Special projects are started with the tmuxinator
if ! tmux has-session -t $session_name 2>/dev/null; then
    if [[ " ${all_projects[@]} " =~ " ${session_name} " ]] 2>/dev/null; then
        # each tmuxinator project is created with attached: false
        tmuxinator start $session_name
    elif [[ " ${WA_projects[@]} " =~ " ${session_name} " ]] 2>/dev/null; then
        session_name="wa" # all wa projects are started with the same tmuxinator project
        tmuxinator start $session_name
    else
        tmux new-session -d -s $session_name -c $session
    fi
fi

# 2. If the session name is in the WPR_projects array, check if the wpr-utils session is running
# - Check if session name is in WPR_projects array
if [[ " ${WPR_projects[@]} " =~ " ${session_name} " ]] 2>/dev/null; then
    # start utils session if it is not running already
    if ! tmux has-session -t "wpr-utils" 2>/dev/null; then
        tmuxinator start wpr-utils
    fi
fi

# 3. Attach to the selected session
if [ -z "$TMUX" ]; then
    # if not inside a tmux session
    tmux attach-session -t $session_name
else
     # if already inside a tmux session
    tmux switch-client -t $session_name
fi
