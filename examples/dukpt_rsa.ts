import { check } from "k6";
import { derivationOfInitialKey, deriveCurrentTransactionKey } from "k6/x/dukpt/des";
import secrets from 'k6/secrets';

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export async function rsaWrapwithIk(baseKey: Uint8Array, ksn: Uint8Array, publicKey: CryptoKey): Promise<ArrayBuffer> {
  const ik = derivationOfInitialKey(baseKey, ksn)

  const encryptedKey = await crypto.subtle.encrypt(
    {
      name: 'RSA-OAEP',
    },
    publicKey,
    new Uint8Array(ik)
  );
  return encryptedKey;
}

// this example uses well-known BDK, but you should replace it with your own.
const bdk = Uint8Array.fromBase64(await secrets.source("file").get("bdk"));
const keyPair: CryptoKeyPair = await crypto.subtle.generateKey(
  {
    name: 'RSA-OAEP',
    modulusLength: 2048,
    publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
    hash: 'SHA-256',
  },
  true,
  ['encrypt', 'decrypt']
);

/**
 * run this example with the following command:
 * `./k6 run --secret-source=file=(yourAbsolutePath)/examples/example.secrets.txt examples/dukpt_rsa.ts`
*/
async function test(){
  const ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  const wrappedik = await rsaWrapwithIk(bdk, ksn, keyPair.publicKey);
  const ik = await crypto.subtle.decrypt(
    {
      name: 'RSA-OAEP',
    },
    keyPair.privateKey,
    wrappedik
  );
  const ck = deriveCurrentTransactionKey(ik, ksn);

  const ikExpected = "asKS+qExW02FirOj19WTOg=="; //"6AC292FAA1315B4D858AB3A3D7D5933A";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";
  
  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => new Uint8Array(ik).toBase64() === ikExpected,
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => new Uint8Array(ck).toBase64() === ckExpected,
  });
}

export default async function () {
  await test()
}
