#!/bin/bash

usage() {
    echo "Helper script that runs benchmarks on two different git branches"
    echo
    echo "Usage: [-h | -n]"
    echo
    echo "Options:"
    echo "    -h Show this message and exit."
    echo "    -n Specify number of benchmark runs."
}

main() {
    local current=$(git branch | grep "*" | cut -d ' ' -f2)

    if [[ $current == "master" ]]; then
        bench master
        benchstat master.txt
        rm master.txt
        exit
    fi

    bench current
    git checkout master > /dev/null 2>&1

    bench master
    git checkout - > /dev/null 2>&1

    benchstat master.txt current.txt
    rm master.txt current.txt
}

bench() {
    local branch=$1
    local file="$branch.txt"

    if [[ $v = true ]]; then
        echo -n "Benchmarking $branch"
    fi
    bench_command > $file
    if [[ $v = true ]]; then
    echo -n "."
    fi
    if [[ $n > 1 ]]; then
        for ((i=1;i<$n;i++)); do
            bench_command >> $file
            if [[ $v = true ]]; then
                echo -n "."
            fi
        done
    fi

    if [[ $v = true ]]; then
        echo
    fi
}

bench_command() {
    go test ./... -bench=$b -run=^$ -benchmem
}

n=5
v=false
b="."

while getopts "lhn:v" opt; do
    case ${opt} in
        v)
            v=true
            ;;
        n)
            n="${OPTARG}"
            ;;
        h)
            usage
            exit 0
            ;;
        *)
            echo ""
            usage
            exit 1
            ;;
    esac
done

main
