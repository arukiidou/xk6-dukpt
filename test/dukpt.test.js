import { check } from "k6";
import { getDesTcFromKsn } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  check(null, {
    'getDesTcFromKsn(ksn)': () => getDesTcFromKsn(ksn) === 1,
  });
}
