#!/bin/bash

eval $(egrep -v '^#' variables.env | xargs) ./main
