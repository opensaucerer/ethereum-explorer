API_PORT=5000
INDEX_PORT=2500

start:
	chmod +x ./start.sh
	./start.sh

services:
	chmod +x ./services.sh
	./services.sh

servicesbg:
	chmod +x ./servicesbg.sh
	./servicesbg.sh

loadtest:
	make servicesbg
	chmod +x ./loadtest.sh
	./loadtest.sh

database:
	sudo docker-compose up -d database

stop:
	sudo docker-compose down
	rm -rf ./build
	kill $$(lsof -t -i:${API_PORT})
	kill $$(lsof -t -i:${INDEX_PORT})