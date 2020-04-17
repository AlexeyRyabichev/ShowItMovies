#!/bin/bash
go build
nohup ./ShowItMovies >>out.log 2>&1 &
