#!/bin/bash
SECRET=$(grep  -Po "secret: (.*)" $(dirname "$0")/../config.yml | cut -d: -f2 | xargs)

uuid() {
    cat /proc/sys/kernel/random/uuid
}

sha1() {
    echo -n "$1" | sha1sum | cut -d ' ' -f 1
}

usage() {
  echo "
    Usage: api.sh [OPTION]

    Call b3lb api functions
    -c, --create         Call create api. It prompts for the room name parameter. Meeting identifier, attendee password and moderator password are generated automatically.
    -j, --join           Call join api. It prompts for the meeting ID parameter, the full name and the password.
  "
}

function urlencode() {
	local LANG=C
	for ((i=0;i<${#1};i++)); do
		if [[ ${1:$i:1} =~ ^[a-zA-Z0-9\.\~\_\-]$ ]]; then
			printf "${1:$i:1}"
		else
			printf '%%%02X' "'${1:$i:1}"
		fi
	done
}

join() {
    echo "Enter meeting ID:"
    read MEETING_ID
    echo "Enter your name:"
    read FULL_NAME
    echo "Enter your password:"
    read PASSWORD
    PARAMETERS="meetingID=$MEETING_ID&fullName=$(urlencode "$FULL_NAME")&password=$PASSWORD"
    CHECKSUM=$(sha1 "join$PARAMETERS$SECRET")
    curl -s -G http://localhost:8090/bigbluebutton/api/join --data-urlencode "meetingID=$MEETING_ID" --data-urlencode "fullName=$FULL_NAME" --data-urlencode "password=$PASSWORD" --data-urlencode "checksum=$CHECKSUM" | tidy -xml -i -q -
}

create() {
    MEETING_ID=$(uuid)
    ATTENDEE_PW=$(uuid)
    MODERATOR_PW=$(uuid)
    echo "Enter the room name:"
    read ROOM_NAME
    PARAMETERS="meetingID=$MEETING_ID&attendeePW=$ATTENDEE_PW&moderatorPW=$MODERATOR_PW&name=$(urlencode "$ROOM_NAME")"
    CHECKSUM=$(sha1 "create$PARAMETERS$SECRET")
    curl -s -G "http://localhost:8090/bigbluebutton/api/create" -H "Accept: application/xml" --data-urlencode "meetingID=$MEETING_ID" --data-urlencode "attendeePW=$ATTENDEE_PW" --data-urlencode "moderatorPW=$MODERATOR_PW" --data-urlencode "name=$ROOM_NAME" --data-urlencode "checksum=$CHECKSUM" | tidy -xml -i -q -
}

for param in "$@"
do
  case $param in
    -c | --create)
      create
      ;;
    -j | --join)
      join
      ;;
    -h | --help)
      usage
      exit 1
      ;;
    *)
      echo "Invalid argument : $param"
  esac
  if [ ! $? -eq 0 ]; then
    exit 1
  fi
done
