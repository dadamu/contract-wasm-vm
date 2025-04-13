// The entry file of your WebAssembly module.
import * as runtime from "./runtime";

export function addOne(state: ArrayBuffer, args: ArrayBuffer): void {
  const id = "test";
  const idBuffer = String.UTF8.encode(id);

  // Load and parse the current count
  const count = runtime.load(idBuffer);
  let newCount = 1;
  if (count.byteLength !== 0) {
    const countView = new DataView(count);
    newCount = countView.getInt32(0, true) + 1;
  }

  // Save the new count
  const newCountBuffer = new ArrayBuffer(4);
  const newCountView = new DataView(newCountBuffer);
  newCountView.setInt32(0, newCount, true);
  runtime.save(idBuffer, newCountBuffer);
}

export function crash(state: ArrayBuffer, args: ArrayBuffer): void {
  throw new Error("crash");
}

export function callback(state: ArrayBuffer, args: ArrayBuffer): void {
  const contractIdBuffer = String.UTF8.encode("contract");
  const methodBuffer = String.UTF8.encode("method");
  const argsBuffer = String.UTF8.encode("args");
  runtime.contractCall(contractIdBuffer, methodBuffer, argsBuffer);
}

export function infiniteLoop(state: ArrayBuffer, args: ArrayBuffer): void {
  while (true) {
    // Do nothing
  }
}
