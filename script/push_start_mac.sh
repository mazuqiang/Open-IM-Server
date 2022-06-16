#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

# shellcheck disable=SC2154
# shellcheck disable=SC2002
list1=$(cat "$config_path" | grep openImPushPort | awk -F '[:]' '{print $NF}')
list_to_string "$list1"
rpc_ports=($ports_array)


#Check if the service exists
#If it is exists,kill this process
# shellcheck disable=SC2154
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${push_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2086
  oldPid=$(ps -afx | grep -w ./${push_name} | grep -v grep | awk '{print $2}')
  kill -9 "$oldPid"
fi
#Waiting port recycling
sleep 1
# shellcheck disable=SC2164
# shellcheck disable=SC2154
cd "${push_binary_root}"

for ((i = 0; i < ${#rpc_ports[@]}; i++)); do
  nohup ./"${push_name}" -port "${rpc_ports[$i]}" >>../logs/push_start_mac.log 2>&1 &
done

sleep 3
#Check launched service process
# shellcheck disable=SC2009
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${push_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  newPid=$(ps -afx | grep -w ./"${push_name}" | grep -v grep | awk '{print $2}')
  # shellcheck disable=SC2086

  #  echo "netstat -netulp | grep -w ${newPid} | awk '{print \$4}' | awk -F '[:]' '{print \$NF}'"
  #  ports=$(netstat -netulp | grep -w ${newPid} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  ports=$(lsof -a -n -iTCP -sTCP:LISTEN -p"${newPid}" -P | awk '{if(FNR!=1) print $9}' | awk -F '[:-]' '{print $2}')
  allPorts=""

  for i in $ports; do
    allPorts=${allPorts}"$i "
  done
  echo -e "${SKY_BLUE_PREFIX}""SERVICE START SUCCESS ""${COLOR_SUFFIX}"
  # shellcheck disable=SC2086
  echo -e "${SKY_BLUE_PREFIX}""SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${push_name}${COLOR_SUFFIX}
  # shellcheck disable=SC2086
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""LISTENING_PORT: ""${COLOR_SUFFIX}""${YELLOW_PREFIX}""${allPorts}""${COLOR_SUFFIX}"
else
  # shellcheck disable=SC2027
  echo -e "${YELLOW_PREFIX}""${push_name}""${COLOR_SUFFIX}""${RED_PREFIX}""SERVICE START ERROR, PLEASE CHECK openIM.log""${COLOR_SUFFIX}"
fi
