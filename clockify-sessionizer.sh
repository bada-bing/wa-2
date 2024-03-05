#!/bin/bash

# I decided to execute this script as part of the tmuxinator projects
# - It should be easy to move it to tmux_sessionizer.sh if I decide to do so

# tmux list-sessions # it is easier with tmux than with vs code
# echo "hello $1" # pwd is in the correct dir inside of tmux session

# 1. identify the task (find its id)
# 2. set current time as start time for the task
# 3. TODO 🚧 set description for the task (read it from linear)
# 4. start a new timer for the task
    # - Clockify API will automatically stop any running timer, so no need to do that


# 1. Identify the task
# 1.1. - extract the ISSUE key from the current branch
ISSUE_KEY=$(git branch | grep \* | cut -d ' ' -f2 |  awk -F'/' '{print $2}')
# cut -d ' ' -f2 - separate by ' ' and return second field

API_KEY=${1:-$CLOCKIFY_API_KEY}
WORKSPACES_URL=https://api.clockify.me/api/v1/workspaces

# 1.2. - get all tasks for the project
TASKS_URL="$WORKSPACES_URL/$CLOCKIFY_MAIN_WORKSPACE_ID/projects/$STORIES_ID/tasks"
RESPONSE=$(curl -s -H "X-Api-Key: $API_KEY" $TASKS_URL)

# 1.3 - find the task with the current issue key and extract its id
# Filter the task with the issue key
TASK=$(echo "$RESPONSE" | jq -r ".[] | select(.name==\"$ISSUE_KEY\") | {id: .id, name: .name}")

# Extract the task id
TASK_ID=$(echo "$TASK" | jq -r '.id')

# 2. Get the current date
CURRENT_DATE=$(date -v-2H +'%Y-%m-%dT%H:%M:%SZ')

# CURRENT_DATE=$(date -u -d '+1 hour' +'%Y-%m-%dT%H:%M:%SZ')
# CURRENT_DATE=$(date +'%Y-%m-%dT%H:%M:%S%z') # - this should give relative timezone but Clockify doesnt accept it

# 3. Set description for the task
DESCRIPTION="✨ Implement"

# 4. Start a new timer for the task
# 4.1 - set the url for the request
ADD_ENTRY_URL="https://api.clockify.me/api/v1/workspaces/$CLOCKIFY_MAIN_WORKSPACE_ID/time-entries"

# 4.2 - set the request body
read -r -d '' REQ_BODY << EOF
{
  "billable": true,
  "description": "$DESCRIPTION",
  "projectId": "$STORIES_ID",
  "taskId": "$TASK_ID",
  "start": "$CURRENT_DATE",
  "type": "REGULAR"
}
EOF

# 4.3 - make curl post request to start a new timer
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-Api-Key: $API_KEY" -H "Content-Type: application/json" -d "$REQ_BODY" $ADD_ENTRY_URL)

echo "Add time entry: Issue: $ISSUE_KEY; Task ID: $TASK_ID; Status: $STATUS_CODE"