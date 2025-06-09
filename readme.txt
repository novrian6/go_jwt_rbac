Login as admin
curl -X POST -H "Content-Type: application/json" -d '{"user_id":1}' http://localhost:8080/login

Login as user

curl -X POST -H "Content-Type: application/json" -d '{"user_id":2}' http://localhost:8080/login


access dashboard
curl -H "Authorization: <token>" http://localhost:8080/dashboard
access admin
curl -H "Authorization: <token>" http://localhost:8080/admin