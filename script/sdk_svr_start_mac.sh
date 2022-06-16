#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000
# shellcheck disable=SC2154
# shellcheck disable=SC2086
# shellcheck disable=SC2002
list1=$(cat $config_path | grep openImApiPort | awk -F '[:]' '{print $NF}')
# shellcheck disable=SC2002
list2=$(cat "$config_path" | grep openImWsPort | awk -F '[:]' '{print $NF}')
# shellcheck disable=SC2002
list3=$(cat "$config_path" | grep openImSdkWsPort | awk -F '[:]' '{print $NF}')
# shellcheck disable=SC2002
logLevel=$(cat "$config_path" | grep remainLogLevel | awk -F '[:]' '{print $NF}')
list_to_string "$list1"
api_ports=($ports_array)
list_to_string "$list2"
ws_ports=($ports_array)
list_to_string "$list3"
sdk_ws_ports=($ports_array)
list_to_string $list4

#echo  "./\"${sdk_server_name}\" -openIM_api_port "${api_ports[0]}" -openIM_ws_port ${ws_ports[0]} -sdk_ws_port ${sdk_ws_ports[0]} -openIM_log_level ${logLevel} >>../logs/openIM.log 2>&1 &"
#exit 0
#Check if the service exists
#If it is exists,kill this process
# shellcheck disable=SC2154
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${sdk_server_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  oldPid=$(ps -afx | grep -w ./"${sdk_server_name}" | grep -v grep | awk '{print $2}')
  kill -9 "${oldPid}"
fi
#Waiting port recycling
sleep 1
# shellcheck disable=SC2164
# shellcheck disable=SC2154
cd "${sdk_server_binary_root}"
nohup ./${sdk_server_name} -openIM_api_port ${api_ports[0]} -openIM_ws_port ${ws_ports[0]} -sdk_ws_port ${sdk_ws_ports[0]} -openIM_log_level ${logLevel} >>../logs/sdk_svr_start_mac.log 2>&1 &

#Check launched service process
sleep 3
# shellcheck disable=SC2009
# shellcheck disable=SC2126

check=$(ps -afx | grep -w ./"${sdk_server_name}" | grep -v grep | wc -l)
allPorts=""
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  allNewPid=$(ps -afx | grep -w ./"${sdk_server_name}" | grep -v grep | awk '{print $2}')
  for i in $allNewPid; do
    #    ports=$(netstat -netulp | grep -w "${i}" | awk '{print $4}' | awk -F '[:]' '{print $NF}')
    ports=$(lsof -a -n -iTCP -sTCP:LISTEN -p"${allNewPid}" -P | awk '{if(FNR!=1) print $9}' | awk -F '[:-]' '{print $2}')
    allPorts=${allPorts}"$ports "
  done
  echo -e "${SKY_BLUE_PREFIX}""SERVICE START SUCCESS ""${COLOR_SUFFIX}"
  # shellcheck disable=SC2086
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${sdk_server_name}${COLOR_SUFFIX}
  # shellcheck disable=SC2086
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allNewPid}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""LISTENING_PORT: ""${COLOR_SUFFIX}""${YELLOW_PREFIX}""${allPorts}""${COLOR_SUFFIX}"
else
  echo -e "${YELLOW_PREFIX}""${sdk_server_name}""${COLOR_SUFFIX}"${RED_PREFIX}"SERVICE START ERROR PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
