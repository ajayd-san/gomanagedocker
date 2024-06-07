#!/bin/bash

for i in {1..100}; do 
    gnome-terminal --full-screen --title "testing" -- bash -c "go run main.go -debug 2>>error.log" &

    sleep 2

    xdotool search --name testing type q
done 

