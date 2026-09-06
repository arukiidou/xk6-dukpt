import base64 from "./dukptBase64.test.js"
import dukpt from "./dukpt.test.js"

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  base64()
  dukpt()
}
