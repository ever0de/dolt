#!/usr/bin/env bats
load $BATS_TEST_DIRNAME/helper/common.bash
load $BATS_TEST_DIRNAME/helper/query-server-common.bash

setup() {
    setup_common
}

teardown() {
    stop_sql_server
    assert_feature_version
    teardown_common
}

@test "sql-load-data: simple load from file into table" {
    cat <<DELIM > 1pk5col-ints.csv
pk||c1||c2||c3||c4||c5
0||1||2||3||4||5
1||1||2||3||4||5
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int, c2 int, c3 int, c4 int, c5 int);
LOAD DATA INFILE '1pk5col-ints.csv' INTO TABLE test CHARACTER SET UTF8MB4 FIELDS TERMINATED BY '||' ESCAPED BY '' LINES TERMINATED BY '\n' IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1,c2,c3,c4,c5" ]
    [ "${lines[1]}" = "0,1,2,3,4,5" ]
    [ "${lines[2]}" = "1,1,2,3,4,5" ]
}

@test "sql-load-data: load into unknown table throws error" {
    run dolt sql -q "LOAD DATA INFILE '1pk5col-ints.csv' INTO TABLE test CHARACTER SET UTF8MB4 FIELDS TERMINATED BY '||' ESCAPED BY '' LINES TERMINATED BY '\n' IGNORE 1 LINES;"
    [ "$status" -eq 1 ]
    [[ "$output" =~  "table not found: test" ]] || false
}

@test "sql-load-data: load with unknown file throws error" {
    skip "Different error msg on windows."
    run dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int, c2 int, c3 int, c4 int, c5 int);
LOAD DATA INFILE 'hello-ints.csv' INTO TABLE test CHARACTER SET UTF8MB4 FIELDS TERMINATED BY '||' ESCAPED BY '' LINES TERMINATED BY '\n' IGNORE 1 LINES;
SQL

    [ "$status" -eq 1 ]
    [[ "$output" =~  "no such file or directory" ]] || false
}

@test "sql-load-data: works with enclosed terms" {
    cat <<DELIM > 1pk5col-ints.csv
pk||c1||c2||c3||c4||c5
"0"||"1"||"2"||"3"||"4"||"5"
"1"||"1"||"2"||"3"||"4"||"5"
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int, c2 int, c3 int, c4 int, c5 int);
LOAD DATA INFILE '1pk5col-ints.csv' INTO TABLE test CHARACTER SET UTF8MB4 FIELDS TERMINATED BY '||' ENCLOSED BY '"'  ESCAPED BY '' LINES TERMINATED BY '\n' IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1,c2,c3,c4,c5" ]
    [ "${lines[1]}" = "0,1,2,3,4,5" ]
    [ "${lines[2]}" = "1,1,2,3,4,5" ]
}

@test "sql-load-data: works with prefixed terms" {
    cat <<DELIM > prefixed.txt
pk
sssHi
sssHello
ignore me
sssYo
DELIM

    dolt sql << SQL
CREATE TABLE test(pk longtext);
LOAD DATA INFILE 'prefixed.txt' INTO TABLE test CHARACTER SET UTF8MB4 LINES STARTING BY 'sss' IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test ORDER BY pk"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk" ]
    [ "${lines[1]}" = "Hello" ]
    [ "${lines[2]}" = "Hi" ]
    [ "${lines[3]}" = "Yo" ]
}

@test "sql-load-data: works when the number of input columns in the file is less than the number of schema columns" {
    cat <<DELIM > 1pk2col-ints.csv
pk,c1
0,1
1,1
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int, c2 int);
LOAD DATA INFILE '1pk2col-ints.csv' INTO TABLE test FIELDS TERMINATED BY ',' IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1,c2" ]
    [ "${lines[1]}" = "0,1," ]
    [ "${lines[2]}" = "1,1," ]
}

