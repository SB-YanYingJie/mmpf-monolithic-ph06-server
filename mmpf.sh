#!/bin/bash

function usage {
    cat <<EOF
Usage:
    $(basename "${0}") [deploy|undeploy] [<options>]

Options:
    --help, -h                                  :display help
    deploy -mmpf                                :deploy
    deploy -mmpf-with-measurement               :deploy and measure
    restart                                     :restart
    undeploy                                    :undeploy
    remove                                      :remove
EOF
}


function deploy_mmpf_with_measurement {
    if [ -r "./ops-utils/sh/container_metrics_exporter.sh" ];then
      bash ./ops-utils/sh/container_metrics_exporter.sh &
      bash ./scripts/start_mmpf_with_redis.sh
    else
      echo "Please grant read permission: chmod +x ./ops-utils/sh/container_metrics_exporter.sh"
      exit 1
    fi
}

function deploy_mmpf {
    bash ./scripts/start_mmpf_with_redis.sh
}

function kill_measurement_process {
    if [ "$(ps -aux | grep "bash ./ops-utils/sh/container_metrics_exporter.sh")" ]; then
      echo "kill measurement process running in the backgroud"
      pid_list="$(ps -aux | grep "bash ./ops-utils/sh/container_metrics_exporter.sh" | awk '{ print $2 }')"
      for pid in "${pid_list[@]}"; do
        echo "[ ${pid} ]"
        kill -9 $pid
      done
    else
      echo "no measurement process running in the backgroud"
    fi
}

function deploy {
    case ${1} in
        -mmpf)
            echo "call kill running measurement process"
            kill_measurement_process
            echo "call remove running containers"
            remove
            echo "call deploy_mmpf"
            deploy_mmpf;;
        -mmpf-with-measurement)
            echo "call kill running measurement process"
            kill_measurement_process
            echo "call remove running containers"
            remove
            echo "call deploy_mmpf_with_measurement"
            deploy_mmpf_with_measurement;;
        *)
            echo "[ERROR] Invalid subcommand '${1}'"
            usage
            exit 1;;
    esac
}

function restart {
    if [ -z "$(docker ps -q --filter "network=monolithic_network")" ]; then
      echo "restart mmpf containers"
      docker start $(docker ps -q --filter network=monolithic_network --filter status=exited)
    else
      echo "no targetting restart containers"
    fi
}

function undeploy {
    if [ -n "$(docker ps -q --filter "network=monolithic_network")" ]; then
      echo "stop mmpf containers"
	    docker stop $(docker ps -q --filter "network=monolithic_network")
    else
      echo "no running mmpf containers"
    fi
}

function remove {
    if [ -n "$(docker ps -aq --filter "network=monolithic_network")" ]; then
        echo "remove mmpf containers"
        docker rm -f $(docker ps -aq --filter "network=monolithic_network")
    else
        echo "no mmpf containers"
    fi
}

case ${1} in
    deploy)
        deploy "${@:2}";;
    undeploy)
        undeploy
        kill_measurement_process
        exit 0;;
    restart)
	restart;;
    remove)
        remove;;
    help|--help|-h)
        usage;;
    *)
        echo "[ERROR] Invalid subcommand '${1}' choose deploy or undeploy"
        usage
        exit 1;;
esac

FAILED_ID=($(docker ps -aq --filter "status=exited" --filter "network=monolithic_network"))
if [ ${#FAILED_ID[@]} -ne 0 ]; then
    echo "Failed to start container."
    for id in "${FAILED_ID[@]}"
    do
        CONTAINER_NAME=$(docker ps -a --filter "id=$id" --format "{{.Names}}")
        echo "failed Container ID: $id Container Name: $CONTAINER_NAME"
        docker logs "$id"
    done
    exit 1;
fi
