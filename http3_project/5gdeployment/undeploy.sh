#!/bin/bash

helm delete central -n oai
helm delete scp1 -n oai
helm delete scp2 -n oai

for ((c=10;c<=$1+9;c++))
do
	helm delete nrf$c -n oai
	helm delete udr$c -n oai
	helm delete udm$c -n oai
	helm delete ausf$c -n oai
	helm delete amf$c -n oai
	helm delete smf$c -n oai
	helm delete upf$c -n oai
    helm delete gnb$c -n oai
    kubectl delete deployment -n oai oai-dnn$c
done

