import { check } from "k6";
import encoding from 'k6/encoding';
import { derivationOfInitialKey, deriveCurrentTransactionKey, encryptPin } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default async function () {
  const ikExpected = "asKS+qExW02FirOj19WTOg==";//"0123456789ABCDEFFEDCBA9876543210";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";
  const pinEncExpected = "G5wYReuZOno="; //"1B9C1845EB993A7A";

  // TODO: wait for merge Uint8Array.fromBase64() to sobek
  const bdk = encoding.b64decode("ASNFZ4mrze/+3LqYdlQyEA=="); //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = encoding.b64decode("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  const ik = derivationOfInitialKey(bdk, ksn);
  const ck = deriveCurrentTransactionKey(ik, ksn);

  check(null, {
    'derivationOfInitialKey(bdk, ksn)': () => encoding.b64encode(ik) === ikExpected,
    'deriveCurrentTransactionKey(ik, ksn)': () => encoding.b64encode(ck) === ckExpected,
  });

  const pin = "1234";
  const pan = "4012345678909";
  const format = "ISO-0";

  const pinEnc = encryptPin(ck, pin, pan, format);

  check(null, {
    'encryptPin(ck, pin, pan, format)': () => encoding.b64encode(pinEnc) === pinEncExpected,
  });
}
