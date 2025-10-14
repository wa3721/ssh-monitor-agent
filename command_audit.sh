#!/bin/bash
# 创建日志记录函数
set +x

AUDIT_LOG="/var/log/ssh_user_audit.log"
if [ ! -f "$AUDIT_LOG" ]; then
    sudo touch "$AUDIT_LOG"
    sudo chmod 666 "$AUDIT_LOG"
fi
SESSION_MARKER_FILE="/home/ssh_session_$$_${RANDOM}"

log_command() {
     local exit_code=$?
   { set +x; } 2>/dev/null

       if [[ ! -f "$SESSION_MARKER_FILE" ]]; then
           touch "$SESSION_MARKER_FILE"
           return 0
       fi
    local last_command
    last_command=$(history 1 | { read -r num  cmd; echo "$cmd"; })
    local timestamp
    timestamp=$(date '+%Y-%m-%d %T')

    if [[ -n "$SSH_CLIENT" ]]; then
        client_ip=$(echo "$SSH_CLIENT" | awk '{print $1}')
        client_port=$(echo "$SSH_CLIENT" | awk '{print $2}')
    else
        client_ip="LOCAL"
        client_port="-"
    fi

   log_entry="| Time: $timestamp | User: $(whoami) | IP: $client_ip | Port: $client_port | PWD: $PWD | Command: $last_command | ExitCode: $exit_code |"

#    echo "$log_entry" >> /var/log/ssh_command_audit.log

    # 2. 发送到远程位置（核心步骤）
    # 示例：使用 curl 发送到 HTTP API
    if curl --max-time 5 -X POST -H "Content-Type: text/plain" -d "$log_entry" http://127.0.0.1:8080/command_log
    then :
    else
       echo "[$timestamp] Failed to send log to remote server." >> "$AUDIT_LOG"
       echo "$log_entry" >> "$AUDIT_LOG"
    fi
}

cleanup_session_marker() {
    rm -f "$SESSION_MARKER_FILE"
}
trap cleanup_session_marker EXIT

export PROMPT_COMMAND="log_command"