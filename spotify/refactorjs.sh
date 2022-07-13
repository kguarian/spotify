#!/usr/bin/bash

echo Building Web Package
cd ../server_root/ui/
npm run build > .compiler_npm_output;
if [ $? -eq 0 ]; then
   echo OK
else
   echo FAIL
   exit 1
fi
cd build;
sed -i 's/="\//="\/ui\//g' *.*;
sed -i 's/https\:\/\/192\.168\.1\.119\:8081\//\$1/g ../src/components/api_consts.js'
# cd ../../dashboard/build;
# sed -i 's/="\//="\/dashboard\//g' *.*;
echo Done Building and Refactoring UI Folder