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
  "
}

create() {
    MEETING_ID=$(uuid)
    ATTENDEE_PW=$(uuid)
    MODERATOR_PW=$(uuid)
    echo "Enter the room name:"
    read ROOM_NAME
    PARAMETERS="meetingID=$MEETING_ID&attendeePW=$ATTENDEE_PW&moderatorPW=$MODERATOR_PW&name=$ROOM_NAME"
    CHEKSUM=$(sha1 "create$PARAMETERS$SECRET")
    http GET :8090/bigbluebutton/api/create meetingID=="$MEETING_ID" attendeePW=="$ATTENDEE_PW" moderatorPW=="$MODERATOR_PW" name=="$ROOM_NAME" checksum=="$CHEKSUM"
}

for param in "$@"
do
  case $param in
    -c | --create)
      create
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
