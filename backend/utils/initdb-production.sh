#!/bin/sh -e

PATH=/go/bin
DB=production
soda create -e $DB || true
soda migrate up -e $DB
