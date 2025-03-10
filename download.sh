#!/bin/bash

os=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$os" in
  darwin*) os="macos" ;;
  linux*) os="linux" ;;
  *)
    echo "Sorry, operating system ${os} not supported by this install script."
    exit 1
    ;;
esac

echo "Identified your operating system as ${os}"

file="handles-server-${os}"

echo "Starting download of ${file} from github.com/prompt/handles-server"

curl -L --progress-bar --output handles-server https://github.com/prompt/handles-server/releases/download/v1/handles-server-${os}
chmod +x handles-server

echo "handles-server is downloaded and ready to run."

echo "Start demo of handles server using the in memory provider? (Y/n)"
read response

function ended_demo() {
  printf "\n\n"
  echo "> hopefully the demo was helpful!"
  printf "\n"
}

case "$response" in
  [yY][eE][sS] | [yY] | "")
    echo "> Starting server for domain example.com with handles alice.example.com and bob.example.com"
    sleep 0.5 # An artificial delay for time to read the message
    echo "> Test it out with: curl --header 'Host: alice.example.com' http://localhost:8888/.well-known/atproto-did"
    sleep 1 # An artificial delay for time to read the message
    echo "> exit with ctrl+c once you're done."
    sleep 0.5
    printf "\n\n"
    trap ended_demo INT
    (
      set -x      
      GIN_MODE=release LOG_LEVEL="debug" PORT=8888 DID_PROVIDER="memory" MEMORY_DOMAINS="example.com" MEMORY_DIDS="alice.example.com@did:plc:001,bob.example.com@did:plc:002" ./handles-server
    )
    ;;
  *)
    echo "> No problem, no demo today :)"
    sleep 1
    ;;
esac

echo "> handles-server is downloaded and ready to run."
sleep 0.5
echo "> any questions? visit https://github.com/prompt/handles-server or email sam@handles.net :)"