#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000

# shellcheck disable=SC2154
# shellcheck disable=SC2002
list1=$(cat "$config_path" | grep openImMessageGatewayPort | awk -F '[:]' '{print $NF}')
# shellcheck disable=SC2002
list2=$(cat "$config_path" | grep openImWsPort | awk -F '[:]' '{print $NF}')
list_to_string "$list1"
rpc_ports=($ports_array)
list_to_string "$list2"
ws_ports=($ports_array)
if [ ${#rpc_ports[@]} -ne ${#ws_ports[@]} ]; then

  echo -e "${RED_PREFIX}""ws_ports does not match push_rpc_ports in quantity!!!""${COLOR_SUFFIX}"
  # shellcheck disable=SC2242
  exit -1

fi
#Check if the service exists
#If it is exists,kill this process
# shellcheck disable=SC2126
# shellcheck disable=SC2154
# shellcheck disable=SC2009
check=$(ps -afx | grep -w ./"${msg_gateway_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  oldPid=$(ps -afx | grep -w ./"${msg_gateway_name}" | grep -v grep | awk '{print $2}')
#  echo "${oldPid}"
  kill -9 "${oldPid}"
fi
#Waiting port recycling
sleep 1
# shellcheck disable=SC2164
# shellcheck disable=SC2154
cd "${msg_gateway_binary_root}"
for ((i = 0; i < ${#ws_ports[@]}; i++)); do
#   shellcheck disable=SC2086
  nohup ./"${msg_gateway_name}" -rpc_port ${rpc_ports[$i]} -ws_port "${ws_ports[$i]}" >>../logs/msg_gateway_start_mac.log 2>&1 &
done

#Check launched service process
sleep 3
# shellcheck disable=SC2009
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${msg_gateway_name}" | grep -v grep | wc -l)
allPorts=""
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2009
  allNewPid=$(ps -afx | grep -w ./"${msg_gateway_name}" | grep -v grep | awk '{print $2}')
  for i in $allNewPid; do
#    echo "netstat -netulp | grep -w "${i}" | awk '{print $4}' | awk -F '[:]' '{print $NF}'"
#    ports=$(netstat -netulp | grep -w "${i}" | awk '{print $4}' | awk -F '[:]' '{print $NF}')
#    echo "lsof -a -n -iTCP -p"${i}" -sTCP:LISTEN |awk '{if(FNR!=1) print \$9}'|awk -F '[:]' '{if(\$2!=\"scp-config\") print \$2}'"
    ports=$(lsof -a -n -iTCP -p"${i}" -sTCP:LISTEN -P|awk '{if(FNR!=1) print $9}'|awk -F '[:]' '{print $2}')
    allPorts=${allPorts}"$ports "
  done
  echo -e "${SKY_BLUE_PREFIX}""SERVICE START SUCCESS""${COLOR_SUFFIX}"
  # shellcheck disable=SC2086
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${msg_gateway_name}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""PID: ""${COLOR_SUFFIX}""${YELLOW_PREFIX}""${allNewPid}""${COLOR_SUFFIX}"
  # shellcheck disable=SC2086
  echo -e "${SKY_BLUE_PREFIX}""LISTENING_PORT: ""${COLOR_SUFFIX}"${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  # shellcheck disable=SC2027
  echo -e "${YELLOW_PREFIX}""${msg_gateway_name}""${COLOR_SUFFIX}""${RED_PREFIX}""SERVICE START ERROR, PLEASE CHECK openIM.log""${COLOR_SUFFIX}"
fi
