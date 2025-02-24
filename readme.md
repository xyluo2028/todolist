# /register 
    curl -X POST \
        http://localhost:7071/register \
        -H 'Content-Type: application/json' \
        -d '{"username":"test","password":"test123"}'

# /welcome
    curl -X GET -u test:test123 http://localhost:7071/welcome

# /addTask
    curl -X GET -u test:test123 http://localhost:7071/addTask?q=go%20travel

# /printTasks
    curl -X GET -u test:test123 http://localhost:7071/printTasks

# /removeTask
    curl -X GET -u test:test123 http://localhost:7071/removeTask?key=key_to_remove

# /updateTask
    curl -X GET -u test:test123 "http://localhost:7071/updateTask?q=go%20to%20concerts&key=key_to_update"

# /deactivate
    curl -X DELETE -u test:test123 http://localhost:7071/deactivate

##  nohup start 
    nohup go run cmd/server/main.go > logs/todolist.log 2>&1 & 
##  stop completely
    pkill -TERM -P <PID>