import { check } from "k6";
import { getDesTcFromKsnAsHex, generateNextDesKsnAsHex } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ksn = "FFFF9876543210E00001";

  check(null, {
    'getDesTcFromKsnAsHex(ksn)': () => getDesTcFromKsnAsHex(ksn) === 1,
    'generateNextDesKsnAsHex(ksn)': () => generateNextDesKsnAsHex(ksn) === "FFFF9876543210E00002",
  });
}
