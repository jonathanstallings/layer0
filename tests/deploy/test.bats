#!/usr/bin/env bats
load ../common/common

@test "deploy create guestbook1" {
    l0 deploy create ./deploy/Guestbook.dockerrun.aws.json guestbook1
}

@test "deploy list" {
    l0 deploy list
}

@test "deploy list --all" {
    l0 deploy list --all
}

@test "deploy get guestbook1" {
    l0 deploy get guestbook1
}

@test "deploy get guest*" {
    l0 deploy get guest\*
}

@test "deploy delete guestbook1" {
    l0 deploy delete guestbook1
}
