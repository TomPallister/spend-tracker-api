#!/bin/sh
### BEGIN INIT INFO
# Provides:          <SCRIPT_NAME>
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       Script to run asp.net 5 application in the background
### END INIT INFO
 
# Author: Ivan Derevianko aka druss <drussilla7@gmail.com>

WWW_USER=mijitt0m
APPROOT=/home/mijitt0m/godutch-api/ #/home/$WWW_USER/vNext/mvc/HelloMvc
 
PIDFILE=$APPROOT/godutch-api.pid
LOGFILE=$APPROOT/godutch-api.log
 
# fix issue with DNX exception in case of two env vars with the same name but different case
TMP_SAVE_runlevel_VAR=$runlevel
unset runlevel
 
start() {
  if [ -f $PIDFILE ] && kill -0 $(cat $PIDFILE); then
    echo 'Service already running' >&2
    return 1
  fi
  echo 'Starting service...' >&2
  su -c "start-stop-daemon -SbmCv -x /usr/bin/nohup -p \"$PIDFILE\" -d \"$APPROOT\" ./api > \"$LOGFILE\"" $WWW_USER
  echo 'Service started' >&2
}
 
stop() {
  if [ ! -f "$PIDFILE" ] || ! kill -0 $(cat "$PIDFILE"); then
    echo 'Service not running' >&2
    return 1
  fi
  echo 'Stopping service...' >&2
  start-stop-daemon -K -p "$PIDFILE"
  rm -f "$PIDFILE"
  echo 'Service stopped' >&2
}
 
case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart}"
esac
 
export runlevel=$TMP_SAVE_runlevel_VAR