@test "sql-load-data: works with fields separated by tabs" {
    skip "This needs to be fixed."
    cat <<DELIM > 1pk2col-ints.csv
pk  c1
0 1
1 1
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int);
LOAD DATA INFILE '1pk2col-ints.csv' INTO TABLE test FIELDS TERMINATED BY '\t' IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1" ]
    [ "${lines[1]}" = "0,1" ]
    [ "${lines[2]}" = "1,1" ]
}

@test "sql-load-data: recognizes certain nulls" {
    cat <<DELIM > 1pk2col-ints.csv
pk
\N
NULL
DELIM

    dolt sql << SQL
CREATE TABLE test(pk longtext);
LOAD DATA INFILE '1pk2col-ints.csv' INTO TABLE test FIELDS IGNORE 1 LINES;
SQL

    run dolt sql -q "select COUNT(*) from test WHERE pk IS NULL"
    [ "$status" -eq 0 ]
    [[ "$output" =~  "2" ]] || false
}

@test "sql-load-data: works when column order is mismatched" {
    skip "This needs to be fixed."

    cat <<DELIM > 1pk2col-ints.csv
pk,c1
"hi","1"
"hello","2"
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int, c1 longtext);
LOAD DATA INFILE '1pk2col-ints.csv' INTO TABLE test FIELDS ENCLOSED BY '"' TERMINATED BY ',' IGNORE 1 LINES (c1,pk);
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1" ]
    [ "${lines[1]}" = "1,hi" ]
    [ "${lines[2]}" = "2,hello" ]
}

@test "sql-load-data: with different column types that uses optionally" {
    skip "This functionality is not present yet."
    cat <<DELIM > complex.csv
1,"a string",100.20
2,"a string containing a , comma",102.20
3,"a string containing a \" quote",102.20
4,"a string containing a \", quote and comma",102.20
DELIM

     dolt sql << SQL
CREATE TABLE test(pk int, c1 longtext, c2 float);
LOAD DATA INFILE 'complex.csv' INTO TABLE test FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"';
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1,c2" ]
    [ "${lines[1]}" = "1,a string,100.20" ]
    [ "${lines[2]}" = "2,a string containing a , comma,100.20" ]
    [ "${lines[3]}" = "3,a string containing a \" quote,100.20" ]
    [ "${lines[4]}" = "4,a string containing a \", quote and comma,100.20" ]
}

@test "sql-load-data: works with escaped columns" {
    skip "This functionality is not present yet."
    cat <<DELIM > escape.txt
"hi"
"\hello"
"Try\\N"
"new\ns"
DELIM

     dolt sql << SQL
CREATE TABLE loadtable(pk longtext);
LOAD DATA INFILE './testdata/test5.txt' INTO TABLE loadtable FIELDS ENCLOSED BY '\"';
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk" ]
    [ "${lines[1]}" = "hi" ]
    [ "${lines[2]}" = "hello" ]
    [ "${lines[3]}" = "TryN" ]
    [ "${lines[4]}" = "new\ns" ]
}

@test "sql-load-data: when the number of input columns in the file is greater than the number of schema columns" {
   skip "This functionality is not present yet."
   cat <<DELIM > 1pk5col-ints.csv
pk||c1||c2||c3||c4||c5
0||1||2||3||4||5||6
1||1||2||3||4||5||6
DELIM

    dolt sql << SQL
CREATE TABLE test(pk int primary key, c1 int, c2 int, c3 int, c4 int, c5 int);
LOAD DATA INFILE '1pk5col-ints.csv' INTO TABLE test
CHARACTER SET UTF8MB4
FIELDS TERMINATED BY '||'
ESCAPED BY ''
LINES TERMINATED BY '\n'
IGNORE 1 LINES;
SQL

    run dolt sql -r csv -q "select * from test"

    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "pk,c1,c2,c3,c4,c5" ]
    [ "${lines[1]}" = "0,1,2,3,4,5" ]
    [ "${lines[2]}" = "1,1,2,3,4,5" ]
}

