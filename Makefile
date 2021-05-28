
bin/data-source:
	go build -o bin/data-source cmd/mock-data-source/mock-data.go

bin/sfn-load-prediction:
	go build -o bin/sfn-load-prediction cmd/stream-fn-realtime-load-prediction/load-prediction.go

bin/sfn-outliers:
	go build -o bin/sfn-outliers cmd/stream-fn-realtime-outliers/outliers.go

build: bin/data-source bin/sfn-load-prediction bin/sfn-outliers

clean:
	@rm -rf bin