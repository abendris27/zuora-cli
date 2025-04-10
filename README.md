# zuora-cli
Zuora command line interface

go build -o zuo

zuo
Usage:
  zuo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  ngrok       run ngrok Server for Zuora
  org         Manage Zuora organizations
  data		  Manage data in bulk

Flags:
  -h, --help   help for zuo

zuo data
Available Commands:
  insert  	  insert a set of element by csv file input
  update      update a set of element by csv file input
  delete      delete a set of element by csv file input
  dataquery	  execute a dataquery

zuo data insert | zuo data update | zuo data delete

Flags:
 -o specify the org where to update
 -f specify the file to insert
 -t specify the object name to insert (Enum)

zuo data dataquery
Flags:
 -o specify the org where to update
 -s query string
 -q query file
 --output outputFile
