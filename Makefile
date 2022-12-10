.PHONY: syncgommon
URL="github.com/halilylm/gommon@v1.1.19"
syncgommon:
	cd orders && go get ${URL}
	cd payments && go get ${URL}
	cd tickets && go get ${URL}