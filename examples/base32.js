import { b32encode } from "k6/x/example";

export default function () {
  console.log(b32encode("Hello, World!"))
}
