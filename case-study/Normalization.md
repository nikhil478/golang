
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