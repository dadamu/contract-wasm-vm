namespace db {
  export declare function save(key: ArrayBuffer, value: ArrayBuffer): void;

  export declare function load(key: ArrayBuffer): ArrayBuffer;
}

namespace contract {
  export declare function call(
    id: ArrayBuffer,
    method: ArrayBuffer,
    args: ArrayBuffer
  ): void;

  export declare function create(
    codeId: i32,
    initArgs: ArrayBuffer
  ): ArrayBuffer;
}

namespace event {
  export declare function emit(event: string, data: string): void;
}

export const save = db.save;
export const load = db.load;
export const contractCall = contract.call;
export const createContract = contract.create;
export const emitEvent = event.emit;
