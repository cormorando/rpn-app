#!./api/bin/python
import json, itertools, logging, socket, time
from flask import abort, Flask, jsonify, make_response, request

app = Flask(__name__)

logFile = logging.FileHandler(time.strftime('%Y-%m-%d') + '.log')
app.logger.addHandler(logFile)
app.logger.setLevel(logging.INFO)

class JSONClient(object):

    def __init__(self, addr):
        self.socket = socket.create_connection(addr)
        self.id_counter = itertools.count()

    def __del__(self):
        self.socket.close()

    def call(self, name, *params):
        request = dict(id = next(self.id_counter), params = list(params), method = name)
        self.socket.sendall(json.dumps(request).encode())

        response = self.socket.recv(4096)
        response = json.loads(response.decode())

        if response.get('error') is not None:
            raise Exception(response.get('error'))

        return response.get('result')

@app.route('/parse', methods = ['POST'])
def api_parse():
    # print request.headers
    # print request.json
    if request.headers['Content-Type'] == 'application/json':
        app.logger.info(time.strftime('[%H:%M:%S]') + ' Parsing request from ' + request.remote_addr)

        jsonData = json.loads(json.dumps(request.json))

        if 'expressions' in jsonData:
            data = {}
            data['results'] = []

            rpc = JSONClient(("localhost", 1234))

            for expression in jsonData['expressions']:
                result = {}
                start = time.time()
                result['value'] = rpc.call("RPCFunc.Parse", expression['value'])
                end = time.time()
                execTime = end - start
                result['time'] = execTime
                data['results'].append(result)

                app.logger.info(time.strftime('[%H:%M:%S]') + ' Request parsed!\tInput: ' + expression['value'] + '\tTime: ' + str(execTime))

            data['count'] = len(data['results'])

            return jsonify(data)
        else:
            message = {
                'status': 400,
                'message': 'Bad Request',
            }
            response = jsonify(message)
            response.status_code = 400

            return response
    else:
        message = {
            'status': 415,
            'message': 'Unsupported Media Type',
        }
        response = jsonify(message)
        response.status_code = 415

        return response

@app.errorhandler(404)
def not_found(error):
    return make_response(jsonify({'error': 'Not found'}), 404)

if __name__ == '__main__':
    app.run()
