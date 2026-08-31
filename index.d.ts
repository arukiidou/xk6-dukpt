/**
 * **Example k6 extension**
 *
 * @module example
 */
export as namespace example;

/**
 * Generate a personalized greeting message.
 *
 * @param name The name to be greeted (default "World")
 */
export declare function greeting(name: string): string

/**
 * Encode the passed `input` string or ArrayBuffer object to base32 encoded string.
 * Available options for encoding parameter are:
 * - **std**: the standard encoding with `=` padding chars. This is the default.
 * - **rawstd**: like **std** but without `=` padding characters
 * - **hex**: the "Extended Hex Alphabet" encoding defined in RFC 4648. It is typically used in DNS.
 * - **rawhex**: like **hex** but without `=` padding characters
 *
 * @param input The input string or ArrayBuffer object to base32 encode.
 * @param encoding The base32 encoding to use.
 */
export declare function b32encode(
  input: string | ArrayBuffer,
  encoding: string
): string;

/**
 * Decode the passed base32 encoded `input` string into the unencoded original input in either binary or string formats.
 * Available options for encoding parameter are:
 * - **std**: the standard encoding with `=` padding chars. This is the default.
 * - **rawstd**: like **std** but without `=` padding characters
 * - **hex**: the "Extended Hex Alphabet" encoding defined in RFC 4648. It is typically used in DNS.
 * - **rawhex**: like **hex** but without `=` padding characters
 *
 * @param input The string to base64 decode.
 * @param encoding The base32 encoding to use.
 * @param format If "s" return the data as a string, otherwise if unspecified an ArrayBuffer object is returned.
 */
export declare function b32decode(
  input: string,
  encoding?: string,
  format?: string
): string | ArrayBuffer;

/**
 * Pseudo random number generator.
 */
export declare class Random {
  /**
   * Initial seed for the generator.
   * If it is changed, the generator will restart.
   */
  seed: number;

  /**
   * Create a new Random instance.
   *
   * The seed parameter is optional, if missing a random initial seed will be used.
   *
   * @param seed Initial seed for the generator
   */
  constructor(seed?: Number);

  /**
   * Generate a non-negative pseudo-random integer.
   *
   * The generated integer will be from the half-open interval [0,n).
   * If the optional parameter n is missing, `Number.MAX_SAFE_INTEGER` will be used instead.
   *
   * @param n The upper bound, `Number.MAX_SAFE_INTEGER` if missing.
   */
  int(n: number): number;

  /**
   * Generate a non-negative pseudo-random float.
   *
   * The generated float will be from the half-open interval [0,n).
   * If the optional parameter n is missing, `1` will be used instead.
   *
   * @param n The upper bound, `1` if missing.
   */
  float(n: number): number;
}
