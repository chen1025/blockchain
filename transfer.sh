#!/bin/bash

./block send  tx1001 李四 10 班长 "tx1001转李四10"
./block send  tx1001 王五 20 班长 "tx1001转王五20"
./checkBalance.sh

echo "======================"
./block send  王五 李四 2 班长 "王五转李四2"
./block send  王五 李四 3 班长 "王五转李四3"
./block send  王五 tx1001 5 班长 "王五转tx10015"
./checkBalance.sh

echo "======================"
./block send  李四 赵六 14 班长 "李四转赵六14"
./checkBalance.sh

