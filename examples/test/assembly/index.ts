// The entry file of your WebAssembly module.
import * as runtime from "./runtime";

export function init(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {}

export function addOne(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  // Load and parse the current count
  const count = runtime.load("test");
  let newCount = 1;
  if (count.byteLength !== 0) {
    const countView = new DataView(count);
    newCount = countView.getInt32(0, true) + 1;
  }

  // Save the new count
  const newCountBuffer = new ArrayBuffer(4);
  const newCountView = new DataView(newCountBuffer);
  newCountView.setInt32(0, newCount, true);
  runtime.save("test", newCountBuffer);
}

export function crash(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  throw new Error("crash");
}

export function callback(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  const argsBuffer = String.UTF8.encode("args");
  runtime.contractCall("contract", "method", argsBuffer);
}

export function infiniteLoop(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  while (true) {
    // Do nothing
  }
}

export function createContract(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  const argsBuffer = String.UTF8.encode("args");
  runtime.createContract(1, argsBuffer);
}

export function emitEvent(
  state: ArrayBuffer,
  sender: ArrayBuffer,
  args: ArrayBuffer
): void {
  runtime.emitEvent("event", "data");
}
