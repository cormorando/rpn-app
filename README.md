# README

How to run this app:

* Copy rpn-worker.go to $GOPATH/src/rpn-worker
* `cd $GOPATH/src/rpn-worker` -> `go build` -> `./rpn-worker` - this will start the worker
* Go back to app dir -> `virtualenv api` -> `source api/bin/activate` -> `pip install -r requirements.txt` -> `./api.py` - this will start the api
* `cd rpn-app` -> `bin/rails server` - this will start the app
* In web browser go to http://localhost:3000 and use the app
