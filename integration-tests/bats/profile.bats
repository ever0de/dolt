#!/usr/bin/env bats
load $BATS_TEST_DIRNAME/helper/common.bash
load $BATS_TEST_DIRNAME/helper/query-server-common.bash

make_repo() {
  mkdir "$1"
  cd "$1"
  dolt init
  dolt sql -q "create table $1_tbl (id int)"
  dolt sql <<SQL
CREATE TABLE table1 (pk int PRIMARY KEY);
CREATE TABLE table2 (pk int PRIMARY KEY);
INSERT INTO dolt_ignore VALUES ('generated_*', 1);
SQL
  dolt add -A && dolt commit -m "tables table1, table2"
  cd ..
}

setup() {
    setup_no_dolt_init
    make_repo defaultDB
    make_repo altDB

    unset DOLT_CLI_PASSWORD
    unset DOLT_SILENCE_USER_REQ_FOR_TESTING
}

teardown() {
    stop_sql_server 1
    teardown_common
}

@test "profile: --profile exists and isn't empty" {
    cd defaultDB
    dolt sql -q "create table test (pk int primary key)"
    dolt sql -q "insert into test values (999)"
    dolt add test
    dolt commit -m "insert initial value into test"
    cd ..

    dolt profile add --use-db defaultDB defaultTest
    run dolt --profile defaultTest sql -q "select * from test"
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "999" ]] || false
}

@test "profile: --profile doesn't exist" {
    dolt profile add --use-db defaultDB defaultTest

    run dolt --profile nonExistentProfile sql -q "select * from altDB_tbl"
    [ "$status" -eq 1 ] || false
    [[ "$output" =~ "Failure to parse arguments: profile nonExistentProfile not found" ]] || false
}

@test "profile: additional flag gets used" {
    cd altDB
    dolt sql -q "create table test (pk int primary key)"
    dolt sql -q "insert into test values (999)"
    dolt add test
    dolt commit -m "insert initial value into test"
    cd ..

    dolt profile add --user dolt --password "" userProfile
    run dolt --profile userProfile --use-db altDB sql -q "select * from test"
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "999" ]] || false
}

@test "profile: duplicate flag overrides correctly" {
    cd altDB
    dolt sql -q "create table test (pk int primary key)"
    dolt sql -q "insert into test values (999)"
    dolt add test
    dolt commit -m "insert initial value into test"
    cd ..

    dolt profile add --use-db defaultDB defaultTest
    run dolt --profile defaultTest --use-db altDB sql -q "select * from test"
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "999" ]] || false
}

@test "profile: duplicate flag with non-duplicate flags in profile overrides correctly" {
    cd altDB
    dolt sql -q "create table test (pk int primary key)"
    dolt sql -q "insert into test values (999)"
    dolt add test
    dolt commit -m "insert initial value into test"
    cd ..

    start_sql_server altDb
    dolt --user dolt --password "" sql -q "CREATE USER 'steph' IDENTIFIED BY 'pass'; GRANT ALL PRIVILEGES ON altDB.* TO 'steph' WITH GRANT OPTION;";
    dolt profile add --user "not-steph" --password "pass" --use-db altDB userWithDBProfile

    run dolt --profile userWithDBProfile --user steph sql -q "select * from test"
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "999" ]] || false
}

@test "profile: duplicate flag with non-duplicate flags overrides correctly" {
    cd altDB
    dolt sql -q "create table test (pk int primary key)"
    dolt sql -q "insert into test values (999)"
    dolt add test
    dolt commit -m "insert initial value into test"
    cd ..

    start_sql_server altDb
    dolt --user dolt --password "" sql -q "CREATE USER 'steph' IDENTIFIED BY 'pass'; GRANT ALL PRIVILEGES ON altDB.* TO 'steph' WITH GRANT OPTION;";
    dolt profile add --user "not-steph" --password "pass" userProfile

    run dolt --profile userProfile --user steph --use-db altDB sql -q "select * from test"
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "999" ]] || false
}

@test "profile: dolt profile add adds a profile" {
    run dolt profile add --use-db altDB altTest
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest" ]] || false
    [[ "$output" =~ "altDB" ]] || false
}

