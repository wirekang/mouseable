#!/bin/bash
cd "$(dirname "$(dirname "$0")}")" || return
makensis.exe nsi.nsi