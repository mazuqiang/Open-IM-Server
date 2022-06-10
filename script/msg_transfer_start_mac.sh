#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg

#Check if the service exists
#If it is exists,kill this process
# shellcheck disable=SC2006
# shellcheck disable=SC2154
check=$(ps -afx | grep -w ./${msg_transfer_name} | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2006
  # shellcheck disable=SC2009
  oldPid=$(ps -afx | grep -w ./"${msg_transfer_name}" | grep -v grep | awk '{print $2}')
  kill -9 "$oldPid"
fi

#Waiting port recycling
sleep 3

# shellcheck disable=SC2164
# shellcheck disable=SC2154
cd "${msg_transfer_binary_root}"
# shellcheck disable=SC2154
# shellcheck disable=SC2004
for ((i = 0; i < ${msg_transfer_service_num}; i++)); do
  # shellcheck disable=SC2086
  nohup ./${msg_transfer_name} >>../logs/openIM.log 2>&1 &
done

#Check launched service process
# shellcheck disable=SC2006
# shellcheck disable=SC2126
check=$(ps -afx | grep -w ./"${msg_transfer_name}" | grep -v grep | wc -l)
if [ "$check" -ge 1 ]; then
  # shellcheck disable=SC2006
  # shellcheck disable=SC2009
  newPid=$(ps -afx | grep -w ./"${msg_transfer_name}" | grep -v grep | awk '{print $2}')
  allPorts=""

  echo -e "${SKY_BLUE_PREFIX}""SERVICE START SUCCESS ""${COLOR_SUFFIX}"
  # shellcheck disable=SC2086
  echo -e "${SKY_BLUE_PREFIX}""SERVICE_NAME: ""${COLOR_SUFFIX}"${YELLOW_PREFIX}${msg_transfer_name}${COLOR_SUFFIX}
  # shellcheck disable=SC2086
  echo -e "${SKY_BLUE_PREFIX}""PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e "${SKY_BLUE_PREFIX}""LISTENING_PORT: ""${COLOR_SUFFIX}""${YELLOW_PREFIX}""${allPorts}""${COLOR_SUFFIX}"
else
  echo -e "${YELLOW_PREFIX}"${msg_transfer_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
