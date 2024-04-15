server=10.10.10.223
echo "deploying api to dev..."
echo ">> making file..."
make build_dev
wait
echo ">> deploying..."
scp -r ./bin/api zmkj@$server:~/woof-prj
wait
ssh zmkj@$server "bash ~/woof-prj/start-api.sh"
wait
echo ">> done!"
exit 0