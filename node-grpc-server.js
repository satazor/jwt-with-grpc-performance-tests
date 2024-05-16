
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import { createSigner, createVerifier } from 'fast-jwt';
import { createPrivateKey, createPublicKey } from 'node:crypto';
import fs from 'node:fs';

const emails = JSON.parse(fs.readFileSync('./emails.json').toString());
const emailsLength = emails.length;
let emailsIdx = 0;

const createHS256 = () => {
  const jwtSecret = 'bmtXKKngXH1HRdshrI7LkJxmyNZyDN1f';

  const signSync = createSigner({ key: jwtSecret, algorithm: 'HS256' });
  const verifySync = createVerifier({ key: jwtSecret, algorithms: ['HS256'] });

  return { signSync, verifySync };
}

const createES256 = () => {
  const publicKey = JSON.parse('{"x": "7-INQ150R-MCWlj5X_wyGLRIRYAA-o8NakJiUq7gOGg", "y": "dM-GsyJvdDOuALE3l-U9lPL8V3gY_5BPjLH539yTdKU", "alg": "ES256", "crv": "P-256", "kid": "cdd2969c-7e49-4a46-bcbe-e8bbdf74c7f3", "kty": "EC"}');
  const privateKey = JSON.parse('{"alg":"ES256","crv":"P-256","d":"h-UIda1elff-qw81gsSQakyzOv8Dozv5RcQqFIV6R1Y","kid":"cdd2969c-7e49-4a46-bcbe-e8bbdf74c7f3","kty":"EC","x":"7-INQ150R-MCWlj5X_wyGLRIRYAA-o8NakJiUq7gOGg","y":"dM-GsyJvdDOuALE3l-U9lPL8V3gY_5BPjLH539yTdKU"}');

  const publicKeyPem = createPublicKey({ format: 'jwk', key: publicKey }).export({ format: 'pem', type: 'spki' });
  const privateKeyPem = createPrivateKey({ format: 'jwk', key: privateKey }).export({ format: 'pem', type: 'pkcs8' });

  const signSync = createSigner({ key: privateKeyPem, algorithm: privateKey.alg });
  const verifySync = createVerifier({ key: publicKeyPem, algorithms: [publicKey.alg] });

  return { signSync, verifySync };
}

const algorithms = {
  HS256: createHS256(),
  ES256: createES256()
}

const generateJwt = async ({ request }, callback) => {
  const { algorithm } = request;

  const nowSeconds = Math.floor(Date.now() / 1000);
  const sub = emails[emailsIdx];

  if (emailsIdx >= emailsLength) {
    emailsIdx = 0;
  }

  const jwt = await algorithms[algorithm].signSync({
    sub,
    iat: nowSeconds,
    exp: nowSeconds + 7200,
  });

  const claims = await algorithms[algorithm].verifySync(jwt);

  if (claims.sub !== sub) {
    throw new Error('sub mismatch');
  }

  emailsIdx++;

  if (emailsIdx >= emailsLength) {
    emailsIdx = 0;
  }

  callback(null, { jwt, sub });
};

// ----------------------------------------------------------------

const packageDefinition = protoLoader.loadSync('./jwt.proto', {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);

const server = new grpc.Server();

server.addService(protoDescriptor.jwt.services.JwtService.service, {
  generateJwt,
});

server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {
  console.log('Serving on 50051');
});
