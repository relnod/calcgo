#!/bin/bash

echo "mode: count" > profile.cov

for pkg in $(go list ./...); do
    go test -covermode=count -coverprofile=profile_temp.cov $pkg
    tail -n +2 profile_temp.cov >> profile.cov
done

rm profile_temp.cov
