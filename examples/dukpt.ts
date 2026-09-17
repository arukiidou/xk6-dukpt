import { check } from "k6";
import { deriveCurrentTransactionKey } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  // this example uses well-known BDK, but you should replace it with your own.
  example("6AC292FAA1315B4D858AB3A3D7D5933A", "FFFF9876543210E00001");
}

function example(ik: string, ksn: string) {
  const ckExpected = "042666B49184CFA368DE9628D0397BC9";
  const ckBase64Expected = "BCZmtJGEz6No3pYo0Dl7yQ==";

  const ck = deriveCurrentTransactionKey(Uint8Array.fromHex(ik), Uint8Array.fromHex(ksn));
  const ckBase64 = new Uint8Array(ck).toBase64();

  check(null, {
    'deriveCurrentTransactionKey(ik, ksn) as hex': () => new Uint8Array(ck).toHex().toUpperCase() === ckExpected,
    'deriveCurrentTransactionKey(ik, ksn) as base64': () => ckBase64 === ckBase64Expected,
    // Uint8Array.fromBase64() converts a base64 key back into a call argument.
    'Uint8Array.fromBase64(ck)': () => Uint8Array.fromBase64(ckBase64).toHex().toUpperCase() === ckExpected,
  });
}