@test "profile: dolt profile add overwrites an existing profile" {
    run dolt profile add --user "not-steph" --password "password123" userProfile
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "userProfile" ]] || false
    [[ "$output" =~ "not-steph" ]] || false
    [[ "$output" =~ "password123" ]] || false

    run dolt profile add --user "steph" --password "password123" userProfile
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "userProfile" ]] || false
    [[ "$output" =~ "steph" ]] || false
    [[ ! "$output" =~ "not-steph" ]] || false
    [[ "$output" =~ "password123" ]] || false
}

@test "profile: dolt profile add overwrites an existing profile with different flags" {
    run dolt profile add --user "steph" --password "password123" userProfile
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "userProfile" ]] || false
    [[ "$output" =~ "steph" ]] || false
    [[ "$output" =~ "password123" ]] || false

    run dolt profile add --user "dolt" --use-db defaultDB userProfile
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "userProfile" ]] || false
    [[ ! "$output" =~ "steph" ]] || false
    [[ ! "$output" =~ "password123" ]] || false
    [[ "$output" =~ "dolt" ]] || false
    [[ "$output" =~ "defaultDB" ]] || false
}

@test "profile: dolt profile add adds a profile with existing profiles" {
    dolt profile add --use-db altDB altTest
    dolt profile add --use-db defaultDB defaultTest

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest:" ]] || false
    [[ "$output" =~ "defaultTest:" ]] || false
}

@test "profile: dolt profile add with multiple names errors" {
    run dolt profile add --use-db altDB altTest altTest2
    [ "$status" -eq 1 ] || false
    [[ "$output" =~ "Only one profile name can be specified" ]] || false
}

@test "profile: dolt profile add locks global config with 0600" {
    run dolt profile add --use-db altDB altTest
    [ "$status" -eq 0 ] || false

    run stat "$BATS_TMPDIR/config-$$/.dolt/config_global.json"
    [[ "$output" =~ "-rw-------" ]] || false
}

@test "profile: dolt profile remove removes a profile" {
    dolt profile add --use-db altDB altTest
    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest" ]] || false

    run dolt profile remove altTest
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ ! "$output" =~ "altTest:" ]] || false
}

@test "profile: dolt profile remove leaves existing profiles" {
    dolt profile add --use-db altDB altTest
    dolt profile add --use-db defaultDB defaultTest
    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest" ]] || false
    [[ "$output" =~ "defaultTest" ]] || false

    run dolt profile remove altTest
    [ "$status" -eq 0 ] || false

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ ! "$output" =~ "altTest:" ]] || false
    [[ "$output" =~ "defaultTest" ]] || false
}

@test "profile: dolt profile remove with no existing profiles errors" {
    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" = "" ]] || false

    run dolt profile remove altTest
    [ "$status" -eq 1 ] || false
    [[ "$output" =~ "no existing profiles" ]] || false
}

@test "profile: dolt profile remove with non-existent profile errors" {
    dolt profile add --use-db altDB altTest
    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest" ]] || false

    run dolt profile remove defaultTest
    [ "$status" -eq 1 ] || false
    [[ "$output" =~ "profile defaultTest does not exist" ]] || false
}

@test "profile: dolt profile remove with multiple names errors" {
    dolt profile add --use-db altDB altTest
    run dolt profile remove altTest altTest2
    [ "$status" -eq 1 ] || false
    [[ "$output" =~ "Only one profile name can be specified" ]] || false
}

@test "profile: dolt profile remove locks global config with 0600" {
    dolt profile add --use-db altDB altTest
    run dolt profile remove altTest
    [ "$status" -eq 0 ] || false

    run stat "$BATS_TMPDIR/config-$$/.dolt/config_global.json"
    [[ "$output" =~ "-rw-------" ]] || false
}

@test "profile: dolt profile lists all profiles" {
    dolt profile add --use-db altDB altTest
    dolt profile add --use-db defaultDB -u "steph" --password "pass" defaultTest

    run dolt profile
    [ "$status" -eq 0 ] || false
    [[ "$output" =~ "altTest:" ]] || false
    [[ "$output" =~ "use-db: altDB" ]] || false
    [[ "$output" =~ "defaultTest:" ]] || false
    [[ "$output" =~ "user: steph" ]] || false
    [[ "$output" =~ "password: pass" ]] || false
    [[ "$output" =~ "use-db: defaultDB" ]] || false
}
