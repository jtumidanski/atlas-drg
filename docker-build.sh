if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache --tag atlas-drg:latest .
else
   docker build --tag atlas-drg:latest .
fi
