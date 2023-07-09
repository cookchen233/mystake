ps -ef | grep mystake | grep -v "grep" | awk '{print $2}' | xargs kill -9
unset http_proxy https_proxy all_proxy
nohup mystake monitor up >/dev/null  2>&1  &
ps -ef |grep mystake