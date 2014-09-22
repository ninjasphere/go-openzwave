(
sleep 2
open http://localhost:8080
) &

cd ./ozwcp/darwin &&
exec ./ozwcp -p 8080 

