#!/usr/bin/env node
'use strict';

const fs = require('fs');

const PROTO_PATH = __dirname + '/helloworld/helloworld.proto';

const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    { keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true
    });

const api = grpc.loadPackageDefinition(packageDefinition);

colors(api);

const client = new api.helloworld.Greeter(
  'localhost:50051',
  // grpc.credentials.createInsecure(),
  grpc.credentials.createSsl(
    // fs.readFileSync('./openssl-certs/ca.crt'),
    fs.readFileSync('./cfssl-certs/ca.pem'),
    // ca, key, cert,
    fs.readFileSync('./cfssl-certs/client-key.pem'),
    // fs.readFileSync('./openssl-certs/client.key'),
    fs.readFileSync('./cfssl-certs/client.pem'),
    // fs.readFileSync('./openssl-certs/client.crt'),
    // { checkServerIdentity() {} },
  ),
  // { 'grpc.ssl_target_name_override' : 'example', 'grpc.default_authority': 'localhost', },
);

colors(client);

client.SayHello({ name: 'abc' }, (err, resp) => {
  console.log(new Date, 'ERROR:', err, resp);
});

client.SayHelloAgain({ name: 'abc2' }, (err, resp) => {
  console.log(new Date, 'ERROR:', err, resp);
});

function colors(data) { console.dir(data, { depth: null, colors: true }); }

