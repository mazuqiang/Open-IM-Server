#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
# shellcheck disable=SC2154
# shellcheck disable=SC2002
switch=$(cat "$config_path" | grep demoswitch |awk -F '[:]' '{print $NF}')

if [ ${switch} != "true" ]; then
      echo -e "${YELLOW_PREFIX}"" demo service switch is false not start demo "${COLOR_SUFFIX}
      exit 0
fi
# shellcheck disable=SC2086
# shellcheck disable=SC2002
list1=$(cat $config_path | grep openImDemoPort | awk -F '[:]' '{print $NF}')
list_to_string "$list1"
api_ports=($ports_array)

#Check if the service exists
#If it is exists,kill this process
# shellcheck disable=SC2154
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${demo_server_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  oldPid=$(ps -afx | grep -w ./"${demo_server_name}" | grep -v grep | awk '{print $2}')
  kill -9 "$oldPid"
fi
#Waiting port recycling
sleep 1
# shellcheck disable=SC2164
# shellcheck disable=SC2154
cd "${demo_server_binary_root}"

for ((i = 0; i < ${#api_ports[@]}; i++)); do
  nohup ./"${demo_server_name}" -port "${api_ports[$i]}" >>../logs/openIM.log 2>&1 &
done

sleep 120
#Check launched service process
# shellcheck disable=SC2126
# shellcheck disable=SC2009
check=$(ps -afx | grep -w ./"${demo_server_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  newPid=$(ps -afx | grep -w ./"${demo_server_name}" | grep -v grep | awk '{print $2}')
#  ports=$(netstat -netulp | grep -w "${newPid}" | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  ports=$(lsof -a -n -iTCP -p"${newPid}" -sTCP:LISTEN  -P |awk '{if(FNR!=1) print $9}'|awk -F '[:]' '{ifprint $2}')
  allPorts=""

  for i in $ports; do
    allPorts=${allPorts}"$i "
  done
  echo -e "${SKY_BLUE_PREFIX}""SERVICE START SUCCESS "${COLOR_SUFFIX}
  # shellcheck disable=SC2086
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${demo_server_name}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${demo_server_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi