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
}

export const save = db.save;
export const load = db.load;
export const contractCall = contract.call;
