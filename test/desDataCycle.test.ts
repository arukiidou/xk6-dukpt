import { check } from "k6";
import { decryptDataAsHex, deriveCurrentTransactionKey, encryptData } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

// The well-known ANSI X9.24-1 test IPEK and KSN, so the BDK is never used.
const ik = Uint8Array.fromHex("6AC292FAA1315B4D858AB3A3D7D5933A");
const ksn = Uint8Array.fromHex("FFFF9876543210E00001");

// Multi-byte UTF-8 (3 byte kana, 4 byte emoji) whose length is not a multiple of 8.
const plainText = "DUKPT サイクルテスト 🔐 4012345678909D987";

export default function () {
  const ck = deriveCurrentTransactionKey(ik, ksn);
  const ckHex = new Uint8Array(ck).toHex();

  for (const action of ["request", "response"] as const) {
    // UTF-8 text -> Uint8Array
    const data = new TextEncoder().encode(plainText);
    // encrypt with the raw API
    const encrypted = encryptData(ck, null, data, action);
    // ArrayBuffer -> hex
    const encryptedHex = new Uint8Array(encrypted).toHex();
    // decrypt with the hex API
    const decryptedHex = decryptDataAsHex(ckHex, encryptedHex, "", action);
    // hex -> Uint8Array, zero padded to a multiple of 8 bytes
    const decrypted = Uint8Array.fromHex(decryptedHex);
    // Uint8Array -> UTF-8 text, cut at the original length to drop the padding
    const result = new TextDecoder().decode(decrypted.subarray(0, data.length));

    check(null, {
      [`encryptData(ck, null, data, "${action}") pads to a multiple of 8 bytes`]: () =>
        encrypted.byteLength % 8 === 0 && encrypted.byteLength >= data.length,
      [`decryptDataAsHex(ck, encrypted, "", "${action}") zero pads the plain text`]: () =>
        decrypted.length === encrypted.byteLength && decrypted.subarray(data.length).every((b) => b === 0),
      [`"${action}" round trips the UTF-8 plain text`]: () => result === plainText,
    });
  }
}
