#!/usr/bin/env sh

pgrep hot_reloaded_app
pgrep hot_reloaded_app |xargs -I {{}} kill {{}}
