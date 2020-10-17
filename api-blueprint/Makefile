PWD := ${CURDIR}

.PHONY: render
render: 
	docker run -it --rm -v $(PWD):/doc quay.io/bukalapak/snowboard html -o rendered Polls.apib
	docker build -t aglio . && docker run -v $(PWD):/tmp -t aglio -i /tmp/Polls.apib --theme-template triple -o /tmp/rendered/Polls.aglio.html
