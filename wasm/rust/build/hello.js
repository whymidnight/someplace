let wasm_bindgen;
(function() {
    const __exports = {};
    let wasm;

    const cachedTextDecoder = new TextDecoder('utf-8', { ignoreBOM: true, fatal: true });

    cachedTextDecoder.decode();

    let cachegetUint8Memory0 = null;
    function getUint8Memory0() {
        if (cachegetUint8Memory0 === null || cachegetUint8Memory0.buffer !== wasm.memory.buffer) {
            cachegetUint8Memory0 = new Uint8Array(wasm.memory.buffer);
        }
        return cachegetUint8Memory0;
    }

    function getStringFromWasm0(ptr, len) {
        return cachedTextDecoder.decode(getUint8Memory0().subarray(ptr, ptr + len));
    }

    const heap = new Array(32).fill(undefined);

    heap.push(undefined, null, true, false);

    let heap_next = heap.length;

    function addHeapObject(obj) {
        if (heap_next === heap.length) heap.push(heap.length + 1);
        const idx = heap_next;
        heap_next = heap[idx];

        heap[idx] = obj;
        return idx;
    }

function getObject(idx) { return heap[idx]; }

let WASM_VECTOR_LEN = 0;

const cachedTextEncoder = new TextEncoder('utf-8');

const encodeString = (typeof cachedTextEncoder.encodeInto === 'function'
    ? function (arg, view) {
    return cachedTextEncoder.encodeInto(arg, view);
}
    : function (arg, view) {
    const buf = cachedTextEncoder.encode(arg);
    view.set(buf);
    return {
        read: arg.length,
        written: buf.length
    };
});

function passStringToWasm0(arg, malloc, realloc) {

    if (realloc === undefined) {
        const buf = cachedTextEncoder.encode(arg);
        const ptr = malloc(buf.length);
        getUint8Memory0().subarray(ptr, ptr + buf.length).set(buf);
        WASM_VECTOR_LEN = buf.length;
        return ptr;
    }

    let len = arg.length;
    let ptr = malloc(len);

    const mem = getUint8Memory0();

    let offset = 0;

    for (; offset < len; offset++) {
        const code = arg.charCodeAt(offset);
        if (code > 0x7F) break;
        mem[ptr + offset] = code;
    }

    if (offset !== len) {
        if (offset !== 0) {
            arg = arg.slice(offset);
        }
        ptr = realloc(ptr, len, len = offset + arg.length * 3);
        const view = getUint8Memory0().subarray(ptr + offset, ptr + len);
        const ret = encodeString(arg, view);

        offset += ret.written;
    }

    WASM_VECTOR_LEN = offset;
    return ptr;
}

function isLikeNone(x) {
    return x === undefined || x === null;
}

let cachegetInt32Memory0 = null;
function getInt32Memory0() {
    if (cachegetInt32Memory0 === null || cachegetInt32Memory0.buffer !== wasm.memory.buffer) {
        cachegetInt32Memory0 = new Int32Array(wasm.memory.buffer);
    }
    return cachegetInt32Memory0;
}

function dropObject(idx) {
    if (idx < 36) return;
    heap[idx] = heap_next;
    heap_next = idx;
}

function takeObject(idx) {
    const ret = getObject(idx);
    dropObject(idx);
    return ret;
}

let cachegetFloat64Memory0 = null;
function getFloat64Memory0() {
    if (cachegetFloat64Memory0 === null || cachegetFloat64Memory0.buffer !== wasm.memory.buffer) {
        cachegetFloat64Memory0 = new Float64Array(wasm.memory.buffer);
    }
    return cachegetFloat64Memory0;
}

function debugString(val) {
    // primitive types
    const type = typeof val;
    if (type == 'number' || type == 'boolean' || val == null) {
        return  `${val}`;
    }
    if (type == 'string') {
        return `"${val}"`;
    }
    if (type == 'symbol') {
        const description = val.description;
        if (description == null) {
            return 'Symbol';
        } else {
            return `Symbol(${description})`;
        }
    }
    if (type == 'function') {
        const name = val.name;
        if (typeof name == 'string' && name.length > 0) {
            return `Function(${name})`;
        } else {
            return 'Function';
        }
    }
    // objects
    if (Array.isArray(val)) {
        const length = val.length;
        let debug = '[';
        if (length > 0) {
            debug += debugString(val[0]);
        }
        for(let i = 1; i < length; i++) {
            debug += ', ' + debugString(val[i]);
        }
        debug += ']';
        return debug;
    }
    // Test for built-in
    const builtInMatches = /\[object ([^\]]+)\]/.exec(toString.call(val));
    let className;
    if (builtInMatches.length > 1) {
        className = builtInMatches[1];
    } else {
        // Failed to match the standard '[object ClassName]'
        return toString.call(val);
    }
    if (className == 'Object') {
        // we're a user defined class or Object
        // JSON.stringify avoids problems with cycles, and is generally much
        // easier than looping through ownProperties of `val`.
        try {
            return 'Object(' + JSON.stringify(val) + ')';
        } catch (_) {
            return 'Object';
        }
    }
    // errors
    if (val instanceof Error) {
        return `${val.name}: ${val.message}\n${val.stack}`;
    }
    // TODO we could test for more things here, like `Set`s and `Map`s.
    return className;
}

function makeMutClosure(arg0, arg1, dtor, f) {
    const state = { a: arg0, b: arg1, cnt: 1, dtor };
    const real = (...args) => {
        // First up with a closure we increment the internal reference
        // count. This ensures that the Rust closure environment won't
        // be deallocated while we're invoking it.
        state.cnt++;
        const a = state.a;
        state.a = 0;
        try {
            return f(a, state.b, ...args);
        } finally {
            if (--state.cnt === 0) {
                wasm.__wbindgen_export_2.get(state.dtor)(a, state.b);

            } else {
                state.a = a;
            }
        }
    };
    real.original = state;

    return real;
}
function __wbg_adapter_28(arg0, arg1) {
    wasm._dyn_core__ops__function__FnMut_____Output___R_as_wasm_bindgen__closure__WasmClosure___describe__invoke__h2924d9b283b99c0e(arg0, arg1);
}

function __wbg_adapter_31(arg0, arg1, arg2) {
    wasm._dyn_core__ops__function__FnMut__A____Output___R_as_wasm_bindgen__closure__WasmClosure___describe__invoke__h51cf90af102bd259(arg0, arg1, addHeapObject(arg2));
}

/**
* @param {string} candy_machine_id
* @returns {Promise<any>}
*/
__exports.fetch_candies = function(candy_machine_id) {
    const ptr0 = passStringToWasm0(candy_machine_id, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ret = wasm.fetch_candies(ptr0, len0);
    return takeObject(ret);
};

function handleError(f, args) {
    try {
        return f.apply(this, args);
    } catch (e) {
        wasm.__wbindgen_exn_store(addHeapObject(e));
    }
}

function notDefined(what) { return () => { throw new Error(`${what} is not defined`); }; }

function _assertClass(instance, klass) {
    if (!(instance instanceof klass)) {
        throw new Error(`expected instance of ${klass.name}`);
    }
    return instance.ptr;
}

function getArrayU8FromWasm0(ptr, len) {
    return getUint8Memory0().subarray(ptr / 1, ptr / 1 + len);
}

let cachegetUint32Memory0 = null;
function getUint32Memory0() {
    if (cachegetUint32Memory0 === null || cachegetUint32Memory0.buffer !== wasm.memory.buffer) {
        cachegetUint32Memory0 = new Uint32Array(wasm.memory.buffer);
    }
    return cachegetUint32Memory0;
}

function passArrayJsValueToWasm0(array, malloc) {
    const ptr = malloc(array.length * 4);
    const mem = getUint32Memory0();
    for (let i = 0; i < array.length; i++) {
        mem[ptr / 4 + i] = addHeapObject(array[i]);
    }
    WASM_VECTOR_LEN = array.length;
    return ptr;
}

const u32CvtShim = new Uint32Array(2);

const uint64CvtShim = new BigUint64Array(u32CvtShim.buffer);
/**
* Initialize Javascript logging and panic handler
*/
__exports.init = function() {
    wasm.init();
};

function __wbg_adapter_140(arg0, arg1, arg2, arg3) {
    wasm.wasm_bindgen__convert__closures__invoke2_mut__h57f69b06e27a4d20(arg0, arg1, addHeapObject(arg2), addHeapObject(arg3));
}

/**
*/
class Hash {

    static __wrap(ptr) {
        const obj = Object.create(Hash.prototype);
        obj.ptr = ptr;

        return obj;
    }

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_hash_free(ptr);
    }
    /**
    * Create a new Hash object
    *
    * * `value` - optional hash as a base58 encoded string, `Uint8Array`, `[number]`
    * @param {any} value
    */
    constructor(value) {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.hash_constructor(retptr, addHeapObject(value));
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var r2 = getInt32Memory0()[retptr / 4 + 2];
            if (r2) {
                throw takeObject(r1);
            }
            return Hash.__wrap(r0);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
    /**
    * Return the base58 string representation of the hash
    * @returns {string}
    */
    toString() {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.hash_toString(retptr, this.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            return getStringFromWasm0(r0, r1);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
            wasm.__wbindgen_free(r0, r1);
        }
    }
    /**
    * Checks if two `Hash`s are equal
    * @param {Hash} other
    * @returns {boolean}
    */
    equals(other) {
        _assertClass(other, Hash);
        const ret = wasm.hash_equals(this.ptr, other.ptr);
        return ret !== 0;
    }
    /**
    * Return the `Uint8Array` representation of the hash
    * @returns {Uint8Array}
    */
    toBytes() {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.hash_toBytes(retptr, this.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var v0 = getArrayU8FromWasm0(r0, r1).slice();
            wasm.__wbindgen_free(r0, r1 * 1);
            return v0;
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
}
__exports.Hash = Hash;
/**
* A directive for a single invocation of a Solana program.
*
* An instruction specifies which program it is calling, which accounts it may
* read or modify, and additional data that serves as input to the program. One
* or more instructions are included in transactions submitted by Solana
* clients. Instructions are also used to describe [cross-program
* invocations][cpi].
*
* [cpi]: https://docs.solana.com/developing/programming-model/calling-between-programs
*
* During execution, a program will receive a list of account data as one of
* its arguments, in the same order as specified during `Instruction`
* construction.
*
* While Solana is agnostic to the format of the instruction data, it has
* built-in support for serialization via [`borsh`] and [`bincode`].
*
* [`borsh`]: https://docs.rs/borsh/latest/borsh/
* [`bincode`]: https://docs.rs/bincode/latest/bincode/
*
* # Specifying account metadata
*
* When constructing an [`Instruction`], a list of all accounts that may be
* read or written during the execution of that instruction must be supplied as
* [`AccountMeta`] values.
*
* Any account whose data may be mutated by the program during execution must
* be specified as writable. During execution, writing to an account that was
* not specified as writable will cause the transaction to fail. Writing to an
* account that is not owned by the program will cause the transaction to fail.
*
* Any account whose lamport balance may be mutated by the program during
* execution must be specified as writable. During execution, mutating the
* lamports of an account that was not specified as writable will cause the
* transaction to fail. While _subtracting_ lamports from an account not owned
* by the program will cause the transaction to fail, _adding_ lamports to any
* account is allowed, as long is it is mutable.
*
* Accounts that are not read or written by the program may still be specified
* in an `Instruction`'s account list. These will affect scheduling of program
* execution by the runtime, but will otherwise be ignored.
*
* When building a transaction, the Solana runtime coalesces all accounts used
* by all instructions in that transaction, along with accounts and permissions
* required by the runtime, into a single account list. Some accounts and
* account permissions required by the runtime to process a transaction are
* _not_ required to be included in an `Instruction`s account list. These
* include:
*
* - The program ID &mdash; it is a separate field of `Instruction`
* - The transaction's fee-paying account &mdash; it is added during [`Message`]
*   construction. A program may still require the fee payer as part of the
*   account list if it directly references it.
*
* [`Message`]: crate::message::Message
*
* Programs may require signatures from some accounts, in which case they
* should be specified as signers during `Instruction` construction. The
* program must still validate during execution that the account is a signer.
*/
class Instruction {

    static __wrap(ptr) {
        const obj = Object.create(Instruction.prototype);
        obj.ptr = ptr;

        return obj;
    }

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_instruction_free(ptr);
    }
}
__exports.Instruction = Instruction;
/**
*/
class Instructions {

    static __wrap(ptr) {
        const obj = Object.create(Instructions.prototype);
        obj.ptr = ptr;

        return obj;
    }

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_instructions_free(ptr);
    }
    /**
    */
    constructor() {
        const ret = wasm.instructions_constructor();
        return Instructions.__wrap(ret);
    }
    /**
    * @param {Instruction} instruction
    */
    push(instruction) {
        _assertClass(instruction, Instruction);
        var ptr0 = instruction.ptr;
        instruction.ptr = 0;
        wasm.instructions_push(this.ptr, ptr0);
    }
}
__exports.Instructions = Instructions;
/**
* A Solana transaction message (legacy).
*
* See the [`message`] module documentation for further description.
*
* [`message`]: crate::message
*
* Some constructors accept an optional `payer`, the account responsible for
* paying the cost of executing a transaction. In most cases, callers should
* specify the payer explicitly in these constructors. In some cases though,
* the caller is not _required_ to specify the payer, but is still allowed to:
* in the `Message` structure, the first account is always the fee-payer, so if
* the caller has knowledge that the first account of the constructed
* transaction's `Message` is both a signer and the expected fee-payer, then
* redundantly specifying the fee-payer is not strictly required.
*/
class Message {

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_message_free(ptr);
    }
    /**
    * The id of a recent ledger entry.
    */
    get recent_blockhash() {
        const ret = wasm.__wbg_get_message_recent_blockhash(this.ptr);
        return Hash.__wrap(ret);
    }
    /**
    * The id of a recent ledger entry.
    * @param {Hash} arg0
    */
    set recent_blockhash(arg0) {
        _assertClass(arg0, Hash);
        var ptr0 = arg0.ptr;
        arg0.ptr = 0;
        wasm.__wbg_set_message_recent_blockhash(this.ptr, ptr0);
    }
}
__exports.Message = Message;
/**
*/
class Pubkey {

    static __wrap(ptr) {
        const obj = Object.create(Pubkey.prototype);
        obj.ptr = ptr;

        return obj;
    }

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_pubkey_free(ptr);
    }
    /**
    * Create a new Pubkey object
    *
    * * `value` - optional public key as a base58 encoded string, `Uint8Array`, `[number]`
    * @param {any} value
    */
    constructor(value) {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.pubkey_constructor(retptr, addHeapObject(value));
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var r2 = getInt32Memory0()[retptr / 4 + 2];
            if (r2) {
                throw takeObject(r1);
            }
            return Pubkey.__wrap(r0);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
    /**
    * Return the base58 string representation of the public key
    * @returns {string}
    */
    toString() {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.pubkey_toString(retptr, this.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            return getStringFromWasm0(r0, r1);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
            wasm.__wbindgen_free(r0, r1);
        }
    }
    /**
    * Check if a `Pubkey` is on the ed25519 curve.
    * @returns {boolean}
    */
    isOnCurve() {
        const ret = wasm.pubkey_isOnCurve(this.ptr);
        return ret !== 0;
    }
    /**
    * Checks if two `Pubkey`s are equal
    * @param {Pubkey} other
    * @returns {boolean}
    */
    equals(other) {
        _assertClass(other, Pubkey);
        const ret = wasm.pubkey_equals(this.ptr, other.ptr);
        return ret !== 0;
    }
    /**
    * Return the `Uint8Array` representation of the public key
    * @returns {Uint8Array}
    */
    toBytes() {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            wasm.pubkey_toBytes(retptr, this.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var v0 = getArrayU8FromWasm0(r0, r1).slice();
            wasm.__wbindgen_free(r0, r1 * 1);
            return v0;
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
    /**
    * Derive a Pubkey from another Pubkey, string seed, and a program id
    * @param {Pubkey} base
    * @param {string} seed
    * @param {Pubkey} owner
    * @returns {Pubkey}
    */
    static createWithSeed(base, seed, owner) {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            _assertClass(base, Pubkey);
            const ptr0 = passStringToWasm0(seed, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
            const len0 = WASM_VECTOR_LEN;
            _assertClass(owner, Pubkey);
            wasm.pubkey_createWithSeed(retptr, base.ptr, ptr0, len0, owner.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var r2 = getInt32Memory0()[retptr / 4 + 2];
            if (r2) {
                throw takeObject(r1);
            }
            return Pubkey.__wrap(r0);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
    /**
    * Derive a program address from seeds and a program id
    * @param {any[]} seeds
    * @param {Pubkey} program_id
    * @returns {Pubkey}
    */
    static createProgramAddress(seeds, program_id) {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            const ptr0 = passArrayJsValueToWasm0(seeds, wasm.__wbindgen_malloc);
            const len0 = WASM_VECTOR_LEN;
            _assertClass(program_id, Pubkey);
            wasm.pubkey_createProgramAddress(retptr, ptr0, len0, program_id.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var r2 = getInt32Memory0()[retptr / 4 + 2];
            if (r2) {
                throw takeObject(r1);
            }
            return Pubkey.__wrap(r0);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
    /**
    * Find a valid program address
    *
    * Returns:
    * * `[PubKey, number]` - the program address and bump seed
    * @param {any[]} seeds
    * @param {Pubkey} program_id
    * @returns {any}
    */
    static findProgramAddress(seeds, program_id) {
        try {
            const retptr = wasm.__wbindgen_add_to_stack_pointer(-16);
            const ptr0 = passArrayJsValueToWasm0(seeds, wasm.__wbindgen_malloc);
            const len0 = WASM_VECTOR_LEN;
            _assertClass(program_id, Pubkey);
            wasm.pubkey_findProgramAddress(retptr, ptr0, len0, program_id.ptr);
            var r0 = getInt32Memory0()[retptr / 4 + 0];
            var r1 = getInt32Memory0()[retptr / 4 + 1];
            var r2 = getInt32Memory0()[retptr / 4 + 2];
            if (r2) {
                throw takeObject(r1);
            }
            return takeObject(r0);
        } finally {
            wasm.__wbindgen_add_to_stack_pointer(16);
        }
    }
}
__exports.Pubkey = Pubkey;

class SystemInstruction {

    __destroy_into_raw() {
        const ptr = this.ptr;
        this.ptr = 0;

        return ptr;
    }

    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_systeminstruction_free(ptr);
    }
    /**
    * @param {Pubkey} from_pubkey
    * @param {Pubkey} to_pubkey
    * @param {BigInt} lamports
    * @param {BigInt} space
    * @param {Pubkey} owner
    * @returns {Instruction}
    */
    static createAccount(from_pubkey, to_pubkey, lamports, space, owner) {
        _assertClass(from_pubkey, Pubkey);
        _assertClass(to_pubkey, Pubkey);
        uint64CvtShim[0] = lamports;
        const low0 = u32CvtShim[0];
        const high0 = u32CvtShim[1];
        uint64CvtShim[0] = space;
        const low1 = u32CvtShim[0];
        const high1 = u32CvtShim[1];
        _assertClass(owner, Pubkey);
        const ret = wasm.systeminstruction_createAccount(from_pubkey.ptr, to_pubkey.ptr, low0, high0, low1, high1, owner.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} from_pubkey
    * @param {Pubkey} to_pubkey
    * @param {Pubkey} base
    * @param {string} seed
    * @param {BigInt} lamports
    * @param {BigInt} space
    * @param {Pubkey} owner
    * @returns {Instruction}
    */
    static createAccountWithSeed(from_pubkey, to_pubkey, base, seed, lamports, space, owner) {
        _assertClass(from_pubkey, Pubkey);
        _assertClass(to_pubkey, Pubkey);
        _assertClass(base, Pubkey);
        const ptr0 = passStringToWasm0(seed, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        uint64CvtShim[0] = lamports;
        const low1 = u32CvtShim[0];
        const high1 = u32CvtShim[1];
        uint64CvtShim[0] = space;
        const low2 = u32CvtShim[0];
        const high2 = u32CvtShim[1];
        _assertClass(owner, Pubkey);
        const ret = wasm.systeminstruction_createAccountWithSeed(from_pubkey.ptr, to_pubkey.ptr, base.ptr, ptr0, len0, low1, high1, low2, high2, owner.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} pubkey
    * @param {Pubkey} owner
    * @returns {Instruction}
    */
    static assign(pubkey, owner) {
        _assertClass(pubkey, Pubkey);
        _assertClass(owner, Pubkey);
        const ret = wasm.systeminstruction_assign(pubkey.ptr, owner.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} pubkey
    * @param {Pubkey} base
    * @param {string} seed
    * @param {Pubkey} owner
    * @returns {Instruction}
    */
    static assignWithSeed(pubkey, base, seed, owner) {
        _assertClass(pubkey, Pubkey);
        _assertClass(base, Pubkey);
        const ptr0 = passStringToWasm0(seed, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        _assertClass(owner, Pubkey);
        const ret = wasm.systeminstruction_assignWithSeed(pubkey.ptr, base.ptr, ptr0, len0, owner.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} from_pubkey
    * @param {Pubkey} to_pubkey
    * @param {BigInt} lamports
    * @returns {Instruction}
    */
    static transfer(from_pubkey, to_pubkey, lamports) {
        _assertClass(from_pubkey, Pubkey);
        _assertClass(to_pubkey, Pubkey);
        uint64CvtShim[0] = lamports;
        const low0 = u32CvtShim[0];
        const high0 = u32CvtShim[1];
        const ret = wasm.systeminstruction_transfer(from_pubkey.ptr, to_pubkey.ptr, low0, high0);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} from_pubkey
    * @param {Pubkey} from_base
    * @param {string} from_seed
    * @param {Pubkey} from_owner
    * @param {Pubkey} to_pubkey
    * @param {BigInt} lamports
    * @returns {Instruction}
    */
    static transferWithSeed(from_pubkey, from_base, from_seed, from_owner, to_pubkey, lamports) {
        _assertClass(from_pubkey, Pubkey);
        _assertClass(from_base, Pubkey);
        const ptr0 = passStringToWasm0(from_seed, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        _assertClass(from_owner, Pubkey);
        _assertClass(to_pubkey, Pubkey);
        uint64CvtShim[0] = lamports;
        const low1 = u32CvtShim[0];
        const high1 = u32CvtShim[1];
        const ret = wasm.systeminstruction_transferWithSeed(from_pubkey.ptr, from_base.ptr, ptr0, len0, from_owner.ptr, to_pubkey.ptr, low1, high1);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} pubkey
    * @param {BigInt} space
    * @returns {Instruction}
    */
    static allocate(pubkey, space) {
        _assertClass(pubkey, Pubkey);
        uint64CvtShim[0] = space;
        const low0 = u32CvtShim[0];
        const high0 = u32CvtShim[1];
        const ret = wasm.systeminstruction_allocate(pubkey.ptr, low0, high0);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} address
    * @param {Pubkey} base
    * @param {string} seed
    * @param {BigInt} space
    * @param {Pubkey} owner
    * @returns {Instruction}
    */
    static allocateWithSeed(address, base, seed, space, owner) {
        _assertClass(address, Pubkey);
        _assertClass(base, Pubkey);
        const ptr0 = passStringToWasm0(seed, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        uint64CvtShim[0] = space;
        const low1 = u32CvtShim[0];
        const high1 = u32CvtShim[1];
        _assertClass(owner, Pubkey);
        const ret = wasm.systeminstruction_allocateWithSeed(address.ptr, base.ptr, ptr0, len0, low1, high1, owner.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} from_pubkey
    * @param {Pubkey} nonce_pubkey
    * @param {Pubkey} authority
    * @param {BigInt} lamports
    * @returns {Array<any>}
    */
    static createNonceAccount(from_pubkey, nonce_pubkey, authority, lamports) {
        _assertClass(from_pubkey, Pubkey);
        _assertClass(nonce_pubkey, Pubkey);
        _assertClass(authority, Pubkey);
        uint64CvtShim[0] = lamports;
        const low0 = u32CvtShim[0];
        const high0 = u32CvtShim[1];
        const ret = wasm.systeminstruction_createNonceAccount(from_pubkey.ptr, nonce_pubkey.ptr, authority.ptr, low0, high0);
        return takeObject(ret);
    }
    /**
    * @param {Pubkey} nonce_pubkey
    * @param {Pubkey} authorized_pubkey
    * @returns {Instruction}
    */
    static advanceNonceAccount(nonce_pubkey, authorized_pubkey) {
        _assertClass(nonce_pubkey, Pubkey);
        _assertClass(authorized_pubkey, Pubkey);
        const ret = wasm.systeminstruction_advanceNonceAccount(nonce_pubkey.ptr, authorized_pubkey.ptr);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} nonce_pubkey
    * @param {Pubkey} authorized_pubkey
    * @param {Pubkey} to_pubkey
    * @param {BigInt} lamports
    * @returns {Instruction}
    */
    static withdrawNonceAccount(nonce_pubkey, authorized_pubkey, to_pubkey, lamports) {
        _assertClass(nonce_pubkey, Pubkey);
        _assertClass(authorized_pubkey, Pubkey);
        _assertClass(to_pubkey, Pubkey);
        uint64CvtShim[0] = lamports;
        const low0 = u32CvtShim[0];
        const high0 = u32CvtShim[1];
        const ret = wasm.systeminstruction_withdrawNonceAccount(nonce_pubkey.ptr, authorized_pubkey.ptr, to_pubkey.ptr, low0, high0);
        return Instruction.__wrap(ret);
    }
    /**
    * @param {Pubkey} nonce_pubkey
    * @param {Pubkey} authorized_pubkey
    * @param {Pubkey} new_authority
    * @returns {Instruction}
    */
    static authorizeNonceAccount(nonce_pubkey, authorized_pubkey, new_authority) {
        _assertClass(nonce_pubkey, Pubkey);
        _assertClass(authorized_pubkey, Pubkey);
        _assertClass(new_authority, Pubkey);
        const ret = wasm.systeminstruction_authorizeNonceAccount(nonce_pubkey.ptr, authorized_pubkey.ptr, new_authority.ptr);
        return Instruction.__wrap(ret);
    }
}
__exports.SystemInstruction = SystemInstruction;

async function load(module, imports) {
    if (typeof Response === 'function' && module instanceof Response) {
        if (typeof WebAssembly.instantiateStreaming === 'function') {
            try {
                return await WebAssembly.instantiateStreaming(module, imports);

            } catch (e) {
                if (module.headers.get('Content-Type') != 'application/wasm') {
                    console.warn("`WebAssembly.instantiateStreaming` failed because your server does not serve wasm with `application/wasm` MIME type. Falling back to `WebAssembly.instantiate` which is slower. Original error:\n", e);

                } else {
                    throw e;
                }
            }
        }

        const bytes = await module.arrayBuffer();
        return await WebAssembly.instantiate(bytes, imports);

    } else {
        const instance = await WebAssembly.instantiate(module, imports);

        if (instance instanceof WebAssembly.Instance) {
            return { instance, module };

        } else {
            return instance;
        }
    }
}

async function init(input) {
    if (typeof input === 'undefined') {
        let src;
        if (typeof document === 'undefined') {
            src = location.href;
        } else {
            src = document.currentScript.src;
        }
        input = src.replace(/\.js$/, '_bg.wasm');
    }
    const imports = {};
    imports.wbg = {};
    imports.wbg.__wbindgen_string_new = function(arg0, arg1) {
        const ret = getStringFromWasm0(arg0, arg1);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_Window_424f0a27f77b8fdd = function(arg0) {
        const ret = getObject(arg0).Window;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_WorkerGlobalScope_212ebfcc0586c20e = function(arg0) {
        const ret = getObject(arg0).WorkerGlobalScope;
        return addHeapObject(ret);
    };
    imports.wbg.__wbindgen_string_get = function(arg0, arg1) {
        const obj = getObject(arg1);
        const ret = typeof(obj) === 'string' ? obj : undefined;
        var ptr0 = isLikeNone(ret) ? 0 : passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        var len0 = WASM_VECTOR_LEN;
        getInt32Memory0()[arg0 / 4 + 1] = len0;
        getInt32Memory0()[arg0 / 4 + 0] = ptr0;
    };
    imports.wbg.__wbindgen_is_undefined = function(arg0) {
        const ret = getObject(arg0) === undefined;
        return ret;
    };
    imports.wbg.__wbindgen_cb_drop = function(arg0) {
        const obj = takeObject(arg0).original;
        if (obj.cnt-- == 1) {
            obj.a = 0;
            return true;
        }
        const ret = false;
        return ret;
    };
    imports.wbg.__wbg_setTimeout_290c28f3580809b6 = function() { return handleError(function (arg0, arg1) {
        const ret = setTimeout(getObject(arg0), arg1);
        return ret;
    }, arguments) };
    imports.wbg.__wbg_clearTimeout_d8b36ad8fa330187 = typeof clearTimeout == 'function' ? clearTimeout : notDefined('clearTimeout');
    imports.wbg.__wbindgen_is_object = function(arg0) {
        const val = getObject(arg0);
        const ret = typeof(val) === 'object' && val !== null;
        return ret;
    };
    imports.wbg.__wbindgen_number_new = function(arg0) {
        const ret = arg0;
        return addHeapObject(ret);
    };
    imports.wbg.__wbindgen_number_get = function(arg0, arg1) {
        const obj = getObject(arg1);
        const ret = typeof(obj) === 'number' ? obj : undefined;
        getFloat64Memory0()[arg0 / 8 + 1] = isLikeNone(ret) ? 0 : ret;
        getInt32Memory0()[arg0 / 4 + 0] = !isLikeNone(ret);
    };
    imports.wbg.__wbg_instruction_new = function(arg0) {
        const ret = Instruction.__wrap(arg0);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_pubkey_new = function(arg0) {
        const ret = Pubkey.__wrap(arg0);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_fetch_ef7a6623d1fcd3b8 = function(arg0, arg1) {
        const ret = getObject(arg0).fetch(getObject(arg1));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_debug_fda1f49ea6af7a1d = function(arg0) {
        console.debug(getObject(arg0));
    };
    imports.wbg.__wbg_error_8ff19d586a987aef = function(arg0) {
        console.error(getObject(arg0));
    };
    imports.wbg.__wbg_info_c8f1b00be4ef10bc = function(arg0) {
        console.info(getObject(arg0));
    };
    imports.wbg.__wbg_log_e8ba7b992c7ad0eb = function(arg0) {
        console.log(getObject(arg0));
    };
    imports.wbg.__wbg_warn_0227db1aa6989248 = function(arg0) {
        console.warn(getObject(arg0));
    };
    imports.wbg.__wbg_set_6884dcc6cdd65022 = function() { return handleError(function (arg0, arg1, arg2, arg3, arg4) {
        getObject(arg0).set(getStringFromWasm0(arg1, arg2), getStringFromWasm0(arg3, arg4));
    }, arguments) };
    imports.wbg.__wbg_fetch_bf56e2a9f0644e3f = function(arg0, arg1) {
        const ret = getObject(arg0).fetch(getObject(arg1));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_headers_0a71906114661592 = function(arg0) {
        const ret = getObject(arg0).headers;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_newwithstrandinit_fd99688f189f053e = function() { return handleError(function (arg0, arg1, arg2) {
        const ret = new Request(getStringFromWasm0(arg0, arg1), getObject(arg2));
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_instanceof_Response_ccfeb62399355bcd = function(arg0) {
        const ret = getObject(arg0) instanceof Response;
        return ret;
    };
    imports.wbg.__wbg_status_600fd8b881393898 = function(arg0) {
        const ret = getObject(arg0).status;
        return ret;
    };
    imports.wbg.__wbg_headers_9e7f2c05a9b962ea = function(arg0) {
        const ret = getObject(arg0).headers;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_arrayBuffer_5a99283a3954c850 = function() { return handleError(function (arg0) {
        const ret = getObject(arg0).arrayBuffer();
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbindgen_object_clone_ref = function(arg0) {
        const ret = getObject(arg0);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_error_09919627ac0992f5 = function(arg0, arg1) {
        try {
            console.error(getStringFromWasm0(arg0, arg1));
        } finally {
            wasm.__wbindgen_free(arg0, arg1);
        }
    };
    imports.wbg.__wbg_new_693216e109162396 = function() {
        const ret = new Error();
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_stack_0ddaca5d1abfb52f = function(arg0, arg1) {
        const ret = getObject(arg1).stack;
        const ptr0 = passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        getInt32Memory0()[arg0 / 4 + 1] = len0;
        getInt32Memory0()[arg0 / 4 + 0] = ptr0;
    };
    imports.wbg.__wbindgen_object_drop_ref = function(arg0) {
        takeObject(arg0);
    };
    imports.wbg.__wbg_new_94fb1279cf6afea5 = function() {
        const ret = new Array();
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_newwithlength_e80fb11cf19c1628 = function(arg0) {
        const ret = new Array(arg0 >>> 0);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_set_561aac756158708c = function(arg0, arg1, arg2) {
        getObject(arg0)[arg1 >>> 0] = takeObject(arg2);
    };
    imports.wbg.__wbg_isArray_6721f2e508996340 = function(arg0) {
        const ret = Array.isArray(getObject(arg0));
        return ret;
    };
    imports.wbg.__wbg_push_40c6a90f1805aa90 = function(arg0, arg1) {
        const ret = getObject(arg0).push(getObject(arg1));
        return ret;
    };
    imports.wbg.__wbg_instanceof_ArrayBuffer_b81b40c2ae0ab898 = function(arg0) {
        const ret = getObject(arg0) instanceof ArrayBuffer;
        return ret;
    };
    imports.wbg.__wbg_values_b1b9e8c63dbe01c2 = function(arg0) {
        const ret = getObject(arg0).values();
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_newnoargs_e23b458e372830de = function(arg0, arg1) {
        const ret = new Function(getStringFromWasm0(arg0, arg1));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_call_ae78342adc33730a = function() { return handleError(function (arg0, arg1) {
        const ret = getObject(arg0).call(getObject(arg1));
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_call_3ed288a247f13ea5 = function() { return handleError(function (arg0, arg1, arg2) {
        const ret = getObject(arg0).call(getObject(arg1), getObject(arg2));
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_next_bf3d83fc18df496e = function() { return handleError(function (arg0) {
        const ret = getObject(arg0).next();
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_next_cabb70b365520721 = function(arg0) {
        const ret = getObject(arg0).next;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_done_040f966faa9a72b3 = function(arg0) {
        const ret = getObject(arg0).done;
        return ret;
    };
    imports.wbg.__wbg_value_419afbd9b9574c4c = function(arg0) {
        const ret = getObject(arg0).value;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_new_36359baae5a47e27 = function() {
        const ret = new Object();
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_iterator_4832ef1f15b0382b = function() {
        const ret = Symbol.iterator;
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_new_37705eed627d5ed9 = function(arg0, arg1) {
        try {
            var state0 = {a: arg0, b: arg1};
            var cb0 = (arg0, arg1) => {
                const a = state0.a;
                state0.a = 0;
                try {
                    return __wbg_adapter_140(a, state0.b, arg0, arg1);
                } finally {
                    state0.a = a;
                }
            };
            const ret = new Promise(cb0);
            return addHeapObject(ret);
        } finally {
            state0.a = state0.b = 0;
        }
    };
    imports.wbg.__wbg_resolve_a9a87bdd64e9e62c = function(arg0) {
        const ret = Promise.resolve(getObject(arg0));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_then_ce526c837d07b68f = function(arg0, arg1) {
        const ret = getObject(arg0).then(getObject(arg1));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_then_842e65b843962f56 = function(arg0, arg1, arg2) {
        const ret = getObject(arg0).then(getObject(arg1), getObject(arg2));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_globalThis_8e275ef40caea3a3 = function() { return handleError(function () {
        const ret = globalThis.globalThis;
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_self_99737b4dcdf6f0d8 = function() { return handleError(function () {
        const ret = self.self;
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_window_9b61fbbf3564c4fb = function() { return handleError(function () {
        const ret = window.window;
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_global_5de1e0f82bddcd27 = function() { return handleError(function () {
        const ret = global.global;
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_instanceof_Uint8Array_edb92795fc0c63b4 = function(arg0) {
        const ret = getObject(arg0) instanceof Uint8Array;
        return ret;
    };
    imports.wbg.__wbg_new_cc9018bd6f283b6f = function(arg0) {
        const ret = new Uint8Array(getObject(arg0));
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_newwithbyteoffsetandlength_ce1e75f0ce5f7974 = function(arg0, arg1, arg2) {
        const ret = new Uint8Array(getObject(arg0), arg1 >>> 0, arg2 >>> 0);
        return addHeapObject(ret);
    };
    imports.wbg.__wbg_length_0acb1cf9bbaf8519 = function(arg0) {
        const ret = getObject(arg0).length;
        return ret;
    };
    imports.wbg.__wbg_set_f25e869e4565d2a2 = function(arg0, arg1, arg2) {
        getObject(arg0).set(getObject(arg1), arg2 >>> 0);
    };
    imports.wbg.__wbg_get_a9cab131e3152c49 = function() { return handleError(function (arg0, arg1) {
        const ret = Reflect.get(getObject(arg0), getObject(arg1));
        return addHeapObject(ret);
    }, arguments) };
    imports.wbg.__wbg_set_93b1c87ee2af852e = function() { return handleError(function (arg0, arg1, arg2) {
        const ret = Reflect.set(getObject(arg0), getObject(arg1), getObject(arg2));
        return ret;
    }, arguments) };
    imports.wbg.__wbindgen_is_function = function(arg0) {
        const ret = typeof(getObject(arg0)) === 'function';
        return ret;
    };
    imports.wbg.__wbg_buffer_7af23f65f6c64548 = function(arg0) {
        const ret = getObject(arg0).buffer;
        return addHeapObject(ret);
    };
    imports.wbg.__wbindgen_debug_string = function(arg0, arg1) {
        const ret = debugString(getObject(arg1));
        const ptr0 = passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        getInt32Memory0()[arg0 / 4 + 1] = len0;
        getInt32Memory0()[arg0 / 4 + 0] = ptr0;
    };
    imports.wbg.__wbindgen_throw = function(arg0, arg1) {
        throw new Error(getStringFromWasm0(arg0, arg1));
    };
    imports.wbg.__wbindgen_memory = function() {
        const ret = wasm.memory;
        return addHeapObject(ret);
    };
    imports.wbg.__wbindgen_closure_wrapper1401 = function(arg0, arg1, arg2) {
        const ret = makeMutClosure(arg0, arg1, 94, __wbg_adapter_28);
        return addHeapObject(ret);
    };
    imports.wbg.__wbindgen_closure_wrapper8808 = function(arg0, arg1, arg2) {
        const ret = makeMutClosure(arg0, arg1, 369, __wbg_adapter_31);
        return addHeapObject(ret);
    };

    if (typeof input === 'string' || (typeof Request === 'function' && input instanceof Request) || (typeof URL === 'function' && input instanceof URL)) {
        input = fetch(input);
    }



    const { instance, module } = await load(await input, imports);

    wasm = instance.exports;
    init.__wbindgen_wasm_module = module;

    return wasm;
}

wasm_bindgen = Object.assign(init, __exports);

})();