@test "sql-load-data: load data ignore" {
    cat <<CSV > in.csv
0,0,0
1,1,1
CSV

    dolt sql -q "create table t (pk int primary key, c1 int, c2 int)"
    dolt sql -q "insert into t values (0,0,0)"
    skip "load data ignore not supported"
    run dolt sql <<SQL
load data infile 'in.csv' ignore into table t
fields terminated by ','
lines terminated by '\n'
SQL
    [ $status -eq 0 ]

    run dolt sql -r csv -q "select * from t"
    [ $status -eq 0 ]
    [[ $output =~ 0,0,0 ]] || false
    [[ $output =~ 1,1,1 ]] || false
    
}

@test "sql-load-data: load data replace" {
    cat <<CSV > in.csv
0,0,1
1,1,1
CSV

    dolt sql -q "create table t (pk int primary key, c1 int, c2 int)"
    dolt sql -q "insert into t values (0,0,0)"
    skip "load data replace not supported"
    run dolt sql <<SQL
load data infile 'in.csv' replace into table t
fields terminated by ','
lines terminated by '\n'
SQL
    [ $status -eq 0 ]
    ! [[ $output =~ 0,0,0 ]] || false
    [[ $output =~ 0,0,1 ]] || false
    [[ $output =~ 1,1,1 ]] || false
}


@test "sql-load-data: keyless table" {

}

@test "sql-load-data: test schema with many types" {

}

@test "sql-load-data: test schema with not null constraints" {

}

@test "sql-load-data: test schema with column defaults" {

}

@test "sql-load-data: test schema with check constraints" {

}

@test "sql-load-data: test schema with foreign keys" {

}

@test "sql-load-data: test load data defaults" {

}

@test "sql-load-data: load data local" {
    cat <<CSV > in.csv
0,0,0
1,1,1
CSV
    
    dolt sql -q "create table t (pk int primary key, c1 int, c2 int)"

    run dolt sql <<SQL
load data local infile 'in.csv' into table t
fields terminated by ','
lines terminated by '\n'
SQL
    [ $status -eq 1 ]
    [[ $output =~ "LOCAL supported only in sql-server mode" ]]

    start_sql_server 

    run dolt sql-client -P $PORT -u dolt -q "
load data local infile 'in.csv' into table t
fields terminated by ','
lines terminated by '\n'
"
    [ $status -ne 0 ]
    [[ $output =~ "local_infile needs to be set to 1 to use LOCAL" ]]

    # This should work but does not because of dolt sql-client
    # mysql -e works locally
    run dolt sql-client -P $PORT -u dolt -q "
set global local_infile=1;
load data local infile 'in.csv' into table t
fields terminated by ','
lines terminated by '\n'
"
    [ $status -ne 0 ]
    [[ $output =~ "local file 'in.csv' is not registered" ]] || false
    
    stop_sql_server

    skip "dolt sql-client does not work with local infile but a mysql client does"
    
    run dolt sql -r csv -q "select * from t"
    [ $status -eq 0 ]
    [[ $output =~ "0,0,0" ]] || false
    [[ $output =~ "1,1,1" ]] || false
}

@test "sql-load-data: sql-server mode" {
    cat <<CSV > in.csv
0,0,0
1,1,1
CSV
    dolt sql -q "create table t (pk int primary key, c1 int, c2 int)"

    start_sql_server

    # File not found errors
    run dolt sql-client -P $PORT -u dolt -q "load data infile 'foo.csv' into table t"
    [ $status -ne 0 ]
    [[ $output =~ "no such file or directory" ]] || false
    
    dolt sql-client -P $PORT -u dolt -q "
load data infile 'in.csv' into table t                                    
fields terminated by ','
lines terminated by '\n'
"

    stop_sql_server
    
    run dolt sql -r csv -q "select * from t"
    [ $status -eq 0 ]
    [[ $output =~ "0,0,0" ]] || false
    [[ $output =~ "1,1,1" ]] || false
}
