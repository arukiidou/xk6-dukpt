import desBase64Test from "./desBase64.test.js"
import desHexTest from "./desHex.test.js"
import desTest from "./des.test.js"
import dukptBase64Test from "./dukptBase64.test.js"
import dukptHexTest from "./dukptHex.test.js"
import dukptTest from "./dukpt.test.js"

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  desBase64Test()
  desHexTest()
  desTest()
  dukptBase64Test()
  dukptHexTest()
  dukptTest()
}
