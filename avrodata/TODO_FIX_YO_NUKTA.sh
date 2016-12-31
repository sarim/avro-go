#!/bin/sh

NUKTA="়"
DA="ড"
DHA="ঢ"
ZA="য"

DA_R="ড়"
DHA_R="ঢ়"
ZYA="য়"

#  one char + nukta
egrep -o .$NUKTA . -R

gsed -i "s/$DA$NUKTA/$DA_R/g" ./avroclassic.json
gsed -i "s/$DHA$NUKTA/$DHA_R/g" ./avroclassic.json
gsed -i "s/$ZA$NUKTA/$ZYA/g" ./avroclassic.json