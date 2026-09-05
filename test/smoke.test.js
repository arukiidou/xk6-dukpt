import dukpt from "./dukpt.test.js"

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  dukpt()
}
