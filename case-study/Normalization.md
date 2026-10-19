The Passenger Who Doesn't Exist - Solution
SkyStack passengers checked in successfully, held valid boarding passes, and still failed at the gate with PASSENGER NOT FOUND. The pass was real. The gate reader and check-in service computed different pax_key values from the same documented name.

1. Reproduce The Failure
Start with the Boarding Pass Tester.

Simple names pass, but accented names can fail:

ANA                  -> passes
ANA MULLER-SCHMIDT  -> passes
ANÄ MULLER-SCHMIDT -> fails
The failure is not that the boarding pass signature is invalid. The reader accepts the pass, computes a lookup key, and then cannot find that key in the manifest.

2. Compare The Logged Keys
Both services log a pax_key:

checkin-api logs the key it wrote to the manifest.
gate-reader logs the key it computed from the boarding pass field.
For affected passengers, those two 12-character hex strings differ. Nothing crashes and nothing logs an error because each service is internally consistent.

3. Identify The Three Derivation Differences
The key is derived from the passenger name field plus the global salt:

sha256(f"{salt}:{field}")[:12]
The salt is:

gate-7f3a91
The two services disagree in three ways.

Check-in derives the manifest key like this:

NFC -> truncate to 20 characters -> hash without case folding
The gate reader derives the lookup key like this:

NFD -> truncate to 20 bytes -> uppercase -> hash
All three differences matter:

Normalization: NFC and NFD can render the same glyph but produce different bytes.
Truncation unit: the vendor field is 20 bytes, not 20 Python characters.
Case folding: the reader uppercases before hashing.
4. Ignore The Visible Off-By-One Decoy
For the reference passenger, the displayed fields look like a simple truncation bug:

check-in field: Lévi/Shay Sofia SR R
gate field:     LÉVI/SHAY SOFIA SR
That is only the visible part. Fixing truncation alone still produces the wrong key because the normalization and case-folding steps also affect the hash input.

5. Reconstruct The Gate Algorithm
Use the vendor strings, crash-core excerpt, and logs to reconstruct the reader side:

import hashlib
import unicodedata

SALT = "gate-7f3a91"

def gate_pax_key(documented_name: str) -> str:
    nfd_name = unicodedata.normalize("NFD", documented_name)
    field_bytes = nfd_name.encode("utf-8")[:20]
    field = field_bytes.decode("utf-8", errors="ignore").upper()
    return hashlib.sha256(f"{SALT}:{field}".encode("utf-8")).hexdigest()[:12]
The exact order is important: normalize first, truncate bytes second, uppercase third, hash last.

6. Calibrate With The Oracle Passenger
Before submitting your own key, validate the implementation against the one passenger whose raw documented name and reader pax_key are both visible:

pax_4417
documented_name = İlgaz/Zeynep Meriç
gate field = İLGAZ/ZEYNEP MERIC
expected gate pax_key = 25d340322140
If your implementation does not produce:

25d340322140
do not submit yet. You are missing one of the three levers or applying them in the wrong order.

7. Compute Your Passenger Key
Run the verified algorithm on your own documented name from the affected-passenger manifest.

For the worked reference name:

Documented name: Lévi/Shay Sofia SR Renee
Gate field:     LÉVI/SHAY SOFIA SR
pax_key:        d38a63a53199
The submitted flag shape is:

SKY{d38a63a53199}
Your own flag uses your own computed 12-hex key.

8. Root Cause
The root cause was not a bad pass, bad signature, scanner outage, or missing passenger record. It was a lookup-key derivation split between two systems:

checkin-api wrote manifest rows under an NFC, character-truncated, mixed-case key.
The vendor gate-reader looked up passengers under an NFD, byte-truncated, uppercased key.
For names with diacritics inside the first 20 bytes, those inputs diverged. The reader computed a valid but absent key and reported PASSENGER NOT FOUND.

The challenge has no deploy step. The proof is reproducing the reader's key and submitting it as SKY{<12 hex>}.


The exact probability depends on the system's passenger data, supported languages, field format, and testing. We can't responsibly assign a percentage from this scenario alone.

## 1. How likely are the different failure modes?

| Scenario                                                     | Production likelihood     | Why                                                                         |
| ------------------------------------------------------------ | ------------------------- | --------------------------------------------------------------------------- |
| ASCII-only names                                             | Low for this specific bug | Most ordinary English names won't expose Unicode normalization differences. |
| Accented or international names                              | Plausible                 | Different Unicode representations and casing rules can expose the mismatch. |
| Names near the 20-byte boundary                              | Plausible                 | Character-based and byte-based truncation can diverge.                      |
| Multiple teams implementing the same algorithm independently | A real risk factor        | Small differences in assumptions accumulate across services.                |
| Different language runtimes or vendor libraries              | Depends on implementation | Unicode and string-processing behavior may differ.                          |

The important distinction is that Unicode itself isn't unreliable. The failure occurs because two systems have different rules for turning Unicode text into an identifier.

## 2. Why this can become a serious incident

Imagine the mismatch affects only 0.1% of passengers. That sounds small, but the absolute impact depends on traffic.

Illustrative scenarios, not measured industry failure rates.

| Passengers processed | At 0.1% failures | At 1% failures |
| -------------------- | ---------------- | -------------- |
| 1,000                | 1                | 10             |
| 10,000               | 10               | 100            |
| 100,000              | 100              | 1,000          |

At an airport, even a small number of failures could create queues, require manual verification, and delay boarding. In other systems, the same class of bug could prevent a payment from matching an account, cause a user lookup to fail, or split records that should belong to one identity.

Airline passenger-name processing also has to deal with transliteration and constrained character fields; international travel-document standards explicitly address these issues.&#x20;

[image](https://www.google.com/s2/favicons?domain=https://www.icao.int\&sz=32)

ICAO

+1



## 3. How would you prevent this in a real production system?

I'd prioritize these four measures:

1. One canonicalization contract: both services must use the same documented normalization, casing, encoding, and truncation rules.
2. Shared test vectors: include accented names, decomposed Unicode, non-Latin scripts, and inputs exactly at the byte limit.
3. Integration tests: generate the key in check-in, then verify that the gate-side implementation retrieves the same manifest entry.
4. Operational fallback: when a key lookup fails, support safe secondary verification and alerting instead of treating the passenger as nonexistent.

Unicode's own guidance emphasizes normalization when equivalent strings need consistent binary representations.&#x20;

[image](https://www.google.com/s2/favicons?domain=https://www.unicode.org\&sz=32)

Normalization

+1



## My engineering assessment

The most concerning aspect isn't the use of SHA-256. It's that two independent services implement identity derivation separately, and their tests apparently don't cover the boundary cases.

This is a classic integration failure: each service can pass its own unit tests while the overall system fails in production.

If you encounter this in a design or production review, I'd ask one question first: Where is the canonical identity contract defined, and how do we prove every producer and consumer follows it? That question gets closer to the root cause than simply asking which line of code is wrong.