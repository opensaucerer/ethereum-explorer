sudo docker-compose up -d database indexer api --force-recreate --build
echo "--------------------------------------------------------------------------------------"
echo "API Service started at http://0.0.0.0:5000"
echo "--------------------------------------------------------------------------------------"