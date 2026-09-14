/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * k6 API for [pkg.GetDesTcFromKsn] port from moov-io,
 * gets the transaction counter from the key serial number.
 * KSN length is not validated.
 * @param ksn The key serial number as a base64 string.
 * @returns The 21 bit transaction counter.
 */
export declare function getDesTcFromKsnAsBase64(ksn: string): number;

/**
 * k6 API for [pkg.GenerateNextDesKsn] port from moov-io,
 * generates the next key serial number as a base64 encoded string.
 * @param ksn The key serial number as a base64 string.
 * @returns The next key serial number as a base64 string.
 */
export declare function generateNextDesKsnAsBase64(ksn: string): string;
